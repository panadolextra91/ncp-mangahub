package http_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/internal/adapters/database"
	"github.com/user/mangahub/internal/application"
	"github.com/user/mangahub/internal/eventbus"
	"github.com/user/mangahub/internal/infrastructure"
	mh_http "github.com/user/mangahub/internal/interfaces/http"
)

func TestHTTPIntegration(t *testing.T) {
	// Setup Infrastructure
	db, _ := infrastructure.NewSQLiteDB(":memory:")
	infrastructure.InitSchema(db)
	bus := eventbus.NewEventBus(10)
	secret := "test-secret"

	// Setup Services
	userRepo := database.NewSqliteUserRepository(db)
	mangaRepo := database.NewSqliteMangaRepository(db)
	progRepo := database.NewSqliteProgressRepository(db)

	authSvc := application.NewAuthService(userRepo)
	mangaSvc := application.NewMangaService(mangaRepo, bus)
	progSvc := application.NewProgressService(progRepo, bus)

	// Setup Handlers & Router
	authH := mh_http.NewAuthHandler(authSvc, secret)
	mangaH := mh_http.NewMangaHandler(mangaSvc)
	progH := mh_http.NewProgressHandler(progSvc)
	// Ports "0" mean no listener is up on those addresses → TCP and gRPC
	// probes fail predictably, exercising the degraded path without needing
	// real protocol servers in the test. startedAt is 30s in the past so
	// uptime_seconds is comfortably positive.
	healthH := mh_http.NewHealthHandler(db, bus, "8080", "0", "0", "0", time.Now().Add(-30*time.Second))

	mux := mh_http.SetupRouter(authH, mangaH, progH, healthH, secret)

	// --- 1. Health Check ---
	t.Run("Health Check", func(t *testing.T) {
		// Capture log output to verify the user-required logging.
		var logBuf bytes.Buffer
		origOut := log.Writer()
		log.SetOutput(&logBuf)
		defer log.SetOutput(origOut)

		req := httptest.NewRequest("GET", "/api/health", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		// No real TCP/gRPC listeners on test ports → degraded → 503.
		assert.Equal(t, http.StatusServiceUnavailable, rr.Code)

		var resp map[string]interface{}
		assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

		// Top-level keys.
		assert.Equal(t, "degraded", resp["status"])
		assert.Contains(t, resp, "checks")
		assert.Contains(t, resp, "bus")
		assert.Contains(t, resp, "uptime_seconds")
		assert.Contains(t, resp, "timestamp")

		// All 6 check keys present.
		checks, ok := resp["checks"].(map[string]interface{})
		assert.True(t, ok, "checks should be a map")
		for _, key := range []string{"http", "tcp", "udp", "ws", "grpc", "db"} {
			assert.Contains(t, checks, key, "missing check: "+key)
		}

		// Self-checks and in-memory SQLite always pass.
		assert.Equal(t, "ok", checks["http"])
		assert.Equal(t, "ok", checks["ws"])
		assert.Equal(t, "ok", checks["db"])

		// Real protocol probes against port "0" must fail (TCP + gRPC are
		// connection-oriented and will not find a listener). UDP is
		// connectionless so we only assert non-emptiness — see note in plan.
		tcpVal, _ := checks["tcp"].(string)
		grpcVal, _ := checks["grpc"].(string)
		assert.True(t, strings.HasPrefix(tcpVal, "error:"), "tcp should error, got: "+tcpVal)
		assert.True(t, strings.HasPrefix(grpcVal, "error:"), "grpc should error, got: "+grpcVal)
		assert.NotEmpty(t, checks["udp"], "udp probe should produce a result")

		// Uptime is a positive number.
		uptime, ok := resp["uptime_seconds"].(float64)
		assert.True(t, ok && uptime > 0, "uptime_seconds should be positive number, got %v", resp["uptime_seconds"])

		// Logging requirement (user explicit).
		logs := logBuf.String()
		assert.Contains(t, logs, "🩺 Health check from", "missing request log line")
		assert.Contains(t, logs, "🩺 Health check result: status=degraded", "missing result log line")
	})

	// --- 2. Registration ---
	t.Run("Register Admin", func(t *testing.T) {
		body := map[string]string{"username": "admin", "password": "password", "role": "admin"}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(b))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	// --- 3. Login & JWT Generation ---
	var token string
	t.Run("Login", func(t *testing.T) {
		body := map[string]string{"username": "admin", "password": "password"}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(b))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		
		var res map[string]string
		json.Unmarshal(rr.Body.Bytes(), &res)
		token = res["token"]
		assert.NotEmpty(t, token)
	})

	// --- 4. Protected Route: Create Manga (Admin Account) ---
	t.Run("Create Manga (Admin)", func(t *testing.T) {
		body := map[string]string{"title": "One Piece", "author": "Oda"}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/api/manga", bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer " + token)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	// --- 5. Protected Route: Unauthorized (No Token) ---
	t.Run("Create Manga (Unauthorized)", func(t *testing.T) {
		body := map[string]string{"title": "Solo Leveling", "author": "Chugong"}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/api/manga", bytes.NewBuffer(b))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	// --- 6. Update Progress ---
	t.Run("Update Progress", func(t *testing.T) {
		body := map[string]int{"manga_id": 1, "current_chapter": 5}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/api/manga/progress", bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// --- 7. Error Cases ---
	t.Run("Register Duplicate", func(t *testing.T) {
		body := map[string]string{"username": "admin", "password": "password", "role": "admin"}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(b))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code) // ErrUserExists results in 500 in current impl
	})

	t.Run("Invalid Login", func(t *testing.T) {
		body := map[string]string{"username": "admin", "password": "wrong"}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(b))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Create Manga (Invalid Body)", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/manga", bytes.NewBufferString("invalid json"))
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Update Progress (Invalid Body)", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/manga/progress", bytes.NewBufferString("invalid json"))
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

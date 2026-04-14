package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
	healthH := mh_http.NewHealthHandler(db, bus)

	mux := mh_http.SetupRouter(authH, mangaH, progH, healthH, secret)

	// --- 1. Health Check ---
	t.Run("Health Check", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/health", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "MangaHub is alive")
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

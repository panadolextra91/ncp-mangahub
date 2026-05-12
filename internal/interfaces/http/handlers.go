package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/user/mangahub/internal/application"
	"github.com/user/mangahub/internal/eventbus"
	"github.com/user/mangahub/internal/middleware"
	"github.com/user/mangahub/pkg/auth"
	"github.com/user/mangahub/pkg/models"
)

type AuthHandler struct {
	authService application.AuthService
	jwtSecret   string
}

func NewAuthHandler(svc application.AuthService, secret string) *AuthHandler {
	return &AuthHandler{authService: svc, jwtSecret: secret}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(req.Username, req.Password, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("👤 [AUTH] New User Registered: %s (Role: %s)", req.Username, req.Role)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	tokenString, err := auth.GenerateToken(user.ID, user.Username, user.Role, h.jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	log.Printf("🔑 [AUTH] User Logged In: %s (Role: %s)", user.Username, user.Role)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"user":  user,
	})
}

type MangaHandler struct {
	mangaService application.MangaService
}

func NewMangaHandler(svc application.MangaService) *MangaHandler {
	return &MangaHandler{mangaService: svc}
}

func (h *MangaHandler) Create(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok {
		http.Error(w, "Unauthorized context", http.StatusUnauthorized)
		return
	}

	var manga models.Manga
	if err := json.NewDecoder(r.Body).Decode(&manga); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.mangaService.CreateManga(role, &manga); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	log.Printf("🏗️ [MANGA] New Manga Created: %s (by User Role: %s)", manga.Title, role)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(manga)
}

func (h *MangaHandler) List(w http.ResponseWriter, r *http.Request) {
	qp := r.URL.Query()
	q := qp.Get("q")
	genresRaw := qp.Get("genres")
	status := qp.Get("status")
	sortBy := qp.Get("sortBy")

	// Parse genres: comma-separated, trim whitespace, drop empties, cap at 10.
	var genres []string
	if genresRaw != "" {
		for _, g := range strings.Split(genresRaw, ",") {
			g = strings.TrimSpace(g)
			if g == "" {
				continue
			}
			genres = append(genres, g)
			if len(genres) == 10 {
				break
			}
		}
	}

	// Backwards compatibility: if NO new filters AND no sort change, route
	// through the existing methods so q-only callers see bit-for-bit identical
	// SQL semantics to the pre-WH7 behavior.
	hasFilters := len(genres) > 0 || status != "" || (sortBy != "" && sortBy != "recent")

	var mangas []*models.Manga
	var err error

	log.Printf("🔍 [SEARCH] q='%s' genres=%v status='%s' sortBy='%s'", q, genres, status, sortBy)

	switch {
	case !hasFilters && q == "":
		mangas, err = h.mangaService.ListMangas()
	case !hasFilters && q != "":
		mangas, err = h.mangaService.SearchMangas(q)
	default:
		mangas, err = h.mangaService.SearchMangasWithFilters(application.SearchFilters{
			Query:  q,
			Genres: genres,
			Status: status,
			SortBy: sortBy,
		})
	}

	if err != nil {
		log.Printf("❌ [SEARCH] Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ [SEARCH] Found %d results", len(mangas))
	json.NewEncoder(w).Encode(mangas)
}

func (h *MangaHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	manga, err := h.mangaService.GetManga(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if manga == nil {
		http.Error(w, "Manga not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(manga)
}

type ProgressHandler struct {
	progressService application.ProgressService
}

func NewProgressHandler(svc application.ProgressService) *ProgressHandler {
	return &ProgressHandler{progressService: svc}
}

func (h *ProgressHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized context", http.StatusUnauthorized)
		return
	}

	var req struct {
		MangaID        int    `json:"manga_id"`
		CurrentChapter int    `json:"current_chapter"`
		Status         string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	progress := &models.UserProgress{
		UserID:         userID,
		MangaID:        req.MangaID,
		CurrentChapter: req.CurrentChapter,
		Status:         req.Status,
	}

	if err := h.progressService.UpdateProgress(progress); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("📖 [PROGRESS] User #%d updated Manga #%d to Chapter %d (Status: %s)", userID, req.MangaID, req.CurrentChapter, req.Status)

	json.NewEncoder(w).Encode(progress)
}

func (h *ProgressHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized context", http.StatusUnauthorized)
		return
	}

	progress, err := h.progressService.GetUserProgress(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(progress)
}

// HealthHandler probes every listener (HTTP self, TCP, UDP, WS, gRPC) plus the
// SQLite DB in parallel and reports an aggregate "ok|degraded" status. Every
// call is logged twice (request + result) per the audit requirement.
type HealthHandler struct {
	db        *sql.DB
	bus       *eventbus.EventBus
	httpPort  string
	tcpPort   string
	udpPort   string
	grpcPort  string
	startedAt time.Time
}

func NewHealthHandler(db *sql.DB, bus *eventbus.EventBus, httpPort, tcpPort, udpPort, grpcPort string, startedAt time.Time) *HealthHandler {
	return &HealthHandler{
		db:        db,
		bus:       bus,
		httpPort:  httpPort,
		tcpPort:   tcpPort,
		udpPort:   udpPort,
		grpcPort:  grpcPort,
		startedAt: startedAt,
	}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	log.Printf("🩺 Health check from %s at %s", r.RemoteAddr, time.Now().UTC().Format(time.RFC3339))

	const probeTimeout = 500 * time.Millisecond
	results := make(map[string]string, 6)
	var mu sync.Mutex
	var wg sync.WaitGroup

	set := func(name, val string) {
		mu.Lock()
		results[name] = val
		mu.Unlock()
	}

	// 1. HTTP self-check — if we are answering this request, HTTP is up.
	wg.Add(1)
	go func() {
		defer wg.Done()
		set("http", "ok")
	}()

	// 2. TCP probe — dial + immediate close.
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+h.tcpPort, probeTimeout)
		if err != nil {
			set("tcp", "error: "+err.Error())
			return
		}
		_ = conn.Close()
		set("tcp", "ok")
	}()

	// 3. UDP probe — best-effort send (UDP is connectionless; success means
	// the kernel let us emit a packet to that port).
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := net.Dial("udp", "127.0.0.1:"+h.udpPort)
		if err != nil {
			set("udp", "error: "+err.Error())
			return
		}
		defer conn.Close()
		_ = conn.SetWriteDeadline(time.Now().Add(probeTimeout))
		if _, err := conn.Write([]byte("PING\n")); err != nil {
			set("udp", "error: "+err.Error())
			return
		}
		set("udp", "ok")
	}()

	// 4. WS self-check — WS upgrade endpoint is mounted on the same HTTP
	// listener; if HTTP serves, the upgrade path is reachable.
	wg.Add(1)
	go func() {
		defer wg.Done()
		set("ws", "ok")
	}()

	// 5. gRPC probe — dial and wait for Ready within probeTimeout.
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), probeTimeout)
		defer cancel()

		conn, err := grpc.NewClient(
			"127.0.0.1:"+h.grpcPort,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			set("grpc", "error: "+err.Error())
			return
		}
		defer conn.Close()

		conn.Connect()
		for {
			s := conn.GetState()
			if s == connectivity.Ready {
				set("grpc", "ok")
				return
			}
			if !conn.WaitForStateChange(ctx, s) {
				set("grpc", "error: timeout (last state "+s.String()+")")
				return
			}
		}
	}()

	// 6. DB probe — keep existing ping semantics.
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := h.db.Ping(); err != nil {
			set("db", "error: "+err.Error())
			return
		}
		set("db", "ok")
	}()

	wg.Wait()

	overall := "ok"
	for _, v := range results {
		if v != "ok" {
			overall = "degraded"
			break
		}
	}

	uptime := math.Round(time.Since(h.startedAt).Seconds()*100) / 100

	log.Printf("🩺 Health check result: status=%s http=%s tcp=%s udp=%s ws=%s grpc=%s db=%s",
		overall, results["http"], results["tcp"], results["udp"], results["ws"], results["grpc"], results["db"])

	w.Header().Set("Content-Type", "application/json")
	if overall == "degraded" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":         overall,
		"checks":         results,
		"bus":            map[string]interface{}{"dropped_events": h.bus.DroppedCount()},
		"uptime_seconds": uptime,
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
	})
}

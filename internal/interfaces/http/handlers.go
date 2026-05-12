package http

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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

type HealthHandler struct {
	db  *sql.DB
	bus *eventbus.EventBus
}

func NewHealthHandler(db *sql.DB, bus *eventbus.EventBus) *HealthHandler {
	return &HealthHandler{db: db, bus: bus}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	dbErr := h.db.Ping()
	dbStatus := "OK"
	if dbErr != nil {
		dbStatus = "ERROR"
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "MangaHub is alive",
		"db":     dbStatus,
		"bus": map[string]interface{}{
			"dropped_events": h.bus.DroppedCount(),
		},
		"timestamp": time.Now(),
	})
}

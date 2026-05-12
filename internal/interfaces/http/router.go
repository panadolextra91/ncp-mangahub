package http

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/user/mangahub/internal/middleware"
)

//go:embed static/*
var staticAssets embed.FS

// SetupRouter configures the Go 1.22+ pattern-based router.
func SetupRouter(
	authH *AuthHandler,
	mangaH *MangaHandler,
	progH *ProgressHandler,
	healthH *HealthHandler,
	jwtSecret string,
) *http.ServeMux {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("GET /api/health", healthH.Check)
	mux.HandleFunc("POST /api/auth/register", authH.Register)
	mux.HandleFunc("POST /api/auth/login", authH.Login)

	// Protected routes
	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	// Manga routes
	mux.Handle("POST /api/manga", authMiddleware(http.HandlerFunc(mangaH.Create)))
	mux.Handle("GET /api/manga", authMiddleware(http.HandlerFunc(mangaH.List)))
	mux.Handle("GET /api/manga/{id}", authMiddleware(http.HandlerFunc(mangaH.Get)))

	// Progress routes
	mux.Handle("PUT /api/manga/progress", authMiddleware(http.HandlerFunc(progH.Update)))
	mux.Handle("GET /api/manga/library", authMiddleware(http.HandlerFunc(progH.List)))

	// Static Assets (Pink Dashboard)
	staticSub, _ := fs.Sub(staticAssets, "static")
	fileServer := http.FileServer(http.FS(staticSub))
	
	mux.Handle("GET /", fileServer)
	mux.Handle("GET /style.css", fileServer)
	mux.Handle("GET /app.js", fileServer)

	return mux
}

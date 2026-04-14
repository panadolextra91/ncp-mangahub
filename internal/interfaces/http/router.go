package http

import (
	"net/http"

	"github.com/user/mangahub/internal/middleware"
)

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

	// Progress routes
	mux.Handle("PUT /api/manga/progress", authMiddleware(http.HandlerFunc(progH.Update)))

	return mux
}

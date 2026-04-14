# Phase 4 Research: HTTP & SQLite Implementation

## 1. Schema Definition (SQLite)
Based on the models, the following schema will be initialized on startup:
```sql
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS mangas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_progress (
    user_id INTEGER NOT NULL,
    manga_id INTEGER NOT NULL,
    current_chapter INTEGER NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, manga_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (manga_id) REFERENCES mangas(id)
);
```

## 2. HTTP Routing (Go 1.22+)
Standard library `http.ServeMux` now supports method-based routing and path parameters:
```go
mux := http.NewServeMux()
mux.HandleFunc("POST /api/auth/register", authHandler.Register)
mux.HandleFunc("GET /api/manga", mangaHandler.List)
mux.HandleFunc("POST /api/manga", authMiddleware(mangaHandler.Create))
mux.HandleFunc("PUT /api/manga/{id}/progress", authMiddleware(progressHandler.Update))
```

## 3. JWT Middleware Pattern
The middleware will extract the `Authorization: Bearer <token>` header, validate it, and inject the claims into the context:
```go
type contextKey string
const roleKey contextKey = "role"
const userIDKey contextKey = "user_id"

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // ... validate JWT ...
            ctx := context.WithValue(r.Context(), roleKey, claims["role"])
            ctx = context.WithValue(ctx, userIDKey, claims["user_id"])
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

## 4. Config Loader
A simple config struct with fallback logic:
```go
type Config struct {
    JWTSecret string
    DBPath    string
}

func LoadConfig() *Config {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "mangahub-fallback-secret-123"
    }
    return &Config{JWTSecret: secret, DBPath: "mangahub.db"}
}
```

## 5. Repository Implementations
The repositories will live in `internal/adapters/database` and implement the interfaces from `internal/domain`. They will require `*sql.DB` via constructor.

## Validation Plan
- Unit tests for Repositories (using SQLite memory mode).
- Integration tests for HTTP handlers (using `net/http/httptest`).
- Manual verification of JWT extraction and context injection.
- Total coverage target: > 80%.

---
wave: 1
depends_on: [03-domain-models-application-services]
files_modified:
  - internal/infrastructure/sqlite.go
  - internal/adapters/database/sqlite_user_repo.go
  - internal/adapters/database/sqlite_manga_repo.go
  - internal/adapters/database/sqlite_progress_repo.go
  - internal/middleware/auth.go
  - internal/interfaces/http/handlers.go
  - internal/interfaces/http/router.go
  - config/config.go
  - cmd/server/main.go
  - test.sh
autonomous: true
---

# Phase 4 Plan: HTTP Protocol Layer (Core API)

## Objective
Implement the HTTP transport layer using Go 1.22+ standard library, integrate JWT authentication with role-based authorization, and provide concrete SQLite repository implementations for system persistence.

## Tasks

### [ ] Task 4.1: Database Schema & Repository Implementation
<read_first>
- .planning/phases/04-http-protocol-layer-core-api/04-RESEARCH.md
</read_first>
<action>
1. Update `internal/infrastructure/sqlite.go` to include a `InitSchema(db *sql.DB)` function that creates the `users`, `mangas`, and `user_progress` tables.
2. Create `internal/adapters/database/sqlite_user_repo.go` implementing `domain.UserRepository`.
3. Create `internal/adapters/database/sqlite_manga_repo.go` implementing `domain.MangaRepository`.
4. Create `internal/adapters/database/sqlite_progress_repo.go` implementing `domain.ProgressRepository`.
</action>
<acceptance_criteria>
- Repositories pass unit tests using a memory-mode SQLite database.
- Database correctly enforces foreign keys and unique constraints.
</acceptance_criteria>

### [ ] Task 4.2: Auth Middleware & Config Loader
<read_first>
- .planning/phases/04-http-protocol-layer-core-api/04-RESEARCH.md
</read_first>
<action>
1. Create `config/config.go` with `LoadConfig()` implementing the secret fallback logic (`mangahub-fallback-secret-123`).
2. Create `internal/middleware/auth.go`. Implement `AuthMiddleware(secret string)` that validates JWTs and injects `user_id` and `role` into `context`.
3. Export custom context keys to avoid collisions.
</action>
<acceptance_criteria>
- Middleware correctly extracts claims from a valid JWT.
- Invalid or missing tokens result in `401 Unauthorized`.
- Config loader correctly prioritizes env variables but falls back safely.
</acceptance_criteria>

### [ ] Task 4.3: HTTP Handlers & Router Setup
<read_first>
- .planning/phases/04-http-protocol-layer-core-api/04-RESEARCH.md
</read_first>
<action>
1. Create `internal/interfaces/http/handlers.go`. Implement `AuthHandler`, `MangaHandler`, and `ProgressHandler`.
2. Handlers must parse JSON requests, call Application Services (from Phase 3), and return JSON responses.
3. Handlers should extract `role` from the request `context` to pass as an argument to services.
4. Create `internal/interfaces/http/router.go` utilizing `http.NewServeMux()` with Go 1.22+ patterns (e.g. `POST /api/manga`).
</action>
<acceptance_criteria>
- API endpoints correctly handle happy paths and return appropriate error codes for failures.
- All 5 protocols plan (Phase 4 focuses on HTTP) is initiated.
</acceptance_criteria>

### [ ] Task 4.4: Main Orchestration & Integration Testing
<action>
1. Update `cmd/server/main.go` to:
   - Load config.
   - Initialize SQLite and run `InitSchema`.
   - Instantiate Repositories, EventBus, and Services.
   - Start the HTTP server in a goroutine (as per `plan.md`).
2. Update `test.sh` to include the new packages.
3. Write integration tests in `internal/interfaces/http/http_test.go` using `httptest`.
</action>
<acceptance_criteria>
- `go run cmd/server/main.go` starts the server and logs "MangaHub started" (or similar).
- Total coverage >= 80% confirmed via `./test.sh`.
</acceptance_criteria>

## Verification
- Run `./test.sh` to confirm all tests pass and coverage is > 80%.
- Manually verify token-based access to `/api/manga` (POST) requires `Admin` role.
- Verify that a `User` role can update progress but not create manga.

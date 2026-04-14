# Phase 4 Validation Strategy
**Status**: Completed

## Objective
Verify the HTTP Core API's capability to route requests, enforce JWT-based security, and persist data via SQLite repositories while maintaining the strict 80% coverage mandate.

## Verification Checklist

### 1. Database Persistence
- [x] Schema initialization creates `users`, `mangas`, and `user_progress` tables.
- [x] `SqliteUserRepository` successfully saves and retrieves users with Bcrypt verification.
- [x] `SqliteMangaRepository` handles CRUD operations correctly.
- [x] `SqliteProgressRepository` manages pivot records without collisions.

### 2. HTTP Routing & Middleware
- [x] `http.ServeMux` routes match Go 1.22+ patterns (Method + Path).
- [x] Auth Middleware correctly extracts and validates JWT from Authorization header.
- [x] Auth Middleware injects `role` and `userID` into the request context.
- [x] Unauthorized requests (missing/invalid token) are rejected with 401.

### 3. Application Integration
- [x] Register/Login flow returns valid JWTs.
- [x] Manga creation is strictly restricted to `Admin` roles.
- [x] Progress updates are permitted for authenticated users.

### 4. Definition of Done (DoD)
- [x] Total package coverage >= 80%.
- [x] "Hell Case": Verify no "database is locked" errors occur under concurrent HTTP requests (simulated).
- [x] `test.sh` passes all tests.

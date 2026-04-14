# Phase 4 Validation Strategy

## Objective
Verify the HTTP Core API's capability to route requests, enforce JWT-based security, and persist data via SQLite repositories while maintaining the strict 80% coverage mandate.

## Verification Checklist

### 1. Database Persistence
- [ ] Schema initialization creates `users`, `mangas`, and `user_progress` tables.
- [ ] `SqliteUserRepository` successfully saves and retrieves users with Bcrypt verification.
- [ ] `SqliteMangaRepository` handles CRUD operations correctly.
- [ ] `SqliteProgressRepository` manages pivot records without collisions.

### 2. HTTP Routing & Middleware
- [ ] `http.ServeMux` routes match Go 1.22+ patterns (Method + Path).
- [ ] Auth Middleware correctly extracts and validates JWT from Authorization header.
- [ ] Auth Middleware injects `role` and `userID` into the request context.
- [ ] Unauthorized requests (missing/invalid token) are rejected with 401.

### 3. Application Integration
- [ ] Register/Login flow returns valid JWTs.
- [ ] Manga creation is strictly restricted to `Admin` roles.
- [ ] Progress updates are permitted for authenticated users.

### 4. Definition of Done (DoD)
- [ ] Total package coverage >= 80%.
- [ ] "Hell Case": Verify no "database is locked" errors occur under concurrent HTTP requests (simulated).
- [ ] `test.sh` passes all tests.

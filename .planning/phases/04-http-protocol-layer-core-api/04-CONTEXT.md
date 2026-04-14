# Phase 4 Context: HTTP Protocol Layer (Core API)

## Domain Boundary
**What this encapsulates:** The external transport layer for the system via RESTful HTTP. It handles request parsing, authentication via JWT, authorization via Role-based middleware, and delegates business logic to Application Services. This phase also includes the concrete implementation of persistence adapters (SQLite) to fuel the API.

## Key Decisions

1. **HTTP Router Selection**
   - **Decision:** Pure `net/http` standard library from Go 1.22+.
   - **Details:** Leveraging the enhanced `http.ServeMux` for pattern-based routing (e.g., `POST /api/manga`) to demonstrate core language proficiency and avoid unnecessary dependencies.

2. **Persistence Implementation**
   - **Decision:** Concrete SQLite Repositories.
   - **Details:** Implementation of `UserRepository`, `MangaRepository`, and `ProgressRepository` using the standard `database/sql` driver and the shared `*sql.DB` connection pool from Phase 1 (enforcing WAL and `SetMaxOpenConns(1)`).

3. **Authentication & Config**
   - **Decision:** JWT-based Auth with Config Fallbacks.
   - **Details:**
     - Use `github.com/golang-jwt/jwt/v5`.
     - Configuration loaded from Environment/YAML with a hardcoded fallback (`mangahub-fallback-secret-123`) for demo safety.
     - Authentication Middleware: Validates JWT, extracts `UserID` and `Role`, and injects them into the request `context` using a typed key (e.g. `context.Context`).

4. **Role Enforcement**
   - **Decision:** Middleware-to-Service context flow.
   - **Details:** Handlers pull Role from context and pass it to Application Services. `MangaService` (from Phase 3) already enforces `Admin` checks.

## Canonical Refs
- `docs/plan.md`
- Phase 1 (SQLite Connection Pool)
- Phase 3 (Application Logic and Model Definitions)

## Deferred Ideas
None.

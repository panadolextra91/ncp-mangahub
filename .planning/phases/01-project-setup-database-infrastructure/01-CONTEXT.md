# Phase 1 Context

## Domain Boundary
**What this encapsulates:** Establishing the foundation of the project using a standard Go layout, initializing the SQLite database with WAL configuration, and setting up repository interfaces to decouple infrastructure from business logic.

## Key Decisions

1. **Tech Stack & Architecture Override**
   - **Decision:** Ignore the Global Rule requiring TS/Prisma + Hexagonal modules. STRICTLY FOLLOW `docs/plan.md` 100%.
   - **Details:** Use Go (1.21+), Go Channels, Goroutines, and the single-process monolith structure (`cmd/`, `internal/domain`, `internal/infrastructure`, `internal/interfaces`, `internal/eventbus`, `pkg/models`) outlined in the plan.
   
2. **Database Tooling**
   - **Decision:** SQLite with WAL mode.
   - **Details:** Manage lock contention by enforcing a strictly single write connection (`db.SetMaxOpenConns(1)`).

3. **Testing Definition of Done**
   - **Decision:** > 80% test coverage for both Unit and Integration tests.
   - **Details:** Must rigorously include "Hell Cases" (e.g., SQLite lock timeouts) and must print the visual coverage profile via `go test -cover` for showcasing.

## Canonical Refs
- `docs/plan.md` (Primary blueprint)

## Deferred Ideas
None for this phase.

# Phase 3 Context: Domain Models & Application Services

## Domain Boundary
**What this encapsulates:** Core business logic structures (Manga, User, Progress) and their governing Application Services (`MangaService`, `AuthService`, `ProgressService`). This layer orchestrates business requirements, entirely agnostic of the HTTP/WS/TCP protocols rendering it perfectly testable.

## Key Decisions

1. **Event Publishing Responsibility**
   - **Decision:** Application Service exclusively manages Event trigger actions.
   - **Details:** Downstream Protocol Handlers (Controllers) are structurally forbidden from dispatching events to the Event Bus. `MangaService` and `ProgressService` act as the undisputed controllers, executing logic and subsequently invoking `bus.Publish()` safely.

2. **Authentication & Roles Mechanics**
   - **Decision:** Standard JWT with Bcrypt hashing coupled with defined Role restrictions.
   - **Details:** The system incorporates a dedicated User entity featuring `Role` and Bcrypt hashed payloads. Critical core pathways (like `manga.new`) mandate `Admin` level authorization rigorously intercepted by the Application layer.

3. **Manga Progress Entity Model**
   - **Decision:** Dedicated Relational Structure (Pivot Mapping).
   - **Details:** Tracking progress directly inside the generic `Manga` model invites catastrophic write-overwrites (Race Conditions) among diverse users. Progress tracking operates explicitly within a decoupled `UserProgress{UserID, MangaID, CurrentChapter, UpdatedAt}` entity.

## Canonical Refs
- `docs/plan.md`

## Deferred Ideas
None.

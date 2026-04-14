---
wave: 1
depends_on: [01-project-setup-database-infrastructure, 02-internal-event-bus-implementation]
files_modified:
  - go.mod
  - go.sum
  - pkg/models/manga.go
  - pkg/models/user.go
  - pkg/models/progress.go
  - internal/domain/repositories.go
  - internal/application/auth_service.go
  - internal/application/manga_service.go
  - internal/application/progress_service.go
autonomous: true
---

# Phase 3 Plan: Domain Models & Application Services

## Objective
Establish perfectly isolated Domain Entities and their orchestrating Application Logic (Auth, Manga, Progress), aggressively wiring Database abstractions and seamlessly fusing them to the local distributed Event Bus to drive pure domain workflow without HTTP pollution.

## Tasks

### [ ] Task 3.1: Complete Domain Models
<read_first>
- .planning/phases/03-domain-models-application-services/03-CONTEXT.md
</read_first>
<action>
1. Relocate the conceptual definition of `Manga` to `pkg/models/manga.go` cleanly. Incorporate standard context `ID`, `Title`, `Author`, and `CreatedAt` parameters. Delete `internal/domain/manga.go` subsequently since `MangaRepository` will move.
2. Formulate `pkg/models/user.go` structured elegantly around `ID`, `Username`, `PasswordHash`, and internal `Role` (string: admin/user).
3. Elaborate `pkg/models/progress.go` explicitly declaring a rigid pivot mapping `UserProgress{UserID, MangaID, CurrentChapter, UpdatedAt}` resolving potential spatial mutation overwrites entirely. 
4. Forge `internal/domain/repositories.go` centralizing all necessary contractual agnostic guarantees: `UserRepository`, `MangaRepository`, and `ProgressRepository` interfaces to specify strict data gateway access.
</action>
<acceptance_criteria>
- `cat pkg/models/progress.go` outlines the requested relational pivot keys safely disconnected from Manga model.
- `cat internal/domain/repositories.go` exhibits `type UserRepository interface`.
</acceptance_criteria>

### [ ] Task 3.2: AuthService Implementation (Bcrypt)
<read_first>
- .planning/phases/03-domain-models-application-services/03-RESEARCH.md
</read_first>
<action>
1. Boot standard module dependencies (`go get golang.org/x/crypto/bcrypt`).
2. Draft `internal/application/auth_service.go`. Provide interface `AuthService` orchestrating fundamental `Register` and `Login`.
3. Construct the concrete logic layer `authService struct` binding internally with `domain.UserRepository`.
4. Command Bcrypt integrations utilizing `bcrypt.GenerateFromPassword` seamlessly upon Registration attempts caching passwords securely and verifying hashes effectively via `CompareHashAndPassword` during Logins.
</action>
<acceptance_criteria>
- `cat internal/application/auth_service.go` correctly maps `golang.org/x/crypto/bcrypt` invocations mitigating clear text bugs.
</acceptance_criteria>

### [ ] Task 3.3: MangaService & ProgressService Logic Orchestrator
<read_first>
- .planning/phases/03-domain-models-application-services/03-CONTEXT.md
</read_first>
<action>
1. Draft `internal/application/manga_service.go` containing `interface MangaService` enclosing a `CreateManga(role string, manga models.Manga)` action. Instantiate the internal implementation utilizing both `MangaRepository` and our specific standard `*eventbus.EventBus`. Command the logic to definitively verify an `Admin` identity parameter before orchestrating the Repo save sequences, triggering `bus.Publish("manga.new")` dynamically afterward!
2. Draft `internal/application/progress_service.go`. Expose a `UpdateProgress(progress models.UserProgress)` gateway ensuring successful entity updates invoke `bus.Publish("progress.updated")` allowing dynamic WebSocket listeners automated real-time alerts.
</action>
<acceptance_criteria>
- File `manga_service.go` strictly possesses `if role != "admin"` rejection paths.
- Code asserts that Application service implementations are explicitly initiating `.Publish()` requests natively!
</acceptance_criteria>

### [ ] Task 3.4: Pure Logic Tests
<action>
1. Engineer targeted functional validations inside `internal/application/application_test.go` covering core logic routing (e.g. Bcrypt auth rejection and Manga restrictions).
2. The testing execution suite automatically relies on native abstractions to rapidly compile and run. Ensure `./test.sh` honors the original `80.0%+` constraints across this critical Business suite layer validating internal robustness!
</action>
<acceptance_criteria>
- `test.sh` output indicates > 80% coverage against standard `internal/application...` bounds successfully without triggering database panics!
</acceptance_criteria>

## Verification
1. Execution of `./test.sh`.
2. Inspect logic source codes directly ensuring the HTTP controllers are safely entirely decoupled from internal publishing behavior mechanisms. All `bus.Publish` calls lie purely inside `.planning/application...` packages!

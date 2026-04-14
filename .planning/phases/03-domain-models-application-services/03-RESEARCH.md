# Phase 3: Domain Models & Application Services Research

## Key Findings

1. **Service Definitions and Clean Architecture Boundaries**
   - Application Services should expose interfaces to the presentation tier: `AuthService`, `MangaService`, and `ProgressService`.
   - Concrete implementations (e.g., `type authService struct`) explicitly harbor the `EventBus` and `Repositories` as private dependencies.
   - By shielding the implementations, unit testing protocol handlers in future stages becomes radically simplified by mocking these exact interfaces.

2. **Entity Modeling Constraints**
   - Without a heavy ORM (as established by the native SQLite plan), entities within `pkg/models/` must remain pure context structs.
   - Avoid littering entities with complex DB-oriented driver directives. Standard `json` mapping tags are ideal for serialized data broadcast events.
   - Models required: `User`, `Manga`, `UserProgress`.

3. **Bcrypt Integration**
   - We must rely on `golang.org/x/crypto/bcrypt` to implement standard secure password encryption safely avoiding cleartext leaks in the upcoming database integration.

## Validation Architecture

1. **Domain Test Validations:** Assert model structural instantiations exist independently.
2. **Business Requirements Checks:** Write logic verifications dictating only distinct `Admin` checks permit manga creations. Ensure Bcrypt hashing passes/fails optimally upon simulated login efforts within isolated `AuthService` integration tests.

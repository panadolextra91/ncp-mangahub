---
wave: 1
depends_on: []
files_modified:
  - cmd/server/main.go
  - internal/domain/manga.go
  - internal/infrastructure/sqlite.go
  - internal/infrastructure/sqlite_test.go
  - test.sh
  - go.mod
  - go.sum
autonomous: true
---

# Phase 1 Plan: Project Setup & Database Infrastructure

## Objective
Establish the foundational Go modular monolith structure and initialize a fully concurrent-safe SQLite database connection using WAL and pooling optimizations.

## Tasks

### [ ] Task 1.1: Initialize Go Module and Project Structure
<read_first>
- .planning/phases/01-project-setup-database-infrastructure/01-CONTEXT.md
</read_first>
<action>
1. Run `go mod init github.com/user/mangahub`.
2. Install necessary dependencies via `go get github.com/mattn/go-sqlite3` and `go get github.com/stretchr/testify/assert`.
3. Create the standard directory tree:
   - `mkdir -p cmd/server`
   - `mkdir -p internal/domain`
   - `mkdir -p internal/infrastructure`
   - `mkdir -p internal/interfaces`
   - `mkdir -p internal/eventbus`
   - `mkdir -p pkg/models`
4. Create a minimal `cmd/server/main.go` that prints "MangaHub started".
</action>
<acceptance_criteria>
- `go.mod` exists with the correct module name and dependencies.
- `ls cmd/server/main.go` returns success.
- Running `go run cmd/server/main.go` outputs "MangaHub started".
</acceptance_criteria>

### [ ] Task 1.2: Define Domain Contracts
<read_first>
- .planning/phases/01-project-setup-database-infrastructure/01-CONTEXT.md
</read_first>
<action>
1. Create `internal/domain/manga.go`.
2. Define a simple struct `type Manga struct { ID int; Title string }`.
3. Define `type MangaRepository interface` with a placeholder function signature `Save(Manga) error`.
4. Do NOT import any external / infrastructure packages in this file to maintain the pure domain boundary.
</action>
<acceptance_criteria>
- `cat internal/domain/manga.go` contains `type MangaRepository interface`.
- The `manga.go` file explicitly has no `internal/infrastructure` or database imports.
</acceptance_criteria>

### [ ] Task 1.3: Implement SQLite Infrastructure
<read_first>
- internal/domain/manga.go
- .planning/phases/01-project-setup-database-infrastructure/01-RESEARCH.md
</read_first>
<action>
1. Create `internal/infrastructure/sqlite.go`.
2. Implement a connection builder `func NewSQLiteDB(dsn string) (*sql.DB, error)`.
3. After opening the connection, aggressively enforce a single connection to serialize writes: run `db.SetMaxOpenConns(1)`.
4. Ensure the Write-Ahead Log is used: execute `_, err = db.Exec("PRAGMA journal_mode = WAL; PRAGMA synchronous = NORMAL;")`.
</action>
<acceptance_criteria>
- `cat internal/infrastructure/sqlite.go` exactly contains `db.SetMaxOpenConns(1)`.
- `cat internal/infrastructure/sqlite.go` exactly contains `"PRAGMA journal_mode = WAL;`.
</acceptance_criteria>

### [ ] Task 1.4: Unit, Integration & Hell Case Testing
<read_first>
- internal/infrastructure/sqlite.go
- .planning/PROJECT.md
</read_first>
<action>
1. Create `internal/infrastructure/sqlite_test.go`.
2. Write a unit test `TestSQLiteConnection` testing the successful creation of a `file:test.db?cache=shared&mode=memory` database and affirming no errors open/pinging it.
3. Write a Hell Case integration test `TestSQLiteHellPath_ConcurrentWriteRace`. Spawn 50 Goroutines attempting to create a test table and insert a row simultaneously. The test MUST inherently pass without timing out or throwing `SQLITE_BUSY` (database is locked) panics thanks to `SetMaxOpenConns(1)`. Use `sync.WaitGroup` to wait for all coroutines to finish.
4. Create a bash wrapper `test.sh` that securely evaluates the required coverage DoD:
   ```bash
   #!/bin/bash
   go test -v -coverprofile=coverage.out ./...
   go tool cover -func=coverage.out
   ```
</action>
<acceptance_criteria>
- `test.sh` exists and is executable.
- `cat internal/infrastructure/sqlite_test.go` exhibits `TestSQLiteHellPath_ConcurrentWriteRace` including a loop spawning multiple goroutines `go func()`.
</acceptance_criteria>

## Verification
- Execute `chmod +x test.sh && ./test.sh`. The test suite must pass without fatal errors or blocking panics.
- The output of `./test.sh` strictly prints out function coverage levels, honoring the Testing DoD established in `PROJECT.md`.

## Must Haves
- WAL enabled.
- `SetMaxOpenConns(1)` executed on boot.
- Over 80% coverage on initialization mechanisms.

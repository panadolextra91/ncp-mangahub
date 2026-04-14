# Phase 1: Project Setup & Database Infrastructure Research

## Key Findings

1. **Go Standard Project Layout**
   - The Go community emphasizes using standard structures for modular monoliths.
   - `cmd/server/main.go` will be the system's entry point.
   - `internal/` directory is critical as the compiler prevents external packages from importing code inside `internal/`.
   - Sub-directories: `internal/domain` (entities and interfaces), `internal/infrastructure` (database implementation), `internal/interfaces` (HTTP/TCP/etc), `internal/eventbus`.

2. **SQLite Go Driver Selection**
   - `github.com/mattn/go-sqlite3` is the standard (requires CGO). Given local development + CLI context, CGO is usually manageable.
   - `modernc.org/sqlite` is a CGO-free alternative, but `mattn/go-sqlite3` is historically more battle-tested for concurrent WAL operations. We will use `mattn/go-sqlite3`.

3. **SQLite Concurrency and WAL**
   - Default SQLite behavior locks the entire database for any write.
   - WAL (Write-Ahead Logging) allows simultaneous readers and a single writer. 
   - Connecting string MUST contain `_journal=WAL`.
   - In Go's `database/sql` pool, connection multiplexing can cause multiple parallel write attempts, triggering `database is locked` error (SQLITE_BUSY). 
   - **Crucial fix:** `db.SetMaxOpenConns(1)` restricts the Go application to only maintain a single connection to SQLite. This effectively queues all Go database operations behind one sequential connection, mitigating busy locks securely.

4. **Repository Pattern Implementation**
   - `internal/domain/manga.go` will define `type MangaRepository interface { ... }`.
   - `internal/infrastructure/sqlite/manga_repo.go` will implement it.
   - The `main.go` file explicitly injects the sqlite repository into the application services.

## Validation Architecture

1. **File System Verification:** Ensure `cmd/`, `internal/domain`, `internal/infrastructure` exist.
2. **Database Verification:** Execute a query to verify `journal_mode=wal` is active on the connected DB. Test `MaxOpenConns == 1`.
3. **Hexagonal Dependency Checks:** Assert that `internal/domain` does NOT import `internal/infrastructure`.
4. **Hell Case Testing:** Write a test simulating concurrent writes to verify no `SQLITE_BUSY` panic occurs when multiple goroutines attempt database inserts.

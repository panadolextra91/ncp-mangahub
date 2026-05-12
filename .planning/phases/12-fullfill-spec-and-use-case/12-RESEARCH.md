# Phase 12: Fullfill Spec and Use Case - Research

## 1. Schema Migration Strategy (SQLite)
Since we are adding columns to existing tables (`mangas` and `user_progress`), we need to handle this without losing data. 
SQLite's `ALTER TABLE ADD COLUMN` is supported but has limitations (e.g., cannot add `NOT NULL` columns without a default value).

**Plan:**
- Update `InitSchema` to include the new columns in the `CREATE TABLE` statements for fresh installs.
- Add individual `ALTER TABLE ... ADD COLUMN IF NOT EXISTS ...` statements for existing installations.
- **New Manga Columns:** `genres` (TEXT), `status` (TEXT), `total_chapters` (INTEGER), `description` (TEXT).
- **New Progress Columns:** `status` (TEXT).

## 2. JSON Seeding Implementation
We need to auto-populate the database with ~40 manga series.
- **File:** `data/manga_seed.json`.
- **Trigger:** In `cmd/server/main.go`, after schema initialization, check if `mangas` count is 0.
- **Logic:** Parse JSON and insert into DB.

## 3. Cobra CLI Structure
To match `docs/cli_manual.md`, the CLI should be structured as follows:
- `mangahub auth [login|register|logout|status]`
- `mangahub manga [search|info|list]`
- `mangahub library [add|list|remove|update]`
- `mangahub progress [update|history|sync]`
- `mangahub server [start|stop|status]`

**Implementation Details:**
- Use `github.com/spf13/cobra`.
- Config file in `~/.mangahub/config.yaml` to store API endpoint and Token.
- Use `net/http` to communicate with the API server.

## 4. Advanced Search (HTTP)
The `GET /api/manga` endpoint needs to handle filters.
- **Implementation:** Use `fmt.Sprintf` or a query builder to construct the SQL query dynamically based on provided genre, status, etc.
- **Path:** `internal/interfaces/http/handlers.go`.

## 5. UI/TUI Updates
- The TUI needs to display the new fields.
- The `PageCreate` in TUI should also be updated to include these fields.

## Validation Architecture
- **Tests:** 
  - Unit tests for the new search logic.
  - Integration tests for the JSON seed data import.
  - CLI "happy path" tests (mocking the server).

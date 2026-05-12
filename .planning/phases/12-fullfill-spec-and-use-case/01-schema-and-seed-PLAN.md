---
phase: 12
plan: 1
wave: 1
autonomous: true
files_modified:
  - pkg/models/manga.go
  - pkg/models/progress.go
  - internal/infrastructure/sqlite.go
  - internal/adapters/database/sqlite_manga_repo.go
  - internal/adapters/database/sqlite_progress_repo.go
---

# Plan: Schema and Seed Data Implementation

## Tasks

<task>
<read_first>
- `pkg/models/manga.go`
- `pkg/models/progress.go`
</read_first>
<action>
Update the domain models to include the missing fields from spec.md.
- `Manga`: Add `Genres` (string), `Status` (string), `TotalChapters` (int), `Description` (string).
- `UserProgress`: Add `Status` (string).
</action>
<acceptance_criteria>
- `pkg/models/manga.go` contains the new fields in the `Manga` struct.
- `pkg/models/progress.go` contains the `Status` field in the `UserProgress` struct.
</acceptance_criteria>
</task>

<task>
<read_first>
- `internal/infrastructure/sqlite.go`
</read_first>
<action>
Update the database schema to include the new columns and handle migrations.
- Update `CREATE TABLE mangas` with: `genres TEXT, status TEXT, total_chapters INTEGER, description TEXT`.
- Update `CREATE TABLE user_progress` with: `status TEXT`.
- Add migration logic in `InitSchema` to `ALTER TABLE` if columns are missing.
</action>
<acceptance_criteria>
- `internal/infrastructure/sqlite.go` has updated `CREATE TABLE` statements.
- `InitSchema` includes `ALTER TABLE` commands for each new column.
- Server starts successfully without errors on an existing database.
</acceptance_criteria>
</task>

<task>
<read_first>
- `internal/adapters/database/sqlite_manga_repo.go`
- `internal/adapters/database/sqlite_progress_repo.go`
</read_first>
<action>
Update the repository implementations to handle the new fields in SQL queries.
- Update `Create`, `GetByID`, `List` in `MangaRepository`.
- Update `UpdateProgress`, `GetByUserID` in `ProgressRepository`.
</action>
<acceptance_criteria>
- Repositories correctly scan/insert the new fields.
- No compilation errors in repository files.
</acceptance_criteria>
</task>

<task>
<action>
Create the `data/manga_seed.json` file with 40+ popular manga series including metadata.
</action>
<acceptance_criteria>
- `data/manga_seed.json` exists and is valid JSON.
- Contains at least 40 entries with "One Piece", "Naruto", etc.
</acceptance_criteria>
</task>

<task>
<read_first>
- `cmd/server/main.go`
</read_first>
<action>
Implement auto-seeding logic in the server startup.
- If `mangas` table is empty, load `data/manga_seed.json` and insert entries using `mangaRepo`.
</action>
<acceptance_criteria>
- Server startup logs "Database is empty. Seeding data..." when appropriate.
- Database contains 40+ entries after first run.
</acceptance_criteria>
</task>

<task>
<read_first>
- `internal/infrastructure/sqlite_test.go`
- `internal/adapters/database/sqlite_manga_repo_test.go`
</read_first>
<action>
Implement Unit and Integration tests for Schema migration and Seeding.
- **Happy Path:** Verify new columns exist and can be read/written.
- **Happy Path:** Verify seed data loads exactly 40+ entries.
- **Edge Case:** Migration on an already updated database (should be idempotent).
- **Edge Case:** Seeding when the file `manga_seed.json` is missing or malformed (should fail gracefully).
</action>
<acceptance_criteria>
- `go test ./internal/infrastructure/...` passes.
- `go test ./internal/adapters/database/...` passes.
- Test coverage for repository is verified.
</acceptance_criteria>
</task>

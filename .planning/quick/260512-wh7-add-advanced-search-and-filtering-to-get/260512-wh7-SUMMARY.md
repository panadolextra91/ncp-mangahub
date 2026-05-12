---
phase: 260512-wh7
plan: 01
subsystem: search-and-filtering
tags: [search, filtering, grpc, http, sqlite, proto, bonus-feature]
type: execute
wave: 1
depends_on: []
requirements:
  - WH7-FILTER-GENRES
  - WH7-FILTER-STATUS
  - WH7-SORT
  - WH7-GRPC-PARITY
  - WH7-BACKCOMPAT

requires:
  - "internal/domain/repositories.go::MangaRepository (existing)"
  - "internal/application/manga_service.go::MangaService (existing)"
  - "internal/adapters/database/sqlite_manga_repo.go::sqliteMangaRepo (existing)"
  - "internal/interfaces/http/handlers.go::MangaHandler.List (existing)"
  - "internal/interfaces/grpc/services.go::MangaService.SearchManga (existing)"
  - "api/proto/mangahub.proto::SearchRequest (existing)"

provides:
  - "domain.SearchFilters struct (Query/Genres/Status/SortBy)"
  - "application.SearchFilters type alias (= domain.SearchFilters)"
  - "MangaService.SearchMangasWithFilters(f) method"
  - "MangaRepository.SearchWithFilters(f) method"
  - "sqliteMangaRepo.SearchWithFilters: dynamic WHERE clause, quoted-token genre LIKE, sort recent|title"
  - "MangaHandler.List parses q/genres/status/sortBy with cap-of-10 + legacy routing"
  - "grpc.MangaService.SearchManga forwards Genres/Status/SortBy with legacy routing"
  - "Proto SearchRequest fields 2/3/4 (genres/status/sort_by); query=1 preserved"

affects:
  - "HTTP GET /api/manga: now accepts genres, status, sortBy query params"
  - "gRPC MangaService.SearchManga: accepts SearchRequest.Genres/Status/SortBy"
  - "Wire-level backwards compat preserved for both protocols"

tech_stack:
  added: []
  patterns:
    - "Type alias re-export (domain → application) to keep handler imports clean while avoiding domain→application cycle"
    - "Dynamic WHERE clause builder via []string clauses + []interface{} args"
    - "Quoted-token LIKE form (\"%\\\"Action\\\"%\") to prevent substring collisions"
    - "Legacy-path routing: backcompat-only calls bypass new filtered method"

key_files:
  created: []
  modified:
    - "api/proto/mangahub.proto"
    - "pkg/pb/mangahub.pb.go"
    - "internal/domain/repositories.go"
    - "internal/application/manga_service.go"
    - "internal/application/application_test.go"
    - "internal/adapters/database/sqlite_manga_repo.go"
    - "internal/adapters/database/database_test.go"
    - "internal/interfaces/http/handlers.go"
    - "internal/interfaces/http/search_test.go"
    - "internal/interfaces/grpc/services.go"

decisions:
  - "SearchFilters in domain package, application re-exports as type alias — only cycle-safe shape"
  - "Sibling method (SearchMangasWithFilters) instead of changing SearchMangas signature — zero call-site churn for q-only callers"
  - "Backcompat split at handler/service layer: q-only calls take the legacy SearchMangas path verbatim (identical SQL to pre-WH7)"
  - "Quoted-token LIKE form (%\\\"Action\\\"%) instead of json_each — simpler, SQLite-portable, avoids 'Action' matching 'Reaction'"
  - "Cap of 10 genres applied at BOTH the handler (early truncation) and the repo (defensive belt-and-suspenders)"
  - "Unknown sortBy values fall back to recent (id DESC) — no 400 error, fully backwards compatible"
  - "Unknown status values produce empty result set — no error (lenient validation)"

metrics:
  duration: "4m 26s"
  completed: "2026-05-12T16:38:00Z"
  tasks: 3
  files_modified: 10
  files_created: 0
  new_tests: 21
  new_test_functions: 2
---

# Quick Task 260512-wh7: Advanced Search & Filtering Summary

Wired multi-criteria search (genres OR, status, sortBy) through proto, application, repo, HTTP and gRPC layers while preserving bit-for-bit backwards compatibility for q-only callers — no schema migration required.

## Tasks Completed

| # | Task | Commit | Files |
|---|------|--------|-------|
| 1 | Extend SearchRequest proto + regenerate Go stubs | `1356372` | `api/proto/mangahub.proto`, `pkg/pb/mangahub.pb.go` |
| 2 | Wire filters through application, repo, HTTP, gRPC | `af51f23` | 7 files (domain, application, repo, http, grpc + 2 mocks) |
| 3 | Add filter tests at HTTP and repo layers | `5a44dcd` | `internal/interfaces/http/search_test.go`, `internal/adapters/database/database_test.go` |

## Architecture Decisions

### `SearchFilters` location: `domain` package

Placed in `internal/domain/repositories.go` (same file as `MangaRepository`). The application package re-exports it via type alias:

```go
// internal/application/manga_service.go
type SearchFilters = domain.SearchFilters
```

**Why:** `domain.MangaRepository.SearchWithFilters(f SearchFilters)` must reference the type. If `SearchFilters` lived in `application`, the domain package would have to import application, which already imports domain — a cycle. The type alias keeps handler/service call-sites referring to `application.SearchFilters` without leaking domain types upward.

### Interface extension: sibling method (Option A)

Added `SearchMangasWithFilters(f SearchFilters) ([]*models.Manga, error)` alongside the unchanged `SearchMangas(query string)`. **Zero call-site churn** for the q-only path:

- `MangaHandler.List` (HTTP): when only `?q=` is present, calls `SearchMangas` exactly as before. Filtered calls route to `SearchMangasWithFilters`.
- `grpc.MangaService.SearchManga`: when only `req.Query` is set, calls `SearchMangas`. Otherwise routes to `SearchMangasWithFilters`.

This is the explicit `<done>` requirement: q-only callers must hit the unchanged code path so their SQL/results are bit-for-bit identical.

### Proto field number preservation

```proto
message SearchRequest {
    string query = 1;          // PRESERVED — wire-level backcompat for existing gRPC clients
    repeated string genres = 2; // NEW
    string status = 3;          // NEW
    string sort_by = 4;         // NEW
}
```

Verified via `grep -E 'repeated string genres = 2;|string status = 3;|string sort_by = 4;' api/proto/mangahub.proto` → 3 matches. Regenerated `pkg/pb/mangahub.pb.go` exposes `Genres []string`, `Status string`, `SortBy string` on `SearchRequest`.

### Genre matching: quoted-token LIKE

Genres are stored as JSON-array text (e.g. `'["Action","Shounen"]'`). Naive substring `LIKE '%Action%'` would match `'["Reaction Time"]'`. The repo wraps each genre as `'%"Action"%'` — the JSON quotes guard against substring collisions. Verified by the `single genre — no substring collision` repo subtest.

### Cap-of-10 applied at two layers

- **Handler:** truncates input list during parsing, breaks at 10.
- **Repo:** defensively slices to 10 if a caller bypasses the handler.

This is belt-and-suspenders so gRPC clients (which don't go through the handler) also get the cap. The repo cap is verified by the `genres over cap of 10 are silently truncated` subtest.

## Impact Analysis Results

GitNexus MCP was not available in this worktree, so analysis used grep-based caller discovery (CLAUDE.md fallback noted in task constraints).

| Symbol | d=1 Callers | Risk | Action |
|--------|-------------|------|--------|
| `MangaHandler.List` | `router.go:34` (registration only) | LOW | None — handler internals changed, signature unchanged |
| `grpc.MangaService.SearchManga` | `cmd/server/main.go:127`, `grpc_test.go:42` (registration) | LOW | None — gRPC method signature unchanged |
| `application.MangaService.SearchMangas` | `handlers.go:117`, `services.go:49` | LOW | None — method unchanged; new sibling added |
| `application.MangaService` interface | implementor `mangaService`; mock `MockMangaService` | MEDIUM | Added `SearchMangasWithFilters` impl + mock stub |
| `sqliteMangaRepo.Search` | `manga_service.go:61`, `seeding_test.go:42` | LOW | None — method unchanged |
| `domain.MangaRepository` interface | impl `sqliteMangaRepo`; mock `mockMangaRepo` | MEDIUM | Added `SearchWithFilters` impl + mock stub |

No HIGH or CRITICAL risk. All d=1 dependents updated in the same commit (Task 2). Post-edit `git status --short` confirmed only the 7 expected files for Task 2 and 2 expected files for Task 3 were modified.

## Test Coverage Delta

**21 new subtests across 2 new test functions:**

- `internal/interfaces/http/search_test.go::TestMangaListFilters` — 9 subtests:
  - Filter by single genre, multi-genre OR, status, sortBy=title
  - Combined filters (AND across types)
  - Genre cap at 10 (12 → 10)
  - Whitespace/empty stripping
  - Backcompat: q-only routes to legacy `SearchMangas` (not filtered path)
  - Backcompat: `q + sortBy=recent` stays on legacy path

- `internal/adapters/database/database_test.go::TestSqliteMangaRepository_SearchWithFilters` — 12 subtests:
  - Single genre with substring-collision guard (`'Action' MUST NOT match 'Reaction Time'`)
  - Multi-genre OR semantics (3-of-4 match)
  - Status filter, unknown-status lenient empty
  - AND across types (genre + status narrows to 1)
  - Sort by title ASC, sort by recent (id DESC), unknown sortBy fallback to recent
  - Query matches title or author
  - No filters returns all
  - 12-genre input silently truncated to 10

**Existing tests untouched and still passing:** `TestMangaListSearch`, `TestHTTPIntegration`, `TestSqliteMangaRepository`, `TestSeeding`, `TestMangaService_CreateRoleBlocks`, `TestGRPCIntegration`, `TestSqliteUserRepository`, `TestSqliteProgressRepository`.

## Verification Results

```bash
go build ./...    # PASS
go vet ./...      # PASS
go test ./... -count=1
# All packages PASS:
#   internal/adapters/database  1.162s
#   internal/application        0.988s
#   internal/interfaces/grpc    2.164s
#   internal/interfaces/http    2.405s
#   internal/middleware         3.108s
#   tests/e2e                  10.296s
#   (...all others)
```

Proto-level verification:

```bash
grep -E 'repeated string genres = 2;|string status = 3;|string sort_by = 4;' api/proto/mangahub.proto | wc -l
# 3 (all three fields present)
grep -E 'Genres\s+\[\]string|SortBy\s+string' pkg/pb/mangahub.pb.go
# Genres []string (rep,name=genres,proto3)
# SortBy string (opt,name=sort_by,json=sortBy,proto3)
```

`gitnexus_detect_changes` was not run (MCP unavailable in worktree). Equivalent `git status` and `git diff --stat` checks confirmed:

- Task 2 modified exactly the 7 files listed in its `<files>` block.
- Task 3 modified exactly the 2 test files listed.

## Backcompat Verification

Both wire-level and SQL-level backwards compatibility are proved by:

1. **HTTP test `Backcompat: q-only routes to SearchMangas (not SearchMangasWithFilters)`:** asserts `MockMangaService.SearchMangas("naruto")` is called and `SearchMangasWithFilters` is **never** called. Same for `q + sortBy=recent` (treats `recent` as the no-op default).
2. **Proto field `query=1` preserved** (grep'd in summary above).
3. **`sqliteMangaRepo.Search(q)` method body is unchanged.** Verified by `git diff` (no edits to that function).
4. **`TestSeeding` (which calls `mangaRepo.Search("One Piece")` against real seed data) still passes** with no modifications — confirms legacy SQL semantics intact.
5. **`TestGRPCIntegration` (full gRPC roundtrip) still passes** with no modifications — confirms gRPC wire format unchanged for existing clients.

## Deviations from Plan

**None — plan executed exactly as written.**

The plan was unusually precise about file locations, expected interface shapes, and required tests. Deviation rules 1–3 did not trigger. One sub-optimal pattern was discovered (the legacy `Search` query lowercases LIKE patterns against all three columns including `genres`, so even a q-only search will surface genre matches — this is intentional pre-existing behavior, not a bug, and preserving it is part of the backcompat contract).

## Authentication Gates

None. Task was fully autonomous; no auth secrets, no manual steps, no external services touched.

## Known Stubs

None. All new code paths are fully wired and exercised by tests:

- `SearchMangasWithFilters` connects HTTP handler → application service → domain repo → real SQL.
- gRPC handler forwards real request fields to real application service.
- No placeholder data, no TODO markers, no mock-only paths shipped to production code.

## Threat Flags

None. No new network endpoints (existing `/api/manga` route is unchanged structurally). No new auth paths. No file access. No schema changes at trust boundaries.

The only behavior surface added is parameter parsing on an already-authenticated route, with all string inputs handled via parameterized SQL (`r.db.Query(query, args...)`) — no string concatenation of user input into the SQL body. Cap of 10 limits DoS surface from arbitrarily long genre lists. Unknown status/sortBy values are silently lenient rather than erroring.

## Self-Check

Verifying claims before finalizing:

**Files exist:**
- `api/proto/mangahub.proto` — FOUND
- `pkg/pb/mangahub.pb.go` — FOUND
- `internal/domain/repositories.go` — FOUND
- `internal/application/manga_service.go` — FOUND
- `internal/application/application_test.go` — FOUND
- `internal/adapters/database/sqlite_manga_repo.go` — FOUND
- `internal/adapters/database/database_test.go` — FOUND
- `internal/interfaces/http/handlers.go` — FOUND
- `internal/interfaces/http/search_test.go` — FOUND
- `internal/interfaces/grpc/services.go` — FOUND

**Commits exist (per `git log --oneline -4`):**
- `1356372` — FOUND
- `af51f23` — FOUND
- `5a44dcd` — FOUND

## Self-Check: PASSED

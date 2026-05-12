---
phase: 12
plan: 2
wave: 2
autonomous: true
files_modified:
  - internal/interfaces/http/handlers.go
  - internal/interfaces/tui/app.go
  - internal/interfaces/tui/model.go
---

# Plan: API and TUI Enhancements

## Tasks

<task>
<read_first>
- `internal/interfaces/http/handlers.go`
</read_first>
<action>
Enhance the `GET /api/manga` handler to support filtering.
- Support query parameters: `genre`, `status`, `title`.
- Update the repository `List` method call to pass these filters.
</action>
<acceptance_criteria>
- `curl "http://localhost:8080/api/manga?genre=Shounen"` returns only Shounen manga.
- `curl "http://localhost:8080/api/manga?status=Ongoing"` works correctly.
</acceptance_criteria>
</task>

<task>
<read_first>
- `internal/interfaces/http/handlers_test.go`
</read_first>
<action>
Write integration tests for the enhanced search API.
- **Happy Path:** Search by valid genre "Shounen".
- **Happy Path:** Search by valid title "One Piece".
- **Edge Case:** Search with non-existent genre (should return empty list, not error).
- **Edge Case:** Search with special characters in query params.
</action>
<acceptance_criteria>
- New tests in `handlers_test.go` cover all filter combinations.
- All tests pass with `go test ./internal/interfaces/http/...`.
</acceptance_criteria>
</task>

<task>
<read_first>
- `internal/interfaces/tui/model.go`
- `internal/interfaces/tui/app.go`
</read_first>
<action>
Update the TUI to display the new Manga metadata and support library status updates.
- Show `Genres`, `Status`, and `Chapters` in the Manga Info/List views.
- Add a field to select library status (Reading, Completed, etc.) when updating progress.
</action>
<acceptance_criteria>
- TUI displays "Genres: ..." and "Status: ..." in the manga details.
- User can see the reading status of their tracked manga.
</acceptance_criteria>
</task>

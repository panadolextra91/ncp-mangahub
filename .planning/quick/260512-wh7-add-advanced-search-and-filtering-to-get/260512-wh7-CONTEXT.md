# Quick Task 260512-wh7: Advanced Search & Filtering - Context

**Gathered:** 2026-05-12
**Status:** Ready for planning

<domain>
## Task Boundary

Add advanced multi-criteria filters to manga search:
- HTTP: `GET /api/manga` extended with new query params
- gRPC: `MangaService.SearchManga` extended with filter fields in `SearchRequest`

Bonus feature target: +5pt "Advanced Search & Filtering" per spec. Must not regress existing `?q=` behavior.

</domain>

<decisions>
## Implementation Decisions

### Migration scope — LOCKED: No schema migration
- Do NOT add `rating` or `year` columns to `mangas`.
- Reason: seed data has no values; filters would be useless at demo. Bonus is 5pt — not worth migration risk.
- Filterable fields are limited to what exists today: `title`, `author`, `genres` (TEXT JSON), `status` (TEXT), `total_chapters` (INT).

### Sort options — LOCKED: `recent` + `title` only
- `sortBy=recent` (default) → `ORDER BY id DESC` (matches current behavior).
- `sortBy=title` → `ORDER BY LOWER(title) ASC`.
- Do NOT implement `chapters` or `popularity` sort.
- If `sortBy` value is anything else, fall back to `recent` (do not error — backwards compatible).

### Multi-genre logic — LOCKED: OR semantics
- Query param `?genres=Action,Romance` returns manga having AT LEAST ONE of the genres.
- SQL approach: `genres` is TEXT JSON array. Use `LIKE '%"Action"%' OR LIKE '%"Romance"%'` for simplicity (avoids json_each table-valued function complexity). Each genre must be looked up with the quoted-token form to avoid substring collisions (e.g. "Action" matching "Reaction").
- Comma-separated parsing; trim whitespace; empty entries ignored; max 10 genres accepted (cap to keep WHERE clause sane).

### gRPC parity — LOCKED: Update proto + regen
- `SearchRequest` gets new fields: `repeated string genres = 2; string status = 3; string sort_by = 4;`
- Field number 1 (`query`) stays — backwards compatible at wire level.
- Regen via `./gen_proto.sh`. Update `internal/interfaces/grpc/services.go::SearchManga` to read new fields and pass through to application service.

### Claude's Discretion
- Pagination: defer (existing endpoint has no pagination — keep parity, don't introduce in this quick task)
- Validation errors: be lenient — unknown sortBy → fallback to recent; unknown status → no manga matches (empty result, not 400)
- Application-service signature: extend existing `SearchMangas(query string)` to `SearchMangas(filters SearchFilters)` OR add a sibling method `SearchMangasWithFilters` — pick whichever causes less call-site churn. Old proto query strings still hit the same code path.
- Test coverage: at least 1 unit test per filter (genres single, genres multi, status, sortBy=title, backwards-compat `?q=` only). HTTP layer test sufficient; gRPC can rely on the same application-service tests.

</decisions>

<specifics>
## Specific Ideas

### HTTP query param contract
```
GET /api/manga?q=naruto&genres=Action,Romance&status=ongoing&sortBy=title
```
All filters optional. All filters combinable. Result = intersection (AND across different filter types, OR within `genres`).

### SQL skeleton (illustrative — planner will refine)
```sql
SELECT * FROM mangas
WHERE 1=1
  AND (LOWER(title) LIKE LOWER(?) OR LOWER(author) LIKE LOWER(?))  -- when q present
  AND (genres LIKE ? OR genres LIKE ?)                              -- one LIKE per genre
  AND status = ?                                                    -- when status present
ORDER BY {id DESC | LOWER(title) ASC}                               -- by sortBy
```

### Proto change (illustrative)
```proto
message SearchRequest {
    string query = 1;
    repeated string genres = 2;
    string status = 3;
    string sort_by = 4;
}
```

</specifics>

<canonical_refs>
## Canonical References

- `docs/spec.md` — "Advanced Search & Filtering (5 pts)" bonus section
- `internal/interfaces/http/handlers.go::MangaHandler.List` — current search endpoint
- `internal/interfaces/grpc/services.go::SearchManga` — current gRPC method
- `api/proto/mangahub.proto::SearchRequest` — proto definition
- `internal/infrastructure/sqlite.go` — schema (line 45+ mangas table)
- `internal/application/manga_service.go` — application service `SearchMangas`
- `gen_proto.sh` — proto regen script
</canonical_refs>

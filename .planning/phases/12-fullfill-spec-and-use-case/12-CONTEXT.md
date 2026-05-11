# Phase 12: Fullfill Spec and Use Case - Context

**Gathered:** 2026-05-11
**Status:** Ready for planning
**Source:** User discussion + `docs/spec.md`, `docs/use_case.md`, `docs/cli_manual.md`

<domain>
## Phase Boundary

This phase bridges all remaining gaps between the current implementation and the project specifications. It focuses on schema completion, data population, and providing the "Standard CLI" interface alongside the existing TUI.

</domain>

<decisions>
## Implementation Decisions

### 1. CLI Interface (Cobra)
- **Decision:** Implement a subcommand-based CLI using `github.com/spf13/cobra`.
- **Command Structure:** Match `docs/cli_manual.md` (e.g., `mangahub auth login`, `mangahub manga search`, `mangahub library add`).
- **Integration:** This CLI will act as a secondary client interface to the existing API server.

### 2. Database & Schema Completion
- **Manga Table:** Add `genres` (JSON/Text), `status`, `total_chapters`, and `description`.
- **UserProgress Table:** Add `status` field (Enum/Text: `reading`, `completed`, `plan_to_read`, `on_hold`, `dropped`).
- **Migration:** Update `InitSchema` in `internal/infrastructure/sqlite.go`.

### 3. Seed Data (JSON)
- **Strategy:** Create `data/manga_seed.json` with ~40-50 popular series.
- **Auto-Import:** The server will check if the `mangas` table is empty on startup and import the seed data automatically.

### 4. Search & Filtering (HTTP)
- **Enhancement:** Update `GET /api/manga` to support query params for `title`, `author`, `genre`, and `status`.

### 5. the agent's Discretion
- **TUI Updates:** Refresh the TUI to include the new fields and allow filtering by library status.
- **Mapping**: Use `Slug` for manga IDs where appropriate if requested by the spec, or stick to consistent integer IDs if simpler.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Project Specs
- `docs/spec.md` — Core requirements and database structure.
- `docs/use_case.md` — Functional requirements for search and library.
- `docs/cli_manual.md` — The blueprint for the new Cobra CLI interface.

### Infrastructure
- `internal/infrastructure/sqlite.go` — Current DB schema and initialization logic.
- `pkg/models/` — Domain entities to be updated.

</canonical_refs>

<specifics>
## Specific Ideas
- "One Piece", "Naruto", "Dragon Ball", "Berserk", "Monster" should be in the seed data.
- The CLI should use `127.0.0.1` as the default server host.

</specifics>

<deferred>
## Deferred Ideas
- Advanced AI Recommendations (Phase 13+).
- Social Friend System (Phase 13+).
</deferred>

---

*Phase: 12-fullfill-spec-and-use-case*
*Context gathered: 2026-05-11 via Discuss Phase*

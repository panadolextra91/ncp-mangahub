# Plan 1: Infrastructure & Static Asset Setup

This plan sets up the necessary dependencies for the TUI and the "Single Binary" delivery mechanism using Go `embed`.

- **Wave**: 1
- **Depends on**: None
- **Files modified**: `go.mod`, `internal/interfaces/http/router.go`, `cmd/server/main.go`
- **New files**: `internal/interfaces/http/static/index.html`, `internal/interfaces/http/static/style.css`, `internal/interfaces/http/static/app.js`

## Tasks

<task>
<action>Add Bubbletea and Lipgloss dependencies to `go.mod`.</action>
<read_first>
- [go.mod](file:///Users/huynhngocanhthu/ncp-mangahub/go.mod)
</read_first>
<acceptance_criteria>
- `go.mod` contains `github.com/charmbracelet/bubbletea` and `github.com/charmbracelet/lipgloss`.
- `go mod tidy` runs successfully.
</acceptance_criteria>
</task>

<task>
<action>Initialize the `static/` directory and create placeholder "Pink" HTML/CSS/JS files.</action>
<read_first>
- [11-RESEARCH.md](file:///Users/huynhngocanhthu/ncp-mangahub/.planning/phases/11-tui-and-html-interface/11-RESEARCH.md)
</read_first>
<acceptance_criteria>
- `internal/interfaces/http/static/index.html` created.
- `internal/interfaces/http/static/style.css` contains the Pink Pastel (#FBCFE8) and Subtle Black (#0B0E11) variables.
</acceptance_criteria>
</task>

<task>
<action>Modify `router.go` to serve the `static/` directory via `embed.FS`.</action>
<read_first>
- [router.go](file:///Users/huynhngocanhthu/ncp-mangahub/internal/interfaces/http/router.go)
</read_first>
<acceptance_criteria>
- `router.go` uses `//go:embed static/*`.
- `SetupRouter` includes a handler for `GET /` that serves the embedded filesystem.
- Accessing `http://localhost:8080/` returns the placeholder HTML.
</acceptance_criteria>
</task>

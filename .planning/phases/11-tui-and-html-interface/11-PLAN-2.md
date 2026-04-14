# Plan 2: Pink TUI Implementation

This plan implements the interactive, pink-themed TUI client with alternating ASCII art.

- **Wave**: 2
- **Depends on**: Plan 1
- **Files modified**: `cmd/client/main.go`
- **New files**: `internal/interfaces/tui/app.go`, `internal/interfaces/tui/model.go`, `internal/interfaces/tui/styles.go`, `internal/interfaces/tui/assets.go`

## Tasks

<task>
<action>Create `internal/interfaces/tui/assets.go` containing the ASCII art templates provided by the user.</action>
<read_first>
- [11-CONTEXT.md](file:///Users/huynhngocanhthu/ncp-mangahub/.planning/phases/11-tui-and-html-interface/11-CONTEXT.md)
</read_first>
<acceptance_criteria>
- File contains `KeroChan`, `Berserk`, and `Evangelion` string constants.
- ASCII art is scaled/trimmed if necessary (approx 50% width).
</acceptance_criteria>
</task>

<task>
<action>Implement `styles.go` using `lipgloss` with Pastel Pink (#FBCFE8) theme.</action>
<read_first>
- [11-RESEARCH.md](file:///Users/huynhngocanhthu/ncp-mangahub/.planning/phases/11-tui-and-html-interface/11-RESEARCH.md)
</read_first>
<acceptance_criteria>
- Styles for `Header`, `Tab`, `ActiveTab`, `Border`, and `Notification` defined.
</acceptance_criteria>
</task>

<task>
<action>Implement the Bubbletea `Model` and `Update` loop in `model.go` and `app.go`.</action>
<read_first>
- [demo/run_show.go](file:///Users/huynhngocanhthu/ncp-mangahub/demo/run_show.go)
</read_first>
<acceptance_criteria>
- App supports Tab switching (Chat, Events, Create).
- A ticker `tea.Tick` updates the ASCII art index every 60s.
- `cmd/client/main.go` starts the Bubbletea program.
</acceptance_criteria>
</task>

<task>
<action>Integrate Auth, WebSocket Chat, and Manga Creation into the TUI using Goroutines and `tea.Cmd`.</action>
<read_first>
- [internal/interfaces/http/router.go](file:///Users/huynhngocanhthu/ncp-mangahub/internal/interfaces/http/router.go)
- [11-CONTEXT.md](file:///Users/huynhngocanhthu/ncp-mangahub/.planning/phases/11-tui-and-html-interface/11-CONTEXT.md)
</read_first>
<acceptance_criteria>
- WebSocket/TCP reading loops run in separate goroutines (via `tea.Cmd`).
- Incoming data is wrapped in custom types (e.g., `ChatMsg`, `EventMsg`) and sent to the Bubbletea `Update` loop.
- The UI remains completely responsive during network operations.
- User can login, chat, and receive pink notifications in-TUI.
</acceptance_criteria>
</task>

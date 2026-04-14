# Phase 11 CONTEXT: TUI and HTML interface

## Implementation Decisions

### 1. Unified Aesthetic: "Pink Professional"
- **TUI**: Use `bubbletea` for state management and `lipgloss` for styling. Theme is Pastel Pink (#FBCFE8) with white/light pink accents.
- **HTML**: Dark theme with a "Subtle Black" background (e.g., #0B0E11). Components will mimic Shadcn/UI but use the Pink Pastel palette.
- **Icons**: Use Lucide (for Web) and Nerd Font symbols/emojis (for TUI).

### 2. TUI (The Pink Terminal)
- **Interactive**: Must support `Login`, `Send Chat`, `View Manga Log`, and `Create Manga`.
- **Non-blocking TUI**: All networking (WS/TCP) MUST run in goroutines. Use `tea.Cmd` to bridge incoming data to the UI loop via custom messages (e.g., `type chatMsg string`).
- **ASCII Art**: Alternating Kero-chan, Berserk, and Evangelion icons every 60 seconds. Scale to fit smaller terminals (approx 50% width).
- **Notifications**: In-TUI popups or status line highlights for new manga releases.

### 3. Web Dashboard (The Shadcn Mimic)
- **Stack**: Pure HTML/CSS/JS (no heavy frameworks) served via Go `embed`.
- **Architecture**: Single Binary delivery. The Go server will serve the static files at `GET /`.
- **Glow Indicators**: 5 LED dots (HTTP, TCP, UDP, WS, gRPC) that pulse/glow Pink (#FBCFE8) with a CSS `drop-shadow` effect whenever an event is received on that protocol.
- **Real-time**: WebSockets for live chat and manga notifications.


### 4. Integration
- **Auth**: Both TUI and HTML must provide a Login form to obtain the JWT.
- **Routing**: Serve at `http://localhost:8080/`.

## Open Questions (Clarified)
- **Scale**: User provided ASCII blocks, con sẽ tinh chỉnh để nó không chiếm quá nhiều không gian TUI.
- **Binary**: Đóng gói tất cả vào 1 file duy nhất cho mẹ dễ nộp bài.

## Ref References
- [run_show.go](file:///Users/huynhngocanhthu/ncp-mangahub/demo/run_show.go) (Old CLI pattern to be replaced)
- [router.go](file:///Users/huynhngocanhthu/ncp-mangahub/internal/interfaces/http/router.go) (Mounting point for static)

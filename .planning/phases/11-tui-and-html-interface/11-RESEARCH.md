# Phase 11 RESEARCH: TUI and HTML interface

## Technical Overview
The goal is to provide a "Pink Professional" visualization layer. This involves a TUI client for terminal lovers and a Web Dashboard for browser users. Both will be served/delivered as part of the MangaHub modular monolith.

## 1. Aesthetic Research
### Color Palette
- **Background (Subtle Black)**: `#0B0E11` (Very dark grey/near-black).
- **Primary Accent (Pink Pastel)**: `#FBCFE8` (Tailwind Pink-200).
- **Secondary Accent**: `#F9A8D4` (Tailwind Pink-300).
- **Text (Secondary)**: `#E2E8F0` (Slate-200).

### Shadcn/UI Style (Vanilla CSS)
- **Radii**: `0.5rem`.
- **Shadows**: Subtle box-shadow with low opacity.
- **Glassmorphism**: `backdrop-filter: blur(8px); background: rgba(11, 14, 17, 0.8)`.

## 2. TUI Implementation (Bubbletea/Lipgloss)
### Component Structure
- **Header**: Banner with system status.
- **Side Panel**: Alternating ASCII Art (Kero-chan, etc.).
- **Main View**: Tabs for `[1] Chat`, `[2] Events`, `[3] Create`.
- **Status Bar**: Live log of bus activity.

### ASCII Alternator
- Use `time.Ticker` in a separate `Msg` stream.
- Every 60s, update the `Model.AsciiIndex`.
- ASCII Art blocks will be hardcoded in a `pkg/tui/assets.go`.

## 3. Web Dashboard (Embed)
### File Structure
- `internal/interfaces/http/static/index.html`
- `internal/interfaces/http/static/style.css`
- `internal/interfaces/http/static/app.js`

### Go Embed Integration
```go
//go:embed static/*
var staticFS embed.FS
```
Mount using `http.FileServer(http.FS(staticFS))`.

## 4. Notification Bridge
- **WebSocket**: Use the existing `/api/chat` route or a new `/api/events` route for system-wide broadcasts.
- **TUI**: Connect as a TCP or WebSocket client internally.

## 5. Verification Strategy
- **Visual Audit**: Manual check of "Pink" consistency.
- **Responsiveness**: Verify TUI handles window resizing (using `tea.WindowSizeMsg`).
- **Concurrency**: Ensure multiple TUIs/Browsers can connect simultaneously without leaking server resources.

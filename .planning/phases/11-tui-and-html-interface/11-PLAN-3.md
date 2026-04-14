# Plan 3: Pink Dashboard Implementation

This plan implements the "Shadcn-style" premium dashboard with a pink pastel theme and dark subtle background.

- **Wave**: 3
- **Depends on**: Plan 1
- **Files modified**: `internal/interfaces/http/static/index.html`, `internal/interfaces/http/static/style.css`, `internal/interfaces/http/static/app.js`

## Tasks

<task>
<action>Refine `style.css` to mimic Shadcn/UI components (Card, Button, Input, Toast).</action>
<read_first>
- [11-RESEARCH.md](file:///Users/huynhngocanhthu/ncp-mangahub/.planning/phases/11-tui-and-html-interface/11-RESEARCH.md)
</read_first>
<acceptance_criteria>
- Background set to `#0B0E11`.
- Primary highlight set to `#FBCFE8`.
- Rounded corners (`0.5rem`) and sleek Inter font used.
</acceptance_criteria>
</task>

<task>
<action>Implement the dashboard layout in `index.html` including Glow LED Indicators.</action>
<read_first>
- [11-CONTEXT.md](file:///Users/huynhngocanhthu/ncp-mangahub/.planning/phases/11-tui-and-html-interface/11-CONTEXT.md)
</read_first>
<acceptance_criteria>
- Sections for: Login (Overlay), Status Grid (HTTP/TCP/UDP/WS/gRPC), Live Event Log, and Live Chat.
- Status Grid includes 5 LED dots with "Pink Glow" (#FBCFE8) CSS effects.
- LEDs pulse/glow when their respective protocol receives an event.
</acceptance_criteria>
</task>

<task>
<action>Implement `app.js` for real-time interactions.</action>
<read_first>
- [tests/e2e/integration_test.go](file:///Users/huynhngocanhthu/ncp-mangahub/tests/e2e/integration_test.go)
</read_first>
<acceptance_criteria>
- WS client connects to `/api/chat` with token.
- Handles incoming messages and updates the UI instantly.
- Triggers "Pink Toasts" for system-wide notifications.
</acceptance_criteria>
</task>

<task>
<action>Add a "Manga Creator" form to the dashboard.</action>
<read_first>
- [internal/interfaces/http/router.go](file:///Users/huynhngocanhthu/ncp-mangahub/internal/interfaces/http/router.go)
</read_first>
<acceptance_criteria>
- Submits `POST /api/manga` with JWT.
- Success/Error reflected via Toasts.
</acceptance_criteria>
</task>

# Phase 14: Add Personal Library - Context & Decisions

## 🎯 Domain Boundary
Implementing a comprehensive User Library system that allows users to save manga, track reading progress, and synchronize states across multiple clients using HTTP and TCP protocols.

## 🛠️ Decisions & Specifics

### 1. User Journey (The "How")
- **Search to Library**: User searches for manga -> Selects result -> Press `Enter` to view Detail page -> Press `a` (Add) to save to Library.
- **Why**: Prevents keyboard input collisions in the Search query box and ensures the user sees details before adding.

### 2. Library Structure (The "What")
- **States**: Every library entry must have one of three statuses: `Reading`, `Completed`, or `Plan to Read`.
- **Default**: New entries default to `Plan to Read`.
- **Data Model**: A many-to-many relationship (or a link table `user_libraries`) between Users and Manga, storing `current_chapter` and `status`.

### 3. Visual Presentation (The "Where")
- **New Tab**: A dedicated **LIBRARY [7]** tab will be added to the TUI.
- **Content**: Displays a list of saved manga with their current progress and status.
- **Styling**: Consistent with the Pink Hub aesthetic, using highlighted rows for navigation.

### 4. Network Synchronization (The "Real-time")
- **Persistence**: Saved via **HTTP REST API** (`POST /api/users/library`).
- **Broadcasting**: Upon successful addition, the server will broadcast a message via **TCP Progress Sync** to all connected clients (e.g., `"USER_ADDED_TO_LIBRARY|{UserID}|{MangaID}"`).
- **Client Side**: Other TUIs receiving this message will show a toast notification and refresh their Library tab if active.

## 🔗 Canonical Refs
- `docs/requirements.md` (User Data Management section)
- `internal/interfaces/http/handlers.go` (For new REST endpoints)
- `internal/interfaces/tcp/server.go` (For broadcast logic)
- `internal/interfaces/tui/app.go` (For UI integration)

## 🚀 Next Steps
- Research/Design the `user_libraries` database schema.
- Implement the HTTP handlers for library management.
- Update the TCP server to handle library-specific broadcasts.
- Build the TUI Detail View and Library Tab.

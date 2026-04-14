# Phase 6 Research: WebSocket Protocol Layer (Chat)

## 1. Database Schema (SQLite)
The SQL to initialize the chat table:
```sql
CREATE TABLE IF NOT EXISTS chat_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    manga_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_chat_manga_id ON chat_messages(manga_id);
```

## 2. Models & Repository
- **Model**: `ChatMessage` struct in `pkg/models/chat.go`.
- **Repository**: `ChatRepository` interface in `internal/domain/repositories.go`.
  - `Save(msg *ChatMessage) error`
  - `GetRecentByManga(mangaID int, limit int) ([]*ChatMessage, error)`

## 3. WebSocket Hub Design
To support partitioned chat rooms while keeping it simple:
- A `ChatHub` will maintain a `map[int]map[*Client]bool` where the first key is `manga_id`.
- `Client` struct will wrap `*websocket.Conn` and have a `send` channel.
- `Hub.Run()` will handle `register`, `unregister`, and `broadcast` channels.
- `broadcast` channel will carry a `ChatMessage` object, and the hub will only send it to clients in the matching `manga_id` bucket.

## 4. Auth & Upgrading
- `http.HandlerFunc` will extract `manga_id` and `token` from URL query.
- Use `pkg/auth.ValidateToken(token, secret)`.
- If valid, upgrade connection using `websocket.Upgrader`.
- Send the last 20 messages from `ChatService` immediately after upgrade.

## 5. EventBus Bridge
- When a `ChatMessage` is saved via `ChatService.SendMessage`, it must also `bus.Publish("chat.message", msg)`.
- The `main.go` event bridge will need an update to broadcast these to TCP clients.

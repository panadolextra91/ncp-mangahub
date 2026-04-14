---
wave: 1
depends_on: [02-internal-event-bus-implementation, 03-domain-models-application-services, 04-http-protocol-layer-core-api]
files_modified:
  - internal/infrastructure/sqlite.go
  - pkg/models/chat.go
  - internal/domain/repositories.go
  - internal/adapters/database/sqlite_chat_repo.go
  - internal/application/chat_service.go
  - internal/interfaces/ws/hub.go
  - internal/interfaces/ws/handlers.go
  - cmd/server/main.go
  - internal/interfaces/ws/ws_test.go
autonomous: true
---

# Phase 6 Plan: WebSocket Protocol Layer (Chat)

## Objective
Implement a real-time chat system using WebSockets, partitioned by Manga IDs. The system must support persistent messages (SQLite), immediate history delivery (last 20 messages), and cross-protocol event broadcasting.

## Tasks

### [ ] Task 6.1: Data Layer (Models & Persistence)
<read_first>
- .planning/phases/06-websocket-protocol-layer-chat/06-RESEARCH.md
</read_first>
<action>
1. Update `internal/infrastructure/sqlite.go` to include the `chat_messages` table and index on `manga_id`.
2. Create `pkg/models/chat.go` with `ChatMessage` struct.
3. Update `internal/domain/repositories.go` to include `ChatRepository` interface.
4. Implement `internal/adapters/database/sqlite_chat_repo.go`.
</action>
<acceptance_criteria>
- SQLite schema initializes correctly.
- `SqliteChatRepository` passes unit tests for Save and Retrieve.
</acceptance_criteria>

### [ ] Task 6.2: Application Layer (Chat Service)
<action>
1. Create `internal/application/chat_service.go`.
2. Implement `SendMessage` (Save to DB then publish to `EventBus`).
3. Implement `GetRecentMessages` (Fetch 20 most recent).
</action>
<acceptance_criteria>
- `ChatService` correctly publishes to the "chat.message" topic.
- Service retrieves history as expected.
</acceptance_criteria>

### [ ] Task 6.3: WebSocket Interface (Hub & Handler)
<action>
1. Implement `internal/interfaces/ws/hub.go`. Manage a registry of clients mapped by `manga_id`.
2. Implement `internal/interfaces/ws/handlers.go`.
   - Upgrader logic.
   - Query parameter Auth (`token` and `manga_id`).
   - On connect: Fetch and send last 20 messages.
   - On receive: Route to `ChatService`.
</action>
<acceptance_criteria>
- WebSocket connection requires valid JWT via query param.
- Clients automatically receive room history upon joining.
</acceptance_criteria>

### [ ] Task 6.4: Integration & Cross-Protocol Sync
<action>
1. Update `cmd/server/main.go`:
   - Setup `ChatHub` and `ChatHandler`.
   - Update EventBus bridge: subscribe to `chat.message` and broadcast to TCP Hub.
2. Write integration tests in `internal/interfaces/ws/ws_test.go`.
3. Verify global coverage >= 80%.
</action>
<acceptance_criteria>
- Chat messages from WS appear on TCP clients (nc/telnet).
- DoD coverage is maintained.
</acceptance_criteria>

## Verification
- Run `./test.sh`.
- Manual test:
  1. Start server.
  2. Connect TCP client (`nc localhost 9090`) and AUTH.
  3. Connect WS client (Browser/Postman) to room 1.
  4. Send WS message.
  5. Verify persistent receipt on TCP terminal and storage in DB.

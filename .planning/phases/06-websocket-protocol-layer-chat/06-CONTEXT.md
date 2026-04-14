# Phase 6: WebSocket Protocol Layer (Chat)

## Overview
Phase 6 introduces the third protocol layer: WebSockets. This layer will facilitate real-time chat rooms tied to specific Manga titles.

## Key Decisions (Context)
- **Library**: `github.com/gorilla/websocket`.
- **Persistence (STRICT)**:
  - Table: `chat_messages` (id, manga_id, user_id, content, created_at).
  - **Index**: Mandatory index on `manga_id`.
  - Behavior: Upon connection, the server must fetch and send the **last 20 messages** for that room.
- **Room Management (STRICT)**:
  - 1 Connection = 1 Room (Manga).
  - Join method: `GET /api/chat?manga_id=X&token=Y`.
- **Authentication**: JWT passed via **Query Parameter** (`token`).
- **Cross-Layer Integration**: Every chat message received MUST be published to the `EventBus` (Topic: `chat.message`) to allow other protocol subscribers (TCP/UDP) to react.
- **Architecture**: A central `Registry` managing room-specific client pools.

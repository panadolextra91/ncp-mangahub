# Phase 6 Validation Strategy
**Status**: Pending

## Objective
To verify that the WebSocket Protocol Layer correctly implements a secured, persistent chat system with room isolation and cross-protocol event broadcasting.

## Verification Checklist

### 1. Persistence & Data Integrity
- [ ] Table `chat_messages` exists with index on `manga_id`.
- [ ] `SqliteChatRepository` correctly stores messages and retrieves the latest 20.

### 2. Synchronization & History
- [ ] Upon WS connection, the client receives a burst of up to 20 historical messages for the selected room.
- [ ] Sending a message via WS correctly saves it to the DB and broadcasts it to other clients *in the same room*.

### 3. Cross-Protocol Magic (Hell Test)
- [ ] **TCP Integration**: Connect a TCP client (`nc localhost 9090`). Send a message from a WS client. The TCP client must receive the JSON notification via the `chat.message` topic.
- [ ] **Event Bus Reliability**: Ensure the `ChatService` publishes to the bus without blocking.

### 4. Security
- [ ] WebSocket handshake (Upgrade) is rejected if the `token` query parameter is missing or invalid.
- [ ] User context (UserID) is correctly identified from the JWT and tagged to chat messages.

### 5. Room Isolation
- [ ] Connect Client A to `manga_id=1`.
- [ ] Connect Client B to `manga_id=2`.
- [ ] Client A sends a message. Client B MUST NOT receive it.

### 6. Definition of Done
- [ ] Total statement coverage remains >= 80%.
- [ ] `test.sh` passes 100%.

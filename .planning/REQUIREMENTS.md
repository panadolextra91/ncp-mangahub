# Requirements

## Scope
**In Scope:**
- Core HTTP REST API (JWT Auth, Manga CRUD, Progress update).
- Real-Time Sync via TCP.
- Real-Time Chat via WebSocket.
- Notification broadcasting via UDP.
- Admin CLI operations via gRPC.
- Internal Event Bus (Go channels) to decouple protocols.
- SQLite Database configured with WAL.

**Out of Scope:**
- Distributed Microservices.
- RabbitMQ/Kafka or other external message brokers.
- DB writing from non-HTTP protocols.

## Definition of Done (DoD)
- **Coverage**: Must be above 80% with visual stats printed for showcase (`go test -coverprofile` etc.).
- **Tests Per Phase**: Every phase requires Unit Tests and Integration Tests.
- **Hell Cases**: Must explicitly test edge cases, catastrophic failures, connection drops, and slow consumers.

## Must Have (Table Stakes)
- **REQ-1:** Internal HTTP Server handling Manga CRUD and Auth.
- **REQ-2:** SQLite DB fully operational with `journal_mode=WAL` and single open connection.
- **REQ-3:** Internal Event Bus that supports subscribing and publishing without blocking.
- **REQ-4:** TCP server for real-time progress updates.
- **REQ-5:** WebSocket server for chat functionality.
- **REQ-6:** UDP server for fast notifications.
- **REQ-7:** Graceful shutdown mechanism that safely closes DB and all goroutines.

## Should Have (Differentiators)
- **REQ-8:** gRPC endpoint for admin internal operations to show strong-typed RPC integration.

## Could Have (Delighters)
- **REQ-9:** Simple web UI to demonstrate WebSocket and HTTP interactions.

## Won't Have
- **REQ-10:** Multi-node state scaling (all state is single process).

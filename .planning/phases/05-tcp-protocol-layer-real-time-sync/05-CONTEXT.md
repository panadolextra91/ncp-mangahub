# Phase 5 Context: TCP Protocol Layer (Real-time Sync)

## Domain Boundary
**What this encapsulates:** A real-time JSON broadcasting layer over TCP. It allows external clients to receive live updates of system events (`manga.new`, `progress.updated`) via a persistent socket connection. It demonstrates the multi-protocol capability of the MangaHub modular monolith.

## Key Decisions

1. **Protocol Configuration**
   - **Decision:** Port 9090 (configurable).
   - **Details:** The TCP port must be loaded from `config/config.go`. Hardcoded port numbers are strictly forbidden.

2. **TCP Authentication (Handshake)**
   - **Decision:** Explicit JWT Handshake (`AUTH <token>`).
   - **Details:** Upon connection, the server waits for a mandatory first line starting with `AUTH `. The token is validated using the shared system secret.
     - **Success:** Respond with `200 OK CONNECTED`.
     - **Failure:** Respond with `401 Unauthorized` and terminate connection.
   - **Deadline:** Every connection must complete the `AUTH` handshake within **5 seconds** (using `conn.SetReadDeadline`). Failures or idle connections are terminated to prevent Goroutine leaks.

3. **Concurrency Pattern (Channel Registry)**
   - **Decision:** Hub/Registry pattern using Go Channels (Idiomatic Go) with **DOS Protection**.
   - **Details:**
     - Use a central `Hub` struct with `register`, `unregister`, and `broadcast` channels.
     - **DOS Protection:** Implement a `MaxTCPClients` limit (default 100). If the Hub is full, new connections are served a `503 Service Unavailable - Server Full` message and closed immediately.
     - A single goroutine manages the state in a `select` loop (no `sync.RWMutex` for broadcasting).
     - **Slow-Consumer Prevention:** If a client's send buffer is full, the Hub must immediately drop the message for that specific client or disconnect them to prevent head-of-line blocking for the entire system.

4. **Event Bus Subscription**
   - **Decision:** Transparent Bridge.
   - **Details:** The TCP Hub subscribes to topics from the `internal/eventbus` and routes them into its internal broadcast channel.

## Canonical Refs
- `docs/plan.md`
- Phase 2 (Event Bus Implementation)
- Phase 4 (JWT Auth and Config Layer)

## Deferred Ideas
None.

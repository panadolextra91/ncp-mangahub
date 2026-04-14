---
wave: 1
depends_on: [03-domain-models-application-services, 04-http-protocol-layer-core-api]
files_modified:
  - config/config.go
  - pkg/auth/jwt.go
  - internal/middleware/auth.go
  - internal/interfaces/http/handlers.go
  - internal/interfaces/tcp/hub.go
  - internal/interfaces/tcp/server.go
  - cmd/server/main.go
  - internal/interfaces/tcp/tcp_test.go
autonomous: true
---

# Phase 5 Plan: TCP Protocol Layer (Real-time Sync)

## Objective
Implement a high-performance, secured TCP protocol layer that broadcasts real-time system events (Manga creation and progress updates) to authenticated clients using an idiomatic Go Channel Registry pattern.

## Tasks

### [ ] Task 5.1: Config & JWT Refactor
<read_first>
- .planning/phases/05-tcp-protocol-layer-real-time-sync/05-RESEARCH.md
</read_first>
<action>
1. Update `config/config.go` to include `TCPPort` (default "9090") and `MaxTCPClients` (default 100).
2. Create `pkg/auth/jwt.go` incorporating `GenerateToken` and `ValidateToken` functions.
3. Update `internal/middleware/auth.go` and `internal/interfaces/http/handlers.go` to utilize the new `pkg/auth` logic ensuring 100% backward compatibility for the HTTP API.
</action>
<acceptance_criteria>
- `test.sh` passes for `config` and `http` packages.
- No hardcoded ports in code.
</acceptance_criteria>

### [ ] Task 5.2: TCP Hub (Channel Registry)
<read_first>
- .planning/phases/05-tcp-protocol-layer-real-time-sync/05-RESEARCH.md
</read_first>
<action>
1. Create `internal/interfaces/tcp/hub.go`. Implement the `Hub` with `register`, `unregister`, and `broadcast` channels.
2. Implement DOS check in Hub registration: If `len(clients) >= MaxTCPClients`, send `503 Service Unavailable - Server Full` and close connection.
3. Ensure the `broadcast` loop uses a non-blocking write strategy (select with default or write timeout) to prevent slow consumers from impacting other clients.
</action>
<acceptance_criteria>
- Hub logic avoids `sync.Mutex` for broadcasting.
- One slow client doesn't freeze the hub message distribution.
</acceptance_criteria>

### [ ] Task 5.3: TCP Server & Auth Handshake
<read_first>
- .planning/phases/05-tcp-protocol-layer-real-time-sync/05-CONTEXT.md
</read_first>
<action>
1. Create `internal/interfaces/tcp/server.go`. Implement a TCP server that listens on `TCPPort`.
2. Implement Handshake Timeout: Set `conn.SetReadDeadline` (5 seconds) immediately after Accept.
3. Implement the Handshake: Every new connection must provide `AUTH <JWT>` as the first line.
4. Send `200 OK CONNECTED\n` on success, or `401 Unauthorized\n` + close on failure.
</action>
<acceptance_criteria>
- TCP server requires valid JWT for entry.
- Correct response codes sent to TCP client.
</acceptance_criteria>

### [ ] Task 5.4: Integration & DoD Verification
<action>
1. Update `cmd/server/main.go` to:
   - Instantiate the TCP Hub.
   - Run the Hub and TCP Server in goroutines.
   - Bridge EventBus topics into the Hub's broadcast channel.
2. Write integration tests in `internal/interfaces/tcp/tcp_test.go` covering the Handshake and Broadcast logic.
3. Verify global coverage >= 80%.
</action>
<acceptance_criteria>
- `telnet localhost 9090` followed by `AUTH <token>` successfully connects and receives JSON updates when HTTP requests occur.
- Total coverage meets requirements.
</acceptance_criteria>

## Verification
- Run `./test.sh`.
- Manual test: Run server, connect via `nc localhost 9090`, AUTH, and trigger a Manga create via HTTP. Verify JSON receipt on TCP console.

---
wave: 1
depends_on: [02-internal-event-bus-implementation, 05-tcp-protocol-layer-real-time-sync]
files_modified:
  - config/config.go
  - internal/interfaces/udp/registry.go
  - internal/interfaces/udp/server.go
  - cmd/server/main.go
  - internal/interfaces/udp/udp_test.go
autonomous: true
---

# Phase 7 Plan: UDP Protocol Layer (Notifications)

## Objective
Implement a lightweight, fire-and-forget notification system via UDP, secured with a Heartbeat/TTL mechanism for resource management.

## Tasks

### [ ] Task 7.1: Infrastructure & Configuration
<action>
1. Update `config/config.go` to include `UDPPort` (Defaults to 9191).
</action>
<acceptance_criteria>
- `UDP_PORT` environment variable is correctly loaded.
</acceptance_criteria>

### [ ] Task 7.2: UDP Registry (The "Brain")
<action>
1. Implement `internal/interfaces/udp/registry.go`.
2. Manage a `sync.Map` of authenticated UDP Peers.
3. Implement `Register(mangaID int, addr *net.UDPAddr)` and `KeepAlive(addr *net.UDPAddr)`.
4. Implement a background **Garbage Collector** routine that removes peers if not seen for > 60 seconds.
</action>
<acceptance_criteria>
- Peers are automatically pruned after 60s of inactivity.
- Thread-safe access to peer list.
</acceptance_criteria>

### [ ] Task 7.3: UDP Server (Protocol Layer)
<action>
1. Implement `internal/interfaces/udp/server.go`.
2. Listen for UDP packets on `UDPPort`.
3. Support two commands:
   - `SUB <manga_id> <jwt>`: Validate JWT, register peer.
   - `PING <jwt>`: Validate JWT, refresh peer TTL.
</action>
<acceptance_criteria>
- Server correctly parses SUB and PING commands.
- Invalid JWTs are ignored (silent drop).
</acceptance_criteria>

### [ ] Task 7.4: Integration & DoD Verification
<action>
1. Update `cmd/server/main.go`:
   - Setup `UDPRegistry` and `UDPServer`.
   - Update EventBus bridge: subscribe to system events (e.g., `manga.new`) and broadcast to `UDPRegistry`.
2. Write integration tests in `internal/interfaces/udp/udp_test.go`.
3. Verify global coverage >= 80%.
</action>
<acceptance_criteria>
- UDP clients receive global notifications.
- Heartbeat logic successfully prunes inactive clients in tests.
</acceptance_criteria>

## Verification
- Run `./test.sh`.
- Manual test:
  1. Start server.
  2. Send UDP packet `SUB 1 <jwt>` to 9191.
  3. Wait 10s. Send `PING <jwt>`.
  4. Trigger `manga.new` via HTTP.
  5. Verify UDP receipt.
  6. Wait 70s without PING. Verify no further notifications received (pruned).

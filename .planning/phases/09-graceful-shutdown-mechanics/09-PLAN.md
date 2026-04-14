---
wave: 1
depends_on: [08-grpc-protocol-layer-admin-cli]
files_modified:
  - internal/interfaces/tcp/hub.go
  - internal/interfaces/ws/hub.go
  - internal/interfaces/tcp/server.go
  - internal/interfaces/udp/server.go
  - cmd/server/main.go
  - internal/interfaces/grpc/services.go
autonomous: true
---

# Phase 9 Plan: Graceful Shutdown Mechanics

## Objective
Implement a robust, orchestrated shutdown sequence that ensures no data loss, resource leaks, or "Thundering Herd" reconnection issues.

## Tasks

### [ ] Task 9.1: Hub & Server Refactoring
<read_first>
- .planning/phases/09-graceful-shutdown-mechanics/08-RESEARCH.md
</read_first>
<action>
1. Update `mh_tcp.Hub.Run()` to accept `context.Context` and `*sync.WaitGroup`.
2. Update `mh_ws.Hub.Run()` to accept `context.Context` and `*sync.WaitGroup`.
3. Implement `Stop()` logic in both hubs:
   - Send `{"type": "system", "action": "shutdown", "reason": "server_maintenance"}`.
   - Wait 5s (Grace Period) or until all clients are gone.
   - Force close connections.
4. Update `mh_udp.Server` and `mh_tcp.Server` to use `context.Context` for their Start loops.
</action>
<acceptance_criteria>
- Hubs return from `Run()` when context is cancelled.
- Clients receive shutdown message before disconnect.
</acceptance_criteria>

### [ ] Task 9.2: Bridge & EventBus Cleanup
<action>
1. Refactor the `bridge` logic in `main.go` to listen for `ctx.Done()`.
2. Ensure `bridge` goroutines call `wg.Done()`.
</action>
<acceptance_criteria>
- Bridge goroutines terminate cleanly upon shutdown.
</acceptance_criteria>

### [ ] Task 9.3: Main Orchestration Rewrite
<action>
1. Use `signal.NotifyContext` in `main.go`.
2. Initialize root `sync.WaitGroup`.
3. Wrap final shutdown in `context.WithTimeout(10s)`.
4. Follow the strict sequence: Listeners -> Notify Clients -> Internal Chains -> wg.Wait() -> db.Close().
</action>
<acceptance_criteria>
- Server shuts down within the 10s window.
- Database is closed last.
- All protocols report clean closure.
</acceptance_criteria>

## Verification
- Run `./test.sh` to ensure no regressions.
- Manual test:
  1. Start server.
  2. Connect TCP and WS clients.
  3. Send `SIGINT` (Ctrl+C).
  4. Verify clients receive the shutdown message and the process exits within 5-10s.

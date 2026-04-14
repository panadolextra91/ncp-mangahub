---
wave: 1
depends_on: [09-graceful-shutdown-mechanics]
files_modified:
  - ARCHITECTURE.md
  - tests/e2e/isolation_test.go [NEW]
  - tests/e2e/resilience_test.go [NEW]
  - tests/e2e/integration_test.go [NEW]
  - demo/run_show.go [NEW]
autonomous: true
---

# Phase 10 Plan: Demo Integration & E2E Verification

## Objective
The "Final Stand". Prove the system's resilience under "Hell Path" conditions, automate a pro-level showcase, and document the architecture for academic excellence.

## Tasks

### [ ] Task 10.1: E2E Isolation & Resilience Testing
<action>
1. Create `tests/e2e/isolation_test.go`:
   - Start server.
   - Connect a TCP client and stop it from reading (simulate lag).
   - Verify other clients (WS, gRPC) still operate at full speed.
   - Verify the laggy client is eventually kicked or handled without blocking the Bus.
2. Create `tests/e2e/resilience_test.go`:
   - Concurrent writes to SQLite under load to verify WAL mode effectiveness.
   - Inject context cancellations to verify no partial data corruption.
3. Create `tests/e2e/integration_test.go`:
   - Full user story: Admin creates manga (gRPC) -> User subscribes (UDP) -> User chats (WS) -> Syncs (TCP).
</action>
<acceptance_criteria>
- All tests pass in a single `go test ./tests/e2e/...` run.
- Slow consumers do not block the global event bus.
</acceptance_criteria>

### [ ] Task 10.2: Master Demo Show (`run_show.go`)
<action>
1. Implement `demo/run_show.go` in Go:
   - Orchestrates a sequence:
     - Register & Login (HTTP).
     - Connect and listen (TCP & UDP).
     - Subscribe to events (gRPC Stream).
     - Create Manga (gRPC Admin).
     - Verify notifications received on all protocols with colored logs.
     - Send a chat message (WS) and see it globally.
</action>
<acceptance_criteria>
- Running `go run demo/run_show.go` produces a readable, impressive log of 5-protocol interaction.
</acceptance_criteria>

### [ ] Task 10.3: Final Documentation (`ARCHITECTURE.md`)
<action>
1. Create `ARCHITECTURE.md` using Mermaid:
   - **System Overview**: High-level block diagram.
   - **Protocol Matrix**: Comparison of 5 protocols.
   - **Event Propagation**: Sequence diagram showing how `bus.Publish` reaches all endpoints.
   - **Design Decisions**: Detail WAL, Hub Isolation, Graceful Shutdown, and Repository Pattern.
</action>
<acceptance_criteria>
- Comprehensive documentation that clearly explains the "Why" behind the "How".
</acceptance_criteria>

## Verification
- `go test ./tests/e2e/...`
- `go run demo/run_show.go`
- Inspect `ARCHITECTURE.md` rendering.

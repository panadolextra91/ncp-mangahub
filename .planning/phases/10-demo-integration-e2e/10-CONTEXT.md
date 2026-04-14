# Phase 10 Context: Demo Integration & E2E Verification

## Key Decisions (Context)
- **E2E Test Location**: Root directory `tests/e2e/`. These are blackbox tests that run against a fully initialized server instance.
- **Hell Path Scenarios**:
  1. **Slow Consumer (Isolation)**: A TCP/WS client that stops reading must be kicked/timed out without stalling the EventBus.
  2. **DB Stress (Locking)**: Graceful handling of SQLite `database is locked` errors during high concurrency using WAL and proper retries.
  3. **Concurrent Race**: Orchestrated test where all 5 protocols (HTTP, TCP, gRPC, WS, UDP) perform operations and hóng events simultaneously.
- **Master Demo Script**: `demo/run_show.sh`. A shell script that automates a full "User Story" showcase (Create Manga -> Notify -> Chat -> Progress).
- **Documentation**: A comprehensive `ARCHITECTURE.md` file featuring:
  - Mermaid Diagrams of the Modular Monolith & Event Propagation.
  - Deep-dive into design decisions (WAL mode, Channel Hubs, UDP TTL, Graceful Shutdown).
  - Manual verification guide.

## Technical Details
- **Tools**: Use `curl`, `nc`, `grpcurl` (if possible), and `websocat`/`wscat` for the script.
- **E2E Test Environment**: Standard Go tests running on separate random ports to avoid conflicts.

## Gray Areas to Discuss
*None. All final chốt received from Mẹ Architect.*

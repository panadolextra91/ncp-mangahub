# Roadmap

## Phase 1: Project Setup & Database Infrastructure
- **Goal:** Establish standard Go layout, set up SQLite with WAL mode, and implement basic repository contracts.
- **Depends on:** None

## Phase 2: Internal Event Bus Implementation
- **Goal:** Build the non-blocking Go channels pub/sub mechanism to decouple upcoming protocols.
- **Depends on:** Phase 1

## Phase 3: Domain Models & Application Services
- **Goal:** Implement the core business logic (Manga, Auth) devoid of protocol-specific knowledge.
- **Depends on:** Phase 1

## Phase 4: HTTP Protocol Layer (Core API)
- **Goal:** Implement Auth and Manga CRUD. This acts as the *only* DB writer.
- **Depends on:** Phase 3, Phase 2

## Phase 5: TCP Protocol Layer (Real-time Sync)
- **Goal:** Implement persistent connections and act as a subscriber to the Event Bus to broadcast JSON updates.
- **Depends on:** Phase 2

## Phase 6: WebSocket Protocol Layer (Chat)
- **Goal:** Implement real-time chat with rooms grouped by `manga_id`.
- **Depends on:** Phase 2

## Phase 7: UDP Protocol Layer (Notifications)
- **Goal:** Implement a fire-and-forget lightweight notification subscriber.
- **Depends on:** Phase 2

## Phase 8: gRPC Protocol Layer (Admin/CLI)
- **Goal:** Implement the gRPC server and basic client for CLI administrative tooling.
- **Depends on:** Phase 3

## Phase 9: Graceful Shutdown Mechanics
- **Goal:** Wire up `context.Context` and `sync.WaitGroup` to orchestrate clean termination for all 5 protocol servers and the DB.
- **Depends on:** Phase 8, Phase 7, Phase 6, Phase 5, Phase 4

## Phase 10: Demo Integration & E2E Verification
- **Goal:** Comprehensive E2E Testing of the full 5-protocol interaction flow. This phase must focus intensely on "hell paths" alongside the happy path (e.g., dropping SQLite lock, stalling TCP readers, malformed UDP packets, network partitioning).
- **Depends on:** Phase 9

### Phase 11: TUI and HTML interface

**Goal:** [To be planned]
**Requirements**: TBD
**Depends on:** Phase 10
**Plans:** 0 plans

Plans:
- [ ] TBD (run /gsd-plan-phase 11 to break down)

### Phase 12: fullfill spec and use case

**Goal:** [To be planned]
**Requirements**: TBD
**Depends on:** Phase 11
**Plans:** 0 plans

Plans:
- [ ] TBD (run /gsd-plan-phase 12 to break down)

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
- **Status:** ✅ COMPLETED
- **Goal:** Comprehensive E2E Testing of the full 5-protocol interaction flow.
- **Depends on:** Phase 9

## Phase 11: PinkHub TUI Interface
- **Status:** ✅ COMPLETED
- **Goal:** Build a premium Terminal UI (TUI) using BubbleTea to visualize all 5 protocols in one unified dashboard.
- **Depends on:** Phase 10

## Phase 12: Feature Completeness (Spec & Use Case)
- **Status:** ✅ COMPLETED
- **Goal:** Align with academic requirements, implementing gRPC Search, advanced filtering, and full use case coverage.
- **Depends on:** Phase 11

## Phase 13: Web Scraping Integration
- **Status:** ✅ COMPLETED
- **Goal:** Scrape motivational quotes from `quotes.toscrape.com` using standard Go libraries to fulfill the scraping requirement.
- **Depends on:** Phase 12

## Phase 14: Personal Library & TCP Sync
- **Status:** ✅ COMPLETED
- **Goal:** Implement a personal library with 3 reading statuses and real-time TCP synchronization across multiple clients.
- **Depends on:** Phase 13

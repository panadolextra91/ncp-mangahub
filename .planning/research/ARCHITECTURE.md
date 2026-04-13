# Architecture Research

**Domain:** Modular Monolith + Event-Driven Systems

## Component Boundaries
1. **Interfaces Layer (`interfaces/`)**: Isolates the handling of HTTP, TCP, WS, UDP, gRPC protocols.
2. **Event Bus (`eventbus/`)**: The central nervous system. A pub/sub mechanism utilizing Go channels to broadcast events across boundaries.
3. **Application/Domain Layer (`application/`, `domain/`)**: Pure business logic (Manga CRUD, Users) independent of protocols.
4. **Infrastructure Layer (`infrastructure/`)**: DB connections, repositories with single-write enforcement.

## Data Flow
- **Writes:** `Client -> HTTP Handler -> Application -> Repository (SQLite WAL)`
- **Eventing:** `Application -> EventBus Publish -> Interfaces (TCP/WS/UDP Subscribers)`

## Build Order Implication
1. Setup Infrastructure (SQLite WAL config / Repositories) & Domain entities.
2. Build Event Bus (Topic routing, subscribe/publish mechanics with non-blocking drops).
3. Build HTTP layer (Core functionality & Auth).
4. Implement TCP / WS / UDP subscribers reading from Event Bus.
5. Add gRPC Admin APIs.

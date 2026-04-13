# Features Research

**Domain:** MangaHub Multi-Protocol System

## Table Stakes (Must-Have)
- **HTTP REST API:** User authentication (JWT) and basic CRUD operations for Manga logic. Must be the single source of truth for DB writes.
- **Real-Time Sync (TCP):** Persistent connection for broadcasting manga progress updates.
- **Chat Rooms (WebSocket):** Low-latency, bi-directional communication multiplexed by manga rooms.
- **Notifications (UDP):** Fire-and-forget lightweight event distribution (e.g. "New manga added").
- **Admin CLI (gRPC):** Internal robust tooling for system operators.

## Differentiators
- **Internal Event Bus:** Decoupled channels allowing HTTP to trigger TCP, UDP, or WS endpoints without direct function calls, maintaining a clean architecture.

## Anti-Features (Do NOT Build)
- Multi-node clustering / Distributed tracking
- Complex retry queues (Slow consumers will be dropped to maintain system life)

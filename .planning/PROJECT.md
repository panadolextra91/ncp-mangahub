# MangaHub

## What This Is

MangaHub is a backend system designed with a Modular Monolith architecture combined with a lightweight internal Event-Driven mechanism. It integrates multiple protocols (HTTP, TCP, WebSocket, UDP, gRPC) to provide RESTful APIs, real-time sync, chat, notifications, and internal administrative operations, keeping the system simple, stable, and scalable.

## Core Value

Ensure 100% stable demonstration of all 5 protocols working together seamlessly.

## Requirements

### Validated

- [x] Implementation of Modular Monolith structure in Go (Phase 1)
- [x] Internal Event Bus using buffered, non-blocking channels (Phase 2)
- [x] HTTP Protocol for Core REST API (Phase 4)
- [x] TCP Protocol for persistent real-time sync (Phase 5, 14)
- [x] WebSocket Protocol for real-time chat (Phase 6)
- [x] UDP Protocol for lightweight notification broadcasting (Phase 7)
- [x] gRPC Protocol for Search & Admin (Phase 8, 12)
- [x] SQLite setup with WAL and single write connection (Phase 1)
- [x] Graceful shutdown implementing strict cleanup (Phase 9)
- [x] Web Scraping from quotes.toscrape.com (Phase 13)
- [x] Personal Library with 3 reading statuses (Phase 14)

### Active

(All core objectives for Milestone 1 achieved)

### Out of Scope

- Distributed Microservices — The system uses a modular monolith to maintain simplicity while simulating distributed scenarios via internal event buses.
- Dedicated Message Brokers (e.g. RabbitMQ/Kafka) — Replaced by Go internal channels to keep it dependency-lite.
- DB Writing outside of HTTP — Strict rule: only HTTP REST API can authoritatively write to the database (to prevent concurrency and locking issues).

## Context

- Building a highly available API and sync architecture for an academic/demonstration context.
- Prioritizing stability during live demo (Show all 5 protocols working together: Login -> Add/Update -> TCP Sync -> Chat -> UDP notification -> GRPC call).
- Single process `main.go`.

## Constraints

- **Architecture**: Modular Monolith + Event-Driven Lite — Easiest to deploy and explain conceptually.
- **Database**: SQLite with WAL — Minimizes complex deployment footprint while managing concurrency using connection tuning (`SetMaxOpenConns(1)`).
- **Concurrency**: Goroutines per protocol (`httpServer`, `tcpServer`, etc.), rigorously tied to context cancellation for clean exits.
- **Testing (DoD)**: Every phase must include Unit & Integration tests with >80% coverage and "hell cases" explicitly covered. Coverage output must be printable for showcase.

## Key Decisions

| Allow event drop for slow consumers | Prevents blocking the entire event bus and system | ✅ IMPLEMENTED |
| Only HTTP writes to DB | Centralizes mutation state, avoids SQLite locking nightmares | ✅ IMPLEMENTED |
| gRPC for Admin/Search | Better performance and strong typing | ✅ IMPLEMENTED |
| TUI as unified dashboard | Visualize 5 protocols in one place for demo | ✅ IMPLEMENTED |
| In-memory Scraping | Compliance with volatile data requirements | ✅ IMPLEMENTED |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-05-12 after Phase 14 completion*

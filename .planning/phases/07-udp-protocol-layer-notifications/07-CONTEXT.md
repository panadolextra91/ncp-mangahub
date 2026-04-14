# Phase 7: UDP Protocol Layer (Notifications)

## Overview
Phase 7 introduces the fourth protocol layer: UDP. This layer is designed for lightweight, fire-and-forget notifications (e.g., system alerts, ticker updates) where perfect reliability is not required but low overhead is essential.

## Key Decisions (Context)
- **Port**: Default `9191` (Configurable via `UDP_PORT`).
- **Mode**: Fire-and-forget Unicast.
- **Handshake (The "SUB" Packet)**: 
  - Clients must send `SUB <manga_id> <token>` to register their `IP:Port`.
  - Server validates JWT once and stores the peer.
- **Heartbeat & TTL (STRICT)**:
  - Peers expire after **60 seconds** of inactivity.
  - Clients should send `PING <token>` every 30s to stay alive.
  - A background "Garbage Collector" routine must prune the registry.
- **Content Filtering**: 
  - ONLY Global/System events (e.g., `manga.new`). 
  - **CẤM** (FORBIDDEN) to send Chat or Progress updates via UDP (to avoid out-of-order confusion).
- **No Reliability Layer**: Pure UDP. No ACKs, no retransmissions.

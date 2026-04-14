# Phase 8: gRPC Protocol Layer (Admin/CLI)

## Overview
Phase 8 introduces the final protocol layer: gRPC. This is designed for high-performance service-to-service communication and administrative CLI integrations.

## Key Decisions (Context)
- **Port**: Default `50051` (Configurable via `GRPC_PORT`).
- **Proto Architecture (STRICT)**:
  - Layout: All generated code must reside in `pkg/pb` (for public importability).
  - Automation: A `gen_proto.sh` script must be provided.
  - Service Separation:
    - `MangaService`: `GetManga`, `UpdateProgress` (User level).
    - `AdminService`: `CreateManga`, `DeleteManga` (Admin level).
- **Server-side Streaming (STRICT)**:
  - Method: `SubscribeEvents` returning a `stream Event`.
  - Integration: Bridge to the internal `EventBus`.
- **Authentication & Authorization**:
  - Interceptors: Unary and Stream interceptors must extract JWT from Metadata.
  - RBAC: `AdminService` access must be gated by `Admin` role checks returning `codes.PermissionDenied` if unauthorized.

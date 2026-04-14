---
wave: 1
depends_on: [02-internal-event-bus-implementation, 03-domain-models-application-services, 04-http-protocol-layer-core-api]
files_modified:
  - config/config.go
  - api/proto/mangahub.proto
  - gen_proto.sh
  - pkg/pb/mangahub.pb.go (Generated)
  - pkg/pb/mangahub_grpc.pb.go (Generated)
  - internal/interfaces/grpc/interceptors.go
  - internal/interfaces/grpc/manga_service.go
  - internal/interfaces/grpc/admin_service.go
  - cmd/server/main.go
  - internal/interfaces/grpc/grpc_test.go
autonomous: true
---

# Phase 8 Plan: gRPC Protocol Layer (Admin/CLI)

## Objective
Implement a high-performance gRPC layer with separated User and Admin services, secured by JWT interceptors and supporting real-time event streaming.

## Tasks

### [ ] Task 8.1: Proto Definition & Generation
<read_first>
- .planning/phases/08-grpc-protocol-layer-admin-cli/08-RESEARCH.md
</read_first>
<action>
1. Create `api/proto/mangahub.proto` defining `MangaService` and `AdminService`.
2. Create `gen_proto.sh` with correct `pkg/pb` output paths.
3. Run generation.
</action>
<acceptance_criteria>
- Proto defines `GetManga`, `UpdateProgress`, `SubscribeEvents`, `CreateManga`, `DeleteManga`.
- Binary code generated in `pkg/pb`.
</acceptance_criteria>

### [ ] Task 8.2: Security Interceptors
<action>
1. Implement `internal/interfaces/grpc/interceptors.go`.
2. Create `AuthUnaryInterceptor` and `AuthStreamInterceptor`.
3. Check role authorization for `AdminService`.
</action>
<acceptance_criteria>
- Non-admin tokens are rejected with `codes.PermissionDenied` for admin methods.
- Invalid tokens return `codes.Unauthenticated`.
</acceptance_criteria>

### [ ] Task 8.3: Service Implementation
<action>
1. Implement `MangaService` in `internal/interfaces/grpc/manga_service.go`.
2. Implement `AdminService` in `internal/interfaces/grpc/admin_service.go`.
3. Specialized logic for `SubscribeEvents` to bridge `EventBus` to the GRPC stream.
</action>
<acceptance_criteria>
- Unary calls correctly interact with application services.
- `SubscribeEvents` correctly pushes real-time events to the stream.
</acceptance_criteria>

### [ ] Task 8.4: Integration & DoD Verification
<action>
1. Update `config/config.go` for `GRPCPort` (50051).
2. Update `cmd/server/main.go` to start the gRPC server.
3. Write integration tests in `internal/interfaces/grpc/grpc_test.go`.
4. Verify global coverage >= 80%.
</action>
<acceptance_criteria>
- gRPC server starts and handles requests correctly.
- Streaming events are verified in tests.
</acceptance_criteria>

## Verification
- Run `./test.sh`.
- Manual test:
  1. Start server.
  2. Use `grpcurl` or a custom client to subscribe to events.
  3. Create manga via HTTP.
  4. Verify event appears on gRPC stream.

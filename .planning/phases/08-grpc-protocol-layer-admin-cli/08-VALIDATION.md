# Phase 8 Validation Strategy
**Status**: Pending

## Objective
To verify that the gRPC Protocol Layer correctly implements a high-performance, secured service-to-service communication layer with proper access control and real-time streaming.

## Verification Checklist

### 1. Proto Generation
- [ ] `gen_proto.sh` executes without errors.
- [ ] `pkg/pb/` contains the generated `.go` files for both services.

### 2. Service Logic (Unary)
- [ ] `GetManga` returns correct manga details.
- [ ] `UpdateProgress` correctly updates the database and publishes to the EventBus.
- [ ] `CreateManga` and `DeleteManga` (AdminService) work as expected for admin users.

### 3. Security & Authorization (RBAC)
- [ ] Unary calls without JWT are rejected with `codes.Unauthenticated`.
- [ ] Stream calls without JWT are rejected with `codes.Unauthenticated`.
- [ ] `AdminService` calls by non-admin users return `codes.PermissionDenied`.

### 4. Server-side Streaming (Hell Path)
- [ ] `SubscribeEvents` maintains a long-lived connection.
- [ ] When a manga is created via HTTP, the gRPC stream receives a notification immediately.
- [ ] When a client disconnects, the server correctly cleans up the EventBus subscription (no goroutine leak).

### 5. Definition of Done
- [ ] Total statement coverage remains >= 80%.
- [ ] `test.sh` passes 100%.

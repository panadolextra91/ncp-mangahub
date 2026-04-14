# Phase 8 Research: gRPC Protocol Layer (Admin/CLI)

## 1. Proto Definitions
```proto
syntax = "proto3";
package mangahub.v1;
option go_package = "github.com/user/mangahub/pkg/pb";

service MangaService {
    rpc GetManga(GetMangaRequest) returns (MangaResponse);
    rpc UpdateProgress(UpdateProgressRequest) returns (UpdateProgressResponse);
    rpc SubscribeEvents(Empty) returns (stream EventNotification);
}

service AdminService {
    rpc CreateManga(CreateMangaRequest) returns (MangaResponse);
    rpc DeleteManga(DeleteMangaRequest) returns (DeleteMangaResponse);
}
```

## 2. Auth Interceptors
- **Unary**: `authUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler)`.
- **Stream**: `authStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler)`.
- Metadata extraction: `md, ok := metadata.FromIncomingContext(ctx)`. Token is usually in `authorization` key.

## 3. Streaming Logic
- Subscribe to `EventBus`.
- Loop and `stream.Send()`.
- **Crucial**: Use `stream.Context().Done()` to detect client disconnection and unsubscribe from the bus to avoid goroutine leaks.

## 4. Automation
- `gen_proto.sh`:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/proto/mangahub.proto
```
Wait, Mẹ wants it in `pkg/pb`. I'll adjust the paths.

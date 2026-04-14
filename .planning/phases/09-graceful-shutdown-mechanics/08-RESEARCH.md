# Phase 9 Research: Graceful Shutdown Mechanics

## 1. signal.NotifyContext
We should replace the manual signal handling with `signal.NotifyContext`:
```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
```

## 2. sync.WaitGroup
We need to track:
- 5 Protocol Servers: HTTP, TCP, WebSocket, UDP, gRPC.
- 5 Hubs/Backends: TCP Hub, WS Hub, UDP GC, gRPC Services (streaming).
- 3 Bridge Goroutines.

## 3. Hub Stop Logic
Both `TCP Hub` and `WS Hub` need a way to break their `Run()` loop and notify clients.
```go
func (h *Hub) Run(ctx context.Context, wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        select {
        case <-ctx.Done():
            // 1. Notify Clients
            // 2. Cleanup
            return
        // ... other cases
        }
    }
}
```

## 4. Sequence in main.go
```go
<-ctx.Done() // Signal received
shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// 1. Stop Listeners
httpServer.Shutdown(shutdownCtx)
grpcServer.GracefulStop()
// ...

// 2. Wait for completion
wg.Wait()

// 3. Final cleanup
db.Close()
```

## 5. Client Messages
- TCP: `{"type": "system", "action": "shutdown", "reason": "server_maintenance"}\n`
- WS: Same JSON.
- UDP: No message needed (fire-and-forget).
- gRPC: Use `GracefulStop()` (it handles it).

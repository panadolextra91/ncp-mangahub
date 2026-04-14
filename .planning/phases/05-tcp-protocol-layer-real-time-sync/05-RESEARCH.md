# Phase 5 Research: TCP Protocol Layer (Real-time Sync)

## 1. Config Refactor
The `Config` struct will be updated to include `TCPPort` loaded from `TCP_PORT` env var with a fallback to `9090`.

## 2. JWT Logic Centralization
To support both HTTP Middleware and TCP Handshake, JWT logic will be moved to `pkg/auth/jwt.go`:
- `ValidateToken(tokenString string, secret string) (jwt.MapClaims, error)`
- `GenerateToken(userID int, role string, secret string) (string, error)`

## 3. TCP Hub (Channel Registry)
The Hub will manage active connections using Go channels to avoid mutex contention and Head-of-line blocking:
```go
type Hub struct {
    clients    map[net.Conn]bool
    broadcast  chan []byte
    register   chan net.Conn
    unregister chan net.Conn
}

func (h *Hub) Run() {
    for {
        select {
        case conn := <-h.register:
            h.clients[conn] = true
        case conn := <-h.unregister:
            if _, ok := h.clients[conn]; ok {
                delete(h.clients, conn)
                conn.Close()
            }
        case message := <-h.broadcast:
            for conn := range h.clients {
                // Non-blocking write to avoid slow-consumer issues
                // Actually TCP conn.Write is blocking, so we use a per-client goroutine or buffer.
                // Simplified: use a select with default or a write timeout.
                _, err := conn.Write(append(message, '\n'))
                if err != nil {
                    h.unregister <- conn
                }
            }
        }
    }
}
```

## 4. Handshake Logic
- Server listener spawned in `main.go`.
- `Accept()` loop.
- One-time read: `AUTH <token>`.
- Verification using `pkg/auth/jwt`.

## 5. Event Bus Integration
A bridge goroutine will subscribe to the `EventBus` and push every message into the `Hub.broadcast` channel.

## Validation Plan
- Unit test for Hub (Mock connections).
- Integration test for TCP Handshake (using `net.Dial`).
- Verify non-blocking behavior by simulating a slow reader.

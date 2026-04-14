# Phase 7 Research: UDP Protocol Layer (Notifications)

## 1. Configuration
- `UDPPort` env: `UDP_PORT`, fallback `9191`.

## 2. Registry (Heartbeat Logic)
Wait, `sync.Map` is good, but iterate-and-delete on it is slightly clunky for precise TTL.
However, for 100-1000 clients, a simple loop every 5-10s is fine.

```go
type Peer struct {
    Addr      *net.UDPAddr
    MangaID   int
    LastSeen  time.Time
}

type Registry struct {
    peers sync.Map // string(addr) -> *Peer
}
```

## 3. Package Format
Mẹ requested:
- `SUB <manga_id> <jwt>`
- `PING <jwt>` (Wait, why JWT on PING? Mẹ said: security. Okay, we validate).

## 4. UDP Server Loop
```go
func (s *Server) Start() {
    conn, _ := net.ListenUDP("udp", s.addr)
    buf := make([]byte, 1024)
    for {
        n, addr, _ := conn.ReadFromUDP(buf)
        data := string(buf[:n])
        // Parse "SUB ..." or "PING ..."
        // On SUB: validate JWT, registry.Register(mangaID, addr)
        // On PING: validate JWT, registry.KeepAlive(addr)
    }
}
```

## 5. Event Delivery
When the EventBus bridge sees a `manga.new`:
```go
registry.Broadcast(mangaID, payload)
```
Broadcast logic:
```go
registry.peers.Range(func(key, value interface{}) bool {
    peer := value.(*Peer)
    if peer.MangaID == targetMangaID || targetMangaID == 0 { // 0 for global?
         conn.WriteToUDP(payload, peer.Addr)
    }
    return true
})
```

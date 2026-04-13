# 🚀 MangaHub – FINAL ARCHITECTURE v2.0

---

# 1. 🧠 Architectural Style

## ✅ Modular Monolith + Event-Driven Lite

* 1 process (`main.go`)
* chia module rõ ràng
* giao tiếp nội bộ qua Event Bus (Go channels)

---

## 🎯 Design Goals

* ✔️ Đáp ứng đủ 5 protocol
* ✔️ Stable khi demo (ưu tiên số 1)
* ✔️ Dễ explain với thầy
* ✔️ Có hướng scale (theoretical)

---

## 🧠 One-liner để defend

> “We use a modular monolith with an internal event-driven mechanism to integrate multiple protocols while keeping the system simple, stable, and scalable.”

---

# 2. 🧱 System Structure

```text
mangahub/
├── cmd/
│   └── server/main.go

├── internal/
│   ├── domain/              # Business logic
│   ├── application/         # Use cases
│   ├── infrastructure/      # DB, logger, config
│   ├── interfaces/          # HTTP, TCP, UDP, WS, gRPC
│   └── eventbus/            # Pub/Sub

├── pkg/models/
```

---

# 3. 🔄 Event Bus (FINAL VERSION)

## 🎯 Purpose

* decouple protocol
* simulate distributed system

---

## ✅ Design

### Buffered + Non-blocking

```go
type Event struct {
    Topic   string
    Payload interface{}
}
```

---

### Subscribe

```go
func Subscribe(topic string) <-chan Event
```

---

### Publish (CRITICAL FIX)

```go
for _, ch := range subscribers {
    select {
    case ch <- event:
    default:
        log.Println("Dropping event (slow consumer)")
    }
}
```

---

## 📌 Topics

| Topic            | Flow        |
| ---------------- | ----------- |
| progress.updated | HTTP → TCP  |
| manga.new        | Admin → UDP |
| chat.message     | WS → WS     |

---

## 🧠 Design decision

> ⚠️ Allow event drop → tránh block hệ thống

---

# 4. 🌐 Protocol Responsibilities (FINAL)

---

## 🌐 HTTP (REST API) – Core API

### Role:

* Auth (JWT)
* Manga CRUD
* Progress update

### Rule:

> 🔒 ONLY HTTP writes to DB

---

## 🔗 TCP – Real-time Sync

* persistent connection
* nhận event từ bus
* broadcast JSON

---

## 💬 WebSocket – Chat

* real-time chat
* room theo manga_id

---

## 📡 UDP – Notification

* fire-and-forget
* broadcast event nhẹ

---

## ⚡ gRPC – CLI / Admin API (FIXED)

### Role:

* phục vụ CLI tool
* internal admin operations

### ❌ NOT dùng internal call

---

## 🧠 Defend line

> “gRPC is used for CLI and admin operations due to better performance and strong typing, while HTTP serves general clients.”

---

# 5. 💾 Database Strategy (FINAL FIX)

---

## ❌ REMOVE writeChan

---

## ✅ Use SQLite + WAL

```sql
PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL;
```

---

## ✅ Connection control

```go
db.SetMaxOpenConns(1)
db.SetMaxIdleConns(1)
```

---

## 🧠 Reasoning

> “We rely on SQLite native locking with WAL and connection tuning instead of building a custom write queue.”

---

## ✅ Repository Pattern

* clean separation
* dễ test
* tránh gọi DB trực tiếp từ handler

---

# 6. 🧵 Concurrency Model

---

## Each protocol = goroutine

```go
go httpServer.Start(ctx)
go tcpServer.Start(ctx)
go udpServer.Start(ctx)
go wsServer.Start(ctx)
go grpcServer.Start(ctx)
```

---

## Tools

* `context.Context` → lifecycle
* `WaitGroup` → shutdown sync

---

# 7. 🛑 Graceful Shutdown

---

## Flow

1. Ctrl+C
2. cancel context
3. stop accepting connection
4. close event bus
5. wait goroutines
6. close DB

---

## 🧠 Defend

> “Ensures no resource leaks and clean termination.”

---

# 8. ⚠️ Critical Safety Fixes

---

## ❗ Event Bus Blocking

✔️ buffered
✔️ non-blocking publish

---

## ❗ SQLite Lock

✔️ WAL
✔️ single connection

---

## ❗ Goroutine Leak

✔️ `defer conn.Close()`
✔️ `ctx.Done()`

---

## ❗ Event Ordering

✔️ last_write_wins (timestamp)

---

## ❗ Config

```yaml
http: 8080
tcp: 9090
udp: 9091
grpc: 9092
ws: 9093
```

---

# 9. 🗺️ Implementation Roadmap

---

## 🥇 Phase 1 (Week 1–3)

* DB + HTTP

---

## 🥈 Phase 2 (Week 4–6)

* TCP + WebSocket

---

## 🥉 Phase 3 (Week 7–8)

* UDP + gRPC

---

## 🧪 Phase 4 (Week 9–11)

* Integration
* CLI test
* bug fix

---

# 10. 🎤 Demo Strategy (IMPORTANT)

---

## Flow demo

1. Login (HTTP)
2. Add manga
3. Update progress
4. 👉 TCP realtime sync
5. Chat (WebSocket)
6. UDP notification
7. gRPC CLI call

---

## 🎯 Goal

> Show all 5 protocols working together

---

# 🏁 FINAL CONCLUSION

---

## 💥 System đạt:

* ✔️ Full spec compliance
* ✔️ Clean architecture
* ✔️ Stable demo
* ✔️ Clear separation
* ✔️ Strong reasoning

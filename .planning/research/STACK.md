# Stack Research

**Domain:** Go Modular Monolith + Event-Driven Multi-Protocol

## Core Stack
- **Language:** Go (1.21+)
  - *Rationale:* Native concurrency (goroutines/channels) makes it the perfect fit for multi-protocol servers. Standard library `net/http` and `net` are robust enough for most protocols without heavy dependencies.
- **Database:** SQLite3 (`github.com/mattn/go-sqlite3`)
  - *Rationale:* Serverless, easy to deploy. `PRAGMA journal_mode=WAL` allows concurrent readers and prevents database locks when configured correctly with `SetMaxOpenConns(1)`.
- **gRPC:** `google.golang.org/grpc` & `google.golang.org/protobuf`
  - *Rationale:* Industry standard for CLI/Admin fast communication.
- **WebSocket:** `github.com/gorilla/websocket`
  - *Rationale:* Most stable and widely used WebSocket implementation in Go, better than using the standard library's baseline.

## What NOT to use
- **Heavy web frameworks (e.g., Gin, Fiber):** The internal HTTP layer is a simple component. Standard `net/http` + `http.ServeMux` (or `chi` for lightweight routing) is sufficient to avoid bloating the monolith.
- **Message Brokers (RabbitMQ/Kafka):** Defeats the purpose of a simple event-driven monolith. Standard Go buffered channels are enough.

# 🌸 MangaHub - The Pink Professional Monolith

MangaHub is a comprehensive manga tracking and management system built on **Clean Architecture** principles. It serves as a practical demonstration of integrating **5 essential network protocols** into a single cohesive system: **HTTP, gRPC, WebSocket, TCP, and UDP**.

Experience the ultimate **"TUI-First"** environment with our stylish Pink Pastel Terminal Interface.

---

## 🚀 Key Features

- **🌸 PinkHub TUI**: A high-performance terminal UI featuring ASCII animations, smooth scrolling lists, and role-aware components.
- **🔐 Multi-Protocol Core**:
    - **HTTP**: Handles JWT Authentication and standard CRUD operations.
    - **gRPC**: High-efficiency internal search and administrative services.
    - **WebSocket**: Real-time community chat hub.
    - **TCP (Raw)**: Instant cross-device reading progress synchronization.
    - **UDP**: Fast, lightweight "Fire-and-forget" release notifications.
- **🛡️ RBAC (Role-Based Access Control)**: Strict identity management. Administrative functions are dynamically gated based on user roles.
- **📚 Personal Library**: Save titles to your personal shelf, view details, and track reading progress across three statuses (`reading`, `completed`, `plan_to_read`).
- **🕸️ Web Scraping**: Integrated quote-of-the-day scraper that pulls inspirational quotes from `quotes.toscrape.com` for the TUI dashboard.
- **🔍 Advanced Search & Filtering**: Multi-criteria search on `GET /api/manga` and gRPC `SearchManga`: filter by multiple genres (OR), status, and sort by `title` or `recent`. Capped at 10 genres per request.
- **🩺 Multi-Protocol Health Checks**: `GET /api/health` actively probes all 5 protocols (HTTP, TCP, UDP, WS, gRPC) and the SQLite DB in parallel with a 500ms per-probe timeout. Returns `200/ok` or `503/degraded` with per-subsystem status. Every request is logged to the server console (🩺 prefix).
- **🌱 Pre-seeded Catalog**: Over 100 popular manga series with rich metadata, ready for instant demo.
- **🔋 Graceful Shutdown**: A 5-step termination protocol ensuring data integrity and clean resource cleanup.

---

## 🛠️ Quick Start

For a quick demo, we provide pre-built binaries:

1.  **Launch the Core Server**:
    ```bash
    ./server
    ```
2.  **Launch the TUI Client**:
    ```bash
    ./tui
    ```

*To run from source:*
```bash
go run cmd/server/main.go  # Start Server
go run cmd/client/main.go  # Start TUI
```

3. **Run with Docker (Recommended for Demo)**:
    ```bash
    docker-compose up --build
    ```
    *Note: Run the TUI Client directly on the host for the best terminal rendering experience.*

---

## 📖 Technical Documentation

Dive deeper into our system design:
- [🏗️ System Architecture](./ARCHITECTURE.md)
- [📜 API Contract & Protocol Flows](./API_CONTRACT.md)

---

## 🧪 Testing

MangaHub includes a rigorous E2E test suite in `tests/e2e/`, verifying system stability and cross-protocol data consistency.

```bash
go test -v ./tests/e2e/...
```

---

Crafted with passion by **Mẹ Architect & Antigravity**. Enjoy the "Pink Professional" experience! 🌸🤖

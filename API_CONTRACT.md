# 📜 MangaHub API Contract & Flows

Tài liệu này định nghĩa tất cả các điểm giao tiếp của hệ thống MangaHub, minh họa bằng các sơ đồ Sequence Diagram để mô tả luồng dữ liệu qua 5 giao thức: **HTTP, gRPC, WebSocket, TCP, UDP**.

---

## 🔐 1. Authentication & Identity (HTTP)

Hệ thống sử dụng cơ chế **Auto-Login/Register** để tối ưu hóa trải nghiệm người dùng trên TUI.

```mermaid
sequenceDiagram
    participant TUI as 🖥️ TUI Client
    participant API as 🌐 HTTP API
    participant DB as 🗄️ SQLite

    TUI->>API: POST /api/auth/login (Username, Password)
    API->>DB: SELECT user WHERE username = ?
    alt User Not Found
        API->>DB: INSERT INTO users (Register)
    end
    API->>API: Generate JWT Token (Role: admin/user)
    API-->>TUI: 200 OK + JWT Token + Role
```

---

## 📚 2. Manga Management (HTTP & gRPC)

### 2.1 Create Manga (gRPC - Admin Only)
Dùng cho các hệ thống quản trị nội bộ hoặc CLI Admin.

```mermaid
sequenceDiagram
    participant Admin as 🔑 Admin CLI
    participant GRPC as ⚡ gRPC Server
    participant Bus as 🚌 Event Bus
    participant UDP as 📡 UDP Notifier

    Admin->>GRPC: AdminService/CreateManga (Proto Request)
    GRPC->>GRPC: Verify Admin Role (JWT Interceptor)
    GRPC->>Bus: Publish "manga.new"
    Bus-->>UDP: Trigger Broadcast
    UDP->>UDP: UDP Sendto (All Clients)
    GRPC-->>Admin: Manga Created (Proto Response)
```

### 2.2 Search Manga (HTTP & gRPC)
Hỗ trợ cả Web (HTTP) và App (gRPC).

- **HTTP**: `GET /api/manga?q={query}`
- **gRPC**: `MangaService/SearchManga(SearchRequest)`

---

## 🔄 3. Progress & Sync (HTTP & TCP)

Đây là sự kết hợp giữa **Stateful (TCP)** và **Stateless (HTTP)**.

```mermaid
sequenceDiagram
    participant TUI as 🖥️ TUI Client
    participant API as 🌐 HTTP API
    participant Bus as 🚌 Event Bus
    participant TCP as 🔄 TCP Sync Server
    participant TUI2 as 📱 Other Device

    TUI->>API: PUT /api/manga/progress (Chapter X)
    API->>Bus: Publish "progress.updated"
    Bus-->>TCP: Forward Update
    TCP->>TUI2: TCP Push (JSON: User Y updated to Chapter X)
    API-->>TUI: 200 OK
```

---

## 💬 4. Real-time Communication (WebSocket)

Hệ thống Chat tập trung (Centralized Hub).

```mermaid
sequenceDiagram
    participant C1 as 📱 User Alice
    participant WS as 🔌 WebSocket Hub
    participant C2 as 📱 User Bob

    C1->>WS: Send Message (WS Frame)
    WS->>WS: Validate Session & Role
    WS->>WS: Log to DB
    WS->>C2: Broadcast (JSON: Alice: "Hello")
    WS->>C1: Echo Broadcast
```

---

## 📊 Summary Protocol Table

| Feature | Protocol | Method/Port | Role Required | Interaction Pattern |
| :--- | :--- | :--- | :--- | :--- |
| **Auth** | HTTP | `POST /login` | Any | Request-Response |
| **Search** | HTTP / gRPC | `GET` / `RPC` | Any | Request-Response |
| **Create** | gRPC / HTTP | `RPC` / `POST` | **Admin** | Request-Response + Broadcast |
| **Sync** | TCP | `Port 9090` | User | Server Push |
| **Notify** | UDP | `Port 9191` | None | Fire-and-Forget |
| **Chat** | WebSocket | `/ws/chat` | User | Full Duplex |

---

## 🛠️ Internal SQL Map

- **Fuzzy Search**: `SELECT * FROM mangas WHERE LOWER(title) LIKE LOWER(?)`
- **Upsert Progress**: `INSERT INTO user_progress (...) ON CONFLICT(...) DO UPDATE`
- **Admin Check**: `SELECT role FROM users WHERE id = ?`

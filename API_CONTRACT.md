# 📜 MangaHub API Contract & Flows

Tài liệu này định nghĩa tất cả các điểm giao tiếp của hệ thống MangaHub, bao gồm luồng xử lý (Flow) và các câu lệnh truy vấn dữ liệu (SQL).

---

## 🔐 1. Authentication Module (HTTP)

### 1.1 Login / Register Auto-flow
- **Endpoint**: `POST /api/auth/login`
- **Logic**: Nếu User chưa tồn tại, hệ thống tự động đăng ký (Register) rồi mới đăng nhập (Login).
- **Flow**:
    1. HTTP Handler nhận Username/Password.
    2. Auth Service kiểm tra trong DB.
    3. Nếu không thấy -> Thực hiện `INSERT` mới.
    4. Nếu thấy -> Kiểm tra Password Hash.
    5. Trả về **JWT Token** chứa `UserID`, `Username`, và `Role`.
- **SQL Query**:
    ```sql
    -- Kiểm tra User
    SELECT id, username, password_hash, role FROM users WHERE username = ?;
    -- Tạo User mới (nếu cần)
    INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?);
    ```

---

## 📚 2. Manga Management (HTTP & gRPC)

### 2.1 Create New Manga
- **Endpoint**: `POST /api/manga` (HTTP) hoặc `AdminService/CreateManga` (gRPC)
- **Security**: Chỉ User có `Role = admin` mới được phép.
- **Flow**:
    1. Handler trích xuất `Role` từ JWT Token.
    2. Manga Service kiểm tra quyền -> Gọi Repository.
    3. Lưu vào DB -> Phát sự kiện `manga.new` vào **Event Bus**.
    4. **UDP Server** bắt sự kiện -> Broadcast đến mọi người dùng.
- **SQL Query**:
    ```sql
    INSERT INTO mangas (title, author, genres, status, total_chapters, description) 
    VALUES (?, ?, ?, ?, ?, ?);
    ```

### 2.2 Search / List Manga
- **Endpoint**: `GET /api/manga?q={query}`
- **Flow**: Query string `q` được chuẩn hóa thành lowercase để tìm kiếm mờ (fuzzy search).
- **SQL Query**:
    ```sql
    SELECT * FROM mangas 
    WHERE LOWER(title) LIKE LOWER(?) OR LOWER(author) LIKE LOWER(?) 
    ORDER BY id DESC;
    ```

---

## 💬 3. Real-time Interactions (WebSocket)

### 3.1 Chat Module
- **Endpoint**: `WS /ws/chat`
- **Flow**:
    1. TUI Client thực hiện Handshake với HTTP Server.
    2. Server nâng cấp kết nối lên **WebSocket**.
    3. Server lưu trữ Connection gắn với `UserID` và `Username` (lấy từ Token).
    4. Khi nhận Message -> Server đóng gói kèm `SenderName` (Username) và phát cho tất cả mọi người.
- **SQL Query**: (Tin nhắn được lưu vào DB để tra cứu sau)
    ```sql
    INSERT INTO chat_messages (manga_id, user_id, sender_name, content) 
    VALUES (?, ?, ?, ?);
    ```

---

## 🔄 4. Synchronization & Notifications

### 4.1 Progress Sync (TCP)
- **Port**: `9090` (Raw TCP)
- **Purpose**: Đảm bảo tiến độ đọc truyện của Mẹ luôn được đồng bộ "ngầm" giữa các thiết bị.
- **Flow**: Client duy trì kết nối TCP -> Khi có update từ bất kỳ đâu, Server đẩy payload qua TCP socket.

### 4.2 Notifications (UDP)
- **Port**: `9191` (UDP)
- **Purpose**: Thông báo "Nóng" về truyện mới.
- **Flow**: Server dùng cơ chế **Fire-and-Forget**. Khi có truyện mới, Server gửi một gói tin UDP đến tất cả IP đang lắng nghe. Không cần handshake, tốc độ cực nhanh.

---

## 📊 5. User Progress (HTTP)

### 5.1 Update Reading Progress
- **Endpoint**: `PUT /api/manga/progress`
- **Payload**: `{ "manga_id": 44, "current_chapter": 10, "status": "reading" }`
- **SQL Query**:
    ```sql
    INSERT INTO user_progress (user_id, manga_id, current_chapter, status) 
    VALUES (?, ?, ?, ?)
    ON CONFLICT(user_id, manga_id) DO UPDATE SET 
    current_chapter = excluded.current_chapter, 
    status = excluded.status, 
    updated_at = CURRENT_TIMESTAMP;
    ```

---

## 💡 Summary Table for Demo

| Feature | Protocol | Auth Required | Key SQL Action |
| :--- | :--- | :--- | :--- |
| **Login** | HTTP | No | SELECT/INSERT user |
| **Search** | HTTP | Yes | SELECT LIKE |
| **Create** | HTTP/gRPC | Admin | INSERT manga |
| **Chat** | WS | Yes | INSERT chat_msg |
| **Notify** | UDP | No | N/A (Broadcast) |
| **Sync** | TCP | Yes | N/A (Push) |

---

## 📍 Endpoint Source Code Locations

Nếu Mẹ muốn soi code thực tế của các API này, Mẹ hãy tìm ở các "tọa độ" sau nhé:

### 🌐 HTTP & WebSocket (Port 8080)
- **Định nghĩa Route**: `cmd/server/main.go` (Tìm đoạn `http.HandleFunc`)
- **Logic xử lý Manga/Auth**: `internal/interfaces/http/handlers.go`
- **Logic WebSocket Chat**: `internal/interfaces/ws/handlers.go`

### 🔐 gRPC Admin (Port 50052)
- **Định nghĩa Service**: `api/proto/manga.proto`
- **Logic xử lý Server**: `internal/interfaces/grpc/server.go`

### 🔄 TCP Sync (Port 9090)
- **Logic Server**: `internal/interfaces/tcp/server.go`

### 🔔 UDP Notifications (Port 9191)
- **Logic Server**: `internal/interfaces/udp/server.go`

### 🌸 TUI Client (CMD)
- **Logic App chính**: `internal/interfaces/tui/app.go`
- **Giao diện & Model**: `internal/interfaces/tui/model.go`

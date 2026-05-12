# 🌸 MangaHub - The Pink Professional Monolith

MangaHub là một hệ thống quản lý và theo dõi truyện tranh (Manga/Comic) được xây dựng trên nền tảng **Clean Architecture**, tập trung vào việc minh họa và phối hợp nhuần nhuyễn **5 giao thức mạng** cốt lõi: **HTTP, gRPC, WebSocket, TCP, và UDP**.

Dự án này không chỉ là một bài tập kỹ thuật, mà là một trải nghiệm **"TUI-First"** với giao diện Terminal màu hồng (Pink Pastel) đầy phong cách và hiệu năng cao.

---

## 🚀 Tính năng nổi bật

- **🌸 PinkHub TUI**: Giao diện dòng lệnh (Terminal) mạnh mẽ, hỗ trợ hiệu ứng ASCII, cuộn danh sách (Scrolling) và phân quyền người dùng.
- **🔐 Multi-Protocol Core**:
    - **HTTP**: Quản lý Authentication (JWT) và các thao tác CRUD dữ liệu.
    - **gRPC**: Giao thức quản trị Admin và tìm kiếm nội bộ hiệu suất cao.
    - **WebSocket**: Hệ thống Chat thời gian thực cho cộng đồng.
    - **TCP (Raw)**: Đồng bộ hóa tiến độ đọc truyện đa thiết bị tức thời.
    - **UDP**: Hệ thống thông báo (Notification) truyện mới cực nhanh.
- **🛡️ RBAC (Role-Based Access Control)**: Phân quyền Admin/User nghiêm ngặt. Chỉ Admin mới thấy và sử dụng được tính năng gRPC/Create Manga.
- **📚 100+ Manga Library**: Kho dữ liệu được Seed sẵn hơn 100 bộ truyện nổi tiếng nhất thế giới.
- **🔋 Graceful Shutdown**: Quy trình tắt máy 5 bước đảm bảo an toàn dữ liệu và giải phóng tài nguyên mạng.

---

## 🛠️ Hướng dẫn khởi chạy nhanh

Dành cho buổi Demo nhanh, chúng tôi đã chuẩn bị sẵn các file nhị phân (Binary):

1.  **Khởi động Server**:
    ```bash
    ./server
    ```
2.  **Khởi động TUI Client**:
    ```bash
    ./tui
    ```

*Nếu muốn chạy từ mã nguồn (Source code):*
```bash
go run cmd/server/main.go  # Start Server
go run cmd/client/main.go  # Start TUI
```

---

## 📖 Tài liệu kỹ thuật chi tiết

Để hiểu sâu hơn về "nội công" của MangaHub, vui lòng tham khảo:
- [🏗️ Kiến trúc hệ thống (Architecture)](./ARCHITECTURE.md)
- [📜 Hợp đồng API & Các giao thức (API Contract)](./API_CONTRACT.md)

---

## 🧪 Kiểm thử (Testing)

Hệ thống đi kèm bộ E2E Test cực kỳ khắt khe tại `tests/e2e/`. Luồng test sẽ kiểm tra khả năng chịu tải và tính nhất quán dữ liệu xuyên giao thức.

```bash
go test -v ./tests/e2e/...
```

---

Dự án được thực hiện với tâm huyết bởi **Mẹ Architect & Antigravity**. Chúc mọi người có trải nghiệm "Pink Professional" tuyệt vời nhất! 🌸🤖

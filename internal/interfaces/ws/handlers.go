package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/user/mangahub/internal/application"
	"github.com/user/mangahub/pkg/auth"
	"github.com/user/mangahub/pkg/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // For demo purposes
	},
}

type ChatHandler struct {
	hub         *Hub
	chatService application.ChatService
	jwtSecret   string
}

func NewChatHandler(hub *Hub, chatSvc application.ChatService, secret string) *ChatHandler {
	return &ChatHandler{
		hub:         hub,
		chatService: chatSvc,
		jwtSecret:   secret,
	}
}

func (h *ChatHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	// 1. Parse Query Params
	mangaIDStr := r.URL.Query().Get("manga_id")
	mangaID, err := strconv.Atoi(mangaIDStr)
	if err != nil {
		http.Error(w, "Invalid manga_id", http.StatusBadRequest)
		return
	}

	token := r.URL.Query().Get("token")
	claims, err := auth.ValidateToken(token, h.jwtSecret)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Upgrade Connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WS Upgrade error: %v", err)
		return
	}

	client := &Client{
		Hub:     h.hub,
		MangaID: mangaID,
		Conn:    conn,
		Send:    make(chan *models.ChatMessage, 256),
	}

	h.hub.Register <- client
	log.Printf("🌐 [WS] New Connection: UserID %d (MangaID: %d)", claims.UserID, mangaID)

	// 3. Send History (Bắt buộc lôi 20 tin gần nhất quăng vào mặt nó)
	history, err := h.chatService.GetHistory(mangaID)
	if err == nil {
		for _, msg := range history {
			msgJSON, _ := json.Marshal(msg)
			conn.WriteMessage(websocket.TextMessage, msgJSON)
		}
	}

	// 4. Start Goroutines for Read/Write
	go h.writePump(client)
	go h.readPump(client, claims.UserID, claims.Username)
}

func (h *ChatHandler) readPump(client *Client, userID int, username string) {
	conn := client.Conn.(*websocket.Conn)
	defer func() {
		h.hub.Unregister <- client
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		chatMsg := &models.ChatMessage{
			MangaID:    client.MangaID,
			UserID:     userID,
			SenderName: username,
			Content:    string(message),
			CreatedAt:  time.Now(),
		}

		log.Printf("💬 [CHAT] Message from UserID %d: %s", userID, chatMsg.Content)
		err = h.chatService.SendMessage(chatMsg)
		if err != nil {
			log.Printf("Failed to save/broadcast chat message: %v", err)
		}
	}
}

func (h *ChatHandler) writePump(client *Client) {
	conn := client.Conn.(*websocket.Conn)
	defer conn.Close()

	for msg := range client.Send {
		msgJSON, _ := json.Marshal(msg)
		err := conn.WriteMessage(websocket.TextMessage, msgJSON)
		if err != nil {
			break
		}
	}
}

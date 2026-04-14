package ws

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/user/mangahub/pkg/models"
)

// Client represents a single WebSocket connection.
type Client struct {
	Hub     *Hub
	MangaID int
	Conn    interface{} // Will be *websocket.Conn
	Send    chan *models.ChatMessage
}

// Hub manages active WebSocket connections partitioned by MangaID.
type Hub struct {
	// Registered clients by MangaID.
	rooms map[int]map[*Client]bool
	mu    sync.RWMutex

	// Inbound messages from the ChatService or other sources.
	Broadcast chan *models.ChatMessage

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan *models.ChatMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		rooms:      make(map[int]map[*Client]bool),
	}
}

func (h *Hub) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			h.stop()
			return
		case client := <-h.Register:
			h.mu.Lock()
			if h.rooms[client.MangaID] == nil {
				h.rooms[client.MangaID] = make(map[*Client]bool)
			}
			h.rooms[client.MangaID][client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if room := h.rooms[client.MangaID]; room != nil {
				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.Send)
					if len(room) == 0 {
						delete(h.rooms, client.MangaID)
					}
				}
			}
			h.mu.Unlock()

		case msg := <-h.Broadcast:
			h.mu.RLock()
			if clients := h.rooms[msg.MangaID]; clients != nil {
				for client := range clients {
					select {
					case client.Send <- msg:
					default:
						// If send buffer is full, unregister the client to prevent blocking
						go func(c *Client) { h.Unregister <- c }(client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) stop() {
	fmt.Println("📡 WebSocket Hub: Shutting down and notifying clients...")
	shutdownMsg := &models.ChatMessage{
		MangaID:    0,
		SenderName: "system",
		Content:    `{"action": "shutdown", "reason": "server_maintenance"}`,
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Notify all clients in all rooms
	for _, clients := range h.rooms {
		for client := range clients {
			select {
			case client.Send <- shutdownMsg:
			default:
				// Skip if buffer full
			}
		}
	}

	// Grace period
	time.Sleep(2 * time.Second)

	// Force close all connections
	for mangaID, clients := range h.rooms {
		for client := range clients {
			delete(clients, client)
			close(client.Send)
		}
		delete(h.rooms, mangaID)
	}
}

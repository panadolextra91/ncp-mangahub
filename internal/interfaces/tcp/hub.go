package tcp

import (
	"fmt"
	"net"
)

// Hub manages the set of active TCP clients and broadcasts messages to them.
type Hub struct {
	// Registered clients.
	clients map[net.Conn]bool

	// Inbound messages from the event bus.
	broadcast chan []byte

	// Register requests from the clients.
	register chan net.Conn

	// Unregister requests from clients.
	unregister chan net.Conn

	// Max number of allowed clients (DOS protection)
	maxClients int
}

func NewHub(maxClients int) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan net.Conn),
		unregister: make(chan net.Conn),
		clients:    make(map[net.Conn]bool),
		maxClients: maxClients,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			// DOS Protection Check
			if len(h.clients) >= h.maxClients {
				fmt.Fprintf(conn, "503 Service Unavailable - Server Full\n")
				conn.Close()
				continue
			}
			h.clients[conn] = true

		case conn := <-h.unregister:
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
			}

		case message := <-h.broadcast:
			for conn := range h.clients {
				// Non-blocking write attempt
				// We wrap in a goroutine or use a per-client buffer to avoid blocking the hub loop
				// if one client is slow.
				go func(c net.Conn, msg []byte) {
					_, err := c.Write(append(msg, '\n'))
					if err != nil {
						h.unregister <- c
					}
				}(conn, message)
			}
		}
	}
}

// Broadcast sends a message to all registered clients.
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

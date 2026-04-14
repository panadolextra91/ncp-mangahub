package tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/user/mangahub/pkg/auth"
)

type Server struct {
	addr      string
	hub       *Hub
	jwtSecret string
}

func NewServer(addr string, hub *Hub, secret string) *Server {
	return &Server{
		addr:      addr,
		hub:       hub,
		jwtSecret: secret,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("📡 TCP Real-time Sync Server listening on %s", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("TCP Accept error: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	// Rule: 5-second deadline for Handshake to prevent Goroutine Leak
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(conn, "401 Unauthorized - Handshake Timeout or Error\n")
		conn.Close()
		return
	}

	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "AUTH ") {
		fmt.Fprintf(conn, "401 Unauthorized - Handshake Required\n")
		conn.Close()
		return
	}

	token := strings.TrimPrefix(line, "AUTH ")
	_, err = auth.ValidateToken(token, s.jwtSecret)
	if err != nil {
		fmt.Fprintf(conn, "401 Unauthorized - Invalid Token\n")
		conn.Close()
		return
	}

	// Handshake Successful
	fmt.Fprintf(conn, "200 OK CONNECTED\n")
	
	// Reset deadline for persistent connection (or set a long one/heartbeat)
	conn.SetReadDeadline(time.Time{})

	// Register with Hub
	s.hub.register <- conn

	// Stay alive to detect disconnects
	for {
		_, err := reader.ReadByte()
		if err != nil {
			s.hub.unregister <- conn
			return
		}
	}
}

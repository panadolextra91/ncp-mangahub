package udp

import (
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/user/mangahub/pkg/auth"
)

type Server struct {
	addr      string
	registry  *Registry
	jwtSecret string
	conn      *net.UDPConn
}

func NewServer(addr string, registry *Registry, secret string) *Server {
	return &Server{
		addr:      addr,
		registry:  registry,
		jwtSecret: secret,
	}
}

func (s *Server) Start() error {
	udpAddr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	s.conn = conn

	log.Printf("📡 UDP Notification Server listening on %s", s.addr)

	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("UDP Read Error: %v", err)
			continue
		}

		data := string(buf[:n])
		go s.handlePacket(clientAddr, data)
	}
}

func (s *Server) handlePacket(addr *net.UDPAddr, data string) {
	parts := strings.Fields(data)
	if len(parts) < 2 {
		return
	}

	cmd := strings.ToUpper(parts[0])
	switch cmd {
	case "SUB":
		if len(parts) < 3 {
			return
		}
		mangaID, err := strconv.Atoi(parts[1])
		if err != nil {
			return
		}
		token := parts[2]
		_, err = auth.ValidateToken(token, s.jwtSecret)
		if err != nil {
			return // Silent drop
		}
		s.registry.Register(mangaID, addr)

	case "PING":
		token := parts[1]
		_, err := auth.ValidateToken(token, s.jwtSecret)
		if err != nil {
			return // Silent drop
		}
		s.registry.KeepAlive(addr)
	}
}

func (s *Server) Broadcast(mangaID int, payload []byte) {
	peers := s.registry.GetPeers(mangaID)
	for _, addr := range peers {
		_, err := s.conn.WriteToUDP(payload, addr)
		if err != nil {
			log.Printf("UDP Broadcast Error to %v: %v", addr, err)
		}
	}
}

func (s *Server) Stop() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

package tcp_test

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/internal/interfaces/tcp"
	"github.com/user/mangahub/pkg/auth"
)

func TestTCPProtocol(t *testing.T) {
	secret := "test-secret"
	hub := tcp.NewHub(2) // Max 2 clients for testing DOS
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go hub.Run(ctx, &wg)

	// Get a free port
	listener, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)
	addr := listener.Addr().String()
	listener.Close() 

	srv := tcp.NewServer(addr, hub, secret)
	go func() {
		_ = srv.Start()
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	t.Run("Handshake Success", func(t *testing.T) {
		conn, err := net.Dial("tcp", addr)
		assert.NoError(t, err)
		defer conn.Close()

		token, _ := auth.GenerateToken(1, "tester", "user", secret)
		fmt.Fprintf(conn, "AUTH %s\n", token)

		reader := bufio.NewReader(conn)
		line, _ := reader.ReadString('\n')
		assert.Equal(t, "200 OK CONNECTED\n", line)
	})

	t.Run("Handshake Fail - Invalid Token", func(t *testing.T) {
		conn, err := net.Dial("tcp", addr)
		assert.NoError(t, err)
		defer conn.Close()

		fmt.Fprintf(conn, "AUTH invalid-token\n")

		reader := bufio.NewReader(conn)
		line, _ := reader.ReadString('\n')
		assert.Contains(t, line, "401 Unauthorized")
	})

	t.Run("Handshake Fail - Timeout", func(t *testing.T) {
		conn, err := net.Dial("tcp", addr)
		assert.NoError(t, err)
		defer conn.Close()

		// Do nothing, wait for 5s timeout
		time.Sleep(5200 * time.Millisecond)
		reader := bufio.NewReader(conn)
		line, _ := reader.ReadString('\n')
		assert.Contains(t, line, "401 Unauthorized")
	})

	t.Run("DOS Protection - Max Clients", func(t *testing.T) {
		token, _ := auth.GenerateToken(1, "tester", "user", secret)
		
		c1, _ := net.Dial("tcp", addr)
		fmt.Fprintf(c1, "AUTH %s\n", token)
		time.Sleep(50 * time.Millisecond) // Give time for registration

		c2, _ := net.Dial("tcp", addr)
		fmt.Fprintf(c2, "AUTH %s\n", token)
		time.Sleep(50 * time.Millisecond) // Give time for registration

		// 3rd client should be rejected
		c3, err := net.Dial("tcp", addr)
		assert.NoError(t, err)
		fmt.Fprintf(c3, "AUTH %s\n", token)
		
		reader := bufio.NewReader(c3)
		// First line after AUTH might be Handshake Success OR rejection
		// Wait... the server sends 200 OK CONNECTED before sending the conn to the register channel.
		// If the hub is full, it should send 503 from the Run() loop AFTER receiving the conn.
		// So the client will see 200 OK CONNECTED, then 503.
		line1, _ := reader.ReadString('\n')
		assert.Equal(t, "200 OK CONNECTED\n", line1)
		
		line2, _ := reader.ReadString('\n')
		assert.Contains(t, line2, "503 Service Unavailable")
		
		c1.Close()
		c2.Close()
		c3.Close()
	})

	t.Run("Broadcasting", func(t *testing.T) {
		conn, _ := net.Dial("tcp", addr)
		token, _ := auth.GenerateToken(1, "tester", "user", secret)
		fmt.Fprintf(conn, "AUTH %s\n", token)
		
		reader := bufio.NewReader(conn)
		reader.ReadString('\n') // Consume "200 OK CONNECTED\n"
		
		time.Sleep(50 * time.Millisecond) // Give time for registration

		msg := []byte(`{"event":"test"}`)
		hub.Broadcast(msg)

		line, _ := reader.ReadString('\n')
		assert.Equal(t, string(msg)+"\n", line)
		conn.Close()
	})
}

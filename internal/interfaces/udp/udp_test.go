package udp_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/internal/interfaces/udp"
	"github.com/user/mangahub/pkg/auth"
)

func TestUDPProtocolIntegration(t *testing.T) {
	secret := "test-secret"
	// Tiny TTL and GC for fast testing
	ttl := 200 * time.Millisecond
	gc := 100 * time.Millisecond
	registry := udp.NewRegistry(ttl, gc)
	
	// Port 0 for random available port (placeholder remove)
	
	// Start server in background
	// We need to capture the bound address
	// I'll add a helper to server or just use a known port if necessary.
	// Actually, let's just listen on a fixed port for the test.
	port := 9192
	serverFixed := udp.NewServer(fmt.Sprintf("127.0.0.1:%d", port), registry, secret)
	go serverFixed.Start()
	time.Sleep(100 * time.Millisecond) // Wait for bind

	token, _ := auth.GenerateToken(1, "tester", "user", secret)
	serverAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", port))
	
	t.Run("SUB and Broadcast", func(t *testing.T) {
		clientConn, _ := net.ListenUDP("udp", nil)
		defer clientConn.Close()

		// 1. Send SUB
		subMsg := fmt.Sprintf("SUB 1 %s", token)
		clientConn.WriteToUDP([]byte(subMsg), serverAddr)
		time.Sleep(100 * time.Millisecond) // Wait for server to process

		// 2. Server Broadcast
		payload := []byte("new-manga-alert")
		serverFixed.Broadcast(1, payload)

		// 3. Receive on Client
		buf := make([]byte, 1024)
		clientConn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, _, err := clientConn.ReadFromUDP(buf)
		assert.NoError(t, err)
		assert.Equal(t, payload, buf[:n])
	})

	t.Run("PING and TTL Pruning", func(t *testing.T) {
		clientConn, _ := net.ListenUDP("udp", nil)
		defer clientConn.Close()

		// 1. SUB
		subMsg := fmt.Sprintf("SUB 2 %s", token)
		clientConn.WriteToUDP([]byte(subMsg), serverAddr)
		time.Sleep(100 * time.Millisecond)

		// Verify registered
		assert.Len(t, registry.GetPeers(2), 1)

		// 2. PING
		pingMsg := fmt.Sprintf("PING %s", token)
		clientConn.WriteToUDP([]byte(pingMsg), serverAddr)
		time.Sleep(100 * time.Millisecond)

		// 3. Wait for TTL to expire
		time.Sleep(500 * time.Millisecond)
		
		// 4. Verify Pruned
		assert.Len(t, registry.GetPeers(2), 0)
	})
}

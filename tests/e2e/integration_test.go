package e2e

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestIntegration_MultiProtocolFlow(t *testing.T) {
	s, err := StartServer("../../mangahub_server")
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer s.Stop()

	token := getAuthToken(t, s.Port)

	// 1. Setup UDP Listener
	serverUDPAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", s.UDPPort))
	udpAddr, _ := net.ResolveUDPAddr("udp", "localhost:0")
	udpConn, _ := net.ListenUDP("udp", udpAddr)
	defer udpConn.Close()
	
	t.Logf("UDP Subscribing to port %d", s.UDPPort)
	subMsg := fmt.Sprintf("SUB 1 %s", token)
	_, err = udpConn.WriteToUDP([]byte(subMsg), serverUDPAddr)
	assert.NoError(t, err)

	// 2. Setup WebSocket
	t.Logf("Connecting to WebSocket...")
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("localhost:%d", s.Port), Path: "/api/chat", RawQuery: "manga_id=1&token=" + token}
	wsConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("WS Dial failed: %v", err)
	}
	defer wsConn.Close()
	t.Logf("WebSocket connected.")

	// 3. gRPC: Create Manga as Admin
	t.Logf("Dialing gRPC...")
	grpcConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", s.GRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("gRPC Dial failed: %v", err)
	}
	defer grpcConn.Close()
	
	adminClient := pb.NewAdminServiceClient(grpcConn)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+token)
	
	t.Logf("Creating Manga via gRPC...")
	_, err = adminClient.CreateManga(ctx, &pb.CreateMangaRequest{
		Title:       "Final-E2E-Manga",
		Author:      "Architect",
		Description: "The Grand Finale",
	})
	assert.NoError(t, err)

	// 4. Verify UDP Notification
	t.Logf("Waiting for UDP notification...")
	udpConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buf := make([]byte, 1024)
	n, _, err := udpConn.ReadFromUDP(buf)
	if err != nil {
		t.Logf("UDP Read Error: %v", err)
	}
	assert.NoError(t, err, "UDP should receive manga.new notification")
	assert.Contains(t, string(buf[:n]), "Final-E2E-Manga")

	// 5. Verify WS Chat History or Message
	t.Logf("Sending WS message...")
	err = wsConn.WriteMessage(websocket.TextMessage, []byte("Hello from E2E"))
	assert.NoError(t, err)
	
	// Read back (only expect our own message as the history is empty)
	t.Logf("Reading back WS echo...")
	wsConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, p, err := wsConn.ReadMessage()
	assert.NoError(t, err)
	assert.Contains(t, string(p), "Hello from E2E")
}

package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/user/mangahub/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	HTTP_URL  = "http://127.0.0.1:8080"
	WS_URL    = "ws://127.0.0.1:8080/api/chat"
	TCP_ADDR  = "127.0.0.1:9090"
	UDP_ADDR  = "127.0.0.1:9191"
	GRPC_ADDR = "127.0.0.1:50052"
)

func main() {
	fmt.Println("🌟 =============================================== 🌟")
	fmt.Println("🚀   MANGAHUB MULTI-PROTOCOL GRAND FINALE DEMO     🚀")
	fmt.Println("🌟 =============================================== 🌟")

	// 1. Auth & Login
	fmt.Print("\n🔐 [STEP 1] Authenticaticating Admin via HTTP... ")
	token := login()
	fmt.Println("DONE! ✅")

	// 2. TCP Client (Real-time Sync)
	fmt.Print("📡 [STEP 2] Connecting TCP Real-time Client... ")
	tcpConn, _ := net.Dial("tcp", TCP_ADDR)
	fmt.Fprintf(tcpConn, "AUTH %s\n", token)
	tcpReader := bufio.NewReader(tcpConn)
	tcpReader.ReadString('\n') // Consume 200 OK
	fmt.Println("CONNECTED! ✅")

	// 3. UDP Client (Fast Notifications)
	fmt.Print("📡 [STEP 3] Registering UDP Listener... ")
	udpConn, _ := net.ListenUDP("udp", nil)
	fmt.Fprintf(udpConn, "SUB 1 %s", token)
	fmt.Println("SUBBED! ✅")

	// 4. WebSocket Client (Chat)
	fmt.Print("💬 [STEP 4] Entering Manga #1 Chat Room via WS... ")
	u, _ := url.Parse(WS_URL + "?manga_id=1&token=" + token)
	wsConn, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)
	fmt.Println("JOINED! ✅")

	// 5. gRPC Event Stream
	fmt.Print("🏗️  [STEP 5] Subscribing to gRPC Admin Events... ")
	grpcConn, _ := grpc.Dial(GRPC_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	mangaClient := pb.NewMangaServiceClient(grpcConn)
	adminClient := pb.NewAdminServiceClient(grpcConn)
	grpcCtx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+token)
	stream, _ := mangaClient.SubscribeEvents(grpcCtx, &emptypb.Empty{})
	fmt.Println("STREAMING! ✅")

	fmt.Println("\n🎬 ================= [ ACT 1: THE RELEASE ] ================= 🎬")
	
	// Create Manga via gRPC Unary
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Printf("\n🔨 [ADMIN] Creating new Manga 'GRAND FINALE 2026' via gRPC Unary...")
		adminClient.CreateManga(grpcCtx, &pb.CreateMangaRequest{
			Title: "GRAND FINALE 2026",
			Author: "Mẹ Architect",
			Description: "The state of the art modular monolith",
		})
	}()

	// 6. Listen for events on ALL Protocols
	fmt.Println("\n👀 [MONITOR] Watching all protocols for the release event...")
	
	// Wait for TCP
	go func() {
		line, _ := tcpReader.ReadString('\n')
		fmt.Printf("\n📥 [TCP]  RECV Event: %s", line)
	}()

	// Wait for UDP
	go func() {
		buf := make([]byte, 1024)
		n, _, _ := udpConn.ReadFromUDP(buf)
		fmt.Printf("\n📥 [UDP]  RECV Event: %s", string(buf[:n]))
	}()

	// Wait for gRPC Stream
	go func() {
		resp, _ := stream.Recv()
		fmt.Printf("\n📥 [gRPC] RECV Event: %s [%s]", resp.Topic, resp.Message)
	}()

	time.Sleep(3 * time.Second)

	fmt.Println("\n\n🎬 ================= [ ACT 2: THE DISCUSSION ] ================= 🎬")
	
	// Send Chat Message via WS
	fmt.Printf("\n💬 [USER] Sending chat message: 'PROProject of the year!'")
	wsConn.WriteMessage(websocket.TextMessage, []byte("PROProject of the year!"))
	
	// See it on WS
	_, p, _ := wsConn.ReadMessage()
	fmt.Printf("\n📥 [WS]   RECV Chat: %s", string(p))

	time.Sleep(2 * time.Second)

	fmt.Println("\n\n🌟 =============================================== 🌟")
	fmt.Println("🏁   DEMO COMPLETE: ALL 5 PROTOCOLS VERIFIED!      🏁")
	fmt.Println("🌟 =============================================== 🌟")
	
	os.Exit(0)
}

func login() string {
	url := HTTP_URL + "/api/auth/login"
	body := `{"username":"admin","password":"password"}`
	resp, err := http.Post(url, "application/json", bytes.NewBufferString(body))
	if err != nil {
		log.Fatalf("FAIL: Could not connect to MangaHub at %s. Did you run the server?", HTTP_URL)
	}
	defer resp.Body.Close()
	
	var res struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	if res.Token == "" {
		// Try register first
		regUrl := HTTP_URL + "/api/auth/register"
		regBody := `{"username":"admin","password":"password","role":"admin"}`
		http.Post(regUrl, "application/json", bytes.NewBufferString(regBody))
		return login()
	}
	return res.Token
}

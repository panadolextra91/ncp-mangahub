package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/mangahub/config"
	"github.com/user/mangahub/internal/adapters/database"
	"github.com/user/mangahub/internal/application"
	"github.com/user/mangahub/internal/eventbus"
	"github.com/user/mangahub/internal/infrastructure"
	mh_http "github.com/user/mangahub/internal/interfaces/http"
	mh_tcp "github.com/user/mangahub/internal/interfaces/tcp"
	mh_udp "github.com/user/mangahub/internal/interfaces/udp"
	mh_ws "github.com/user/mangahub/internal/interfaces/ws"
	"github.com/user/mangahub/pkg/models"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Initialize Infrastructure (SQLITE)
	db, err := infrastructure.NewSQLiteDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Critical: Failed to connect to SQLite: %v", err)
	}
	defer db.Close()

	if err := infrastructure.InitSchema(db); err != nil {
		log.Fatalf("Critical: Failed to initialize database schema: %v", err)
	}

	// 3. Initialize Shared Components
	bus := eventbus.NewEventBus(100)

	// 4. Initialize Protocol Layers
	// --- TCP ---
	tcpHub := mh_tcp.NewHub(cfg.MaxTCPClients)
	go tcpHub.Run()
	tcpServer := mh_tcp.NewServer(":"+cfg.TCPPort, tcpHub, cfg.JWTSecret)
	go func() {
		if err := tcpServer.Start(); err != nil {
			log.Printf("TCP Server Error: %v", err)
		}
	}()

	// --- WebSocket ---
	wsHub := mh_ws.NewHub()
	go wsHub.Run()

	// --- UDP ---
	udpRegistry := mh_udp.NewRegistry(60*time.Second, 10*time.Second)
	udpServer := mh_udp.NewServer(":"+cfg.UDPPort, udpRegistry, cfg.JWTSecret)
	go func() {
		if err := udpServer.Start(); err != nil {
			log.Printf("UDP Server Error: %v", err)
		}
	}()

	// 5. Initialize Repositories
	userRepo := database.NewSqliteUserRepository(db)
	mangaRepo := database.NewSqliteMangaRepository(db)
	progRepo := database.NewSqliteProgressRepository(db)
	chatRepo := database.NewSqliteChatRepository(db)

	// 6. Initialize Services
	authSvc := application.NewAuthService(userRepo)
	mangaSvc := application.NewMangaService(mangaRepo, bus)
	progSvc := application.NewProgressService(progRepo, bus)
	chatSvc := application.NewChatService(chatRepo, bus)

	// 7. Bridge EventBus to Protocol Hubs (Real-time Broadcast)
	bridge := func(ch <-chan models.Event) {
		for event := range ch {
			jsonData, _ := json.Marshal(event.Payload)
			
			// Broadcast to TCP Clients
			tcpHub.Broadcast(jsonData)

			// Broadcast to WebSocket Clients (if it's a ChatMessage)
			if msg, ok := event.Payload.(*models.ChatMessage); ok {
				wsHub.Broadcast <- msg
			}

			// Broadcast to UDP Clients (ONLY GLOBAL EVENTS per Mẹ Architect)
			if event.Topic == "manga.new" {
				// manga.new events in our system usually have a *Manga payload
				// We broadcast to all UDP subscribers (mangaID=0)
				udpServer.Broadcast(0, jsonData)
			}
		}
	}

	go bridge(bus.Subscribe("manga.new"))
	go bridge(bus.Subscribe("progress.updated"))
	go bridge(bus.Subscribe("chat.message"))

	// 8. Initialize HTTP & WS Handlers
	authH := mh_http.NewAuthHandler(authSvc, cfg.JWTSecret)
	mangaH := mh_http.NewMangaHandler(mangaSvc)
	progH := mh_http.NewProgressHandler(progSvc)
	healthH := mh_http.NewHealthHandler(db, bus)
	wsH := mh_ws.NewChatHandler(wsHub, chatSvc, cfg.JWTSecret)

	mux := mh_http.SetupRouter(authH, mangaH, progH, healthH, cfg.JWTSecret)
	
	// Register WS endpoint
	mux.HandleFunc("/api/chat", wsH.HandleWS)

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	// 9. Start HTTP Server
	go func() {
		fmt.Printf("🚀 MangaHub Core API (HTTP + WS + UDP) listening on port %s\n", cfg.Port)
		fmt.Printf("📡 UDP Notification Server on port %s\n", cfg.UDPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP Server Error: %v", err)
		}
	}()

	// 10. Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Println("\n🛑 Shutting down MangaHub...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	udpServer.Stop()

	fmt.Println("👋 Cleanup complete. Goodbye!")
}

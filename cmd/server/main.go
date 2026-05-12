package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/user/mangahub/config"
	"github.com/user/mangahub/internal/adapters/database"
	"github.com/user/mangahub/internal/application"
	"github.com/user/mangahub/internal/eventbus"
	"github.com/user/mangahub/internal/infrastructure"
	mh_grpc "github.com/user/mangahub/internal/interfaces/grpc"
	mh_http "github.com/user/mangahub/internal/interfaces/http"
	mh_tcp "github.com/user/mangahub/internal/interfaces/tcp"
	mh_udp "github.com/user/mangahub/internal/interfaces/udp"
	mh_ws "github.com/user/mangahub/internal/interfaces/ws"
	"github.com/user/mangahub/pkg/models"
	"github.com/user/mangahub/pkg/pb"
	"google.golang.org/grpc"
)

func main() {
	// 1. Initialize Root Context and Signal Handling (LOUDSPEAKER)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. Initialize Roll Call (WAIT GROUP)
	var wg sync.WaitGroup

	// 3. Load Config
	cfg := config.LoadConfig()

	// 4. Initialize Infrastructure (SQLITE)
	db, err := infrastructure.NewSQLiteDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Critical: Failed to connect to SQLite: %v", err)
	}
	defer db.Close() // Safety fallback

	if err := infrastructure.InitSchema(db); err != nil {
		log.Fatalf("Critical: Failed to initialize database schema: %v", err)
	}

	// 5. Initialize Shared Components
	bus := eventbus.NewEventBus(100)

	// 6. Initialize Repositories
	userRepo := database.NewSqliteUserRepository(db)
	mangaRepo := database.NewSqliteMangaRepository(db)
	progRepo := database.NewSqliteProgressRepository(db)
	chatRepo := database.NewSqliteChatRepository(db)

	// 6.5 Auto-seeding
	mangas, _ := mangaRepo.List()
	if len(mangas) == 0 {
		log.Println("🌱 Database is empty. Seeding initial manga data...")
		seedFile, err := os.Open("data/manga_seed.json")
		if err == nil {
			var seedData []*models.Manga
			if err := json.NewDecoder(seedFile).Decode(&seedData); err == nil {
				for _, m := range seedData {
					_ = mangaRepo.Save(m)
				}
				log.Printf("✅ Successfully seeded %d manga records.", len(seedData))
			} else {
				log.Printf("⚠️ Failed to decode seed file: %v", err)
			}
			seedFile.Close()
		} else {
			log.Printf("⚠️ No seed file found at data/manga_seed.json: %v", err)
		}
	}

	// 7. Initialize Services
	authSvc := application.NewAuthService(userRepo)
	mangaSvc := application.NewMangaService(mangaRepo, bus)
	progSvc := application.NewProgressService(progRepo, bus)
	chatSvc := application.NewChatService(chatRepo, bus)

	// 8. Initialize Protocol Layers

	// --- TCP ---
	tcpHub := mh_tcp.NewHub(cfg.MaxTCPClients)
	wg.Add(1)
	go tcpHub.Run(ctx, &wg)
	tcpServer := mh_tcp.NewServer(":"+cfg.TCPPort, tcpHub, cfg.JWTSecret)
	wg.Add(1)
	go func() {
		defer wg.Done()
		go tcpServer.Start() // Start listener
		<-ctx.Done()
		log.Println("🛑 TCP Server: Stopping listener...")
		// Listener close is handled by Start loop usually, but here we just wait for hub to close conns
	}()

	// --- WebSocket ---
	wsHub := mh_ws.NewHub()
	wg.Add(1)
	go wsHub.Run(ctx, &wg)

	// --- UDP ---
	udpRegistry := mh_udp.NewRegistry(60*time.Second, 10*time.Second)
	udpServer := mh_udp.NewServer(":"+cfg.UDPPort, udpRegistry, cfg.JWTSecret)
	wg.Add(1)
	go func() {
		defer wg.Done()
		go udpServer.Start()
		<-ctx.Done()
		log.Println("🛑 UDP Server: Stopping server...")
		udpServer.Stop()
	}()

	// --- gRPC ---
	authInterceptor := mh_grpc.NewAuthInterceptor(cfg.JWTSecret)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)
	pb.RegisterMangaServiceServer(grpcServer, mh_grpc.NewMangaService(mangaSvc, progSvc, bus))
	pb.RegisterAdminServiceServer(grpcServer, mh_grpc.NewAdminService(mangaSvc))

	wg.Add(1)
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Printf("gRPC Listen Error: %v", err)
			return
		}
		log.Printf("📡 gRPC Admin Server listening on :%s", cfg.GRPCPort)
		go grpcServer.Serve(lis)
		<-ctx.Done()
		log.Println("🛑 gRPC Server: GracefulStop...")
		grpcServer.GracefulStop()
	}()

	// 9. Bridge EventBus to Protocol Hubs (Real-time Broadcast)
	bridge := func(topic string) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch := bus.Subscribe(topic)
			defer bus.Unsubscribe(topic, ch)
			for {
				select {
				case <-ctx.Done():
					log.Printf("🛑 Bridge [%s]: Stopping...", topic)
					return
				case event := <-ch:
					jsonData, _ := json.Marshal(event.Payload)
					tcpHub.Broadcast(jsonData)
					if msg, ok := event.Payload.(*models.ChatMessage); ok {
						wsHub.Broadcast <- msg
					}
					if event.Topic == "manga.new" {
						if m, ok := event.Payload.(*models.Manga); ok {
							udpServer.Broadcast(m.ID, jsonData)
						}
					}
				}
			}
		}()
	}

	bridge("manga.new")
	bridge("progress.updated")
	bridge("chat.message")

	// 10. HTTP & WS Setups
	authH := mh_http.NewAuthHandler(authSvc, cfg.JWTSecret)
	mangaH := mh_http.NewMangaHandler(mangaSvc)
	progH := mh_http.NewProgressHandler(progSvc)
	healthH := mh_http.NewHealthHandler(db, bus)
	wsH := mh_ws.NewChatHandler(wsHub, chatSvc, cfg.JWTSecret)
	mux := mh_http.SetupRouter(authH, mangaH, progH, healthH, cfg.JWTSecret)
	mux.HandleFunc("GET /api/chat", wsH.HandleWS)

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("🚀 MangaHub Core API listening on %s", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP Server Error: %v", err)
		}
	}()

	// ---------------------------------------------------------
	// 11. STRATEGIC SHUTDOWN SEQUENCE (ENTREPRISE STANDARD)
	// ---------------------------------------------------------
	<-ctx.Done() // Block until SIGINT/SIGTERM
	log.Println("\n🛑 SHUTDOWN SIGNAL RECEIVED. Initiating 5-step exit protocol...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// STEP 1: Reject New Work (Stop HTTP Listener)
	log.Println("Step 1: Closing HTTP Listener (No more new requests)...")
	httpServer.Shutdown(shutdownCtx)

	// STEP 2 & 3: Evict Active Clients & Stop internal chains
	// (This is triggered by the context propagation to Hubs and Bridges)
	log.Println("Step 2 & 3: Notifying active clients and stopping internal event chains...")

	// STEP 4: Roll Call (Wait for everything to finish)
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Step 4: All background workers reported clean completion.")
	case <-shutdownCtx.Done():
		log.Println("Step 4: Graceful shutdown timed out. Some connections may be forced closed.")
	}

	// STEP 5: Final Cleanup (DB Close)
	log.Println("Step 5: Closing Database connection...")
	db.Close()

	log.Println("👋 MangaHub successfully decommissioned. Goodbye!")
}

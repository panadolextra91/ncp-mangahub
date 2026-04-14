package main

import (
	"context"
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
	bus := eventbus.NewEventBus(100) // 100 buffer size for demo

	// 4. Initialize Repositories (Adapters)
	userRepo := database.NewSqliteUserRepository(db)
	mangaRepo := database.NewSqliteMangaRepository(db)
	progRepo := database.NewSqliteProgressRepository(db)

	// 5. Initialize Services (Application)
	authSvc := application.NewAuthService(userRepo)
	mangaSvc := application.NewMangaService(mangaRepo, bus)
	progSvc := application.NewProgressService(progRepo, bus)

	// 6. Initialize HTTP Transport (Interfaces)
	authH := mh_http.NewAuthHandler(authSvc, cfg.JWTSecret)
	mangaH := mh_http.NewMangaHandler(mangaSvc)
	progH := mh_http.NewProgressHandler(progSvc)
	healthH := mh_http.NewHealthHandler(db, bus)

	mux := mh_http.SetupRouter(authH, mangaH, progH, healthH, cfg.JWTSecret)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	// 7. Start Server (Goroutine as per protocol plan)
	go func() {
		fmt.Printf("🚀 MangaHub Core API listening on port %s\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP Server Error: %v", err)
		}
	}()

	// 8. Graceful Shutdown (Partial implementation for Phase 4)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Println("\n🛑 Shutting down MangaHub...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("👋 Cleanup complete. Goodbye!")
}

package grpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/internal/adapters/database"
	"github.com/user/mangahub/internal/application"
	"github.com/user/mangahub/internal/eventbus"
	"github.com/user/mangahub/internal/infrastructure"
	mh_grpc "github.com/user/mangahub/internal/interfaces/grpc"
	"github.com/user/mangahub/pkg/auth"
	"github.com/user/mangahub/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGRPCIntegration(t *testing.T) {
	secret := "test-secret"
	db, _ := infrastructure.NewSQLiteDB(":memory:")
	infrastructure.InitSchema(db)
	bus := eventbus.NewEventBus(10)
	
	mangaRepo := database.NewSqliteMangaRepository(db)
	mangaSvc := application.NewMangaService(mangaRepo, bus)
	progRepo := database.NewSqliteProgressRepository(db)
	progSvc := application.NewProgressService(progRepo, bus)

	// Interceptor
	interceptor := mh_grpc.NewAuthInterceptor(secret)
	
	// Server
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)
	pb.RegisterMangaServiceServer(srv, mh_grpc.NewMangaService(mangaSvc, progSvc, bus))
	pb.RegisterAdminServiceServer(srv, mh_grpc.NewAdminService(mangaSvc))

	lis, _ := net.Listen("tcp", "127.0.0.1:0")
	go srv.Serve(lis)
	defer srv.Stop()

	// Client
	conn, _ := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	mangaClient := pb.NewMangaServiceClient(conn)
	adminClient := pb.NewAdminServiceClient(conn)

	// Tokens
	adminToken, _ := auth.GenerateToken(1, "admin", secret)
	userToken, _ := auth.GenerateToken(2, "user", secret)

	t.Run("Unary: Create Manga (Admin Auth)", func(t *testing.T) {
		ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+adminToken)
		resp, err := adminClient.CreateManga(ctx, &pb.CreateMangaRequest{
			Title: "G-Manga",
			Author: "G-Author",
		})
		assert.NoError(t, err)
		assert.Equal(t, "G-Manga", resp.Title)
	})

	t.Run("Unary: Create Manga (Unauthorized)", func(t *testing.T) {
		ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+userToken)
		_, err := adminClient.CreateManga(ctx, &pb.CreateMangaRequest{Title: "Bad"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "admin role required")
	})

	t.Run("Stream: Subscribe Events", func(t *testing.T) {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		
		ctx := metadata.AppendToOutgoingContext(timeoutCtx, "authorization", "Bearer "+userToken)
		stream, err := mangaClient.SubscribeEvents(ctx, &emptypb.Empty{})
		assert.NoError(t, err)

		// Small delay to ensure server goroutine has subscribed to EventBus
		time.Sleep(200 * time.Millisecond)

		// Trigger event via Unary Admin call
		adminCtx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+adminToken)
		_, err = adminClient.CreateManga(adminCtx, &pb.CreateMangaRequest{Title: "Stream-Manga"})
		assert.NoError(t, err)

		// Receive on stream
		resp, err := stream.Recv()
		assert.NoError(t, err)
		assert.Equal(t, "manga.new", resp.Topic)
		assert.Contains(t, resp.Message, "Stream-Manga")
	})
}

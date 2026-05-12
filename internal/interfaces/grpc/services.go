package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/user/mangahub/internal/application"
	"github.com/user/mangahub/internal/eventbus"
	"github.com/user/mangahub/pkg/auth"
	"github.com/user/mangahub/pkg/models"
	"github.com/user/mangahub/pkg/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MangaService struct {
	pb.UnimplementedMangaServiceServer
	mangaSvc application.MangaService
	progSvc  application.ProgressService
	bus      *eventbus.EventBus
}

func NewMangaService(mSvc application.MangaService, pSvc application.ProgressService, bus *eventbus.EventBus) *MangaService {
	return &MangaService{
		mangaSvc: mSvc,
		progSvc:  pSvc,
		bus:      bus,
	}
}

func (s *MangaService) GetManga(ctx context.Context, req *pb.GetMangaRequest) (*pb.Manga, error) {
	manga, err := s.mangaSvc.GetManga(int(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.Manga{
		Id:     int32(manga.ID),
		Title:  manga.Title,
		Author: manga.Author,
	}, nil
}

func (s *MangaService) UpdateProgress(ctx context.Context, req *pb.UpdateProgressRequest) (*pb.UpdateProgressResponse, error) {
	claims, ok := ctx.Value("claims").(*auth.Claims)
	if !ok {
		return &pb.UpdateProgressResponse{Success: false, Message: "unauthorized"}, nil
	}

	progress := &models.UserProgress{
		UserID:         claims.UserID,
		MangaID:        int(req.MangaId),
		CurrentChapter: int(req.Chapter),
	}

	err := s.progSvc.UpdateProgress(progress)
	if err != nil {
		return &pb.UpdateProgressResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.UpdateProgressResponse{Success: true, Message: "Progress updated"}, nil
}

func (s *MangaService) SubscribeEvents(empty *emptypb.Empty, stream pb.MangaService_SubscribeEventsServer) error {
	chNew := s.bus.Subscribe("manga.new")
	chProg := s.bus.Subscribe("progress.updated")
	
	defer s.bus.Unsubscribe("manga.new", chNew)
	defer s.bus.Unsubscribe("progress.updated", chProg)

	for {
		select {
		case <-stream.Context().Done():
			log.Printf("📡 gRPC Stream: Client disconnected")
			return nil
		case event := <-chNew:
			s.sendEvent(stream, "manga.new", event.Payload)
		case event := <-chProg:
			s.sendEvent(stream, "progress.updated", event.Payload)
		}
	}
}

func (s *MangaService) sendEvent(stream pb.MangaService_SubscribeEventsServer, topic string, payload interface{}) {
	msg := ""
	switch topic {
	case "manga.new":
		if m, ok := payload.(*models.Manga); ok {
			msg = fmt.Sprintf("New release: %s by %s", m.Title, m.Author)
		}
	case "progress.updated":
		if p, ok := payload.(*models.UserProgress); ok {
			msg = fmt.Sprintf("User #%d updated progress on Manga #%d to Chapter %d", p.UserID, p.MangaID, p.CurrentChapter)
		}
	default:
		msg = fmt.Sprintf("%v", payload)
	}

	notification := &pb.EventNotification{
		Topic:     topic,
		Message:   msg,
		Timestamp: timestamppb.Now(),
	}
	if err := stream.Send(notification); err != nil {
		log.Printf("gRPC Stream Send error: %v", err)
	}
}

// --- Admin Service ---

type AdminService struct {
	pb.UnimplementedAdminServiceServer
	mangaSvc application.MangaService
}

func NewAdminService(mSvc application.MangaService) *AdminService {
	return &AdminService{mangaSvc: mSvc}
}

func (s *AdminService) CreateManga(ctx context.Context, req *pb.CreateMangaRequest) (*pb.Manga, error) {
	claims, ok := ctx.Value("claims").(*auth.Claims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	manga := &models.Manga{
		Title:  req.Title,
		Author: req.Author,
	}
	err := s.mangaSvc.CreateManga(claims.Role, manga)
	if err != nil {
		return nil, err
	}
	return &pb.Manga{
		Id:     int32(manga.ID),
		Title:  manga.Title,
		Author: manga.Author,
	}, nil
}

func (s *AdminService) DeleteManga(ctx context.Context, req *pb.DeleteMangaRequest) (*pb.DeleteMangaResponse, error) {
	return &pb.DeleteMangaResponse{Success: true}, nil
}

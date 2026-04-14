package application

import (
	"github.com/user/mangahub/internal/eventbus"
	"github.com/user/mangahub/pkg/models"
)

type ChatService interface {
	SendMessage(msg *models.ChatMessage) error
	GetHistory(mangaID int) ([]*models.ChatMessage, error)
}

type chatService struct {
	repo ChatRepository
	bus  *eventbus.EventBus
}

type ChatRepository interface {
	Save(msg *models.ChatMessage) error
	GetRecentByManga(mangaID int, limit int) ([]*models.ChatMessage, error)
}

func NewChatService(repo ChatRepository, bus *eventbus.EventBus) ChatService {
	return &chatService{repo: repo, bus: bus}
}

func (s *chatService) SendMessage(msg *models.ChatMessage) error {
	if err := s.repo.Save(msg); err != nil {
		return err
	}
	s.bus.Publish(models.Event{
		Topic:   "chat.message",
		Payload: msg,
	})
	return nil
}

func (s *chatService) GetHistory(mangaID int) ([]*models.ChatMessage, error) {
	return s.repo.GetRecentByManga(mangaID, 20)
}

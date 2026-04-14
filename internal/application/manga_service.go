package application

import (
	"errors"

	"github.com/user/mangahub/internal/domain"
	"github.com/user/mangahub/internal/eventbus"
	"github.com/user/mangahub/pkg/models"
)

var (
	ErrUnauthorizedCreate = errors.New("unauthorized: strictly admin role is required to construct new global manga entity")
)

// MangaService bounds standard interactions shielding them from external APIs.
type MangaService interface {
	CreateManga(role string, manga *models.Manga) error
}

type mangaService struct {
	repo domain.MangaRepository
	bus  *eventbus.EventBus
}

// NewMangaService abstracts DB manipulations coupling explicitly with out-bound Event Bus parameters.
func NewMangaService(repo domain.MangaRepository, bus *eventbus.EventBus) MangaService {
	return &mangaService{repo: repo, bus: bus}
}

func (s *mangaService) CreateManga(role string, manga *models.Manga) error {
	// Logic interception asserting global privileges inherently ignoring HTTP parameters 
	if role != "admin" {
		return ErrUnauthorizedCreate
	}

	if err := s.repo.Save(manga); err != nil {
		return err
	}

	// Triggers Event via decoupling Pub/Sub standard protocols
	s.bus.Publish(models.Event{
		Topic:   "manga.new",
		Payload: manga,
	})

	return nil
}

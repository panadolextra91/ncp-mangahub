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

// SearchFilters is a type alias re-exporting domain.SearchFilters so HTTP/gRPC
// handlers can import a single "application" package without referencing domain.
// Defined in domain to keep MangaRepository interface free of cyclic imports.
type SearchFilters = domain.SearchFilters

// MangaService bounds standard interactions shielding them from external APIs.
type MangaService interface {
	CreateManga(role string, manga *models.Manga) error
	GetManga(id int) (*models.Manga, error)
	ListMangas() ([]*models.Manga, error)
	SearchMangas(query string) ([]*models.Manga, error)
	SearchMangasWithFilters(f SearchFilters) ([]*models.Manga, error)
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

func (s *mangaService) GetManga(id int) (*models.Manga, error) {
	return s.repo.FindByID(id)
}

func (s *mangaService) ListMangas() ([]*models.Manga, error) {
	return s.repo.List()
}

func (s *mangaService) SearchMangas(query string) ([]*models.Manga, error) {
	return s.repo.Search(query)
}

// SearchMangasWithFilters routes through the repo's SearchWithFilters path so
// HTTP/gRPC handlers can express multi-criteria queries. The legacy
// SearchMangas(query string) path stays untouched to preserve bit-for-bit
// backwards compatibility for q-only callers.
func (s *mangaService) SearchMangasWithFilters(f SearchFilters) ([]*models.Manga, error) {
	return s.repo.SearchWithFilters(f)
}

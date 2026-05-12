package application

import (
	"github.com/user/mangahub/internal/domain"
	"github.com/user/mangahub/internal/eventbus"
	"github.com/user/mangahub/pkg/models"
)

// ProgressService dictates operations around separated Relational Pivot maps tracking spatial contexts.
type ProgressService interface {
	UpdateProgress(progress *models.UserProgress) error
	GetUserProgress(userID int) ([]*models.UserProgress, error)
}

type progressService struct {
	repo domain.ProgressRepository
	bus  *eventbus.EventBus
}

func NewProgressService(repo domain.ProgressRepository, bus *eventbus.EventBus) ProgressService {
	return &progressService{repo: repo, bus: bus}
}

func (s *progressService) UpdateProgress(progress *models.UserProgress) error {
	if err := s.repo.Save(progress); err != nil {
		return err
	}

	// Trigger immediate WebSocket feedback globally safely across Channels
	s.bus.Publish(models.Event{
		Topic:   "progress.updated",
		Payload: progress,
	})

	return nil
}

func (s *progressService) GetUserProgress(userID int) ([]*models.UserProgress, error) {
	return s.repo.GetByUserID(userID)
}

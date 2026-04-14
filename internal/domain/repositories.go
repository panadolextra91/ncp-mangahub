package domain

import "github.com/user/mangahub/pkg/models"

// UserRepository controls persistence gateway for User contexts.
type UserRepository interface {
	Save(user *models.User) error
	FindByUsername(username string) (*models.User, error)
}

// MangaRepository controls persistence gateway for Manga inventory.
type MangaRepository interface {
	Save(manga *models.Manga) error
	FindByID(id int) (*models.Manga, error)
}

// ProgressRepository controls persistence gateway orchestrating relational mapping states.
type ProgressRepository interface {
	Save(progress *models.UserProgress) error
	FindByUserAndManga(userID, mangaID int) (*models.UserProgress, error)
}

type ChatRepository interface {
	Save(msg *models.ChatMessage) error
	GetRecentByManga(mangaID int, limit int) ([]*models.ChatMessage, error)
}

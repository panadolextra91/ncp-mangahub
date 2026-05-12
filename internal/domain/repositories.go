package domain

import "github.com/user/mangahub/pkg/models"

// SearchFilters carries the optional, combinable criteria for advanced manga search.
//
// Lives in the domain package so the MangaRepository interface can reference it
// without creating a domain→application import cycle. The application package
// re-exports it as a type alias for ergonomic call-sites.
type SearchFilters struct {
	Query  string   // optional substring match against title/author
	Genres []string // optional OR-match against genres JSON column (cap 10 enforced at handler + repo)
	Status string   // optional equality match
	SortBy string   // "recent" (default, id DESC) or "title" (LOWER(title) ASC); unknown values fall back to recent
}

// UserRepository controls persistence gateway for User contexts.
type UserRepository interface {
	Save(user *models.User) error
	FindByUsername(username string) (*models.User, error)
}

// MangaRepository controls persistence gateway for Manga inventory.
type MangaRepository interface {
	Save(manga *models.Manga) error
	FindByID(id int) (*models.Manga, error)
	List() ([]*models.Manga, error)
	Search(query string) ([]*models.Manga, error)
	SearchWithFilters(f SearchFilters) ([]*models.Manga, error)
}

// ProgressRepository controls persistence gateway orchestrating relational mapping states.
type ProgressRepository interface {
	Save(progress *models.UserProgress) error
	FindByUserAndManga(userID, mangaID int) (*models.UserProgress, error)
	GetByUserID(userID int) ([]*models.UserProgress, error)
}

type ChatRepository interface {
	Save(msg *models.ChatMessage) error
	GetRecentByManga(mangaID int, limit int) ([]*models.ChatMessage, error)
}

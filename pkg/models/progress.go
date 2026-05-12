package models

import "time"

// UserProgress acts as the spatial pivot table tracking temporal chapter reading states
// explicitly disconnected from global Manga entity states to prevent write-overwrites.
type UserProgress struct {
	UserID         int       `json:"user_id"`
	MangaID        int       `json:"manga_id"`
	CurrentChapter int       `json:"current_chapter"`
	Status         string    `json:"status"` // Reading, Completed, Plan to Read
	UpdatedAt      time.Time `json:"updated_at"`
}

package models

import "time"

// UserProgress acts as the spatial pivot table tracking temporal chapter reading states
// explicitly disconnected from global Manga entity states to prevent write-overwrites.
type UserProgress struct {
	UserID         int
	MangaID        int
	CurrentChapter int
	UpdatedAt      time.Time
}

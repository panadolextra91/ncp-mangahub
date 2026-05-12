package models

import "time"

// Manga represents the core content entity.
type Manga struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Genres        string    `json:"genres"`         // Comma separated or JSON string
	Status        string    `json:"status"`         // Ongoing, Completed
	TotalChapters int       `json:"total_chapters"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

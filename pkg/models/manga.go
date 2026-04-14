package models

import "time"

// Manga represents the core content entity.
type Manga struct {
	ID        int
	Title     string
	Author    string
	CreatedAt time.Time
}

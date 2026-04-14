package models

import "time"

// ChatMessage represents a single message in a manga chat room.
type ChatMessage struct {
	ID         int       `json:"id"`
	MangaID    int       `json:"manga_id"`
	UserID     int       `json:"user_id"`
	SenderName string    `json:"sender_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

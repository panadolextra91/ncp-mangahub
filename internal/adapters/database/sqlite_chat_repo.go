package database

import (
	"database/sql"
	"github.com/user/mangahub/pkg/models"
)

type SqliteChatRepository struct {
	db *sql.DB
}

func NewSqliteChatRepository(db *sql.DB) *SqliteChatRepository {
	return &SqliteChatRepository{db: db}
}

func (r *SqliteChatRepository) Save(msg *models.ChatMessage) error {
	query := `INSERT INTO chat_messages (manga_id, user_id, sender_name, content) VALUES (?, ?, ?, ?)`
	res, err := r.db.Exec(query, msg.MangaID, msg.UserID, msg.SenderName, msg.Content)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	msg.ID = int(id)
	return nil
}

func (r *SqliteChatRepository) GetRecentByManga(mangaID int, limit int) ([]*models.ChatMessage, error) {
	query := `SELECT id, manga_id, user_id, sender_name, content, created_at FROM chat_messages WHERE manga_id = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.Query(query, mangaID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.ChatMessage
	for rows.Next() {
		msg := &models.ChatMessage{}
		err := rows.Scan(&msg.ID, &msg.MangaID, &msg.UserID, &msg.SenderName, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		// prepend to keep chronological order for client
		messages = append([]*models.ChatMessage{msg}, messages...)
	}
	return messages, nil
}

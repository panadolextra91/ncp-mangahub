package database_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/internal/adapters/database"
	"github.com/user/mangahub/internal/infrastructure"
	"github.com/user/mangahub/pkg/models"
)

func TestSqliteChatRepository(t *testing.T) {
	db, _ := infrastructure.NewSQLiteDB(":memory:")
	infrastructure.InitSchema(db)
	repo := database.NewSqliteChatRepository(db)

	t.Run("Save and Get Recent", func(t *testing.T) {
		msg := &models.ChatMessage{
			MangaID:    1,
			UserID:     10,
			SenderName: "Alice",
			Content:    "Hello World",
		}

		err := repo.Save(msg)
		assert.NoError(t, err)
		assert.NotZero(t, msg.ID)

		// Save some more
		for i := 0; i < 25; i++ {
			repo.Save(&models.ChatMessage{MangaID: 1, UserID: 10, SenderName: "Alice", Content: "Msg"})
		}
		// Save to another manga
		repo.Save(&models.ChatMessage{MangaID: 2, UserID: 10, SenderName: "Alice", Content: "Other"})

		recent, err := repo.GetRecentByManga(1, 20)
		assert.NoError(t, err)
		assert.Len(t, recent, 20)
		assert.Equal(t, 1, recent[0].MangaID)
	})
}

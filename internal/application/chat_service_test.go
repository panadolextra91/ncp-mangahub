package application_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/internal/application"
	"github.com/user/mangahub/internal/eventbus"
	"github.com/user/mangahub/pkg/models"
)

type mockChatRepo struct {
	savedCount int
}

func (m *mockChatRepo) Save(msg *models.ChatMessage) error {
	m.savedCount++
	msg.ID = m.savedCount
	return nil
}

func (m *mockChatRepo) GetRecentByManga(mangaID int, limit int) ([]*models.ChatMessage, error) {
	return []*models.ChatMessage{{MangaID: mangaID}}, nil
}

func TestChatService(t *testing.T) {
	bus := eventbus.NewEventBus(10)
	repo := &mockChatRepo{}
	svc := application.NewChatService(repo, bus)

	t.Run("SendMessage and Publish", func(t *testing.T) {
		ch := bus.Subscribe("chat.message")

		msg := &models.ChatMessage{MangaID: 1, Content: "Hello"}
		err := svc.SendMessage(msg)
		assert.NoError(t, err)
		assert.Equal(t, 1, repo.savedCount)

		// Check event bus
		received := <-ch
		assert.Equal(t, "chat.message", received.Topic)
		assert.Equal(t, msg, received.Payload)
	})

	t.Run("GetHistory", func(t *testing.T) {
		history, err := svc.GetHistory(1)
		assert.NoError(t, err)
		assert.Len(t, history, 1)
	})
}

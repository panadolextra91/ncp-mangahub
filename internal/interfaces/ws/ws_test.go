package ws_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/internal/application"
	"github.com/user/mangahub/internal/eventbus"
	"github.com/user/mangahub/internal/interfaces/ws"
	"github.com/user/mangahub/pkg/auth"
	"github.com/user/mangahub/pkg/models"
)

type mockChatRepo struct{}
func (m *mockChatRepo) Save(msg *models.ChatMessage) error { return nil }
func (m *mockChatRepo) GetRecentByManga(mangaID int, limit int) ([]*models.ChatMessage, error) {
	return []*models.ChatMessage{{Content: "History 1", MangaID: mangaID}}, nil
}

func TestWebSocketProtocol(t *testing.T) {
	secret := "test-secret"
	bus := eventbus.NewEventBus(10)
	repo := &mockChatRepo{}
	svc := application.NewChatService(repo, bus)
	hub := ws.NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go hub.Run(ctx, &wg)
	handler := ws.NewChatHandler(hub, svc, secret)

	s := httptest.NewServer(http.HandlerFunc(handler.HandleWS))
	defer s.Close()

	url := "ws" + strings.TrimPrefix(s.URL, "http")

	t.Run("Handshake Success & History", func(t *testing.T) {
		token, _ := auth.GenerateToken(1, "user", secret)
		fullURL := url + "?manga_id=1&token=" + token

		conn, _, err := websocket.DefaultDialer.Dial(fullURL, nil)
		assert.NoError(t, err)
		defer conn.Close()

		// Read history
		var msg models.ChatMessage
		_, p, err := conn.ReadMessage()
		assert.NoError(t, err)
		json.Unmarshal(p, &msg)
		assert.Equal(t, "History 1", msg.Content)
	})

	t.Run("Handshake Fail - Invalid Token", func(t *testing.T) {
		fullURL := url + "?manga_id=1&token=bad"
		_, resp, err := websocket.DefaultDialer.Dial(fullURL, nil)
		assert.Error(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Chat & Broadcast", func(t *testing.T) {
		token, _ := auth.GenerateToken(1, "user", secret)
		fullURL := url + "?manga_id=1&token=" + token

		conn, _, _ := websocket.DefaultDialer.Dial(fullURL, nil)
		defer conn.Close()
		
		// Consume history
		conn.ReadMessage()

		// Send message
		content := "Realtime message"
		err := conn.WriteMessage(websocket.TextMessage, []byte(content))
		assert.NoError(t, err)

		// Hub should broadcast back through the bridge (mocking bridge here by manually pushing to hub)
		hub.Broadcast <- &models.ChatMessage{MangaID: 1, Content: content, SenderName: "user"}

		_, p, _ := conn.ReadMessage()
		var received models.ChatMessage
		json.Unmarshal(p, &received)
		assert.Equal(t, content, received.Content)
	})
}

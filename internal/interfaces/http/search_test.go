package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/mangahub/pkg/models"
)

type MockMangaService struct {
	mock.Mock
}

func (m *MockMangaService) CreateManga(role string, manga *models.Manga) error {
	args := m.Called(role, manga)
	return args.Error(0)
}

func (m *MockMangaService) GetManga(id int) (*models.Manga, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Manga), args.Error(1)
}

func (m *MockMangaService) ListMangas() ([]*models.Manga, error) {
	args := m.Called()
	return args.Get(0).([]*models.Manga), args.Error(1)
}

func (m *MockMangaService) SearchMangas(query string) ([]*models.Manga, error) {
	args := m.Called(query)
	return args.Get(0).([]*models.Manga), args.Error(1)
}

func TestMangaListSearch(t *testing.T) {
	svc := new(MockMangaService)
	handler := NewMangaHandler(svc)

	t.Run("List all mangas", func(t *testing.T) {
		svc.On("ListMangas").Return([]*models.Manga{{Title: "Manga 1"}, {Title: "Manga 2"}}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/manga", nil)
		rr := httptest.NewRecorder()
		handler.List(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var result []*models.Manga
		json.Unmarshal(rr.Body.Bytes(), &result)
		assert.Equal(t, 2, len(result))
	})

	t.Run("Search mangas", func(t *testing.T) {
		svc.On("SearchMangas", "One Piece").Return([]*models.Manga{{Title: "One Piece"}}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/manga?q=One+Piece", nil)
		rr := httptest.NewRecorder()
		handler.List(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var result []*models.Manga
		json.Unmarshal(rr.Body.Bytes(), &result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "One Piece", result[0].Title)
	})

	t.Run("Empty search result", func(t *testing.T) {
		svc.On("SearchMangas", "NonExistent").Return([]*models.Manga{}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/manga?q=NonExistent", nil)
		rr := httptest.NewRecorder()
		handler.List(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var result []*models.Manga
		json.Unmarshal(rr.Body.Bytes(), &result)
		assert.Equal(t, 0, len(result))
	})
}

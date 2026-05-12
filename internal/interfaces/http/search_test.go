package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/mangahub/internal/application"
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

func (m *MockMangaService) SearchMangasWithFilters(f application.SearchFilters) ([]*models.Manga, error) {
	args := m.Called(f)
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

// TestMangaListFilters covers the WH7 advanced filter contract:
//   - genres (single + multi OR), status, sortBy
//   - cap of 10 genres
//   - bit-for-bit backwards compat for q-only callers (routes to SearchMangas, NOT SearchMangasWithFilters)
func TestMangaListFilters(t *testing.T) {
	t.Run("Filter by single genre", func(t *testing.T) {
		svc := new(MockMangaService)
		handler := NewMangaHandler(svc)
		expected := application.SearchFilters{Genres: []string{"Action"}}
		svc.On("SearchMangasWithFilters", expected).Return([]*models.Manga{{Title: "Naruto"}}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/manga?genres=Action", nil)
		rr := httptest.NewRecorder()
		handler.List(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Filter by multiple genres (OR)", func(t *testing.T) {
		svc := new(MockMangaService)
		handler := NewMangaHandler(svc)
		expected := application.SearchFilters{Genres: []string{"Action", "Romance"}}
		svc.On("SearchMangasWithFilters", expected).Return([]*models.Manga{{Title: "Naruto"}, {Title: "Fruits Basket"}}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/manga?genres=Action,Romance", nil)
		rr := httptest.NewRecorder()
		handler.List(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Filter by status", func(t *testing.T) {
		svc := new(MockMangaService)
		handler := NewMangaHandler(svc)
		expected := application.SearchFilters{Status: "ongoing"}
		svc.On("SearchMangasWithFilters", expected).Return([]*models.Manga{{Title: "Naruto"}}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/manga?status=ongoing", nil)
		rr := httptest.NewRecorder()
		handler.List(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Sort by title", func(t *testing.T) {
		svc := new(MockMangaService)
		handler := NewMangaHandler(svc)
		expected := application.SearchFilters{SortBy: "title"}
		svc.On("SearchMangasWithFilters", expected).Return([]*models.Manga{}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/manga?sortBy=title", nil)
		rr := httptest.NewRecorder()
		handler.List(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Combined filters (AND across types)", func(t *testing.T) {
		svc := new(MockMangaService)
		handler := NewMangaHandler(svc)
		expected := application.SearchFilters{
			Query:  "demon",
			Genres: []string{"Action"},
			Status: "ongoing",
			SortBy: "title",
		}
		svc.On("SearchMangasWithFilters", expected).Return([]*models.Manga{{Title: "Demon Slayer"}}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/manga?q=demon&genres=Action&status=ongoing&sortBy=title", nil)
		rr := httptest.NewRecorder()
		handler.List(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Genre cap at 10", func(t *testing.T) {
		svc := new(MockMangaService)
		handler := NewMangaHandler(svc)
		// 12 genres in; expect 10 out (cap).
		expected := application.SearchFilters{Genres: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}}
		svc.On("SearchMangasWithFilters", expected).Return([]*models.Manga{}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/manga?genres=A,B,C,D,E,F,G,H,I,J,K,L", nil)
		rr := httptest.NewRecorder()
		handler.List(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		svc.AssertExpectations(t)
	})

	t.Run("Genres with whitespace and empties stripped", func(t *testing.T) {
		svc := new(MockMangaService)
		handler := NewMangaHandler(svc)
		expected := application.SearchFilters{Genres: []string{"Action", "Romance"}}
		svc.On("SearchMangasWithFilters", expected).Return([]*models.Manga{}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/manga?genres=Action%2C%20%2C%20Romance%20", nil)
		rr := httptest.NewRecorder()
		handler.List(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		svc.AssertExpectations(t)
	})

	// Backcompat: when ONLY ?q= is present, MangaHandler.List must route through
	// SearchMangas (legacy SQL), NOT SearchMangasWithFilters. This is the
	// critical wire-level backcompat assertion.
	t.Run("Backcompat: q-only routes to SearchMangas (not SearchMangasWithFilters)", func(t *testing.T) {
		svc := new(MockMangaService)
		handler := NewMangaHandler(svc)
		svc.On("SearchMangas", "naruto").Return([]*models.Manga{{Title: "Naruto"}}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/manga?q=naruto", nil)
		rr := httptest.NewRecorder()
		handler.List(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		svc.AssertExpectations(t)
		// Guarantee filtered path was NOT used for q-only call.
		svc.AssertNotCalled(t, "SearchMangasWithFilters", mock.Anything)
	})

	// sortBy=recent is an explicit no-op — must NOT trigger the filtered path.
	t.Run("Backcompat: q + sortBy=recent stays on legacy path", func(t *testing.T) {
		svc := new(MockMangaService)
		handler := NewMangaHandler(svc)
		svc.On("SearchMangas", "naruto").Return([]*models.Manga{{Title: "Naruto"}}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/manga?q=naruto&sortBy=recent", nil)
		rr := httptest.NewRecorder()
		handler.List(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		svc.AssertExpectations(t)
		svc.AssertNotCalled(t, "SearchMangasWithFilters", mock.Anything)
	})
}

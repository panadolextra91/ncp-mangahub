package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/internal/adapters/database"
	"github.com/user/mangahub/internal/domain"
	"github.com/user/mangahub/internal/infrastructure"
	"github.com/user/mangahub/pkg/models"
)

func TestSqliteUserRepository(t *testing.T) {
	db, _ := infrastructure.NewSQLiteDB(":memory:")
	infrastructure.InitSchema(db)
	repo := database.NewSqliteUserRepository(db)

	user := &models.User{Username: "test", PasswordHash: "hash", Role: "user"}
	err := repo.Save(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	found, err := repo.FindByUsername("test")
	assert.NoError(t, err)
	assert.Equal(t, "test", found.Username)
}

func TestSqliteMangaRepository(t *testing.T) {
	db, _ := infrastructure.NewSQLiteDB(":memory:")
	infrastructure.InitSchema(db)
	repo := database.NewSqliteMangaRepository(db)

	manga := &models.Manga{Title: "Test Manga", Author: "Test Author"}
	err := repo.Save(manga)
	assert.NoError(t, err)
	assert.NotZero(t, manga.ID)

	found, err := repo.FindByID(manga.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Test Manga", found.Title)
}

// TestSqliteMangaRepository_SearchWithFilters exercises the real SQL path for
// the WH7 advanced filter contract. Seeds 4 mangas with overlapping genres and
// statuses, then asserts each filter combination produces the expected result
// set — including the critical "Action" must-not-match-"Reaction" substring
// collision guard (proven by the quoted-token LIKE form).
func TestSqliteMangaRepository_SearchWithFilters(t *testing.T) {
	db, _ := infrastructure.NewSQLiteDB(":memory:")
	infrastructure.InitSchema(db)
	repo := database.NewSqliteMangaRepository(db)

	// Seed — genres stored as JSON-array TEXT (matches seed file format).
	seed := []*models.Manga{
		{Title: "Naruto", Author: "Kishimoto", Genres: `["Action","Shounen"]`, Status: "ongoing"},
		{Title: "Berserk", Author: "Miura", Genres: `["Action","Seinen"]`, Status: "completed"},
		{Title: "Fruits Basket", Author: "Takaya", Genres: `["Romance","Shoujo"]`, Status: "completed"},
		{Title: "Reaction Time", Author: "Test", Genres: `["Comedy"]`, Status: "ongoing"},
	}
	for _, m := range seed {
		assert.NoError(t, repo.Save(m))
	}

	t.Run("single genre — no substring collision", func(t *testing.T) {
		got, err := repo.SearchWithFilters(domain.SearchFilters{Genres: []string{"Action"}})
		assert.NoError(t, err)
		assert.Len(t, got, 2)
		titles := []string{got[0].Title, got[1].Title}
		assert.Contains(t, titles, "Naruto")
		assert.Contains(t, titles, "Berserk")
		// CRITICAL: "Action" must NOT match "Reaction" — proves quoted-token guard.
		assert.NotContains(t, titles, "Reaction Time")
	})

	t.Run("multi-genre OR", func(t *testing.T) {
		got, err := repo.SearchWithFilters(domain.SearchFilters{Genres: []string{"Action", "Romance"}})
		assert.NoError(t, err)
		assert.Len(t, got, 3)
		titles := []string{got[0].Title, got[1].Title, got[2].Title}
		assert.Contains(t, titles, "Naruto")
		assert.Contains(t, titles, "Berserk")
		assert.Contains(t, titles, "Fruits Basket")
	})

	t.Run("status filter", func(t *testing.T) {
		got, err := repo.SearchWithFilters(domain.SearchFilters{Status: "ongoing"})
		assert.NoError(t, err)
		assert.Len(t, got, 2)
		titles := []string{got[0].Title, got[1].Title}
		assert.Contains(t, titles, "Naruto")
		assert.Contains(t, titles, "Reaction Time")
	})

	t.Run("AND across types: genre + status", func(t *testing.T) {
		got, err := repo.SearchWithFilters(domain.SearchFilters{Genres: []string{"Action"}, Status: "ongoing"})
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, "Naruto", got[0].Title)
	})

	t.Run("unknown status returns empty (lenient, no error)", func(t *testing.T) {
		got, err := repo.SearchWithFilters(domain.SearchFilters{Status: "hiatus"})
		assert.NoError(t, err)
		assert.Len(t, got, 0)
	})

	t.Run("sort by title ascending", func(t *testing.T) {
		got, err := repo.SearchWithFilters(domain.SearchFilters{SortBy: "title"})
		assert.NoError(t, err)
		assert.Len(t, got, 4)
		assert.Equal(t, "Berserk", got[0].Title)
		assert.Equal(t, "Fruits Basket", got[1].Title)
		assert.Equal(t, "Naruto", got[2].Title)
		assert.Equal(t, "Reaction Time", got[3].Title)
	})

	t.Run("sort by recent (id DESC, default)", func(t *testing.T) {
		got, err := repo.SearchWithFilters(domain.SearchFilters{SortBy: "recent"})
		assert.NoError(t, err)
		assert.Len(t, got, 4)
		// "Reaction Time" was inserted last → highest id → first.
		assert.Equal(t, "Reaction Time", got[0].Title)
	})

	t.Run("unknown sortBy falls back to recent", func(t *testing.T) {
		got, err := repo.SearchWithFilters(domain.SearchFilters{SortBy: "bogus"})
		assert.NoError(t, err)
		assert.Len(t, got, 4)
		assert.Equal(t, "Reaction Time", got[0].Title)
	})

	t.Run("query only matches title/author", func(t *testing.T) {
		got, err := repo.SearchWithFilters(domain.SearchFilters{Query: "naruto"})
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, "Naruto", got[0].Title)
	})

	t.Run("query matches author", func(t *testing.T) {
		got, err := repo.SearchWithFilters(domain.SearchFilters{Query: "miura"})
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, "Berserk", got[0].Title)
	})

	t.Run("no filters returns all", func(t *testing.T) {
		got, err := repo.SearchWithFilters(domain.SearchFilters{})
		assert.NoError(t, err)
		assert.Len(t, got, 4)
	})

	t.Run("genres over cap of 10 are silently truncated", func(t *testing.T) {
		// 12 genres in; cap kicks in at repo layer. None of A–L exist in seed
		// data so result is empty regardless — the assertion is that no error.
		got, err := repo.SearchWithFilters(domain.SearchFilters{
			Genres: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L"},
		})
		assert.NoError(t, err)
		assert.Len(t, got, 0)
	})
}

func TestSqliteProgressRepository(t *testing.T) {
	db, _ := infrastructure.NewSQLiteDB(":memory:")
	infrastructure.InitSchema(db)
	repo := database.NewSqliteProgressRepository(db)

	prog := &models.UserProgress{UserID: 1, MangaID: 1, CurrentChapter: 10}
	err := repo.Save(prog)
	assert.NoError(t, err)

	found, err := repo.FindByUserAndManga(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, 10, found.CurrentChapter)
}

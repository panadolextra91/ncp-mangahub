package database_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/internal/adapters/database"
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

package database

import (
	"encoding/json"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/internal/infrastructure"
	"github.com/user/mangahub/pkg/models"
)

func TestSeeding(t *testing.T) {
	dbPath := "test_seeding.db"
	defer os.Remove(dbPath)

	db, _ := infrastructure.NewSQLiteDB(dbPath)
	defer db.Close()
	infrastructure.InitSchema(db)

	mangaRepo := NewSqliteMangaRepository(db)

	// Mock seeding logic
	seedFile, err := os.Open("../../../data/manga_seed.json")
	assert.NoError(t, err)
	defer seedFile.Close()

	var seedData []*models.Manga
	err = json.NewDecoder(seedFile).Decode(&seedData)
	assert.NoError(t, err)

	for _, m := range seedData {
		err := mangaRepo.Save(m)
		assert.NoError(t, err)
	}

	// Verify count
	all, err := mangaRepo.List()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(all), 40)

	// Verify fields
	onePiece, err := mangaRepo.Search("One Piece")
	assert.NoError(t, err)
	assert.NotEmpty(t, onePiece)
	assert.Equal(t, "Eiichiro Oda", onePiece[0].Author)
	assert.Equal(t, "Ongoing", onePiece[0].Status)
	assert.Contains(t, onePiece[0].Genres, "Action")
}

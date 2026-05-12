package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/user/mangahub/internal/adapters/database"
	"github.com/user/mangahub/pkg/models"
)

func main() {
	log.Println("🌱 Manual Seeding Process Started...")

	// 1. Kết nối DB
	db, err := database.NewSQLiteDB("mangahub.db")
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	defer db.Close()

	repo := database.NewSQLiteMangaRepo(db)

	// 2. Đọc file JSON
	data, err := os.ReadFile("data/manga_seed.json")
	if err != nil {
		log.Fatalf("❌ Failed to read seed file: %v", err)
	}

	var seedMangas []models.Manga
	if err := json.Unmarshal(data, &seedMangas); err != nil {
		log.Fatalf("❌ Failed to parse JSON: %v", err)
	}

	// 3. Insert từng bộ truyện
	count := 0
	for _, m := range seedMangas {
		err := repo.Save(&m)
		if err != nil {
			log.Printf("⚠️ Skip '%s': %v", m.Title, err)
			continue
		}
		count++
	}

	log.Printf("✅ Successfully seeded %d new manga records!", count)
}

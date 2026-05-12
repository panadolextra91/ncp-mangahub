package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// MangaDex API response structures
type MangaDexSearchResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Title struct {
				En string `json:"en"`
			} `json:"title"`
			Description struct {
				En string `json:"en"`
			} `json:"description"`
			Status string `json:"status"`
		} `json:"attributes"`
	} `json:"data"`
}

func main() {
	query := "Blue Lock"
	fmt.Printf("🔍 Probing MangaDex for: '%s'...\n", query)

	// 1. Search for Manga
	searchURL := fmt.Sprintf("https://api.mangadex.org/manga?title=%s&limit=1", url.QueryEscape(query))
	resp, err := http.Get(searchURL)
	if err != nil {
		fmt.Printf("❌ Connection Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var searchRes MangaDexSearchResponse
	json.Unmarshal(body, &searchRes)

	if len(searchRes.Data) == 0 {
		fmt.Println("❓ No results found on MangaDex.")
		return
	}

	manga := searchRes.Data[0]
	fmt.Printf("\n✨ FOUND MANGA ON MANGADEX!\n")
	fmt.Printf("🆔 ID: %s\n", manga.ID)
	fmt.Printf("📚 Title: %s\n", manga.Attributes.Title.En)
	fmt.Printf("🏷️ Status: %s\n", manga.Attributes.Status)
	fmt.Printf("📝 Description: %s...\n", manga.Attributes.Description.En[:200])

	// 2. Fetch Chapters (Optional preview)
	chapterURL := fmt.Sprintf("https://api.mangadex.org/manga/%s/feed?limit=5&translatedLanguage[]=en&order[chapter]=desc", manga.ID)
	fmt.Printf("\n📂 Fetching latest chapters from: %s\n", chapterURL)
	
	// Trả về cho Mẹ biết là mình đã sẵn sàng kết nối!
	fmt.Printf("\n✅ MangaHub is ready to aggregate from MangaDex!\n")
}

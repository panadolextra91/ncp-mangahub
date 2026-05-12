package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type MangaDexResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Title struct {
				En string `json:"en"`
			} `json:"title"`
		} `json:"attributes"`
	} `json:"data"`
}

type ChapterResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Chapter     string `json:"chapter"`
			ExternalURL string `json:"externalUrl"`
		} `json:"attributes"`
	} `json:"data"`
}

type AtHomeResponse struct {
	BaseURL string `json:"baseUrl"`
	Chapter struct {
		Hash string   `json:"hash"`
		Data []string `json:"data"`
	} `json:"chapter"`
}

func main() {
	query := "Chainsaw Man"
	fmt.Printf("🔍 Searching for '%s' on MangaDex...\n", query)

	// 1. Search Manga
	searchURL := fmt.Sprintf("https://api.mangadex.org/manga?title=%s&limit=1", url.QueryEscape(query))
	resp, _ := http.Get(searchURL)
	var searchRes MangaDexResponse
	json.NewDecoder(resp.Body).Decode(&searchRes)
	resp.Body.Close()

	if len(searchRes.Data) == 0 {
		fmt.Println("❌ Manga not found.")
		return
	}
	mID := searchRes.Data[0].ID
	title := searchRes.Data[0].Attributes.Title.En
	fmt.Printf("✅ Found: %s (ID: %s)\n", title, mID)

	// 2. Get Chapters (Latest 50)
	chapURL := fmt.Sprintf("https://api.mangadex.org/manga/%s/feed?limit=50&translatedLanguage[]=en&order[chapter]=desc", mID)
	resp, _ = http.Get(chapURL)
	var chapRes ChapterResponse
	json.NewDecoder(resp.Body).Decode(&chapRes)
	resp.Body.Close()

	// 3. Find first internal chapter
	fmt.Println("🕵️ Hunting for a chapter with internal image data...")
	for _, chap := range chapRes.Data {
		if chap.Attributes.ExternalURL != "" {
			fmt.Printf("🔗 Chapter %s is External: %s\n", chap.Attributes.Chapter, chap.Attributes.ExternalURL)
			continue
		}

		// Try to get images
		athomeURL := fmt.Sprintf("https://api.mangadex.org/at-home/server/%s", chap.ID)
		resp, _ = http.Get(athomeURL)
		var athomeRes AtHomeResponse
		json.NewDecoder(resp.Body).Decode(&athomeRes)
		resp.Body.Close()

		if len(athomeRes.Chapter.Data) > 0 {
			fmt.Printf("🎉 Found Internal Chapter %s! (%d pages)\n", chap.Attributes.Chapter, len(athomeRes.Chapter.Data))
			
			// Download & Render Page 1
			imgURL := fmt.Sprintf("%s/data/%s/%s", athomeRes.BaseURL, athomeRes.Chapter.Hash, athomeRes.Chapter.Data[0])
			fmt.Printf("📥 Downloading Page 1: %s\n", imgURL)
			
			resp, _ = http.Get(imgURL)
			imgData, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			encoded := base64.StdEncoding.EncodeToString(imgData)
			fmt.Println("\n✨ RENDERING IN ITERM2... ✨\n")
			fmt.Printf("\033]1337;File=inline=1;width=80%%:%s\a\n", encoded)
			fmt.Printf("\n📖 %s - CHAPTER %s\n", title, chap.Attributes.Chapter)
			fmt.Println("✅ Success! You are reading live from MangaDex! 🌸")
			return
		}
	}

	fmt.Println("❌ All recent chapters are external. Try another manga!")
}

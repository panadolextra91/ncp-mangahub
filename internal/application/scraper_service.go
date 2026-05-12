package application

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/user/mangahub/pkg/models"
)

type ScraperService struct {
	targetURL string
}

func NewScraperService() *ScraperService {
	return &ScraperService{
		targetURL: "https://quotes.toscrape.com",
	}
}

// FetchQuotes fetches and parses quotes from quotes.toscrape.com using standard library
func (s *ScraperService) FetchQuotes() ([]models.Quote, error) {
	resp, err := http.Get(s.targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return s.parseHTML(string(body)), nil
}

// parseHTML uses regex to extract quotes, authors and tags
func (s *ScraperService) parseHTML(html string) []models.Quote {
	var quotes []models.Quote

	// Pattern for quote block
	quoteBlockRegex := regexp.MustCompile(`(?s)<div class="quote"[^>]*>(.*?)</div>`)
	textRegex := regexp.MustCompile(`<span class="text"[^>]*>“(.+?)”</span>`)
	authorRegex := regexp.MustCompile(`<small class="author"[^>]*>(.+?)</small>`)
	tagRegex := regexp.MustCompile(`<a class="tag"[^>]*>(.+?)</a>`)

	blocks := quoteBlockRegex.FindAllStringSubmatch(html, -1)
	for _, block := range blocks {
		if len(block) < 2 {
			continue
		}
		content := block[1]

		textMatch := textRegex.FindStringSubmatch(content)
		authorMatch := authorRegex.FindStringSubmatch(content)
		tagMatches := tagRegex.FindAllStringSubmatch(content, -1)

		if len(textMatch) < 2 || len(authorMatch) < 2 {
			continue
		}

		var tags []string
		for _, tm := range tagMatches {
			if len(tm) >= 2 {
				tags = append(tags, tm[1])
			}
		}

		quotes = append(quotes, models.Quote{
			Text:   s.cleanText(textMatch[1]),
			Author: authorMatch[1],
			Tags:   tags,
		})
	}

	return quotes
}

func (s *ScraperService) cleanText(text string) string {
	// Simple HTML entity decoding for common ones
	r := strings.NewReplacer("&ldquo;", "“", "&rdquo;", "”", "&rsquo;", "'", "&amp;", "&")
	return r.Replace(text)
}

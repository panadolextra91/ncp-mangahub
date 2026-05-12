package models

// Quote represents a web-scraped quote for educational practice
type Quote struct {
	Text   string   `json:"text"`
	Author string   `json:"author"`
	Tags   []string `json:"tags"`
}

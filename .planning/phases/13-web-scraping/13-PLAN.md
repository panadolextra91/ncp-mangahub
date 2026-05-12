# Phase 13: Web Scraping - Plan

## 🏁 Phase Goal
Implement a web scraping module to fetch inspirational quotes from `quotes.toscrape.com` and display them on the PinkHub TUI, adding a "Discover" feature and an enhanced "Quote-powered" dashboard.

## 📝 Tasks

### Wave 1: Foundation & Backend (The Scraper)
- [ ] **Task 1: Add Quote Model**
  - Create `pkg/models/quote.go` with `Quote` struct (Text, Author, Tags).
  - `read_first`: `pkg/models/manga.go` (as template).
- [ ] **Task 2: Implement Scraper Service**
  - Create `internal/application/scraper_service.go`.
  - Implement `FetchQuotes() ([]models.Quote, error)` using `net/http` and `regexp`.
  - `read_first`: `internal/application/manga_service.go`.
- [ ] **Task 2.1: Add Unit Tests for Scraper**
  - Create `internal/application/scraper_service_test.go`.
  - Add test cases with mock HTML to verify extraction logic.
  - Ensure 100% coverage for the parsing logic.

### Wave 2: TUI Integration (The Display)
- [ ] **Task 3: Update TUI Model**
  - Add `Quotes []models.Quote` and `CurrentQuoteIndex int` to `internal/interfaces/tui/model.go`.
  - Add `Discover` page constant to `Page` enum.
- [ ] **Task 4: Implement Discover View**
  - Create rendering logic for the "Discover" tab in `internal/interfaces/tui/app.go`.
  - Display list of quotes with authors and tags.
- [ ] **Task 5: Enhance Dashboard with Quotes**
  - Update `renderDashboard()` to show a random quote from the scraped list.
  - Position it creatively below the ASCII logo as per user request.

### Wave 3: Interaction & Polish
- [ ] **Task 6: Add Refresh Keybinding**
  - Update `Update()` in `app.go` to handle the `r` key.
  - Key `r` should trigger a background fetch and update the model.
- [ ] **Task 7: Initial Boot Scrape**
  - Trigger `FetchQuotes()` in `main.go` or TUI `Init()` so the app starts with data.

## 🧪 Verification Criteria
- [ ] `go test ./...` passes 100% (All system tests).
- [ ] New test cases in `internal/application/scraper_service_test.go` pass.
- [ ] TUI opens and shows a Quote on the dashboard.
- [ ] Pressing `5` switches to the Discover tab.
- [ ] Pressing `r` updates the quote on screen.

## 🎯 Must Haves
- [ ] Successful scraping from `quotes.toscrape.com`.
- [ ] In-memory storage only (no DB changes).
- [ ] Visible Quote on TUI dashboard.
- [ ] Functional "r" key for manual refresh.
- [ ] New feature covered by unit tests.

# Phase 13: Web Scraping - Context & Decisions

## 🎯 Domain Boundary
Implementing an educational web scraping module to fetch content from practice sites (`quotes.toscrape.com`) and display it beautifully on the PinkHub TUI.

## 🛠️ Decisions & Specifics

### 1. Content Strategy (The "What")
- **Target**: `https://quotes.toscrape.com` (as per instructor requirements).
- **Data**: Fetching "Quotes of the day" including text, author, and tags.
- **Persistence**: **In-memory only**. Scraped data is fresh per session and not stored in the SQLite database to keep the system lightweight.

### 2. Visual Presentation (The "Where")
- **Dashboard Integration**: A random quote will be displayed prominently on the main dashboard, right below the ASCII art header, serving as a decorative and inspirational element.
- **Dedicated Tab**: A new **DISCOVER [5]** tab will be added to the TUI to list the current batch of scraped quotes.
- **Styling**: Using "Pink Glow" lipgloss styles to ensure consistency with the rest of the app.

### 3. User Interaction (The "How")
- **Automation**: The scraping service will trigger automatically upon server/client startup to ensure data is available immediately.
- **Manual Refresh**: Users can press the `r` key while in the Discover tab (or dashboard) to re-trigger the scraper and fetch a fresh set of quotes.

## 🔗 Canonical Refs
- `docs/requirements.md` (Data Collection section)
- `internal/interfaces/tui/app.go` (For UI integration)
- `internal/application/manga_service.go` (As pattern for new services)

## 📋 Folded Todos
- N/A

## 🚀 Next Steps
- Research lightweight Go scraping libraries (or just use standard `net/http` + `regex/strings` to avoid heavy dependencies).
- Implement `ScraperService` in `internal/application`.
- Update TUI model and view to include the Discover tab and Quote display.

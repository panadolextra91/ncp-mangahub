# Phase 13: Web Scraping - Research

## 🔍 Investigation Results

### 1. Target: `https://quotes.toscrape.com`
The site is designed for scraping practice and has a very stable, semantic HTML structure:
- **Each Quote block**: `<div class="quote" ...>`
- **Quote Text**: `<span class="text" itemprop="text">“...”</span>`
- **Author**: `<small class="author" itemprop="author">Name</small>`
- **Tags**: `<a class="tag" ...>tag</a>`

### 2. Technical Approach
- **Protocol**: Standard HTTP GET.
- **Library**: We will use Go's standard `net/http` for fetching and `regexp` or `strings` for extraction to minimize external dependencies, fulfilling the "educational practice" requirement with raw Go skills.
- **Pattern**:
  - `regexp.MustCompile("<span class=\"text\"[^>]*>“(.+?)”</span>")`
  - `regexp.MustCompile("<small class=\"author\"[^>]*>(.+?)</small>")`

### 3. Integration Points
- **Backend Service**: Create a `ScraperService` in `internal/application`.
- **TUI Model**: Add `Quotes []models.Quote` to the TUI state.
- **TUI Keybindings**: Map `r` to trigger the `ScraperService`.
- **TUI Render**: Update `DashboardView` and create `DiscoverView`.

## 🛠️ Reusable Assets
- `pkg/models`: Add a simple `Quote` struct.
- `internal/interfaces/tui/styles.go`: Use existing `PinkGlow` styles.

## 🏁 Validation Architecture
- **Unit Test**: Test the parser logic with a mock HTML string.
- **Integration**: Verify the TUI displays a non-empty quote after "r" is pressed.

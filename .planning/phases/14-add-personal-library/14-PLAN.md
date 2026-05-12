# Phase 14: Add Personal Library - Plan

## 🏁 Phase Goal
Implement a personal library system in the PinkHub TUI, allowing users to save manga to their profile, track reading progress across 3 statuses (Reading, Completed, Plan to Read), and see real-time updates via TCP.

## 📝 Tasks

### Wave 1: TUI Navigation & UI Structure
- [ ] **Task 1: Update TUI Constants & Model**
  - Add `PageDetail` and `PageLibrary` to `Page` enum in `model.go`.
  - Add `LibraryResults []*models.UserProgress` and `SelectedManga *models.Manga` to `Model`.
  - Add `LibraryIndex int` for scrolling in library.
- [ ] **Task 2: Implement Manga Detail View**
  - Create `renderDetail()` in `app.go`.
  - Display full manga metadata: Title, Author, Genres, Status, Total Chapters, and Description.
  - Add "Press 'a' to Add to Library" prompt.
- [ ] **Task 3: Implement Library View**
  - Create `renderLibrary()` in `app.go`.
  - Show a table of saved manga with their current status and progress.

### Wave 2: Commands & Interaction
- [ ] **Task 4: Implement Detail Navigation**
  - Update `Update()` in `app.go`: Pressing `Enter` on a search result sets `SelectedManga` and switches to `PageDetail`.
- [ ] **Task 5: Implement Library Commands**
  - Add `AddLibraryCmd()`: `POST /api/manga/progress` with status "Plan to Read" and chapter 0.
  - Add `FetchLibraryCmd()`: `GET /api/manga/progress`.
  - Update `Init()` to fetch library on startup.
- [ ] **Task 6: Handle Library Keys**
  - In `PageDetail`: Pressing `a` triggers `AddLibraryCmd()`.
  - Global: Pressing `7` switches to `PageLibrary`.

### Wave 3: Polish & Sync Verification
- [ ] **Task 7: TCP Sync Integration**
  - Ensure TUI `Update()` handles incoming TCP messages related to `progress.updated`.
  - Trigger a library refresh when a sync message is received.
- [ ] **Task 8: Unit Testing**
  - Add tests for `ProgressService.GetUserProgress` if not already covered.

## 🧪 Verification Criteria
- [ ] `go test ./...` passes 100%.
- [ ] Search for a manga -> Press `Enter` -> Detail view appears with correct data.
- [ ] Press `a` in Detail view -> "Added to Library" status appears.
- [ ] Press `7` -> Library tab shows the added manga.
- [ ] Open 2 TUI clients -> Add manga in TUI 1 -> TUI 2 shows notification/update.

## 🎯 Must Haves
- [ ] Detail view for Manga.
- [ ] Functional Library tab (Tab [7]).
- [ ] Persistent storage in DB via HTTP.
- [ ] Real-time sync via TCP.
- [ ] All 3 statuses supported.

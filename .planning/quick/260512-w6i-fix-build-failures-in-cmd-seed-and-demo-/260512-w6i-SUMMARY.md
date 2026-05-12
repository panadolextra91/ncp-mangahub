---
phase: 260512-w6i
plan: 01
status: complete
date: 2026-05-12
---

# Quick Task 260512-w6i — Summary

## Goal

Fix build failures in `cmd/seed` and `demo/` so `go build ./...` and `go vet ./...` succeed and CI goes green.

## Outcome

**`go build ./...` exits 0** ✅
**`go vet ./...` exits 0** ✅

## Commits

| Commit | Task | Scope |
|--------|------|-------|
| `17d2b55` | 260512-w6i-01 | `fix: use correct infrastructure/database constructors in cmd/seed` |
| `bfed747` | 260512-w6i-02 | `refactor: split demo/ files into per-command subdirectories` |

## Changes

### Task 1 — `cmd/seed/main.go`
- Replaced undefined `database.NewSQLiteDB(...)` with `infrastructure.NewSQLiteDB(...)` (correct package).
- Replaced undefined `database.NewSQLiteMangaRepo(...)` with `database.NewSqliteMangaRepository(...)` (actual exported name — note `Sqlite` casing, full `Repository` suffix).
- Added missing `infrastructure.InitSchema(db)` call to match `cmd/server/main.go` wiring so the seed works on a fresh DB.

### Task 2 — `demo/` split (Option A)
Three files moved into per-command subdirectories so each becomes its own `main` package:

| Before | After |
|--------|-------|
| `demo/mangadex_probe.go` | `demo/mangadex-probe/main.go` |
| `demo/mangadex_read.go` | `demo/mangadex-read/main.go` |
| `demo/run_show.go` | `demo/run-show/main.go` |

Moves done via `git mv` to preserve file history.

### Task 3 — Verification
- `go build ./...` → 0
- `go vet ./...` → 0
- All 4 binaries build individually: `cmd/seed`, `demo/mangadex-probe`, `demo/mangadex-read`, `demo/run-show`
- Scope check: only `cmd/seed/` and `demo/<name>/` paths touched. Zero spillover into `internal/`, `pkg/`, `config/`, or `cmd/server/`.

## Deviations

1 auto-fix (Rule 3 — blocking):
- Fixed a pre-existing `go vet` warning in `demo/mangadex-read/main.go` (`fmt.Println` with trailing newline → swapped for `fmt.Print` or removed `\n`) so `go vet ./...` could exit 0 as required by the success criteria. Folded into commit `bfed747`.

## Must-haves verified

- [x] `go build ./...` exits 0
- [x] `go vet ./...` exits 0
- [x] `cmd/seed` compiles against actual exported constructors (`infrastructure.NewSQLiteDB`, `database.NewSqliteMangaRepository`)
- [x] Each demo program compiles as its own buildable command in its own subdirectory
- [x] No refactor of unrelated code

## Next step

Push to `origin/main` — CI workflow `MangaHub CI` will rerun with Go 1.25 and the now-buildable repo.

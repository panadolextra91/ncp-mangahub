# Phase 12: Fullfill Spec and Use Case - Validation

**Date:** 2026-05-11
**Phase:** 12-fullfill-spec-and-use-case

## 1. Schema Integrity
- **Test:** Verify `mangas` table has all new columns after migration.
- **Test:** Verify `user_progress` table has the `status` column.
- **Criteria:** `PRAGMA table_info(mangas)` returns the expected column list.

## 2. Seed Data Correctness
- **Test:** Start server and check if 40+ manga are present in the DB.
- **Test:** Verify "One Piece" exists with correct metadata (genres, etc.).
- **Criteria:** `GET /api/manga` returns a non-empty list of enriched objects.

## 3. CLI Command Parity
- **Test:** Run `mangahub auth login` and check if `~/.mangahub/config.yaml` is updated.
- **Test:** Run `mangahub manga search "One Piece"` and verify output matches `cli_manual.md`.
- **Criteria:** CLI commands exit with code 0 and print expected tables/info.

## 4. Search Functionality
- **Test:** Filter by genre "Shounen" and verify only Shounen manga are returned.
- **Test:** Filter by status "Ongoing".
- **Criteria:** Search results are correctly filtered in the DB query.

## 5. TUI Sync
- **Test:** Verify the new `status` field is visible in the TUI Library/Progress tab.
- **Criteria:** TUI UI reflects the updated schema.

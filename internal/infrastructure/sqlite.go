package infrastructure

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// NewSQLiteDB creates a new SQLite connection enforcing single write access and WAL mode.
func NewSQLiteDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Crucial for mitigating SQLITE_BUSY (database is locked) panics by keeping writes sequential.
	db.SetMaxOpenConns(1)

	// Enable Write-Ahead Logging for better read/write concurrency despite max open conns = 1
	// WAL ensures readers do not block writers and writers do not block readers.
	_, err = db.Exec("PRAGMA journal_mode = WAL; PRAGMA synchronous = NORMAL;")
	if err != nil {
		return nil, fmt.Errorf("failed to set PRAGMA: %w", err)
	}

	return db, nil
}

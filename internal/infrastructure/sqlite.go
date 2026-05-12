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

// InitSchema creates the necessary tables for MangaHub if they do not already exist.
func InitSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS mangas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		genres TEXT,
		status TEXT,
		total_chapters INTEGER,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS user_progress (
		user_id INTEGER NOT NULL,
		manga_id INTEGER NOT NULL,
		current_chapter INTEGER NOT NULL,
		status TEXT DEFAULT 'reading',
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, manga_id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (manga_id) REFERENCES mangas(id)
	);

	CREATE TABLE IF NOT EXISTS chat_messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		manga_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		sender_name TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_chat_manga_id ON chat_messages(manga_id);
	`
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Migrations for existing databases
	migrations := []string{
		"ALTER TABLE mangas ADD COLUMN genres TEXT",
		"ALTER TABLE mangas ADD COLUMN status TEXT",
		"ALTER TABLE mangas ADD COLUMN total_chapters INTEGER",
		"ALTER TABLE mangas ADD COLUMN description TEXT",
		"ALTER TABLE user_progress ADD COLUMN status TEXT DEFAULT 'reading'",
	}

	for _, m := range migrations {
		// Ignore errors (e.g., "duplicate column name") as SQLite doesn't have "IF NOT EXISTS" for ADD COLUMN
		_, _ = db.Exec(m)
	}

	return nil
}

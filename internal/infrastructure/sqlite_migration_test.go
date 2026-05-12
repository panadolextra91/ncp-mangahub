package infrastructure

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSchemaMigration(t *testing.T) {
	dbPath := "test_migration.db"
	defer os.Remove(dbPath)

	db, err := NewSQLiteDB(dbPath)
	assert.NoError(t, err)
	defer db.Close()

	// Initial Init
	err = InitSchema(db)
	assert.NoError(t, err)

	// Check if new columns exist in mangas
	rows, err := db.Query("PRAGMA table_info(mangas)")
	assert.NoError(t, err)
	defer rows.Close()

	columns := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, dtype string
		var notnull int
		var dfltValue interface{}
		var pk int
		err := rows.Scan(&cid, &name, &dtype, &notnull, &dfltValue, &pk)
		assert.NoError(t, err)
		columns[name] = true
	}

	assert.True(t, columns["genres"])
	assert.True(t, columns["status"])
	assert.True(t, columns["total_chapters"])
	assert.True(t, columns["description"])

	// Check if new columns exist in user_progress
	rows2, err := db.Query("PRAGMA table_info(user_progress)")
	assert.NoError(t, err)
	defer rows2.Close()

	columns2 := make(map[string]bool)
	for rows2.Next() {
		var cid int
		var name, dtype string
		var notnull int
		var dfltValue interface{}
		var pk int
		err := rows2.Scan(&cid, &name, &dtype, &notnull, &dfltValue, &pk)
		assert.NoError(t, err)
		columns2[name] = true
	}
	assert.True(t, columns2["status"])

	// Verify idempotency
	err = InitSchema(db)
	assert.NoError(t, err)
}

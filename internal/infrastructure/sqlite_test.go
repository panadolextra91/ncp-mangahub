package infrastructure_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/internal/infrastructure"
)

func TestSQLiteConnection(t *testing.T) {
	// memory mode DB for fast unit test
	db, err := infrastructure.NewSQLiteDB("file:conn_test.db?mode=memory&cache=shared")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = db.Close()
	assert.NoError(t, err)
}

func TestInitSchema(t *testing.T) {
	db, _ := infrastructure.NewSQLiteDB(":memory:")
	defer db.Close()

	err := infrastructure.InitSchema(db)
	assert.NoError(t, err)

	// Verify tables exist
	var count int
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestSQLiteConnection_Error(t *testing.T) {
	// A completely bad path should fail during creation or ping
	db, err := infrastructure.NewSQLiteDB("file:/this/path/does/not/exist.db")
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestSQLiteHellPath_ConcurrentWriteRace(t *testing.T) {
	// The Hell Path test: Spam 50 goroutines inserting data simultaneously.
	// If the connection pool wasn't tightly limited to 1 via SetMaxOpenConns,
	// SQLite would throw "database is locked" (SQLITE_BUSY) exceptions.
	
	db, err := infrastructure.NewSQLiteDB("file:hell_test.db?mode=memory&cache=shared")
	assert.NoError(t, err)
	defer db.Close()

	// Initial table setup
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS stress_test (id INTEGER PRIMARY KEY AUTOINCREMENT, val TEXT);")
	assert.NoError(t, err)

	goroutines := 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Channel to capture errors safely from goroutines
	errCh := make(chan error, goroutines)

	// Unleash the concurrent race
	for i := 0; i < goroutines; i++ {
		go func(val int) {
			defer wg.Done()
			
			// Deliberately induce different timings to maximize collision chances
			time.Sleep(time.Duration(val%5) * time.Millisecond)

			// Try to insert concurrently
			_, execErr := db.Exec("INSERT INTO stress_test (val) VALUES (?)", "stress_test_row")
			if execErr != nil {
				errCh <- execErr
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	// Check if any goroutine triggered a lock error
	for e := range errCh {
		assert.NoError(t, e, "Database threw an error during concurrent writes. This means WAL or max connections setup failed.")
	}

	// Verify that absolutely no inserts were dropped or lost
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM stress_test").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, goroutines, count, "Not all concurrent inserts were successfully recorded!")
}

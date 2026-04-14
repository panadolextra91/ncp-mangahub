package e2e

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResilience_ConcurrentDBWrites(t *testing.T) {
	s, err := StartServer("../../mangahub")
	assert.NoError(t, err)
	defer s.Stop()

	token := getAuthToken(t, s.Port)

	// Goal: Hammer the DB with concurrent CreateManga requests from multiple goroutines
	// to ensure WAL mode handles it without "database is locked" crashes.
	const concurrentUsers = 10
	var wg sync.WaitGroup
	wg.Add(concurrentUsers)

	for i := 0; i < concurrentUsers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				createManga(t, s.Port, token, fmt.Sprintf("Heavy-Manga-%d-%d", id, j))
			}
		}(i)
	}

	wg.Wait()
	// If we reached here without panic/error, WAL mode is working.
}

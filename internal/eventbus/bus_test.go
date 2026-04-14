package eventbus_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/internal/eventbus"
	"github.com/user/mangahub/pkg/models"
)

func TestEventBus_HappyPath(t *testing.T) {
	bus := eventbus.NewEventBus(10)
	ch := bus.Subscribe("test.topic")

	bus.Publish(models.Event{Topic: "test.topic", Payload: "hello"})

	select {
	case ev := <-ch:
		assert.Equal(t, "test.topic", ev.Topic)
		assert.Equal(t, "hello", ev.Payload)
	case <-time.After(1 * time.Second):
		t.Fatal("Expected event but got timeout")
	}

	assert.Equal(t, uint64(0), bus.DroppedCount())
}

func TestEventBus_Unsubscribe(t *testing.T) {
	bus := eventbus.NewEventBus(10)
	ch1 := bus.Subscribe("test.topic")
	ch2 := bus.Subscribe("test.topic")

	bus.Unsubscribe("test.topic", ch1)

	bus.Publish(models.Event{Topic: "test.topic", Payload: "test"})

	// ch2 should receive
	select {
	case <-ch2:
		// success
	case <-time.After(1 * time.Second):
		t.Fatal("Expected event on ch2")
	}

	// ch1 should be closed
	_, ok := <-ch1
	assert.False(t, ok, "Expected ch1 to be closed")
}

func TestEventBusHellPath_SlowConsumers(t *testing.T) {
	// Initialize bus with tiny buffer to ensure immediate dropping
	bus := eventbus.NewEventBus(1)
	
	// Create subscriber but we refuse to read from its channel
	_ = bus.Subscribe("spam_topic")

	concurrency := 100
	eventsPerRoutine := 1000
	totalEvents := concurrency * eventsPerRoutine

	var wg sync.WaitGroup
	wg.Add(concurrency)

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < eventsPerRoutine; j++ {
				bus.Publish(models.Event{Topic: "spam_topic", Payload: "spam"})
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	// Since buffer is only 1 and nobody is reading, the bus MUST drop exactly (totalEvents - 1) events.
	// This proves that `Publish` handles the overload instantly and safely increments metrics without locking up.
	dropped := bus.DroppedCount()
	assert.Greater(t, dropped, uint64(0), "Dropped events must be greater than zero")
	assert.Equal(t, uint64(totalEvents-1), dropped, "Exactly all blocked events must be tracked")

	// Extremely fast processing assures publishers are never stuck behind a dead consumer
	assert.Less(t, duration, 2*time.Second, "Hell test took too long, publishing is suspected of blocking!")
}

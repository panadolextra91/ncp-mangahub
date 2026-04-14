package eventbus

import (
	"sync"
	"sync/atomic"

	"github.com/user/mangahub/pkg/models"
)

// EventBus is the core concurrency-safe backbone for cross-protocol communication.
type EventBus struct {
	subscribers   map[string][]chan models.Event
	mu            sync.RWMutex
	bufferSize    int
	droppedEvents uint64
}

// NewEventBus instantiates a new internal Pub/Sub mechanism.
// Parameter bufferSize controls how much pressure a single subscriber can take before
// the bus starts dropping its messages to prevent global System deadlocks.
func NewEventBus(bufferSize int) *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan models.Event),
		bufferSize:  bufferSize,
	}
}

// Subscribe opens a new subscription queue for a particular topic.
func (b *EventBus) Subscribe(topic string) <-chan models.Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan models.Event, b.bufferSize)
	b.subscribers[topic] = append(b.subscribers[topic], ch)

	return ch
}

// Unsubscribe attempts to remove a given listener to prevent memory leaks and orphaned loops.
func (b *EventBus) Unsubscribe(topic string, targetCh <-chan models.Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	channels, exists := b.subscribers[topic]
	if !exists {
		return
	}

	for i, ch := range channels {
		if ch == targetCh {
			// Fast slice deletion by swapping with the last element
			b.subscribers[topic][i] = b.subscribers[topic][len(b.subscribers[topic])-1]
			b.subscribers[topic] = b.subscribers[topic][:len(b.subscribers[topic])-1]

			// Close the channel to clean up the consumer routine
			close(ch)
			break
		}
	}

	// Optimize map size if empty
	if len(b.subscribers[topic]) == 0 {
		delete(b.subscribers, topic)
	}
}

// Publish distributes an event aggressively to its topic subscribers.
// Any subscriber whose buffer is completely full will be ruthlessly skipped (event dropped).
// This rigidly guarantees Publish completes instantly every single time.
func (b *EventBus) Publish(event models.Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	channels, exists := b.subscribers[event.Topic]
	if !exists {
		return
	}

	for _, ch := range channels {
		select {
		case ch <- event:
			// Success
		default:
			// Non-blocking fallback if the subscriber's channel buffer is fully blocked.
			// Instead of locking up the system, we quietly log the drop directly into metrics.
			atomic.AddUint64(&b.droppedEvents, 1)
		}
	}
}

// DroppedCount returns the precise metric of all skipped events due to slow consumers.
func (b *EventBus) DroppedCount() uint64 {
	return atomic.LoadUint64(&b.droppedEvents)
}

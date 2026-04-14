# Phase 2: Internal Event Bus Implementation Research

## Key Findings

1. **Pub/Sub Mechanism in Go**
   - Given the strict constraints of the project avoiding heavy 3rd-party toolkits for core utilities unless necessary, we will build a native Go Pub/Sub mechanism.
   - The fundamental structure requires a central struct mapping topics (strings) to multiple subscriber channels (`map[string][]chan Event`).
   - Using a `sync.RWMutex` is required when adding/removing subscribers or iterating over them during publishing, as these operations will happen concurrently across multiple protocols (goroutines).

2. **Non-Blocking Publish Action**
   - The core defense against cascading system failures (where one slow Websocket client blocks the entire application) is a non-blocking `Publish` method.
   - Go's `select` statement paired with a `default` case handles this perfectly:
     ```go
     select {
     case ch <- ev:
         // Success
     default:
         // Channel buffer full -> Drop immediately
         atomic.AddUint64(&b.droppedEvents, 1)
     }
     ```
   - Incrementing the metric relies on `sync/atomic` for rapid concurrent metric increments without engaging the global mutex loop.

3. **Subscription Lifecycle Management**
   - The `Subscribe` function needs to issue a buffered channel: `make(chan Event, bufferSize)` to the caller.
   - We must also provide an `Unsubscribe(topic string, ch <-chan Event)` method to gracefully severe ties with disconnected clients (e.g., terminated WS connections), avoiding memory leaks in the `map`.

## Validation Architecture

1. **File System Verification:** Ensure `internal/eventbus/bus.go` exists.
2. **Behavioral Checks:**
   - Verify `Publish` executes without delay even if consumer deliberately sleeps.
   - Verify that `DroppedCount()` accurately increments tracking lost events.
   - Verify `Unsubscribe` reclaims array memory.
3. **Hell Case Testing:** Write a dedicated integration test that spawns 100 super-slow consumers (sleeping for 1 second each) while a single publisher fires 1,000,000 events instantly. The publisher MUST NOT block or fail, and `DroppedCount()` MUST reflect the exact dropped amount.

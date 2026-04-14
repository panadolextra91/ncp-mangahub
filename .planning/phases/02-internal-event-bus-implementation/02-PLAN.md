---
wave: 1
depends_on: []
files_modified:
  - internal/eventbus/bus.go
  - internal/eventbus/bus_test.go
  - pkg/models/event.go
autonomous: true
---

# Phase 2 Plan: Internal Event Bus Implementation

## Objective
Implement a thread-safe, non-blocking Pub/Sub Event Bus architecture across Go channels to securely decouple internal protocol communication, guaranteeing high availability even under extreme subscriber latency.

## Tasks

### [ ] Task 2.1: Define Event Model
<read_first>
- .planning/phases/02-internal-event-bus-implementation/02-CONTEXT.md
</read_first>
<action>
1. Create file `pkg/models/event.go`.
2. Define the pure global structure `type Event struct` containing two attributes exactly: `Topic string` and `Payload interface{}`.
</action>
<acceptance_criteria>
- `cat pkg/models/event.go` possesses `type Event struct`.
- `cat pkg/models/event.go` possesses `Topic string`.
</acceptance_criteria>

### [ ] Task 2.2: Implement Event Bus Core Service
<read_first>
- .planning/phases/02-internal-event-bus-implementation/02-CONTEXT.md
- pkg/models/event.go
</read_first>
<action>
1. Create `internal/eventbus/bus.go`.
2. Define the state-management struct: 
   ```go
   type EventBus struct {
       subscribers   map[string][]chan models.Event
       mu            sync.RWMutex
       bufferSize    int
       droppedEvents uint64
   }
   ```
3. Establish `func NewEventBus(bufferSize int) *EventBus` that prepares the map structures.
4. Implement `func (b *EventBus) Subscribe(topic string) <-chan models.Event`. It generates `make(chan models.Event, b.bufferSize)` and safely writes to the map using `b.mu.Lock()`.
5. Implement `func (b *EventBus) Publish(event models.Event)`. Ensure it utilizes an `RLock()` snapshot and executes non-blocking delivery precisely:
   ```go
   select {
   case ch <- event:
   default:
       atomic.AddUint64(&b.droppedEvents, 1)
   }
   ```
6. Implement `func (b *EventBus) Unsubscribe(topic string, ch <-chan models.Event)` to remove inactive listeners dynamically to prevent memory bleeding. It also closes the channel gracefully.
7. Integrate `func (b *EventBus) DroppedCount() uint64` returning the metric via `atomic.LoadUint64(&b.droppedEvents)`.
</action>
<acceptance_criteria>
- `cat internal/eventbus/bus.go` verifies `atomic.AddUint64(&b.droppedEvents, 1)` execution inside a standard `default:` case.
- `cat internal/eventbus/bus.go` utilizes `sync.RWMutex` locks.
</acceptance_criteria>

### [ ] Task 2.3: Verification via Hell Path
<read_first>
- internal/eventbus/bus.go
- .planning/phases/02-internal-event-bus-implementation/02-VALIDATION.md
</read_first>
<action>
1. Create the testing harness at `internal/eventbus/bus_test.go`.
2. Engineer `TestEventBus_HappyPath` for base routing checks.
3. Engineer `TestEventBus_Unsubscribe` preventing memory leaks when explicitly removed.
4. Forge the `TestEventBusHellPath_SlowConsumers` integration block:
   - Command a bus instantiation specifying `bufferSize=1`.
   - Hook a single subscriber and deliberately delay loop processing (`time.Sleep`).
   - Concurrently spam 100,000 requests to `Publish` using an atomic `sync.WaitGroup` deployment.
   - Assert definitively that `b.DroppedCount() > 0` directly confirming zero IO drops without latency anomalies.
5. `test.sh` currently executes `go test -v -coverprofile=coverage.out ./internal/...`. Add this module transparently.
</action>
<acceptance_criteria>
- `cat internal/eventbus/bus_test.go` exhibits `TestEventBusHellPath_SlowConsumers` successfully deploying the drop metrics assertions alongside a concurrent blast.
- The standard `./test.sh` succeeds with `80.0%`+ aggregate statement execution coverages.
</acceptance_criteria>

## Verification
- Invoke `chmod +x test.sh && ./test.sh`. Ensure total coverage transcends the 80% mark, specifically confirming `bus.go` is safely guarded.
- Manually review `internal/eventbus/bus.go` visually to assert zero logic logs natively executed inside the Select Default block. Validating standard IO purity.

## Must Haves
- Buffer sizing injected dynamically avoiding strict singletons.
- Dropped metric is intrinsically atomic.
- Select `default` case drops precisely without locking.

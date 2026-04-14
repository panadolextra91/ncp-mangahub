# Phase 2 Context: Internal Event Bus Implementation

## Domain Boundary
**What this encapsulates:** Building the core Pub/Sub Event Bus utilizing Go Channels. This acts as the isolated transport backbone connecting all 5 future protocols, enforcing non-blocking event broadcasting to prevent cascading system failures.

## Key Decisions

1. **Bus Instantiation Strategy**
   - **Decision:** Dependency Injection via Constructor (`NewEventBus()`).
   - **Details:** The system will initialize the EventBus instance strictly in `cmd/server/main.go` and explicitly pass it as a dependency into relevant service handlers, rather than leveraging a global singleton. This immensely increases isolated testability.

2. **Channel Buffer Capacity**
   - **Decision:** Dynamic Buffer Config.
   - **Details:** The `NewEventBus(bufferSize int)` function will take a dynamic buffer size upon creation. This allows the system to configure resilience profiles (e.g. high capacity for heavy workloads while retaining smaller limits for testing specific rapid-drop edge cases).

3. **Slow Consumer Drop Strategy**
   - **Decision:** Silent Event Drops with Atomic Metrics Counters.
   - **Details:** When an event is dropped due to full subscriber buffer (non-blocking select), the Bus MUST NOT log to console; doing so would bottleneck IO in Hell Cases. Instead, it must safely increment an internal, atomic/mutex-guarded counter (DropMetric) tracking the total number of discarded events for subsequent health check APIs.

## Canonical Refs
- `docs/plan.md` (Primary logical blueprint)

## Deferred Ideas
None.

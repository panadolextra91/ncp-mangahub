# Phase 9 Context: Graceful Shutdown Mechanics

## Key Decisions (Context)
- **Shutdown Strategy**: **Strict with 5s Grace Period**. Inform clients, then force close.
- **Concurrency Model**: **Context + sync.WaitGroup**.
  - `signal.NotifyContext` for the root context.
  - `WaitGroup` to track all background processes (Servers, Hubs, Bridges).
- **Shutdown Sequence (STRICT)**:
  1. **Reject New Work**: Stop Listeners (HTTP, TCP, UDP, gRPC).
  2. **Evict Active Clients**: Send shutdown message `{"type": "system", "action": "shutdown", "reason": "server_maintenance"}` to WS/TCP clients, then close hubs.
  3. **Stop Internal Chains**: Terminate EventBus bridge goroutines.
  4. **Wait for Completion**: `wg.Wait()`.
  5. **Final Cleanup**: `db.Close()`.
- **Infrastructure**: Wrap the final shutdown logic in a 10s `context.WithTimeout` to prevent indefinite hangs.
- **Client Handling**: Explicitly notify clients to prevent "Thundering Herd" reconnection loops.

## Gray Areas to Discuss
*None. All decisions finalized by Mẹ Architect.*

# Pitfalls Research

**Domain:** Go Multi-Protocol Monoliths & SQLite Eventing

## Critical Mistakes & Prevention
1. **SQLite `database is locked` Errors**
  - *Warning Sign:* Intermittent 500 errors on HTTP writes or test failures.
  - *Prevention:* Enforce `db.SetMaxOpenConns(1)` for writes strictly. Use `PRAGMA journal_mode=WAL`. Only HTTP handles writes.
  - *Phase:* Addressed in initial DB + HTTP phase.

2. **Goroutine Leaks**
  - *Warning Sign:* Memory usage climbing over time.
  - *Prevention:* Always use `defer conn.Close()`. Bind every protocol listener and event consumer to a parent `context.Context` and ensure `WaitGroup` is utilized on shutdown.

3. **Event Bus Deadlocks (Slow Consumers)**
  - *Warning Sign:* Entire system freezes when a single WebSocket or TCP client reads slowly.
  - *Prevention:* Buffered channels & select-default drops:
    ```go
    select {
    case ch <- event:
    default:
      // Drop event
    }
    ```

4. **Dirty Shutdowns**
  - *Warning Sign:* DB corruption or hanging processes.
  - *Prevention:* Signal interceptor -> cancel context -> stop accepting new HTTP/TCP -> wait for current requests to finish -> close DB.

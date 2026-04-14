# Phase 5 Validation Strategy

## Objective
To verify that the TCP Protocol Layer correctly implements a secured, real-time JSON broadcasting system using an idiomatic Channel Registry pattern and mandatory JWT handshake.

## Verification Checklist

### 1. Configuration Check
- [ ] `TCP_PORT` env variable is respected.
- [ ] Default `9090` fallback is functional.

### 2. JWT Refactor Integrity
- [ ] HTTP Middleware still correctly validates tokens using the new `pkg/auth/jwt` package.
- [ ] Auth login flow still generates valid JWTs.

### 3. TCP Protocol Compliance
- [ ] Server greets with nothing until `AUTH ` is sent.
- [ ] `AUTH <invalid>` results in `401 Unauthorized` and immediate disconnection.
- [ ] `AUTH <valid>` results in `200 OK CONNECTED`.
- [ ] Connected clients receive JSON payloads of `manga.new` and `progress.updated` immediately after HTTP trigger.

### 4. Concurrency & Stability (Hell Path)
- [ ] **Channel Registry Test**: Hub correctly registers and unregisters clients without blocking other active connections.
- [ ] **Slow Consumer Test**: A client that connects but does not read should not block the system or the broadcast loop.
- [ ] **Total Coverage**: Global package statement coverage remains >= 80%.

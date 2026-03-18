# OpenClaw SDK Go Migration - Implementation Plan

> **Master Plan** - Links to all phase implementation plans

**Goal:** Implement a feature-complete Go SDK equivalent to the TypeScript openclaw-sdk, following Go idioms with Context + Channel hybrid architecture.

**Architecture:** The SDK uses the Option pattern for configuration, goroutines + channels for event handling, and context.Context for timeout/cancellation. All managers implement graceful shutdown patterns with proper resource cleanup.

**Tech Stack:**
- Go 1.21+
- gorilla/websocket (WebSocket)
- Go standard library (context, sync, net/http, etc.)

---

## Phase Documents

| Phase | Description | Document |
|-------|-------------|----------|
| **Phase 1** | Project Setup and Foundation | [Phase 1: Setup & Foundation](./phase-1-setup-foundation.md) |
| **Phase 2** | Authentication Module | [Phase 2: Auth Module](./phase-2-auth-module.md) |
| **Phase 3** | Protocol Module | [Phase 3: Protocol Module](./phase-3-protocol-module.md) |
| **Phase 4** | Transport Module | [Phase 4: Transport Module](./phase-4-transport-module.md) |
| **Phase 5** | Connection Module | [Phase 5: Connection Module](./phase-5-connection-module.md) |
| **Phase 6** | Events Module | [Phase 6: Events Module](./phase-6-events-module.md) |
| **Phase 7** | Managers Module | [Phase 7: Managers Module](./phase-7-managers-module.md) |
| **Phase 8** | Utils Module | [Phase 8: Utils Module](./phase-8-utils-module.md) |
| **Phase 9** | Main Client | [Phase 9: Main Client](./phase-9-main-client.md) |
| **Phase 10** | Examples | [Phase 10: Examples](./phase-10-examples.md) |

---

## Execution Order

1. **Phase 1** - Project Setup (go.mod, types, errors, logger)
2. **Phase 2** - Auth (CredentialsProvider, AuthHandler)
3. **Phase 3** - Protocol (types, validation)
4. **Phase 4** - Transport (WebSocket)
5. **Phase 5** - Connection (state machine, negotiator, policies, TLS)
6. **Phase 6** - Events (tick, gap)
7. **Phase 7** - Managers (event, request, connection, reconnect)
8. **Phase 8** - Utils (timeout)
9. **Phase 9** - Main Client
10. **Phase 10** - Examples

---

## Notes

- Each phase is self-contained and can be implemented independently
- All phases follow TDD workflow (write test first, then implementation)
- Each phase should end with a working, compilable state
- Commit after completing each phase

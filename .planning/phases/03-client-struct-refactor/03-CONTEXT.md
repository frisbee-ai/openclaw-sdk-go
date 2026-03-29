# Phase 3: Client Struct Refactor - Context

**Gathered:** 2026-03-29
**Status:** Ready for planning

<domain>
## Phase Boundary

Phase 3 delivers **maintainable client struct organization and clear API semantics** for the OpenClaw SDK.

Two requirements:
- **API-01**: `client` struct refactor — group oversized struct into logical sub-structs (`core`, `protocol`, `health`, `api`)
- **API-02**: `Close()` vs `Disconnect()` disambiguation — merge or document semantics

</domain>

<decisions>
## Implementation Decisions

### API-01: Client Struct Sub-structs

- **D-01 (API namespace sub-struct):** Group 15 API namespace fields (`chatAPI`, `agentsAPI`, ..., `usageAPI`) into a nested `api struct { ... }`. Accessor methods on the `OpenClawClient` interface remain unchanged — they delegate to `c.api.chat`, etc. Consistent with existing `managers` sub-struct pattern. Backward compatible.

- **D-02 (Remaining fields grouping):** Final struct layout:
  ```go
  type client struct {
      config   *ClientConfig   // Core: configuration
      managers struct {        // Core: managers (already exists)
          event       *managers.EventManager
          request     *managers.RequestManager
          connection  *managers.ConnectionManager
          reconnect   *managers.ReconnectManager
      }
      api struct {            // API-01: 15 namespace fields
          chat          *api.ChatAPI
          agents        *api.AgentsAPI
          sessions      *api.SessionsAPI
          config        *api.ConfigAPI
          cron          *api.CronAPI
          nodes         *api.NodesAPI
          skills        *api.SkillsAPI
          devicePairing *api.DevicePairingAPI
          browser       *api.BrowserAPI
          channels      *api.ChannelsAPI
          push          *api.PushAPI
          execApprovals *api.ExecApprovalsAPI
          system        *api.SystemAPI
          secrets       *api.SecretsAPI
          usage         *api.UsageAPI
      }
      protocol struct {       // API-01: protocol negotiation
          negotiator  *connection.ProtocolNegotiator
          policy      *connection.PolicyManager
          serverInfo  *connection.HelloOk
          snapshot    *connection.Snapshot
      }
      health struct {         // API-01: connection health
          tickMonitor     *events.TickMonitor
          gapDetector     *events.GapDetector
          tickHandlerUnsub func()
      }
      // Core fields (flat, not grouped)
      requestFn api.RequestFn
      ctx       context.Context
      cancel    context.CancelFunc
      mu        sync.Mutex
  }
  ```

### API-02: Close vs Disconnect

- **D-03 (Semantics):** Keep both methods. Update `OpenClawClient` interface documentation:
  - `Close()` = "Shuts down the entire client and releases all resources. No further operations are valid after calling Close."
  - `Disconnect()` = "Disconnects from the server without shutting down the client. Stops reconnection attempts and cleans up connection state. Call Connect() to reconnect."
- Both methods exist in `pkg/client.go` — no code change needed, only documentation comments.

### Additional Decisions

- **Interface stability:** All existing accessor methods (GetServerInfo, GetSnapshot, GetPolicy, GetTickMonitor, GetGapDetector, all 15 API namespaces) remain on the interface. Sub-struct refactoring is an internal reorganization only.
- **Backward compatibility:** No breaking changes to the public API surface. All `client.api.*` accesses are internal to `pkg/client.go`.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Phase Decisions (must be consistent)
- `.planning/phases/01-foundation/CONTEXT.md` — Phase 1 patterns (option pattern, channel+mutex, sub-struct for managers)
- `.planning/phases/02-observability/02-CONTEXT.md` — Phase 2 patterns (GetMetrics, EventManager structure)

### Project Requirements
- `.planning/ROADMAP.md` §Phase 3 — Phase 3 goal, success criteria (API-01, API-02)
- `.planning/REQUIREMENTS.md` §API-01, API-02 — Acceptance criteria for each requirement

### Codebase Architecture
- `pkg/client.go` — Main client: `OpenClawClient` interface (line 347), `client` struct (line 396), `NewClient()` (line 438), `Disconnect()` (line 590), `Close()` (line 975)

### No external specs
No external specs — requirements fully captured in decisions above.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- **`managers` sub-struct pattern** (`pkg/client.go:398-403`): Already working example of sub-struct grouping — use as template for `api`, `protocol`, `health` sub-structs
- **Functional options for initialization**: `NewClient()` uses `ClientOption` pattern — sub-structs should be initialized inline in `NewClient()`

### Established Patterns
- **Nested sub-struct with embedded types**: Go permits `sync.Mutex` in nested structs. No reflection/encoding issues.
- **Accessor delegation**: Interface methods already exist delegating to fields — `Chat() *api.ChatAPI` stays unchanged; implementation changes to `return c.api.chat`

### Integration Points
- **NewClient initialization** (`pkg/client.go:447-510`): All 15 API namespaces initialized here — must update to `c.api = apiNamespaceStruct{...}` pattern
- **GetMetrics()** (`pkg/client.go:820-859`): References `c.tickMonitor`, `c.managers.reconnect` — tickMonitor becomes `c.health.tickMonitor`
- **Disconnect()** (`pkg/client.go:590`): References `c.tickMonitor`, `c.gapDetector`, `c.protocolNegotiator`, `c.serverInfo`, `c.snapshot`, `c.tickHandlerUnsub` — all move to new sub-structs
- **Interface accessors** (`pkg/client.go:347-391`): 24 interface methods that return internal state — unchanged

</code_context>

<specifics>
## Specific Ideas

- **Interface docs update**: `Disconnect()` comment should say "disconnects without shutting down client; call Connect() to reconnect" (not "closes the connection gracefully" which is vague)
- **Struct comment cleanup**: Remove stale comment "// New fields for Phase 6.1" from line 404
- **TickMonitor references after refactor**: `GetMetrics()`, `Disconnect()`, `Connect()` all reference `c.tickMonitor` → becomes `c.health.tickMonitor`

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 03-client-struct-refactor*
*Context gathered: 2026-03-29*

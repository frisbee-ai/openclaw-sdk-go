# Phase 2: Observability - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-28
**Phase:** 02-observability
**Areas discussed:** ConnectionMetrics (OBS-01), Per-request timeout (OBS-02), Event priority levels (OBS-03)

---

## OBS-01: ConnectionMetrics

### Area: Latency measurement

| Option | Description | Selected |
|--------|-------------|----------|
| Tick-based 估算（推荐） | 用两次 tick 间隔来估算，不额外打点。粗略但零成本，足够判断连接健康状态 | ✓ |
| 端到端 RTT 测量 | 通过心跳 ping/pong 测真实 RTT。更精确但需要改动 transport 层的 ping/pong 逻辑 | |
| 不暴露 Latency 字段 | Metrics 只含 LastTickAge 和 ReconnectCount；Latency 留到 Phase 4 benchmark 再测 | |

**User's choice:** Tick-based 估算（推荐）
**Notes:** Zero-cost, good enough for v1.0 observability. Real RTT via ping/pong deferred to Phase 4.

### Area: Metrics API surface

| Option | Description | Selected |
|--------|-------------|----------|
| 加到接口（推荐） | 用户直接 client.GetMetrics()，最方便。ConnectionMetrics struct 公开但版本稳定即可 | ✓ |
| 不加入接口，内部方法 | GetMetrics() 在 client struct 上但不在接口上。需要 Metrics 的用户做类型断言 | |

**User's choice:** 加到接口（推荐）
**Notes:** All users should have equal access. Simpler API surface.

### Additional Decisions (not in gray area questions):

| Decision | Value |
|----------|-------|
| Metrics struct fields | Latency, LastTickAge, ReconnectCount, IsStale |
| ReconnectCount source | New AttemptCount() method on ReconnectManager |

---

## OBS-02: Per-request timeout

### Area: SendRequest signature

| Option | Description | Selected |
|--------|-------------|----------|
| variadic options（推荐） | SendRequest(ctx, req, WithRequestTimeout(5*time.Second)) — 向后兼容，不传则用 context deadline | ✓ |
| 新增 SendRequestWithOptions | 保留原 SendRequest 不变，新增 SendRequestWithOptions(ctx, req, opts)，不影响现有调用方 | |

**User's choice:** variadic options（推荐）
**Notes:** Backward compatible. No new method name. Existing callers unchanged.

### Additional Decisions:

| Decision | Value |
|----------|-------|
| RequestOption design | Functional option type: `type RequestOption func(*requestConfig)` |
| Timeout precedence | WithRequestTimeout wraps incoming ctx with WithTimeout — overwrites existing deadline (explicit wins) |
| Placement | Options parsed in client.SendRequest, wrapped context passed to RequestManager.SendRequest |

---

## OBS-03: Event priority levels

### Area: Priority level count

| Option | Description | Selected |
|--------|-------------|----------|
| 3 级：高/中/低（推荐） | HIGH=错误/断开/状态变化，中=tick/connect，低=message/request。简单实用 | ✓ |
| 2 级：高/低 | 只有关键和非关键两类，最小化决策 | |
| 5 级（类似 log level） | DEBUG/INFO/WARN/ERROR/CRITICAL。细粒度但 Phase 2 过度设计 | |

**User's choice:** 3 级：高/中/低（推荐）
**Notes:** Balanced. Enough granularity without over-engineering.

### Area: Emit API design

| Option | Description | Selected |
|--------|-------------|----------|
| Event 里加 Priority 字段（推荐） | Event{Priority EventPriority, ...}。已有 Emit() 保持兼容，低优先级事件默认 MEDIUM | ✓ |
| 新增 EmitWithPriority 方法 | 保留 Emit() 不变，新增 EmitWithPriority(e Event, p EventPriority)。更显式但 API 碎片化 | |

**User's choice:** Event 里加 Priority 字段（推荐）
**Notes:** Backward compatible. No new Emit method needed.

### Area: Priority default assignment

| Option | Description | Selected |
|--------|-------------|----------|
| Tick + Response = MEDIUM（推荐） | HIGH=Error/Disconnect/StateChange/Gap；MEDIUM=Tick/Response；LOW=Message/Request | ✓ |
| 全 MEDIUM 仅 Error/Disconnect HIGH | HIGH=Error/Disconnect；MEDIUM=其他所有（Tick/Connect/Response/StateChange/Gap）；LOW=Message/Request | |
| 全 HIGH（保守） | 除 EventMessage 外全 HIGH，EventMessage=LOW。关键事件永不丢，可能 buffer 堆积 | |

**User's choice:** Tick + Response = MEDIUM（推荐）
**Notes:** Tick and Response are important for observability but not as critical as errors. Connect at MEDIUM seems fine.

---

## OBS-04: EventBufferSize

**Status:** Already implemented. `ClientConfig.EventBufferSize` + `WithEventBufferSize()` exist in `pkg/client.go`. No discussion needed.

---

## Claude's Discretion

Delegated to researcher/planner:
- Internal channel structure for priority-based event dispatch (separate per-priority channels vs. single channel with priority tagging)
- How ReconnectManager.AttemptCount() is implemented (atomic vs. mutex-protected)
- Whether GetMetrics() returns a copy or direct reference
- Whether TickMonitor needs a GetTickInterval() method

---

## Deferred Ideas

- **Real RTT via ping/pong**: Deferred to Phase 4 (Benchmarking)
- **5-level priority**: Deferred to future if real use cases emerge
- **Prometheus/OpenTelemetry export**: Deferred to Phase 5

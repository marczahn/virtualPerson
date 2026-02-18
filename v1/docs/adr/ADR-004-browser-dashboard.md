# ADR-004: Browser Dashboard — Transport Multiplexing and Build-Free Frontend

**Status:** Accepted
**Date:** 2026-02-18

---

## Context

The simulation runs in `--server` mode and exposes a WebSocket endpoint (`/ws`) consumed by a Bubbletea TUI client. The TUI is cumbersome for observing the simulation over time: it is a flat text stream with no visual differentiation between biological, psychological, and conscious events.

A browser dashboard is needed that shows live biology (all 20 variables) and psychology (affect dimensions, Big Five, coping/distortions) panels with time-series charts and appropriate visual grouping.

Two protocol questions arose:
1. Should the browser use the existing `/ws` endpoint or a separate one?
2. Should the frontend use a framework (React, Vue) with a build pipeline, or vanilla JS?

One layering question arose:
- The simulation loop (`simulation` package) must not import the `server` package to avoid a circular dependency. How does it trigger WebSocket broadcasts?

---

## Decisions

### 1. Multiplex on the existing `/ws` endpoint

New message types (`"bio_state"`, `"psych_state"`) are added to the existing `ServerMessage` type using optional pointer fields with `omitempty`. The `type` field already serves as a discriminant.

**Rejected alternative:** A separate `/ws/state` endpoint for bio/psych data.
- Rejected because: the browser would need two simultaneous WebSocket connections; the server would need to manage two separate hub instances; the protocol is already type-discriminated and scales to additional types without new endpoints.

**Consequence:** The TUI client receives `bio_state` and `psych_state` messages it did not previously see. A type guard was added to `client/model.go:addThought()` to filter out non-`"thought"` messages. Future message types need the same treatment.

### 2. `go:embed` + vanilla JS, no build pipeline

Static web files (`index.html`, `dashboard.js`, `dashboard.css`) are embedded in the binary via `go:embed` and served at `/` by the existing HTTP server. Chart.js is loaded from CDN. No npm, no bundler, no build step.

**Rejected alternative:** A React/Vue SPA with a build pipeline.
- Rejected because: adds significant toolchain complexity (npm, bundler, transpiler), a separate build step before every deploy, and ongoing dependency maintenance — none of which is justified for a single-developer simulation monitoring tool.

**Consequence:** Frontend changes require recompiling the Go binary. Iteration on CSS/JS requires a rebuild (no hot-reload). This is an acceptable trade-off given the development context.

### 3. Callback injection to avoid circular import

The `simulation.Loop` calls bio/psych state broadcasts via an `OnStateSnapshot func(*biology.State, psychology.State)` field on `simulation.Config`. This callback is wired in `main.go` to call `hub.Broadcast(...)`. The `simulation` package never imports `server`.

**Consequence:** This follows the same pattern as `display.SetListener()`, which is already used for thought broadcasts. The pattern should be continued for any future server-oriented callbacks from the loop.

### 4. Bio/psych snapshots throttled to ~2 seconds

The simulation ticks every 100ms. Broadcasting full state on every tick would be wasteful (20 variables × 10/s). The loop rate-limits snapshots via a configurable `StateSnapshotInterval` (default 2s, overridable for tests).

---

## Known Limitations

- **Browser reconnection loses rolling chart history.** The 60-point history buffer lives in the browser. On reconnect, charts start from empty. Acceptable for a monitoring tool; a future enhancement could replay the last N snapshots from the SQLite store on connect.
- **Thresholds are re-evaluated in the payload constructor.** `BioStatePayloadFromState` calls `biology.EvaluateThresholds()` separately from the loop's own threshold evaluation. This is a negligible double-computation (a handful of comparisons) but means snapshot thresholds are evaluated at broadcast time, not tick time.
- **Personality is re-sent every 2 seconds.** Big Five traits are stable (fixed for the person's lifetime) but are included in every `psych_state` message for simplicity. A connection handshake/init message would be more correct but adds protocol complexity.

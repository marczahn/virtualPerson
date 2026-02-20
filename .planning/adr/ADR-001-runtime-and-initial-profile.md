# ADR-001: Runtime Entrypoint And Initial Profile Contract

## Status
Accepted (2026-02-20)

## Context
`v2` currently has strong domain-level contracts and tests but no complete runtime path to execute the virtual person end-to-end. This blocks the core objective: the person must act virtually, continuously, and with visible motivation-driven behavior. In addition, startup conditions are underspecified. The concept baseline requires explicit initial situation data and person characterization, including stress sensitivity that shapes cortisol/stress response.

## Decision
1. Add an explicit runtime entrypoint requirement (`INF-08`) in Phase 5.
2. Add explicit initial profile requirements (`PRF-01..PRF-03`) in Phase 6.
3. Treat startup profile as a typed configuration contract (validated upfront), not implicit defaults spread across packages.

## Consequences
- Positive: v2 becomes runnable as a coherent simulation, not only as isolated packages.
- Positive: startup behavior is reproducible and tunable.
- Positive: stress-reactivity characterization is first-class and testable.
- Tradeoff: adds one additional runtime plan slice and one profile/config slice before v2 planning is complete.

## Alternatives Considered
- Keep runtime wiring implicit in ad-hoc scripts: rejected (unreliable, poor traceability).
- Keep startup profile as hardcoded defaults in `main`: rejected (non-testable, weak configurability, hidden assumptions).

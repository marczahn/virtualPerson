---
phase: 05-simulation-loop
plan: "03"
subsystem: infrastructure
tags: [go, infrastructure, scenario, input-drainer, inf-06, tdd]

requires:
  - phase: 05-02
    provides: "Deterministic external input drain adapter and environment effect mapping"
provides:
  - "Runtime scenario registration and activation contract"
  - "Deterministic latest-wins scenario switch semantics"
  - "Scenario descriptor injection into TickInput pre-bio effects via existing environment mapper"
affects: [05-04]

tech-stack:
  added: []
  patterns:
    - "Infrastructure wrapper around InputDrainer to preserve inward dependencies"
    - "Scenario-to-environment-descriptor mapping reuses existing effect rules"
    - "Strict tests-first for registration, switching, and next-drain application"

key-files:
  created:
    - ".planning/phases/05-simulation-loop/05-03-SUMMARY.md"
    - "v2/internal/infrastructure/scenario.go"
    - "v2/internal/infrastructure/scenario_test.go"
  modified:
    - ".planning/REQUIREMENTS.md"
    - ".planning/STATE.md"
    - ".planning/CONTINUITY.md"
    - ".planning/NEXT_SESSION_START.md"
    - ".planning/ROADMAP.md"

key-decisions:
  - "Kept scenario logic fully in infrastructure; biology only receives BioRate/BioPulse contracts."
  - "Implemented scenario injection as an InputDrainer wrapper so orchestration and tick order remain unchanged."
  - "Applied scenario effects by feeding descriptors through existing environment mapping (`applyEnvironmentInput`) to avoid duplicated rules."
  - "Exposed runtime switching with deterministic latest activation winning the next drain cycle."

requirements-completed: [INF-06]

completed: 2026-02-20
---

# Phase 5 Plan 03: Summary

Implemented `INF-06` scenario injection with explicit biological effects using strict TDD.

## Red phase (tests first)

Added failing tests for:

- scenario registration and successful activation
- deterministic scenario switching where latest activation wins
- scenario effects appearing on the next drain cycle after runtime activation

## Green phase implementation

Created `v2/internal/infrastructure/scenario.go`:

- added `ScenarioInjector` as an `InputDrainer` wrapper
- added scenario registration and activation APIs
- injected active scenario descriptors into drained `TickInput`
- reused existing environment mapping to produce pre-bio rates/action gates
- appended active scenario context to `ExternalText`

## Verification

- `cd v2 && GOCACHE=/tmp/go-build go test ./internal/infrastructure -count=1` -> pass
- `cd v2 && GOCACHE=/tmp/go-build go test ./... -count=1` -> pass

## Architecture check

- Layer impact: infrastructure only.
- Boundary preserved: domain remains infrastructure-agnostic.
- No new dependencies added.

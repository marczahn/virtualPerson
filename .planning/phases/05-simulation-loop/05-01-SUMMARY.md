---
phase: 05-simulation-loop
plan: "01"
subsystem: infrastructure+biology+motivation+consciousness
tags: [go, infrastructure, simulation-loop, orchestration, inf-07, tdd]

requires:
  - phase: 04-03
    provides: "End-of-tick feedback buffer/apply contracts and pulse/rate split"
provides:
  - "Sequential loop contract: drain input -> biology -> motivation -> consciousness -> feedback"
  - "Infrastructure orchestration interfaces for input, biology, motivation, and mind stages"
  - "Tick-level state handoff contract with parsed response/cooldown persistence"
affects: [05-02]

tech-stack:
  added: []
  patterns:
    - "Single orchestrator entrypoint for deterministic one-tick progression"
    - "Tests-first stage-order contract verification"
    - "Inward dependencies only from infrastructure to domain packages"

key-files:
  created:
    - ".planning/phases/05-simulation-loop/05-01-PLAN.md"
    - ".planning/phases/05-simulation-loop/05-01-SUMMARY.md"
    - "v2/internal/infrastructure/simulation_loop.go"
    - "v2/internal/infrastructure/simulation_loop_test.go"
  modified:
    - ".planning/REQUIREMENTS.md"
    - ".planning/STATE.md"
    - ".planning/CONTINUITY.md"
    - ".planning/NEXT_SESSION_START.md"
    - ".planning/ROADMAP.md"

key-decisions:
  - "Defined explicit infrastructure interfaces (`InputDrainer`, `BioEngine`, `MotivationComputer`, `MindResponder`) to prevent cross-layer leakage."
  - "Applied drained pre-bio effects before biology tick so injected input can influence same-tick motivation computation."
  - "Used `TickFeedbackBuffer` at orchestration layer so consciousness feedback mutates biology only at tick end."
  - "Persisted prior parsed output and cooldown state in `SimulationState` as loop-owned cross-tick state."

requirements-completed: [INF-07]

completed: 2026-02-20
---

# Phase 5 Plan 01: Summary

Implemented the initial simulation-loop integration contract (`INF-07`) in `v2/internal/infrastructure` with tests-first workflow.

## Red phase (tests first)

Added failing contract tests for:

- strict stage order and exactly-once stage invocation
- drained input effects applied before motivation compute input is captured
- end-of-tick-only feedback mutation (no pre-consciousness mutation)

## Green phase implementation

Created `v2/internal/infrastructure/simulation_loop.go`:

- added orchestration interfaces and dependency wiring (`SimulationLoopDeps`)
- added tick input/request/state/result contracts
- implemented sequential `Tick(state, dt)` orchestration:
  - drain input
  - apply pre-bio input effects
  - run biology tick
  - compute motivation
  - build consciousness prompt + mind response
  - parse response + resolve cooldown-aware action outcome
  - accumulate pulses and apply feedback at tick end
  - persist prior parse/cooldown/continuity state

## Verification

- `cd v2 && GOCACHE=/tmp/go-build go test ./... -count=1` -> pass

## Architecture check

- Layer impact: infrastructure orchestrates Biology, Motivation, Consciousness contracts.
- Boundary preserved: no reverse dependency from domain layers to infrastructure.
- No new dependencies added.

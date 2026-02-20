---
phase: 05-simulation-loop
plan: "02"
subsystem: infrastructure+sense
tags: [go, infrastructure, sense, input-adapter, inf-05, tdd]

requires:
  - phase: 05-01
    provides: "Simulation loop orchestration contract and TickInput shape"
provides:
  - "Deterministic external input drain adapter for speech/action/environment conventions"
  - "Sense parser contract for v1-compatible input classification"
  - "TickInput bridge with pre-bio pulses/rates, action gates, and external text aggregation"
affects: [05-03]

tech-stack:
  added: []
  patterns:
    - "Thread-safe queued input drain once per tick"
    - "Convention parser kept pure and deterministic"
    - "Infrastructure-only mapping from parsed input to bio/action gate effects"

key-files:
  created:
    - ".planning/phases/05-simulation-loop/05-02-SUMMARY.md"
    - "v2/internal/sense/parser.go"
    - "v2/internal/sense/parser_test.go"
    - "v2/internal/infrastructure/input_adapter.go"
    - "v2/internal/infrastructure/input_adapter_test.go"
  modified:
    - ".planning/REQUIREMENTS.md"
    - ".planning/STATE.md"
    - ".planning/CONTINUITY.md"
    - ".planning/NEXT_SESSION_START.md"
    - ".planning/ROADMAP.md"

key-decisions:
  - "Preserved v1 input conventions exactly for this slice: plain=speech, `*...*`=action, `~...`=environment."
  - "Kept parser in `sense` side-effect free; effect mapping stays in infrastructure adapter."
  - "Initialized default allowed-action gates to permissive values and let ordered environment inputs override deterministically in one drain cycle."
  - "Aggregated all parsed raw inputs into newline-delimited `ExternalText` so one tick receives complete operator context."

requirements-completed: [INF-05]

completed: 2026-02-20
---

# Phase 5 Plan 02: Summary

Implemented `INF-05` external input handling with strict tests-first workflow.

## Red phase (tests first)

Added failing tests for:

- input convention parsing for speech/action/environment (`sense` parser)
- no-input drain no-op behavior
- deterministic multi-input drain ordering and gate override behavior
- speech/action/environment bridge into `TickInput` effects and text

## Green phase implementation

Created `v2/internal/sense/parser.go`:

- pure parser contract returning typed inputs
- v1-compatible conventions:
  - plain text -> speech
  - `*text*` -> action
  - `~text` -> environment

Created `v2/internal/infrastructure/input_adapter.go`:

- thread-safe input queue (`Enqueue` + one-shot `Drain`)
- deterministic drain order preservation
- bridge from parsed inputs to `TickInput`:
  - `ExternalText` aggregation
  - `AllowedActions` gate map (with default allowed actions)
  - pre-bio pulses from action keywords
  - pre-bio rates and action gate changes from environment keywords

## Verification

- `cd v2 && GOCACHE=/tmp/go-build go test ./... -count=1` -> pass

## Architecture check

- Layer impact: infrastructure + sense.
- Boundary preserved: no domain package depends on infrastructure.
- No new dependencies added.

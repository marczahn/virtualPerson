---
phase: 01-biological-foundation
plan: "06"
subsystem: bio
tags: [go, bio, engine, integration-tests, tdd]

# Dependency graph
requires:
  - phase: 01-04
    provides: "Baseline unit test safety net"
  - phase: 01-05
    provides: "Engine tick pipeline + seeded constructor"
provides:
  - "Scenario-based engine integration tests for long-run degradation"
  - "Fast-mode parity test proving accelerated degradation"
  - "In-range invariant checks under extended stress simulation"
  - "dt=0 no-change invariant for numeric bio variables"
  - "Threshold detection during near-collapse degradation"
affects: [02-01]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Deterministic simulation tests via NewEngineWithSeed"
    - "Long-horizon tick loops for phase acceptance criteria"

key-files:
  created:
    - ".planning/phases/01-biological-foundation/01-06-SUMMARY.md"
  modified:
    - "v2/bio/engine_test.go"

key-decisions:
  - "Interpreted slow-path requirement as full-engine net degradation (not isolated ApplyDecay path)"
  - "Implemented all 5 requested scenario tests plus explicit 60-tick slow-path divergence check"
  - "Used strict equality for dt=0 numeric invariants"

requirements-completed: [BIO-03, BIO-04]

# Metrics
duration: current-session
completed: 2026-02-19
---

# Phase 1 Plan 06: Summary

Implemented scenario integration tests for `Engine.Tick` in `v2/bio/engine_test.go` and validated phase degradation criteria.

## Tests added

- `TestEngine_TenMinutesUnattended`
- `TestEngine_FastMode_TwoMinutes`
- `TestEngine_AllVariablesInRange`
- `TestEngine_DtZero_NoDecay`
- `TestEngine_ThresholdsDetectedDuringDegradation`
- `TestEngine_SlowMode_SixtyTicksDiffersFromStart` (explicit slow-path truth check)

## Verification

- `cd v2 && GOCACHE=/tmp/go-build go test ./bio/... -run TestEngine -v -count=1` -> pass
- `cd v2 && GOCACHE=/tmp/go-build go test ./bio/... -v -count=1` -> pass
- `cd v2 && GOCACHE=/tmp/go-build go test ./bio/... -run TestEngine_TenMinutesUnattended -v -count=1` -> pass
- `cd v2 && GOCACHE=/tmp/go-build go test ./bio/... -race` -> pass

## Observed scenario outputs

- 10-minute unattended: `energy=0.4369`, `hunger=0.5617`, `mood=0.2764`
- Fast mode 2-minute: `energy=0.4209`, `hunger=0.6100`

These satisfy the required degradation thresholds while keeping all variables in their defined ranges.

---
phase: 01-biological-foundation
plan: "02"
subsystem: bio
tags: [go, bio, decay, interactions, v2]

# Dependency graph
requires:
  - phase: 01-biological-foundation/01-01
    provides: "bio.State struct with 8 variables, clampAll helper, VarRange ranges"
provides:
  - "ApplyDecay(s *State, cfg DecayConfig, dt float64) — linear decay for 5 autonomous variables"
  - "DecayConfig struct with DecayMultiplier and HomeostasisEnabled fields"
  - "DefaultDecayConfig() returning 5x speed for development"
  - "ApplyInteractions(s *State, dt float64) []Delta — 22 snapshot-based interaction rules"
  - "Delta struct {Field string; Amount float64} for rule deltas and logging"
  - "Rule struct {Name string; Condition func; Apply func} for data-driven rules"
affects:
  - "02-motivation-layer"
  - "01-biological-foundation (plans 03+)"

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Pre-tick snapshot pattern in ApplyInteractions: snap := *s before condition evaluation"
    - "Unexported decay rate constants at package level for test calibration access"
    - "dt cap (60s) in ApplyDecay to prevent pause-recovery explosions"
    - "applyDelta field-name dispatch (switch string) keeps rule Apply funcs pure"

key-files:
  created:
    - "v2/bio/decay.go"
    - "v2/bio/interactions.go"
  modified: []

key-decisions:
  - "dt capped at 60s in ApplyDecay — fast decay (5x multiplier) makes >60s pause produce unrealistic state jumps"
  - "Snapshot pattern (snap := *s) evaluates all 22 conditions from pre-tick state — rules cannot cascade within a single tick"
  - "CognitiveCapacity inversion: rules fire at <0.2 / <0.3 (not >0.8 / >0.7) because V2 models available capacity (high=good) not load"
  - "Decay rate constants are package-level const, not buried in function body — required for test calibration assertions"

patterns-established:
  - "Caller calls clampAll after both ApplyDecay AND ApplyInteractions — neither function clamps internally"
  - "ApplyInteractions returns []Delta for logging/debugging; caller may ignore the slice"
  - "Rule.Apply reads from snap (pre-tick), not from live state — no intra-tick feedback"

requirements-completed: [BIO-03, BIO-04]

# Metrics
duration: 2min
completed: 2026-02-19
---

# Phase 1 Plan 02: Bio Decay and Interactions Summary

**Linear decay for 5 autonomous variables (dt-capped) plus 22 snapshot-evaluated interaction rules generating compound bio spirals**

## Performance

- **Duration:** ~2 min
- **Started:** 2026-02-19T09:18:57Z
- **Completed:** 2026-02-19T09:20:25Z
- **Tasks:** 2
- **Files modified:** 2

## Accomplishments

- Implemented `ApplyDecay` with 5-variable autonomous degradation (Energy, Hunger, CognitiveCapacity, Mood, SocialDeficit) and dt-cap at 60s — Stress, PhysicalTension, BodyTemp intentionally excluded
- Implemented 22 motivation-relevant interaction rules in `motivationRules` slice with snapshot evaluation pattern preventing intra-tick feedback explosions
- Established `Delta` + `applyDelta` for typed, loggable bio state changes; rules return named deltas enabling future debug output per rule

## Task Commits

Each task was committed atomically:

1. **Task 1: Implement linear decay** - `d7591a6` (feat)
2. **Task 2: Implement data-driven interaction rules** - `a0daf3b` (feat)

## Files Created/Modified

- `v2/bio/decay.go` - DecayConfig, DefaultDecayConfig, ApplyDecay; 5 const decay rates
- `v2/bio/interactions.go` - Delta, Rule, 22-entry motivationRules, ApplyInteractions, applyDelta

## Decisions Made

- dt cap set to 60s (not the V1 value of 300s) — with 5x default multiplier, a 300s dt would produce 1500 equivalent seconds of decay, producing unrealistic step-change state jumps after any pause
- CognitiveCapacity inversion maintained throughout: rules 11 and 12 use `< 0.2` and `< 0.3` thresholds (low capacity = high load), consistent with the field naming decision from Plan 01

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- `v2/bio/decay.go` and `v2/bio/interactions.go` are the two autonomous pressure generators
- Ready for Plan 03 (thresholds) and Plan 04 (noise) — both will compose with decay and interactions in a tick function
- No blockers
- Caller pattern established: `ApplyDecay(s, cfg, dt)` → `ApplyInteractions(s, dt)` → `clampAll(s)`

## Self-Check: PASSED

- v2/bio/decay.go: FOUND
- v2/bio/interactions.go: FOUND
- 01-02-SUMMARY.md: FOUND
- Commit d7591a6 (Task 1): FOUND
- Commit a0daf3b (Task 2): FOUND
- `go build ./bio/...` exits 0: VERIFIED
- `go vet ./bio/...` exits 0: VERIFIED
- ApplyDecay touches exactly 5 fields: VERIFIED
- Stress/PhysicalTension/BodyTemp absent from ApplyDecay body: VERIFIED (comments only)
- 22 Name fields in motivationRules: VERIFIED
- `snap := *s` snapshot pattern in ApplyInteractions: VERIFIED
- dt cap at 60.0 in ApplyDecay: VERIFIED

---
*Phase: 01-biological-foundation*
*Completed: 2026-02-19*

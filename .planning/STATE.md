# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-02-18)

**Core value:** The person must exhibit intrinsic motivation — visible drives, desires, and frustrations that emerge from the interplay between biological state, a reward/motivation system, and LLM-driven consciousness.
**Current focus:** Phase 1: Biological Foundation

## Current Position

Phase: 1 of 6 (Biological Foundation)
Plan: 2 of TBD in current phase
Status: In progress
Last activity: 2026-02-19 — Completed 01-02-PLAN.md (bio decay + interaction rules)

Progress: [░░░░░░░░░░] ~4%

## Performance Metrics

**Velocity:**
- Total plans completed: 2
- Average duration: 2 min
- Total execution time: 0.07 hours

**By Phase:**

| Phase | Plans | Total | Avg/Plan |
|-------|-------|-------|----------|
| 1. Biological Foundation | 2 | 4 min | 2 min |

**Recent Trend:**
- Last 5 plans: 2 min, 2 min
- Trend: Stable at 2 min/plan

*Updated after each plan completion*

## Accumulated Context

### Decisions

Decisions are logged in PROJECT.md Key Decisions table.
Recent decisions affecting current work:

- Clean rebuild in v2/ directory — no code sharing with v1
- 7 personality factors replace Big Five (stress sensitivity, energy resilience, curiosity, self-observation, patience/frustration tolerance, risk aversion, social factor)
- Bio model reduced to 8-10 motivationally-relevant variables
- Hybrid motivation: code computes drives, LLM interprets under pressure
- Bio degradation slow-path must exceed homeostasis recovery rate when needs go unmet
- [Phase 01]: Go 1.22 used instead of 1.24 — toolchain go1.22.2 is installed, 1.24 not available on execution machine
- [Phase 01]: BodyTemp range {25,43} — V1's {34,42} was too narrow, clamped hypothermia reversal temps; wider range required for physiologically meaningful thresholds
- [Phase 01-02]: dt capped at 60s in ApplyDecay — 5x multiplier makes larger pauses produce unrealistic state jumps
- [Phase 01-02]: Snapshot pattern (snap := *s) in ApplyInteractions — all 22 conditions evaluated from pre-tick state, no intra-tick cascades

### Pending Todos

None yet.

### Blockers/Concerns

- [Phase 3 risk]: Phenomenological translation language is the highest-risk surface — exact phrasing that shifts LLM behavioral register at high drive intensities must be validated empirically. Budget iteration time.
- [Phase 4 risk]: dt=0 edge cases and feedback loop coupling errors are silent — plan explicit integration tests at dt=0 and dt=large on all feedback paths.

## Session Continuity

Last session: 2026-02-19
Stopped at: Completed 01-02-PLAN.md (bio decay + interaction rules)
Resume file: .planning/phases/01-biological-foundation/01-03-PLAN.md (next plan)

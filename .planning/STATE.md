# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-02-19)

**Core value:** The person must exhibit intrinsic motivation — visible drives, desires, and frustrations that emerge from the interplay between biological state, a reward/motivation system, and LLM-driven consciousness.
**Current focus:** Phase 5 prep after finishing Phase 4 feedback-loop contracts

## Current Position

Phase: 5 of 6 (Simulation Loop Integration) — ready
Plan: Phase 4 complete (`04-01`, `04-02`, `04-03`)
Status: Ready for next phase planning
Last activity: 2026-02-20 — Completed `04-03` (end-of-tick feedback application + BioRate/BioPulse contract split)

Progress: [████████░░] ~67%

## Performance Metrics

**Velocity:**
- Total plans completed: 14
- Phase 1 plans completed: 6/6
- Phase 2 plans completed: 2

**By Phase:**

| Phase | Plans | Total | Avg/Plan |
|-------|-------|-------|----------|
| 1. Biological Foundation | 6 | 11 min | 2 min |
| 2. Motivation | 2 | current-session | current-session |
| 3. Consciousness | 3 | current-session | current-session |
| 4. Feedback Loop | 3 | current-session | current-session |

*Updated after each plan completion*

## Accumulated Context

### Decisions

Decisions are logged in PROJECT.md Key Decisions table.
Recent decisions affecting current work:

- Clean rebuild in v2/ directory — no code sharing with v1
- 7 personality factors replace Big Five
- Bio model reduced to 8-10 motivationally-relevant variables
- Hybrid motivation: code computes drives, LLM interprets under pressure
- [Phase 01-02]: dt capped at 60s in ApplyDecay
- [Phase 01-02]: Snapshot pattern in ApplyInteractions
- [Phase 01-03]: Most-severe-only thresholds via switch/case
- [Phase 01-05]: Engine pipeline fixed to decay -> interactions -> noise -> clamp -> thresholds/cascades -> clamp
- [Phase 01-05]: Noise uses injected RNG with deterministic seeded constructor
- [Phase 01-05]: Module target updated to Go 1.26
- [Phase 01-06]: Scenario integration tests confirm unattended degradation, fast-mode equivalence, in-range safety, dt=0 invariants, and threshold detection under degradation
- [Phase 02-01]: Biology package migrated from `v2/bio` to `v2/internal/biology` with package name `biology`; `v2/bio` retired
- [Phase 02-02]: Motivation core implemented in `v2/internal/motivation` with deterministic `Compute`, tie-break goal selection, and pure action candidates
- [Phase 03-01]: Added `v2/internal/consciousness` contracts for phenomenological prompt context and defensive parser fallback for `[STATE]`, `[ACTION]`, optional `[DRIVE]` tags
- [Phase 03-02]: Added consciousness action outcome gating contract and drive override application with floor clamp `max(raw*0.5, reported)` and deterministic goal recomputation
- [Phase 03-03]: Added spontaneous thought ticker contracts with urgency-dominant selection, baseline associative drift path, and bounded continuity buffer prompt inclusion
- [Phase 04-01]: Added absolute emotional pulse mapping from parsed state and successful action-to-biology delta mapping with blocked-action no-op contract
- [Phase 04-02]: Added parsed-drive override bridge for next-tick perception and per-action cooldown gating contract with repeated-action rejection
- [Phase 04-03]: Added end-of-tick feedback buffer/apply contracts and explicit `BioRate` (dt-scaled) vs `BioPulse` (one-shot) type separation

### Pending Todos

- Start Phase 5 `05-01` simulation loop integration planning slice.

### Blockers/Concerns

- Phenomenological translation quality remains highest-risk and will need iterative prompt validation.

## Session Continuity

Last session: 2026-02-19
Stopped at: Phase 4 complete through `04-03`
Resume file: `.planning/CONTINUITY.md`
Next plan: Phase 5 `05-01`

# Phase 1: Biological Foundation - Context

**Gathered:** 2026-02-19
**Status:** Ready for planning

<domain>
## Phase Boundary

Build the slim bio model (8-10 variables) with decay, degradation, thresholds, and noise. This is the biological substrate that the motivation layer (Phase 2) will compute drives from. No LLM integration, no consciousness, no persistence — just the bio engine.

</domain>

<decisions>
## Implementation Decisions

### Variable Selection
- Keep all 8 proposed variables: energy, stress, cognitive capacity, mood, physical tension, hunger/satiation, social deficit, body temperature
- Use physiological ranges (not normalized 0-1) — e.g., body temp 25-43, stress 0-1 where that makes physiological sense
- Drive mapping approved:
  - Energy drive ← energy + hunger
  - Social connection drive ← social deficit
  - Stimulation/novelty drive ← cognitive capacity + mood
  - Safety drive ← stress + tension + body temperature
  - Identity coherence drive ← mood + cognitive capacity
- Slim interaction rules from V1's 76 down to ~20-30, keeping only rules that meaningfully affect drives

### Degradation Model
- Linear decay rates (constant per second) — no accelerating curves
- Degradation visible within 3-5 minutes of neglect (configurable, start fast for development)
- No automatic homeostasis — vars only change from explicit causes (decay, actions, feedback). Configurable in case homeostasis is needed later.
- Recovery is gradual over time when needs are met (not instant partial relief) — action triggers recovery process that plays out over multiple ticks

### V1 Code Reuse
- Complete rewrite from scratch — V1 as reference only, no porting
- No V1 patterns carried over (state-as-value, Tick method signature, ThresholdResult struct) — design fresh
- Module path: same as V1 (`github.com/marczahn/person`), V2 lives in `v2/` subdirectory
- Testing: both table-driven tests for unit logic AND scenario-based tests for degradation behavior over time

### Threshold Behavior
- Three severity levels: mild, warning, critical
- Thresholds both flag conditions AND trigger cascading bio effects (e.g., extreme stress → tension spike + cognitive capacity drop)
- Terminal states configurable: vars can reach lethal extremes, but toggle defaults to off for development
- Bio noise magnitude: Claude's discretion — calibrate for variability vs readability

### Claude's Discretion
- Bio noise magnitude calibration (~2-5% range)
- Specific physiological ranges per variable
- Which V1 interaction rules are motivation-relevant (within the ~20-30 target)
- Internal data structures and API design (fresh design, no V1 constraints)

</decisions>

<specifics>
## Specific Ideas

- V1's body temp range {34,42} was too narrow — clamped hypothermia reversal temps up to 34. Use wider range.
- V1's cortisol half-life of 4500s was physiologically correct but killed urgency. With no auto-homeostasis in V2, this is addressed by design.
- The 10-minute unattended degradation test from research — this is the key validation: run bio engine alone for 10+ minutes and verify vars visibly worsen.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 01-biological-foundation*
*Context gathered: 2026-02-19*

# Requirements: VirtualPerson V2

**Defined:** 2026-02-19  
**Updated:** 2026-02-20  
**Core Value:** The person must exhibit intrinsic motivation with visible, behavior-driving pressure (not narrative-only emotions).

## v1 Requirements

### Biology

- [x] **BIO-01**: Bio model reduced to 8-10 motivationally-relevant variables (energy, stress, cognitive capacity, mood, physical tension, hunger/satiation, social deficit, body temperature)
- [x] **BIO-02**: Every bio variable connects to at least one drive in the motivation system
- [x] **BIO-03**: Bio variables decay toward degraded states without engagement (slow-path degradation that accumulates over time)
- [x] **BIO-04**: Slow-path degradation rates exceed homeostasis recovery rates when needs go unmet
- [x] **BIO-05**: Gaussian noise applied to bio state each tick to prevent deterministic stagnation
- [x] **BIO-06**: Bio state clamped within valid ranges per variable
- [x] **BIO-07**: Threshold system detects critical bio conditions and surfaces them

### Motivation

- [x] **MOT-01**: 5 core drives are computed from bio state via explicit formulas: energy, social connection, stimulation/novelty, safety, identity coherence
- [x] **MOT-02**: Each drive urgency is a scalar in [0,1] and monotonic with unmet need severity at fixed personality
- [x] **MOT-03**: 7 personality factors modulate drive sensitivity as multipliers: stress sensitivity, energy resilience, curiosity, self-observation, patience/frustration tolerance, risk aversion, social factor
- [x] **MOT-04**: Personality multipliers are applied exactly once in motivation computation; downstream layers consume outputs, not re-apply multipliers
- [x] **MOT-05**: Motivation processor is deterministic and side-effect free: `Compute(bio, personality, chronic) MotivationState` uses no RNG, time, I/O, env, or network
- [x] **MOT-06**: Active goal is selected from highest-urgency drive with deterministic tie-break order
- [x] **MOT-07**: Action candidates are generated from active goal and current constraints using pure rule logic (no LLM dependency)

### Consciousness

- [x] **CON-01**: Drive intensities translated to phenomenological language in prompts (felt experience, never raw numbers)
- [x] **CON-02**: Only top 1-2 drives injected as primary prompt context; lower drives as background texture
- [x] **CON-03**: Active goal framed as implicit pull, not command syntax
- [x] **CON-04**: LLM output parsed for `[STATE: arousal=X, valence=Y]` tag
- [x] **CON-05**: LLM output parsed for `[ACTION: type]` tag
- [x] **CON-06**: LLM output parsed for `[DRIVE: name=value]` tag (optional)
- [x] **CON-07**: All tags stripped from narrative content before display/storage
- [x] **CON-08**: Defensive parsing with fallback to prior known-good state on malformed/missing tags
- [x] **CON-09**: Action execution gated by environment — action fails if environment does not allow it
- [x] **CON-10**: Failed actions leave drives unsatisfied (drive compounds)
- [x] **CON-11**: Drive overrides clamped: `effective_drive = max(raw_drive*0.5, llm_reported_drive)`
- [x] **CON-12**: Spontaneous thought system fires on a ticker, drive-weighted
- [x] **CON-13**: Spontaneous thoughts fire even at baseline bio state (associative drift always available)
- [x] **CON-14**: Thought continuity buffer included in all prompts

### Feedback Loop

- [x] **FBK-01**: Emotional pulse path: `[STATE]` arousal/valence to calibrated bio changes (absolute pulses, not dt-scaled)
- [x] **FBK-02**: Action execution path: successful `[ACTION]` to bio state changes
- [x] **FBK-03**: Drive override path: `[DRIVE]` to clamped perception update for next tick
- [x] **FBK-04**: Action effects have cooldowns
- [x] **FBK-05**: Bio changes from feedback applied at end of tick, not mid-tick
- [x] **FBK-06**: Explicit type split between `BioRate` and `BioPulse`

### Infrastructure

- [x] **INF-01**: CLI output tagged by source layer: BIO / DRIVES / MIND
- [x] **INF-02**: Drive state changes displayed when significant
- [ ] **INF-03**: Configuration struct for all tunable parameters
- [ ] **INF-04**: All rates, weights, and thresholds adjustable without code changes
- [x] **INF-05**: External input handling: speech, actions, environment changes (same conventions as v1)
- [x] **INF-06**: Scenario injection with explicit bio effects
- [x] **INF-07**: Simulation loop sequencing contract: drain input -> biology -> motivation -> consciousness -> feedback
- [ ] **INF-08**: Executable runtime entrypoint with tick scheduler, input loop, and graceful shutdown

### Initial Profile (discussion-derived)

- [ ] **PRF-01**: Startup requires an explicit profile contract containing baseline bio state, personality factors, and identity seed fields
- [ ] **PRF-02**: Profile includes stress/cortisol reactivity factor that measurably modulates stress-related bio changes
- [ ] **PRF-03**: Required start-situation context fields are validated before runtime start (environment, operator relation stance, immediate constraints)

## Discussion-Derived Basis (2026-02-20)

The current roadmap explicitly incorporates these discussion constraints:
- Consciousness requires both motivation and reflection to avoid passivity.
- Intrinsic motivation must create pressure that changes behavior.
- A defined starting situation is required for meaningful simulation onset.
- Person characterization must include stress-reactivity traits (for example, cortisol sensitivity).
- Generative AI may assist interpretation, but core state transitions remain code-contract based.

## v2 Requirements (deferred)

Deferred to V2.1 and beyond:
- Identity erosion dynamics (`IDE-*`)
- Behavioral drift phases (`DRF-*`)
- Goal formation (`GOL-*`)
- Extended LLM drive modulation (`MOD-*`)
- Satisfaction reward events (`SAT-*`)
- SQLite persistence expansion (`PER-*`)
- Reviewer extensions (`REV-*`)

## Out of Scope

- Circadian rhythm system
- 20-variable biological model
- Big Five personality model
- Multi-agent social interaction
- Web dashboard / WebSocket server
- i18n / localization
- Scenario scripting engine

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| BIO-01..BIO-07 | 1 | Complete |
| MOT-01..MOT-07 | 2 | Complete |
| CON-01..CON-14 | 3 | Complete |
| FBK-01..FBK-06 | 4 | Complete |
| INF-05 | 5 (`05-02`) | Complete |
| INF-06 | 5 (`05-03`) | Complete |
| INF-01 | 5 (`05-04`) | Complete |
| INF-02 | 5 (`05-04`) | Complete |
| INF-07 | 5 (`05-01`) | Complete |
| INF-08 | 5 (`05-05`) | Pending |
| INF-03 | 6 (`06-01`) | Pending |
| INF-04 | 6 (`06-02`) | Pending |
| PRF-01..PRF-03 | 6 (`06-03`) | Pending |

**Coverage:**
- Active v1 requirements: 45 total
- Mapped to phases: 45
- Unmapped: 0

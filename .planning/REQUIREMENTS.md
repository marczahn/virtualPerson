# Requirements: VirtualPerson V2

**Defined:** 2026-02-19
**Core Value:** The person must exhibit intrinsic motivation — visible drives, desires, and frustrations that emerge from the interplay between biological state, a reward/motivation system, and LLM-driven consciousness.

## v1 Requirements

Requirements for initial V2 release. Each maps to roadmap phases.

### Biology

- [x] **BIO-01**: Bio model reduced to 8-10 motivationally-relevant variables (energy, stress, cognitive capacity, mood, physical tension, hunger/satiation, social deficit, body temperature)
- [x] **BIO-02**: Every bio variable connects to at least one drive in the motivation system
- [x] **BIO-03**: Bio variables decay toward degraded states without engagement (slow-path degradation that accumulates over time)
- [x] **BIO-04**: Slow-path degradation rates exceed homeostasis recovery rates when needs go unmet
- [ ] **BIO-05**: Gaussian noise applied to bio state each tick to prevent deterministic stagnation
- [x] **BIO-06**: Bio state clamped within valid ranges per variable
- [ ] **BIO-07**: Threshold system detects critical bio conditions and surfaces them

### Motivation

- [ ] **MOT-01**: 5 core drives computed from bio state: energy, social connection, stimulation/novelty, safety, identity coherence
- [ ] **MOT-02**: Drive urgency is a scalar (0-1) that rises without satisfaction
- [ ] **MOT-03**: 7 personality factors modulate drive sensitivity as multipliers: stress sensitivity, energy resilience, curiosity, self-observation, patience/frustration tolerance, risk aversion, social factor
- [ ] **MOT-04**: Personality multipliers applied exactly once (in motivation processor), not duplicated across layers
- [ ] **MOT-05**: Motivation processor is a pure function: Compute(bio, personality, chronic) returns MotivationState with no side effects
- [ ] **MOT-06**: Active goal selected from highest-urgency drive
- [ ] **MOT-07**: Action candidates generated for the active goal

### Consciousness

- [ ] **CON-01**: Drive intensities translated to phenomenological language in prompts (felt experience, never raw numbers)
- [ ] **CON-02**: Only top 1-2 drives injected as primary prompt context; lower drives as background texture
- [ ] **CON-03**: Active goal framed as implicit pull ("find something to eat"), not command ("GOAL: hunger")
- [ ] **CON-04**: LLM output parsed for [STATE: arousal=X, valence=Y] tag (carried from V1)
- [ ] **CON-05**: LLM output parsed for [ACTION: type] tag — what the person tries to do
- [ ] **CON-06**: LLM output parsed for [DRIVE: name=value] tag — optional drive perception override
- [ ] **CON-07**: All tags stripped from narrative content before display/storage
- [ ] **CON-08**: Defensive parsing with fallback to prior known-good state on malformed/missing tags
- [ ] **CON-09**: Action execution gated by environment — action fails if environment doesn't allow it
- [ ] **CON-10**: Failed actions leave drives unsatisfied (drive compounds)
- [ ] **CON-11**: Drive overrides clamped: effective_drive = max(raw_drive * 0.5, llm_reported_drive)
- [ ] **CON-12**: Spontaneous thought system fires on a ticker, drive-weighted (high-urgency drives dominate thought selection)
- [ ] **CON-13**: Spontaneous thoughts fire even at baseline bio state (associative drift category always available)
- [ ] **CON-14**: Thought continuity buffer included in all prompts (recent thought history)

### Feedback Loop

- [ ] **FBK-01**: Emotional pulse path: [STATE] arousal/valence → calibrated bio changes (absolute pulses, not dt-scaled)
- [ ] **FBK-02**: Action execution path: successful [ACTION] → bio state changes (e.g., eat → hunger reduction)
- [ ] **FBK-03**: Drive override path: [DRIVE] → clamped perception update for next tick
- [ ] **FBK-04**: Action effects have cooldowns (e.g., eat cannot fire again for N seconds)
- [ ] **FBK-05**: Bio changes from feedback applied at end of tick, not mid-tick
- [ ] **FBK-06**: Type distinction between BioRate (per-tick, dt-multiplied) and BioPulse (one-time, not dt-multiplied)

### Infrastructure

- [ ] **INF-01**: CLI output tagged by source layer: BIO / DRIVES / MIND
- [ ] **INF-02**: Drive state changes displayed when significant
- [ ] **INF-03**: Configuration struct for all tunable parameters (drive rates, degradation slopes, feedback multipliers, personality defaults, thresholds)
- [ ] **INF-04**: All rates, weights, and thresholds adjustable without code changes
- [ ] **INF-05**: External input handling: speech, actions, environment changes (same conventions as V1)
- [ ] **INF-06**: Scenario injection with explicit bio effects (cold room → body temp decay)
- [ ] **INF-07**: Simulation loop: sequential tick — drain input → biology → motivation → consciousness → feedback

## v2 Requirements

Deferred to V2.1. Tracked but not in current roadmap.

### Identity Erosion

- **IDE-01**: Identity coherence weakens under isolation and extreme bio degradation
- **IDE-02**: Low coherence produces fragmented/contradictory self-narrative in prompts
- **IDE-03**: Erosion reversible through engagement (not permanent psychological damage)

### Behavioral Drift

- **DRF-01**: Drive frustration produces phase transitions: seeking → demanding → collapsed/resigned
- **DRF-02**: Collapsed phase qualitatively different from seeking (not just intensity scaling)

### Goal Formation

- **GOL-01**: Person spontaneously forms goals from drive state ("I want to do X")
- **GOL-02**: Goals compete with drives in thought queue

### LLM Drive Modulation

- **MOD-01**: LLM can modulate drive intensity through reasoning (push through fatigue, ignore safety)
- **MOD-02**: Modulation bounded — suppressed drive still decays slowly; body eventually overrides

### Satisfaction Events

- **SAT-01**: Drive satisfaction produces reward pulse proportional to drive intensity at time of satisfaction
- **SAT-02**: Greater deprivation → more potent satisfaction (dopamine-style learning curve)

### Persistence

- **PER-01**: SQLite persistence for bio state, personality, identity, memories
- **PER-02**: ChronicState persisted across sessions
- **PER-03**: Motivation history and drive log tables
- **PER-04**: Bio state and identity core persisted atomically in same transaction

### Reviewer

- **REV-01**: Psychologist reviewer at 3-minute tick rate with motivation context
- **REV-02**: Reviewer includes drive trajectory and predictions in analysis
- **REV-03**: Reviewer emits identity-consistency assessments

## Out of Scope

| Feature | Reason |
|---------|--------|
| Circadian rhythm system | Adds bio-realism without motivational richness; inject as scenario context if needed |
| 20-variable biological model | V1 proved complexity without proportional behavioral value |
| Big Five personality model | Replaced by 7 motivation-serving factors |
| Multi-agent social interaction | Operator IS the social relationship; enormous complexity for no hypothesis validation |
| Web dashboard / WebSocket server | Focus on core simulation validation; CLI sufficient for V2 |
| i18n / localization | English only for V2 |
| Multi-session identity arc | Requires long-term pattern infrastructure; defer to V3 |
| Scenario scripting / narrative engine | Product feature, not simulation validation feature |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| BIO-01 | Phase 1 | Complete |
| BIO-02 | Phase 1 | Complete |
| BIO-03 | Phase 1 | Complete |
| BIO-04 | Phase 1 | Complete |
| BIO-05 | Phase 1 | Pending |
| BIO-06 | Phase 1 | Complete |
| BIO-07 | Phase 1 | Pending |
| MOT-01 | Phase 2 | Pending |
| MOT-02 | Phase 2 | Pending |
| MOT-03 | Phase 2 | Pending |
| MOT-04 | Phase 2 | Pending |
| MOT-05 | Phase 2 | Pending |
| MOT-06 | Phase 2 | Pending |
| MOT-07 | Phase 2 | Pending |
| CON-01 | Phase 3 | Pending |
| CON-02 | Phase 3 | Pending |
| CON-03 | Phase 3 | Pending |
| CON-04 | Phase 3 | Pending |
| CON-05 | Phase 3 | Pending |
| CON-06 | Phase 3 | Pending |
| CON-07 | Phase 3 | Pending |
| CON-08 | Phase 3 | Pending |
| CON-09 | Phase 3 | Pending |
| CON-10 | Phase 3 | Pending |
| CON-11 | Phase 3 | Pending |
| CON-12 | Phase 3 | Pending |
| CON-13 | Phase 3 | Pending |
| CON-14 | Phase 3 | Pending |
| FBK-01 | Phase 4 | Pending |
| FBK-02 | Phase 4 | Pending |
| FBK-03 | Phase 4 | Pending |
| FBK-04 | Phase 4 | Pending |
| FBK-05 | Phase 4 | Pending |
| FBK-06 | Phase 4 | Pending |
| INF-01 | Phase 5 | Pending |
| INF-02 | Phase 5 | Pending |
| INF-03 | Phase 6 | Pending |
| INF-04 | Phase 6 | Pending |
| INF-05 | Phase 5 | Pending |
| INF-06 | Phase 5 | Pending |
| INF-07 | Phase 5 | Pending |

**Coverage:**
- v1 requirements: 41 total
- Mapped to phases: 41
- Unmapped: 0

---
*Requirements defined: 2026-02-19*
*Last updated: 2026-02-19 after roadmap creation — all 41 requirements mapped*

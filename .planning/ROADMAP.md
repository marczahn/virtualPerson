# Roadmap: VirtualPerson V2

## Overview

V2 is a clean rebuild that adds intrinsic motivation to V1's proven pipeline architecture. The build follows the strict dependency chain: biological foundation first (nothing else can be computed without bio state), then drive computation (pure functions, fully testable without LLM), then consciousness extension (the highest-risk surface — prompt translation must shift LLM behavioral register), then feedback path wiring (closes the loop), then loop integration (assembles all components into a running simulation), and finally configuration and output (all tunable values reachable without code changes). Each phase delivers one complete, independently testable capability.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [ ] **Phase 1: Biological Foundation** - Slim bio model (8-10 vars), decay, thresholds, noise
- [ ] **Phase 2: Motivation Layer** - Pure-function drive computation from bio state and personality
- [ ] **Phase 3: Consciousness Extension** - Prompt translation, multi-tag feedback parser, thought queue
- [ ] **Phase 4: Feedback Loop** - Bidirectional bio/drive changes from LLM output
- [ ] **Phase 5: Simulation Loop Integration** - Sequential tick wiring, input handling, scenario injection, CLI display
- [ ] **Phase 6: Configuration** - All tunable parameters in config struct, no buried constants

## Phase Details

### Phase 1: Biological Foundation
**Goal**: The biological substrate exists and degrades realistically, creating the pressure that drives exist to relieve
**Depends on**: Nothing (first phase)
**Requirements**: BIO-01, BIO-02, BIO-03, BIO-04, BIO-05, BIO-06, BIO-07
**Success Criteria** (what must be TRUE):
  1. Running the simulation unattended for 10 minutes without LLM input shows measurable degradation accumulation (energy, stress, cognitive capacity all visibly worsen)
  2. Slow-path degradation rates visibly outpace homeostasis recovery when needs go unmet (verified by unit test with degradation-only scenario over N ticks)
  3. Each of the 8-10 bio variables connects to at least one drive in the motivation system (traced by documentation or test)
  4. Bio state is always clamped within valid ranges — no variable goes negative or exceeds its max under any combination of inputs
  5. Critical bio conditions (threshold crossings) are detected and surfaced with enough information to act on
**Plans**: TBD

### Phase 2: Motivation Layer
**Goal**: Drive urgencies are computed from bio state and personality as a pure function — testable without any LLM involvement
**Depends on**: Phase 1
**Requirements**: MOT-01, MOT-02, MOT-03, MOT-04, MOT-05, MOT-06, MOT-07
**Success Criteria** (what must be TRUE):
  1. Given a fixed bio state and personality configuration, Compute() always returns the same MotivationState (referential transparency, verifiable by unit test)
  2. Each of the 5 drives produces a scalar 0-1 urgency that rises as its associated bio variable degrades (verified at bio extremes)
  3. The 7 personality factors change which drives dominate without overriding the bio-derived baseline (e.g., high curiosity raises stimulation drive urgency, not energy drive urgency)
  4. Personality multipliers are applied exactly once — no drive double-counts personality during downstream processing
  5. Goal selection and action candidates reflect the highest-urgency drive, producing different outputs at different bio states
**Plans**: TBD

### Phase 3: Consciousness Extension
**Goal**: Drive state is translated into phenomenological language that actually shifts LLM behavioral register — the person feels needs, not numbers
**Depends on**: Phase 2
**Requirements**: CON-01, CON-02, CON-03, CON-04, CON-05, CON-06, CON-07, CON-08, CON-09, CON-10, CON-11, CON-12, CON-13, CON-14
**Success Criteria** (what must be TRUE):
  1. At drive urgency 0.9, the LLM response register is observably more urgent and less analytical than at drive urgency 0.2 (validated by manual inspection at drive extremes)
  2. No raw numbers appear in LLM prompts for drive state — only felt-experience language ("you cannot stop thinking about food" not "hunger: 0.8")
  3. [STATE], [ACTION], and [DRIVE] tags are parsed correctly from well-formed output, and malformed/missing tags fall back to prior known-good state without crashing
  4. Spontaneous thoughts fire even when bio state is at baseline (associative drift category always available), and thought frequency skews toward highest-urgency drives
  5. Thought continuity buffer is included in every prompt, giving the person consistent self-narrative across ticks
**Plans**: TBD

### Phase 4: Feedback Loop
**Goal**: LLM output closes the loop — emotional pulses, actions, and drive overrides all modify biological state through distinct, auditable paths
**Depends on**: Phase 3
**Requirements**: FBK-01, FBK-02, FBK-03, FBK-04, FBK-05, FBK-06
**Success Criteria** (what must be TRUE):
  1. [STATE] arousal/valence tags produce absolute bio pulses (not dt-scaled) that are measurably different from tick-rate-scaled changes
  2. Successful [ACTION] tags produce the expected bio state change for that action type (e.g., eat reduces hunger variable by the specified amount)
  3. [DRIVE] override tags are clamped — LLM cannot suppress a drive below 50% of its bio-derived value, and override is visible in logs
  4. Action execution has per-action cooldowns — a repeated action within cooldown window is rejected and the associated drive remains unsatisfied
  5. All bio changes from feedback are applied at end-of-tick, not mid-tick (verified by inspecting state at intermediate points in a test scenario)
**Plans**: TBD

### Phase 5: Simulation Loop Integration
**Goal**: All components run together in a sequential tick — the person thinks, receives input, and degrades over time in a single coherent simulation
**Depends on**: Phase 4
**Requirements**: INF-01, INF-02, INF-05, INF-06, INF-07
**Success Criteria** (what must be TRUE):
  1. Running the binary produces a live simulation where BIO, DRIVES, and MIND output streams appear tagged by source layer and update each tick
  2. External speech input produces a visible thought response from the person within the next tick
  3. Scenario injection (e.g., "cold room") produces the specified bio effects within the same tick it is injected
  4. Drive state changes that cross a significance threshold produce visible DRIVES output in the CLI — minor fluctuations do not produce noise
**Plans**: TBD

### Phase 6: Configuration
**Goal**: Every behavioral parameter is reachable and adjustable without touching source code — tuning is data, not surgery
**Depends on**: Phase 5
**Requirements**: INF-03, INF-04
**Success Criteria** (what must be TRUE):
  1. A single configuration struct contains all drive decay rates, degradation slopes, feedback multipliers, personality defaults, and thresholds
  2. Changing a value in the configuration changes the simulation's behavior on next run without any code modification or recompilation
**Plans**: TBD

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4 → 5 → 6

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Biological Foundation | 0/TBD | Not started | - |
| 2. Motivation Layer | 0/TBD | Not started | - |
| 3. Consciousness Extension | 0/TBD | Not started | - |
| 4. Feedback Loop | 0/TBD | Not started | - |
| 5. Simulation Loop Integration | 0/TBD | Not started | - |
| 6. Configuration | 0/TBD | Not started | - |

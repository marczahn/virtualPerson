# Plan: Person Simulation (Revised with Advisory Feedback)

## Context

We're building a Go program that simulates a human being with an externally injected consciousness. The person emerges from four layers: sensory input, biological processing, psychological processing, and an LLM-driven consciousness. The simulated person is unaware of the machinery — like The 13th Floor.

Advisory reports from biologist, philosopher, and psychologist are saved in `docs/advisory/`.

---

## Architecture Overview

### Four-Layer Pipeline with Feedback Loop

```
                    ┌─────────────┐
   User types  ───>│  CLI Input   │
                    └──────┬──────┘
                           │ raw text
                    ┌──────▼──────┐
                    │   Sensory   │  classifies into sensory channels
                    │  Processor  │  (LLM-assisted parsing)
                    └──────┬──────┘
                           │ SensoryEvent
                    ┌──────▼──────┐
                    │  Biological │  updates 20 state variables,
                    │  Processor  │  applies interaction rules,
                    │             │──── tick (time-based decay/circadian)
                    └──────┬──────┘
                           │ BiologicalState
                    ┌──────▼──────┐
                    │Psychological│  computes affect dimensions,
                    │   Layer     │  applies personality, regulation,
                    │             │  coping, distortions, isolation
                    └──────┬──────┘
                           │ PsychologicalState
              ┌────────────▼────────────┐
              │     Consciousness       │
              │  (salience-gated        │
              │   reactive + priority   │
              │   queue spontaneous)    │
              └────────┬────────────────┘
                       │ Thought / Emotion
                       ├──────> Display (tagged output)
                       ├──────> Memory (persist)
                       └──────> Biology FEEDBACK
                                (rumination→cortisol,
                                 catastrophizing→adrenaline,
                                 acceptance→cortisol reduction)

              ┌─────────────────────────┐
              │  Psychologist Reviewer   │  (background, periodic)
              └─────────────────────────┘
```

### Directory Structure

```
person/
├── cmd/person/main.go
├── internal/
│   ├── sense/
│   │   ├── channels.go          # SensoryChannel enum, SensoryEvent struct
│   │   └── parser.go            # LLM-assisted text→SensoryEvent parsing
│   ├── biology/
│   │   ├── state.go             # BiologicalState (20 vars), StateChange
│   │   ├── processor.go         # Tick + stimulus processing
│   │   ├── interactions.go      # Interaction rules (data-driven)
│   │   ├── circadian.go         # Circadian modulation formulas
│   │   └── thresholds.go        # Cascading failure thresholds
│   ├── psychology/
│   │   ├── state.go             # PsychologicalState (affect dimensions)
│   │   ├── processor.go         # Bio→Psych transformation
│   │   ├── personality.go       # Big Five traits + modulation
│   │   ├── regulation.go        # Emotional regulation capacity
│   │   ├── coping.go            # Coping strategy selection
│   │   ├── distortions.go       # Cognitive distortion activation
│   │   ├── emotional_memory.go  # Associative emotional memory
│   │   └── isolation.go         # Isolation effects over time
│   ├── consciousness/
│   │   ├── engine.go            # LLM-based consciousness loop
│   │   ├── prompt.go            # Prompt construction from PsychState + memory
│   │   ├── thought.go           # Thought types, priority queue
│   │   ├── salience.go          # Salience calculator (when to trigger)
│   │   └── feedback.go          # Consciousness→Biology feedback rules
│   ├── memory/
│   │   ├── store.go             # Persistence interface
│   │   ├── sqlite.go            # SQLite implementation
│   │   └── context.go           # Memory selection for LLM context
│   ├── simulation/
│   │   ├── loop.go              # Goroutine orchestration
│   │   └── clock.go             # Simulation time
│   ├── reviewer/
│   │   └── psychologist.go      # Runtime reviewer (optional)
│   └── output/
│       ├── display.go           # CLI formatting
│       └── labels.go            # Source tags (SENSE/BIO/PSYCH/MIND)
├── docs/advisory/               # Advisory reports (reference)
├── go.mod
└── go.sum
```

---

## Biological State Model (20 variables)

| Variable | Baseline | Range | Unit | Key behavior |
|---|---|---|---|---|
| Body temperature | 36.6 | 34-42 | °C | Regulates toward baseline; circadian variation +-0.5°C |
| Heart rate | 70 | 40-200 | bpm | Fast response (seconds); heavily mental-state-influenced |
| Blood pressure (sys) | 120 | 80-200 | mmHg | Follows HR + stress; morning surge |
| Respiratory rate | 15 | 8-40 | /min | Fast response; partially conscious-controllable |
| Hunger | 0.0 | 0-1 | ratio | Derived from blood sugar + glycogen (not independent) |
| Thirst | 0.0 | 0-1 | ratio | Derived from hydration level (not independent) |
| Fatigue | 0.0 | 0-1 | ratio | Accumulates ~0.05/hr waking; circadian afternoon dip |
| Pain | 0.0 | 0-1 | ratio | Partially mental-state-influenced (+-30-40%) |
| Muscle tension | 0.0 | 0-1 | ratio | Follows stress + cold + pain |
| Blood sugar | 90 | 50-200 | mg/dL | Buffered by glycogen; cortisol raises it |
| Cortisol | 0.1 | 0-1 | ratio | Peaks 15-30min after stressor; half-life 60-90min; strong circadian |
| Adrenaline | 0.0 | 0-1 | ratio | Peaks in seconds; half-life 2-3min; mental-state-driven |
| Serotonin | 0.5 | 0-1 | ratio | Slowest changes (hours-days); circadian daylight boost |
| Dopamine | 0.3 | 0-1 | ratio | Phasic spikes (seconds) + tonic baseline (hours) |
| Immune response | 0.1 | 0-1 | ratio | Slowest system (hours-days); suppressed by cortisol load |
| Circadian phase | 8.0 | 0-24 | hours | Advances 1hr/hr; cannot be "reset" instantly |
| SpO2 | 98 | 70-100 | % | Fast-moving; drops in seconds without ventilation |
| Hydration | 0.8 | 0-1 | ratio | Depletes ~0.001/min resting; rehydration ~30-60min |
| Glycogen | 0.7 | 0-1 | ratio | Buffers blood sugar; depletes 12-24hr fasting |
| Endorphins | 0.1 | 0-1 | ratio | Released after 20-30min sustained stress; half-life 20-30min |

**Interaction rules:** See `docs/advisory/biologist.md` for full interaction map with magnitudes. Key rules include:
- Cortisol lag (does NOT spike/clear instantly)
- Hypothermia reversal at 33°C (shivering stops = worsening)
- Cortisol load as time-integral for immune suppression
- Hunger/thirst derived from blood sugar and hydration

**Circadian formulas:** Sine-wave approximations for cortisol, body temp, BP, immune, serotonin, fatigue alertness (see `docs/advisory/biologist.md` section 4).

---

## Psychological State Model

The psychology layer transforms raw biology into structured affect that the consciousness layer interprets. It does NOT generate emotions — it provides the raw material.

### Output: PsychologicalState

```go
type PsychologicalState struct {
    Arousal            float64   // 0-1, from adrenaline/HR/cortisol
    Valence            float64   // -1 to 1, from serotonin/dopamine/endorphins/pain
    Energy             float64   // 0-1, from fatigue/blood sugar/circadian
    CognitiveLoad      float64   // 0-1, from cortisol duration/fatigue/blood sugar
    RegulationCapacity float64   // 0-1, depletable resource
    ActiveDistortions  []string  // e.g., "catastrophizing", "emotional_reasoning"
    LikelyCoping       []string  // e.g., "rumination", "reappraisal"
    EmotionalMemories  []EmotionalMemoryActivation
    IsolationEffects   IsolationState
}
```

### Personality (Big Five, fixed per person)

```go
type Personality struct {
    Openness          float64 // 0-1, mean 0.5
    Conscientiousness float64
    Extraversion      float64
    Agreeableness     float64
    Neuroticism       float64
}
```

Modulates: negative emotion intensity (neuroticism), isolation distress rate (extraversion), conflict response (agreeableness), reappraisal ability (openness), self-regulation bonus (conscientiousness).

### Feedback Loop: Psychology → Biology

Per simulation cycle, active coping/distortions modify biological state:
- Rumination: cortisol +0.02, serotonin -0.01
- Catastrophizing: adrenaline +0.03, cortisol +0.02
- Acceptance/Reappraisal: cortisol -0.01, serotonin +0.005
- Multiple active distortions: cortisol +0.01

---

## Consciousness Model

### Salience Calculator (when to trigger reactive consciousness)

```
salience = (rate_of_change * novelty_weight * attention_modifier) + threshold_breach_bonus
```
- `rate_of_change`: sudden changes score high
- `novelty_weight`: inverse of how recently this variable was in awareness
- `attention_modifier`: inward focus amplifies, outward suppresses
- `threshold_breach_bonus`: large bonus when variable enters extreme range
- Dynamic threshold: lower when idle/anxious, higher when engaged

### Spontaneous Thought Priority Queue

1. **Unresolved prediction errors** (highest priority) — "something unexpected happened"
2. **Active biological needs** — hunger, pain, thermal discomfort
3. **Goal rehearsal** — upcoming tasks, unfinished plans
4. **Social modeling** — "what did they think of me?"
5. **Associative drift** (lowest priority) — daydreaming, mind-wandering

### Identity Core (persistent, fed to every consciousness prompt)

Must contain:
- Dispositional traits (behavioral tendencies with context)
- Core self-narrative (2-3 sentences, may be biased/idealized)
- Relational identity markers
- 3-5 key autobiographical memories
- Emotional patterns and habits
- Values and commitments

### Memory Eviction (when context window is full)

Do NOT truncate. Compress:
- Extract identity-relevant information → update identity core
- Retain emotional residue as affect modifiers
- Summarize dropped period in one sentence

### Prompt receives PsychologicalState, NOT raw biology

The consciousness LLM gets: affect dimensions + context + memory + identity. It *constructs* the emotional experience from these — it doesn't receive pre-labeled emotions.

---

## Implementation Order

### Level 1: Project Foundation

#### 1.1 Initialize Go module
- `go mod init` + create directory structure (including `internal/psychology/`)
- Add dependencies: anthropic SDK, modernc.org/sqlite

#### 1.2 Core domain types
- `internal/sense/channels.go` — SensoryChannel enum, SensoryEvent struct
- `internal/biology/state.go` — BiologicalState (20 vars), NewDefaultState(), StateChange
- `internal/psychology/state.go` — PsychologicalState, Personality, EmotionalMemoryActivation, IsolationState
- `internal/consciousness/thought.go` — ThoughtType, Thought, ThoughtPriority
- `internal/output/labels.go` — Source enum (Sense/Bio/Psych/Mind), OutputEntry

### Level 2: Biological Layer

#### 2.1 Interaction rules
- `internal/biology/interactions.go` — data-driven rules with magnitudes from biologist advisory
- Chain resolution with max depth guard
- Rules for all 20 variables (see `docs/advisory/biologist.md` section 2)

#### 2.2 Circadian modulation
- `internal/biology/circadian.go` — sine-wave formulas for cortisol, body temp, BP, immune, serotonin, fatigue alertness

#### 2.3 Cascading failure thresholds
- `internal/biology/thresholds.go` — nonlinear breakpoints (hypothermia reversal, hypoglycemia, SpO2 cascade, cortisol load integral)

#### 2.4 Biological processor
- `internal/biology/processor.go` — Tick() + ProcessStimulus(), significance filtering

#### 2.5 Biology tests
- Cold chain, hunger chain, adrenaline decay, infinite loop guard, hypothermia reversal, cortisol load accumulation

### Level 3: Psychology Layer

#### 3.1 Affect dimension computation
- `internal/psychology/processor.go` — Bio→PsychState transformation (arousal, valence, energy, cognitive_load formulas from psychologist advisory)

#### 3.2 Personality model
- `internal/psychology/personality.go` — Big Five struct, modulation formulas (neuroticism multiplier, extraversion isolation rate, etc.)

#### 3.3 Emotional regulation
- `internal/psychology/regulation.go` — depletable capacity model, stress depletion curve, recovery rates

#### 3.4 Coping mechanism selection
- `internal/psychology/coping.go` — decision tree based on stress, personality, resources, controllability

#### 3.5 Cognitive distortions
- `internal/psychology/distortions.go` — probability-based activation with stress multiplier and trait multiplier per distortion type

#### 3.6 Emotional memory
- `internal/psychology/emotional_memory.go` — associative model with negativity bias (1.5x), power-law recency decay, trauma handling

#### 3.7 Isolation effects
- `internal/psychology/isolation.go` — timeline model (0-2hr, 2-8hr, 8-24hr, 1-3d, 3-7d, 7d+) with personality vulnerability

#### 3.8 Psychology tests
- Affect dimension computation from known bio states
- Personality modulation (high vs low neuroticism on same stimulus)
- Regulation depletion under sustained stress
- Coping selection under different resource levels
- Distortion activation at various stress levels

### Level 4: Persistence Layer

#### 4.1 Memory store interface
- `internal/memory/store.go` — Store interface (save/load bio state, save/query memories, identity core, personality, emotional memories)

#### 4.2 SQLite implementation
- `internal/memory/sqlite.go` — Schema for: biological_state (20 cols), memories, identity, personality, emotional_associations, metadata

#### 4.3 Memory context selection
- `internal/memory/context.go` — Select relevant memories for consciousness prompt, including somatic similarity retrieval (similar biological state)

#### 4.4 Persistence tests

### Level 5: Consciousness Engine

#### 5.1 Salience calculator
- `internal/consciousness/salience.go` — rate_of_change, novelty_weight, attention_modifier, threshold_breach, dynamic threshold

#### 5.2 Prompt construction
- `internal/consciousness/prompt.go` — builds prompt from PsychologicalState + memories + identity core + trigger
- System prompt: consciousness role without simulation awareness
- Manages token budget across tiers

#### 5.3 Consciousness engine
- `internal/consciousness/engine.go` — React() and Spontaneous() methods, Claude API calls, rate limiting

#### 5.4 Spontaneous thought priority queue
- `internal/consciousness/thought.go` — priority queue with weighted random selection, absorption mechanic

#### 5.5 Consciousness→Biology feedback
- `internal/consciousness/feedback.go` — maps active coping strategies and distortions to biological variable changes

#### 5.6 Consciousness tests
- Prompt construction within token budget
- Salience calculation for known scenarios
- Priority queue selection distribution
- Feedback rule application

### Level 6: Simulation Loop & CLI

#### 6.1 Simulation clock
- `internal/simulation/clock.go` — sim-time tracking, 1:1 with real time initially

#### 6.2 Simulation loop
- `internal/simulation/loop.go` — goroutine orchestration:
  1. Input goroutine (stdin → sensory channel)
  2. Sensory goroutine (parse → SensoryEvent)
  3. Biology goroutine (process stimuli + tick → state changes)
  4. Psychology goroutine (transform bio state → psych state)
  5. Consciousness reactive goroutine (salience-gated)
  6. Consciousness spontaneous goroutine (priority queue on timer)
  7. Feedback goroutine (consciousness output → biology adjustments)
  8. Display goroutine (formatted output)
- Graceful shutdown with state persistence

#### 6.3 CLI output display
- `internal/output/display.go` — color-coded by source, timestamps, buffered output
- Tags: `[SENSE]` cyan, `[BIO]` yellow, `[PSYCH]` blue, `[MIND]` green, `[REVIEW]` magenta

#### 6.4 Main entry point
- `cmd/person/main.go` — config, wiring, signal handling

#### 6.5 Integration tests

### Level 7: Psychologist Runtime Reviewer

#### 7.1 Reviewer implementation
- `internal/reviewer/psychologist.go` — periodic LLM-based review of recent outputs

#### 7.2 Integration into simulation loop
- Optional background goroutine, configurable interval

#### 7.3 Reviewer tests

---

## Verification

### Unit-level
- Biology interaction chains, circadian formulas, threshold behaviors
- Psychology affect computation, personality modulation, regulation depletion
- Persistence round-trips, memory context selection
- Salience calculation, prompt token budgets

### Scenario-based
- **Cold scenario:** temp drops → shivering → cortisol → consciousness feels cold → sustained cold → distress → if glycogen depletes → hypothermia reversal
- **Hunger scenario:** no food → blood sugar drops → glycogen depletes → fatigue → irritability → consciousness becomes desperate
- **Pain scenario:** injury → adrenaline → endorphins (after 20min) → coping strategies emerge → pain-stress amplification if sustained
- **Idle scenario:** no input → spontaneous thoughts via priority queue → loneliness over time → isolation effects based on personality
- **Feedback loop test:** induce stress → verify rumination sustains cortisol → verify acceptance reduces it

### Persistence
- Stop/restart preserves bio state, memories, identity core, personality

### Advisory review
- Share sample outputs with philosopher/biologist/psychologist agents for coherence check

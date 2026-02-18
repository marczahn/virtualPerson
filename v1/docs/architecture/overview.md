# Architecture Overview

A Go program that simulates a human being with externally injected consciousness. The "person" experiences a body, emotions, and thoughts without knowing they are simulated. An LLM provides the subjective experience; everything beneath it — sensory parsing, biological processes, psychological computation — runs invisibly, just as neurons and hormones do in a real human.

Inspired by the movie *The 13th Floor*: if a human's consciousness emerges from biological processes they don't control, and this person's consciousness is injected from computational processes they don't control, the structural relationship is analogous. The project deliberately leaves the metaphysical question open. See `../advisory/philosopher-followup-consciousness.md` for the full discussion.

## The Four-Layer Pipeline

```
                            ┌─────────────────────────────────┐
                            │        EXTERNAL INPUT           │
                            │  ~env  *action*  plain speech   │
                            └──────────────┬──────────────────┘
                                           │
                    ┌──────────────────────▼──────────────────────┐
                    │                   SENSE                      │
                    │     keyword parser → sensory events          │
                    │  channels: visual, auditory, thermal, pain,  │
                    │  tactile, olfactory, gustatory, vestibular,  │
                    │  interoceptive                                │
                    └──────────────────────┬──────────────────────┘
                                           │ sense.Event
                    ┌──────────────────────▼──────────────────────┐
                    │                 BIOLOGY                      │
                    │   20-variable state model                    │
                    │   76 interaction rules, circadian rhythm     │
                    │   decay toward baselines, threshold system   │
                    │   hypothermia reversal, cortisol load        │
                    └──────────────────────┬──────────────────────┘
                                           │ biology.State
                    ┌──────────────────────▼──────────────────────┐
                    │                PSYCHOLOGY                    │
                    │   bio → affect dimensions (arousal, valence, │
                    │   energy, cognitive load)                    │
                    │   Big Five personality modulation            │
                    │   regulation, coping, distortions, memory   │
                    │   isolation timeline                         │
                    └──────────────────────┬──────────────────────┘
                                           │ psychology.State
                    ┌──────────────────────▼──────────────────────┐
                    │              CONSCIOUSNESS                   │
                    │   salience-gated reactive thoughts           │
                    │   priority-queued spontaneous thoughts       │
                    │   direct response to speech/actions          │
                    │   LLM generates first-person experience      │
                    └──────────────┬──────────────────────────────┘
                                   │
                    ┌──────────────▼──────────────────────────────┐
                    │           FEEDBACK LOOP                      │
                    │  thought content → keyword detection →       │
                    │  biological state changes                    │
                    │  (rumination sustains cortisol,              │
                    │   catastrophizing spikes adrenaline,        │
                    │   acceptance lowers cortisol)                │
                    └──────────────┬──────────────────────────────┘
                                   │
                                   └──────────► back to BIOLOGY
```

The feedback loop is what makes this more than a feedforward pipeline. What the person thinks affects their body, and their body affects what they think. Rumination keeps cortisol high, which keeps arousal high, which keeps the person thinking about the problem. Acceptance lowers cortisol, which lowers arousal, which lets the mind rest. This produces emergent behavior that neither layer generates alone.

## The Meta-Observer (Reviewer)

A separate LLM call acts as a clinical psychologist behind a one-way mirror. It reads the person's thoughts and psychological state but cannot influence them. It produces observations like "The subject is showing signs of rumination-driven cortisol elevation." This is observation-only — the person never knows they're being watched. The reviewer is optional and runs on a rate-limited schedule (default: once per 60 seconds).

## Package Dependency Graph

```
cmd/person/main.go
    │
    └─► simulation.Loop
            │
            ├─► sense.Parser (interface)
            │       └── sense.KeywordParser
            │
            ├─► biology.Processor
            │       └── biology.State, Rules, Circadian, Thresholds
            │
            ├─► psychology.Processor
            │       └── psychology.State, Personality, Regulation,
            │           Coping, Distortions, EmotionalMemory, Isolation
            │
            ├─► consciousness.Engine
            │       ├── consciousness.LLM (interface)
            │       │       └── consciousness.ClaudeAdapter
            │       ├── consciousness.SalienceCalculator
            │       ├── consciousness.ThoughtQueue
            │       ├── consciousness.PromptBuilder
            │       └── consciousness.Feedback
            │
            ├─► reviewer.Reviewer
            │       └── reviewer.PromptBuilder
            │
            ├─► memory.Store (interface)
            │       └── memory.SQLiteStore
            │
            └─► output.Display
```

Dependencies always point inward. Biology knows nothing about psychology. Psychology knows nothing about consciousness. The simulation loop orchestrates them all.

## Directory Structure

```
person/
├── cmd/person/
│   └── main.go                  # entry point, config, wiring, signal handling
├── docs/
│   ├── advisory/                # pre-implementation LLM advisory reports
│   │   ├── biologist.md         # 20-var model design, interaction rules
│   │   ├── philosopher.md       # consciousness architecture, ethics
│   │   ├── philosopher-followup-consciousness.md
│   │   └── psychologist.md      # psych layer design, coping, distortions
│   ├── architecture/            # post-implementation documentation (this)
│   └── plan/
│       ├── decisions.md         # all design decisions with reasoning
│       └── implementation-plan.md
├── internal/
│   ├── biology/
│   │   ├── circadian.go         # sine-wave circadian rhythm formulas
│   │   ├── interactions.go      # 76 interaction rules, variable ranges
│   │   ├── processor.go         # tick cycle, decay, stimulus processing
│   │   ├── state.go             # 20-var State struct, Variable enum, Get/Set
│   │   └── thresholds.go        # Normal→Impaired→Critical→Unconscious→Lethal
│   ├── consciousness/
│   │   ├── claude.go            # ClaudeAdapter: LLM interface implementation
│   │   ├── engine.go            # Engine: React, Spontaneous, Respond
│   │   ├── feedback.go          # thought→bio changes (keyword detection)
│   │   ├── prompt.go            # PromptBuilder: system/user prompt construction
│   │   ├── queue.go             # ThoughtQueue: 5-level priority queue
│   │   ├── salience.go          # SalienceCalculator: when to become aware
│   │   └── thought.go           # Thought, ThoughtType, Priority types
│   ├── memory/
│   │   ├── context.go           # ContextSelector: somatic similarity retrieval
│   │   ├── sqlite.go            # SQLiteStore: persistence implementation
│   │   └── store.go             # Store interface, IdentityCore, EpisodicMemory
│   ├── output/
│   │   ├── display.go           # Display: colored CLI output
│   │   └── labels.go            # Source enum (SENSE/BIO/PSYCH/MIND/REVIEW)
│   ├── psychology/
│   │   ├── coping.go            # 7 coping strategies, decision tree
│   │   ├── distortions.go       # 6 cognitive distortions, activation probabilities
│   │   ├── emotional_memory.go  # associative store, negativity bias, similarity
│   │   ├── isolation.go         # 6-phase isolation timeline
│   │   ├── personality.go       # Big Five trait modulation functions
│   │   ├── processor.go         # Processor: bio→psych transformation (9 steps)
│   │   ├── regulation.go        # depletable emotional regulation model
│   │   └── state.go             # State, Personality, Distortion, CopingStrategy types
│   ├── reviewer/
│   │   ├── prompt.go            # reviewer system/user prompt construction
│   │   └── psychologist.go      # Reviewer: rolling buffer, rate limiting
│   ├── sense/
│   │   ├── channels.go          # Channel enum, Event struct
│   │   └── parser.go            # KeywordParser: text→sensory events
│   └── simulation/
│       ├── clock.go             # Clock: simulation time tracking, pause/resume
│       └── loop.go              # Loop: main tick cycle, input routing, shutdown
├── go.mod
└── go.sum
```

32 source files, 28 test files, 327 tests.

## Tech Stack

| Dependency | Purpose |
|---|---|
| Go 1.24 | Language, standard library concurrency |
| `modernc.org/sqlite` | Pure Go SQLite (no CGO) for state persistence |
| `anthropic-sdk-go` v1.22.1 | Claude API for consciousness engine and reviewer |

No web framework, no ORM, no dependency injection library. The standard library handles everything else.

## How to Run

### Configuration

The simulation reads configuration from `config.json` (or `$PERSON_CONFIG`):

```json
{
  "anthropic_api_key": "sk-ant-...",
  "model": "claude-haiku-4-5-20251001",
  "db_path": "person.db"
}
```

Environment variables override the config file:
- `ANTHROPIC_API_KEY` — required, the Claude API key
- `PERSON_MODEL` — model to use (default: `claude-haiku-4-5`)
- `PERSON_DB` — SQLite database path (default: `person.db`)
- `PERSON_CONFIG` — path to config file (default: `config.json`)

### Running

```bash
go run ./cmd/person
```

### Input Conventions

| Format | Type | Routing |
|---|---|---|
| `Hello, how are you?` | Speech | sense display → consciousness.Respond |
| `*pushes you*` | Action | sense parser → bio effects → consciousness.Respond |
| `~a freezing wind blows` | Environment | sense parser → bio effects only (consciousness reacts via salience) |

### Output Tags

All output is tagged by source layer with color:
- `[SENSE]` cyan — sensory event parsing
- `[BIO]` yellow — biological state changes
- `[PSYCH]` blue — psychological state changes (not yet wired to output)
- `[MIND]` green — consciousness thoughts
- `[REVIEW]` magenta — psychologist reviewer observations

### Shutdown

Ctrl+C (SIGINT/SIGTERM) triggers graceful shutdown: biological state and identity core are persisted to SQLite before exit. On next startup, the person resumes from where they left off.

## Related Documentation

- [Data Flow](data-flow.md) — how data moves through the system each tick
- [Biology](biology.md) — the 20-variable biological model
- [Psychology](psychology.md) — affect dimensions, personality, coping, distortions
- [Consciousness](consciousness.md) — salience, prompts, identity, reviewer
- [Design Decisions](../plan/decisions.md) — all architectural choices with reasoning
- [Implementation Plan](../plan/implementation-plan.md) — the 7-level build order

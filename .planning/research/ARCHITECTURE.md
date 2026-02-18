# Architecture Research

**Domain:** Motivation-driven consciousness simulation (Go, LLM-backed)
**Researched:** 2026-02-18
**Confidence:** HIGH (based on deep V1 codebase analysis + domain knowledge of motivation theory, agent architectures, and LLM integration patterns)

---

## Standard Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         EXTERNAL INPUT                              │
│              speech / *action* / ~environment / @scenario           │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │
┌─────────────────────────────────▼───────────────────────────────────┐
│                      SENSE LAYER (unchanged)                        │
│          keyword parser → typed sense events per channel            │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │ sense.Event
┌─────────────────────────────────▼───────────────────────────────────┐
│                      BIOLOGY LAYER (8-10 vars)                      │
│   Reduced from 20 → keep only motivationally-relevant vars          │
│   Core: Hunger, Thirst, Fatigue, Pain, Arousal proxy (HR),          │
│         Cortisol, Dopamine, Serotonin                               │
│   Degradation: persistent bio_damage field (irreversible)           │
│   Still: circadian clock, decay toward baselines, thresholds        │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │ biology.State
┌─────────────────────────────────▼───────────────────────────────────┐
│                    MOTIVATION LAYER (new)                           │
│                                                                     │
│   Drive Computation:                                                │
│     raw_drive[i] = f(bio_vars relevant to drive_i)                 │
│     modulated_drive[i] = raw_drive[i] * personality_factor[i]      │
│                                                                     │
│   Drives: hunger, thirst, sleep, pain_relief, stimulation,         │
│            connection, safety, autonomy, meaning                    │
│                                                                     │
│   Goal Selection: highest-urgency drive → active goal              │
│   Action Candidates: rule-based actions tied to goals              │
│   Reward Signal: Δ(drive_satisfaction) after action                │
│                                                                     │
│   Identity Erosion:  erosion += f(chronic_unmet_drives)            │
└──────────┬────────────────────────┬────────────────────────────────┘
           │ motivation.State       │ motivation.GoalContext
┌──────────▼────────────────────────▼────────────────────────────────┐
│                   CONSCIOUSNESS LAYER (LLM)                         │
│                                                                     │
│   Prompt = identity + scenario + bio_summary + motivation_state     │
│             + active_goal + action_candidates + recent_thoughts     │
│             + cognitive_style (from personality factors)            │
│                                                                     │
│   LLM outputs:                                                      │
│     - First-person thought/speech (cleaned content)                 │
│     - [STATE: arousal=X, valence=Y] (kept from V1)                 │
│     - [ACTION: eat|drink|rest|seek|avoid|ask|...] (new)            │
│     - [DRIVE: hunger=X, connection=Y] (optional override)          │
│                                                                     │
│   Three modes (unchanged from V1):                                  │
│     Reactive (salience-gated), Spontaneous (timer), Respond (input) │
└──────────┬─────────────────────────────────────────────────────────┘
           │ consciousness.Thought (with action tag + drive tag)
┌──────────▼─────────────────────────────────────────────────────────┐
│                    FEEDBACK LAYER (expanded)                        │
│                                                                     │
│   Path A — Emotional pulse (V1 unchanged):                         │
│     [STATE: arousal, valence] → calibrated bio changes             │
│                                                                     │
│   Path B — Drive modulation (new):                                 │
│     [DRIVE: ...] → motivation layer updates drive weights           │
│     Models: rationalization, denial, resignation                    │
│     Clamped: LLM cannot zero-out a drive below bio minimum         │
│                                                                     │
│   Path C — Action execution (new):                                  │
│     [ACTION: eat] → bio state: blood_sugar += meal_delta           │
│     [ACTION: rest] → bio state: fatigue -= rest_delta              │
│     [ACTION: seek_connection] → → isolation reset                   │
│     Gate: action is probabilistic — succeeds only if environment   │
│            allows (no food → eat fails → drive unsatisfied)        │
└──────────┬───────────────────────────────────────────────────────  ┘
           │ applied bio changes
           └──────────────────────────────────────────► back to BIOLOGY
```

---

## Component Responsibilities

| Component | Responsibility | Typical Implementation |
|-----------|----------------|------------------------|
| `biology` | Maintain 8-10 physiological vars; decay, circadian, thresholds, bio_damage accumulation | Go struct + `Tick(dt)` method; rule table slimmed from V1 |
| `motivation` | Compute drive urgencies from bio; modulate by personality factors; select goal; yield action candidates; accumulate identity erosion | Pure functions — stateless transform, no goroutines |
| `consciousness.Engine` | Wrap LLM calls; build prompts from bio+motivation; parse multi-tag output; maintain recent thought buffer | Unchanged structure from V1; new tag parsers added |
| `consciousness.PromptBuilder` | Construct system+user prompt from all input signals | New sections for motivation state and goal context |
| `consciousness.FeedbackParser` | Parse [STATE], [ACTION], [DRIVE] tags from LLM output | Extend V1 ParseFeedback; new tag regex per tag type |
| `simulation.Loop` | Orchestrate tick: input→bio→motivation→consciousness→feedback | Single goroutine main loop + input reader goroutine (same as V1) |
| `memory.Store` | Persist bio state, identity core, motivation history, drive log | SQLite; extend V1 schema |
| `reviewer.Reviewer` | Psychologist meta-observer; unchanged role | Extend prompt to include motivation state |

---

## Recommended Project Structure

```
v2/
├── cmd/person/
│   └── main.go                   # entry point, wiring, signal handling
├── internal/
│   ├── biology/
│   │   ├── state.go              # BioState struct (8-10 vars), Variable enum
│   │   ├── processor.go          # Tick(), ProcessStimulus(), decay, circadian
│   │   ├── thresholds.go         # threshold evaluation (adapted from V1)
│   │   ├── interactions.go       # slimmed rule table (~30 rules vs 76)
│   │   ├── damage.go             # BioDamage: irreversible degradation model
│   │   └── circadian.go          # circadian phase, alertness (adapted from V1)
│   ├── motivation/
│   │   ├── drives.go             # Drive enum, DriveState struct, urgency formulas
│   │   ├── personality.go        # 7 personality factors as float64 multipliers
│   │   ├── goals.go              # GoalContext: active goal, action candidates
│   │   ├── erosion.go            # IdentityErosion: chronic unmet drive tracking
│   │   └── processor.go          # Processor.Compute(bio, personality) → MotivState
│   ├── consciousness/
│   │   ├── engine.go             # Engine: React, Spontaneous, Respond (adapted)
│   │   ├── prompt.go             # PromptBuilder (extended: motivation + goal sections)
│   │   ├── feedback.go           # FeedbackParser: [STATE], [ACTION], [DRIVE] tags
│   │   ├── action.go             # ActionTag type, action validation
│   │   ├── salience.go           # SalienceCalculator (adapted from V1)
│   │   ├── queue.go              # ThoughtQueue (adapted from V1)
│   │   ├── claude.go             # ClaudeAdapter (unchanged)
│   │   └── thought.go            # Thought, ThoughtType (extended)
│   ├── memory/
│   │   ├── store.go              # Store interface (extended schema)
│   │   ├── sqlite.go             # SQLiteStore (extended with motivation tables)
│   │   └── context.go            # ContextSelector (unchanged from V1)
│   ├── sense/
│   │   ├── channels.go           # Channel enum, Event (unchanged from V1)
│   │   └── parser.go             # KeywordParser (unchanged from V1)
│   ├── simulation/
│   │   ├── clock.go              # Clock (unchanged from V1)
│   │   └── loop.go               # Loop: new tick order with motivation step
│   ├── reviewer/
│   │   ├── psychologist.go       # Reviewer (extended: motivation section)
│   │   └── prompt.go             # reviewer PromptBuilder (extended)
│   └── output/
│       ├── display.go            # Display (add MOTIV source label)
│       └── labels.go             # Source enum (add Motivation)
└── go.mod
```

### Structure Rationale

- **`motivation/` as a standalone package:** Keeps clear boundary — it reads `biology.State` and `personality`, outputs `MotivationState`. Nothing in biology knows about motivation. Nothing in consciousness knows about biology directly.
- **`biology/damage.go` isolated:** Bio damage is irreversible, so it needs its own update path separate from the decay/circadian pipeline. Isolating it prevents accidental reversal by other rules.
- **`consciousness/action.go` isolated:** Action tag parsing and validation is complex enough to warrant its own file. Validation (does the environment allow this action?) happens here, not in the loop.
- **Sense and clock unchanged:** They have no motivation-specific concerns and are clean in V1.

---

## Architectural Patterns

### Pattern 1: Pure-Function Motivation Processor

**What:** `motivation.Processor.Compute(bio biology.State, personality Personality) MotivationState` is a pure function — no internal state, no goroutines, no side effects.

**When to use:** The motivation layer sits between bio and consciousness in the tick cycle. It must be deterministic and testable without an LLM.

**Trade-offs:** Drive urgency cannot persist across ticks without external state. Solution: pass accumulated chronic metrics (e.g., `HoursHungry`) as part of the bio state or as a separate `ChronicState` struct passed in from the loop.

**Example:**
```go
// motivation/processor.go
type Processor struct {
    personality Personality
}

func (p *Processor) Compute(bio biology.State, chronic ChronicState) MotivationState {
    drives := computeDrives(bio, chronic)
    modulated := p.applyPersonality(drives)
    goal := selectGoal(modulated)
    candidates := actionCandidates(goal, bio)
    return MotivationState{
        Drives:     modulated,
        ActiveGoal: goal,
        Candidates: candidates,
        Erosion:    computeErosion(chronic),
    }
}
```

### Pattern 2: Multi-Tag LLM Output Protocol

**What:** The LLM appends multiple structured machine-readable tags to its response. Each tag is stripped from `Thought.Content` before storage. Tags are parsed independently and applied through separate feedback paths.

**When to use:** Whenever the LLM must communicate structured state back to the simulation engine without that state contaminating the narrative content the person "says."

**Trade-offs:** Tag order matters for stripping. The system prompt must specify tag format exactly. Tests must cover malformed tags, missing tags, partial tags, and tag-only responses.

**Example tags:**
```
[STATE: arousal=0.7, valence=-0.4]         ← retained from V1
[ACTION: eat]                               ← new: what the person tries to do
[DRIVE: hunger=0.8, connection=0.3]        ← new: optional LLM drive override
```

Strip order: parse all tags, strip all from content, store clean content. Never strip selectively.

**Example:**
```go
// consciousness/feedback.go
type ParsedFeedback struct {
    StateTag  EmotionalTag     // arousal, valence
    ActionTag ActionTag        // intended action
    DriveTag  DriveOverride    // optional drive modulation
    Content   string           // cleaned narrative text
}

func ParseAllTags(output string) ParsedFeedback {
    content := output
    state, content  := parseStateTag(content)
    action, content := parseActionTag(content)
    drives, content := parseDriveTag(content)
    return ParsedFeedback{state, action, drives, strings.TrimSpace(content)}
}
```

### Pattern 3: Gated Action Execution

**What:** When the LLM emits `[ACTION: eat]`, the simulation checks whether the action is possible in the current environment before applying biological effects. If impossible, the drive remains unsatisfied, which is fed back to the motivation layer.

**When to use:** Prevents the LLM from bypassing environmental constraints. The person can want to eat, but if there is no food in the scenario, the hunger drive persists and compounds.

**Trade-offs:** The gate logic must be expressed in the scenario state — the loop needs to know what resources are available. Start simple: actions are always possible unless scenario explicitly says otherwise. Elaborate as needed.

**Example:**
```go
// simulation/loop.go
func (l *Loop) executeAction(action consciousness.ActionTag) ActionResult {
    if !l.environment.Allows(action.Type) {
        // Action failed — drive remains unsatisfied
        return ActionResult{Success: false, Reason: "environment blocked"}
    }
    changes := actionEffects(action, l.cfg.BioState)
    l.applyBioChanges(changes)
    return ActionResult{Success: true}
}
```

### Pattern 4: Damped LLM Drive Overrides

**What:** The LLM can emit `[DRIVE: hunger=0.3]` to modulate its own drive perception. But the simulation clamps this override: `effective_drive = max(bio_minimum, llm_reported_drive)`. The LLM can reduce perceived drive (denial, resignation) but cannot ignore biological minimums.

**When to use:** This models the psychological reality that people rationalize away needs — but the body eventually overrides rationalization. A starving person can report feeling fine for a while, but bio state eventually dominates.

**Trade-offs:** Requires defining bio_minimum per drive (e.g., hunger minimum = f(blood_sugar)). This is a new mapping that needs calibration. Keep it simple at first: bio_minimum = current_raw_drive * 0.5.

---

## Data Flow

### Main Tick Cycle (V2)

```
[100ms ticker fires]
        │
        ▼
1. DRAIN INPUT
   read inputCh → classifyInput → route:
     speech → consciousness.Respond (same as V1)
     action → bio effects + consciousness.Respond
     environment → bio effects only
     scenario → update environment state
        │
        ▼
2. BIOLOGY TICK
   BioProcessor.Tick(bioState, dt)
     → decay, circadian, interactions, thresholds
     → apply bio_damage accumulation (new)
     → display significant changes
        │
        ▼
3. MOTIVATION COMPUTE (new step)
   MotivationProcessor.Compute(bioState, chronic)
     → drive urgencies (with personality multipliers)
     → active goal + action candidates
     → identity erosion update
   Display: [MOTIV] active goal + top drives if changed significantly
        │
        ▼
4. CONSCIOUSNESS: REACTIVE
   Engine.React(ctx, motivState, dt)
     → salience on motivation state (not just bio)
     → build prompt: identity + scenario + bio_summary
                    + motivation_state + goal_context
                    + recent thoughts
     → LLM call → ParseAllTags(response)
     → apply feedback: emotional pulse + action + drive override
        │
        ▼
5. CONSCIOUSNESS: SPONTANEOUS
   Engine.Spontaneous(ctx, motivState)
     → goal-seeded spontaneous thoughts (replaces needs-seeded in V1)
     → LLM call → ParseAllTags(response)
     → apply feedback
        │
        ▼
6. REVIEWER (optional, 3-min intervals)
   Reviewer.Review(ctx, motivState, personality)
     → includes motivation state in review prompt
        │
        ▼
7. SNAPSHOT (server mode, 2s interval)
   OnStateSnapshot(bio, motiv)
```

### Feedback Loop Detail (V2)

The three feedback paths run sequentially after every LLM call:

```
LLM response
    │
    ▼
ParseAllTags(response)
    │
    ├─► EmotionalPulse: [STATE: arousal=X, valence=Y]
    │       → bio changes (cortisol, HR, muscle tension, serotonin, dopamine)
    │       → SAME as V1 — no change
    │
    ├─► ActionExecution: [ACTION: eat]
    │       → environment.Allows(action)?
    │             YES → bio changes (blood_sugar, fatigue, etc.)
    │                   drive partial satisfaction recorded in chronic state
    │             NO  → failure logged, drive unsatisfied, chronic hunger++
    │
    └─► DriveOverride: [DRIVE: hunger=0.3]
            → clamp: effective = max(bio_minimum(hunger), 0.3)
            → update motivation.DrivePerception for next tick
            → models rationalization / denial
```

### How Motivation State Feeds Prompt Construction

```
MotivationState
    ├── Drives[]         → sorted by urgency, top 3 included in prompt
    │                       phrased as felt sensations, not labels
    │                       "your body is demanding food" not "hunger=0.9"
    │
    ├── ActiveGoal       → framed as an implicit pull
    │                       "find something to eat" not "GOAL: hunger"
    │
    ├── Candidates[]     → 2-3 possible actions the person could take
    │                       phrased as options, not commands
    │                       "you could ask for food" not "ACTION: eat"
    │
    └── Erosion.Level    → influences identity section of system prompt
                            HIGH erosion → "you feel less like yourself lately"
                            LOW erosion → identity section unchanged
```

The LLM never sees raw drive numbers. It sees phenomenological descriptions derived from those numbers. This preserves the fiction that the person is not a simulation.

---

## Goroutines vs Sequential: The V2 Evaluation

V1 used goroutines + channels for layer communication. The CLAUDE.md for V1 describes this intent, but the actual implementation in `simulation/loop.go` is **sequential within a single goroutine** — only the input reader runs concurrently. The channel is just `inputCh chan string` for input buffering.

**Verdict for V2: Keep the sequential single-goroutine approach.**

Reasons:
1. Each step needs the output of the prior step (bio → motivation → consciousness). True parallelism would require either locking or duplication.
2. LLM calls are rate-limited anyway — parallelizing the fast computation stages around slow LLM calls gives no real benefit.
3. The input reader goroutine is the one genuine parallelism need — the loop must not block waiting for keyboard input. Keep this goroutine.
4. Go's `select` + channels for input draining is idiomatic and clean. No reason to change it.

The reviewer can optionally run as a goroutine if it has its own LLM and independent state (it does in V1). Keep reviewer on its own timer goroutine.

**Recommendation:**
- Main loop: single goroutine, sequential tick steps
- Input reader: goroutine (same as V1)
- Reviewer: goroutine with independent rate limiter (same as V1)
- No new goroutines in V2

---

## Feedback Loop Stability: Preventing Runaway States

The V2 motivation loop creates new feedback cycle risks not present in V1:

### Risk 1: Drive Escalation Spiral

**Path:** Bio hunger rises → motivation urgency spikes → LLM writes about desperate hunger → EmotionalPulse: high arousal + negative valence → cortisol rises → blood sugar regulation disrupted → hunger rises further.

**Prevention:**
- Drive urgency has a hard cap at 1.0. Bio vars are already clamped.
- EmotionalPulse magnitudes are calibrated for single-thought effects (V1 calibration carries over).
- Biology decay runs every tick — cortisol half-life 4500s is the natural circuit breaker.
- Do NOT make drive urgency directly raise bio vars. Motivation → prompt only. Bio → motivation is the only causal direction.

### Risk 2: LLM Drive Override Manipulation

**Path:** LLM reports `[DRIVE: hunger=0.0]` (denial) → motivation perceives no hunger → no food-seeking behavior → bio hunger keeps rising → eventual crisis.

**Prevention:**
- Clamp: `effective_drive = max(bio_minimum, llm_reported_override)`
- Bio minimum = `raw_drive * 0.5` — LLM can suppress perception by at most 50%
- If a drive stays above 0.8 for >N ticks with no action success, bypass the override cap and report full urgency to consciousness regardless of what LLM said before.
- This models the point where rationalization fails and the body takes over.

### Risk 3: Action Loop (Eat → Full → Eat Again)

**Path:** Action `eat` succeeds → blood sugar rises → hunger drops → dopamine spike → LLM emits positive valence → next tick: dopamine decay → hunger still low → no food action → stable.

This is not a runaway — it's the intended equilibrium. But watch for: action `eat` fires every tick because the environment always provides food, causing blood sugar to max and cortisol to drop artificially.

**Prevention:**
- Action effects must have a cooldown: action `eat` cannot fire again for N seconds (models meal duration).
- Action result is stored in chronic state: `lastAte time.Time`. Cooldown check in `environment.Allows()`.

### Risk 4: Identity Erosion Runaway

**Path:** Multiple drives chronically unmet → erosion accumulates → identity section in prompt says "you feel lost" → LLM generates nihilistic thoughts → negative valence → cortisol sustained → regulatory capacity depleted → more drives unmet → more erosion.

**Prevention:**
- Erosion is monotonically increasing but capped at 1.0 (irreversible design intent).
- Erosion affects the _framing_ of the identity section in the prompt, not the drive computation. It does not feed back into biology.
- If all drives are chronically unmet, this spiral IS the intended consequence (bio degradation + identity erosion is a design goal).
- The runaway is only a problem if it is NOT intended. Since it IS intended as a consequence of extreme neglect, let it run but ensure the bio thresholds (lethal levels) terminate the simulation naturally.

### Risk 5: Spontaneous Thought Flooding

**Path:** Multiple high-urgency drives → thought queue floods → LLM calls fire too frequently → API cost spiral + feedback loop amplification.

**Prevention:**
- `MinCallInterval` from V1 applies to all LLM calls globally. Do not relax this.
- Spontaneous thought interval must be at least 30s (V1 default). Do not lower for urgency.
- High drive urgency should influence which thought fires next (priority queue), not how often thoughts fire.

---

## Integration Points

### External Services

| Service | Integration Pattern | Notes |
|---------|---------------------|-------|
| Claude API | LLM interface — one adapter, multiple call sites | Rate limiting in Engine, not in adapter |
| SQLite | Store interface behind implementation | Extend schema with motivation tables, drive log |

### Internal Boundaries

| Boundary | Communication | Notes |
|----------|---------------|-------|
| sense → biology | `sense.Event` struct | No change from V1 |
| biology → motivation | `biology.State` value (not pointer) | Motivation reads; never writes to bio directly |
| motivation → consciousness | `motivation.MotivationState` struct | New in V2; replaces `psychology.State` in prompt context |
| consciousness → feedback | `consciousness.ParsedFeedback` | Expanded from V1 `ThoughtFeedback` |
| feedback → biology | `[]biology.StateChange` applied in loop | Same as V1 |
| feedback → motivation | drive overrides applied to `ChronicState` | New path; clamped |
| loop → reviewer | `thought`, `motivState`, `personality` | Reviewer extended with motivation section |
| loop → display | `output.Entry` with source label | Add `output.Motivation` source |

**Rule:** Motivation never calls biology. Biology never calls motivation. Both read only from the layer below them. Loop orchestrates all writes.

---

## Build Order Implications

The following order minimizes integration pain:

1. **`biology` package (slim)** — Start with 8-10 var bio state stripped from V1. Remove vars not needed for motivation. Add `damage.go`. All existing bio tests adapt easily.

2. **`motivation` package** — Pure functions, no dependencies on consciousness or LLM. Fully testable in isolation. Build drives, personality multipliers, goal selection, erosion.

3. **`consciousness/feedback.go` (multi-tag parser)** — Extend V1 parser to handle [ACTION] and [DRIVE] tags. Testable without LLM via mock responses.

4. **`consciousness/prompt.go` (motivation sections)** — Add motivation state and goal context to prompt builder. Testable with mock motivation state.

5. **`simulation/loop.go` (new tick order)** — Wire motivation step between bio and consciousness. Add action execution. Add drive override application. This is the most integration-heavy step.

6. **`memory/sqlite.go` (schema extension)** — Add motivation history and drive log tables. Persist `ChronicState` across sessions.

7. **Reviewer extension** — Add motivation context to reviewer prompt. Last step because it's observational and non-critical.

---

## Anti-Patterns

### Anti-Pattern 1: Motivation Writes Directly to Biology

**What people do:** When the person "eats" (action), write to `bioState.BloodSugar` directly from the motivation processor.

**Why it's wrong:** Breaks the dependency rule (motivation → bio direction). Creates hidden state mutation outside the tick cycle, making it impossible to trace what changed and why.

**Do this instead:** Return `[]biology.StateChange` from action execution in the loop. Apply them through the standard `BioProcessor` path with a named source ("action_eat"). This keeps all bio changes auditable.

### Anti-Pattern 2: Raw Drive Numbers in LLM Prompts

**What people do:** Include `hunger=0.87, thirst=0.34` in the user message to consciousness.

**Why it's wrong:** Breaks the fiction. The person is not supposed to know their internal parameters. Numbers also invite the LLM to reason instrumentally about maximizing/minimizing values rather than experiencing embodiment.

**Do this instead:** Convert drive urgency to phenomenological language in `PromptBuilder`. `hunger=0.87` → "your stomach is cramping; it's been too long since you've eaten." The LLM constructs emotion from description, not from parameters.

### Anti-Pattern 3: Personality Factors Applied in Multiple Places

**What people do:** Apply personality multipliers in motivation, then also adjust them in the consciousness prompt, then also reference them in the reviewer.

**Why it's wrong:** Personality influence becomes untrackable and double-counted. High curiosity person becomes unrealistically different from low curiosity because every layer adds a multiplier.

**Do this instead:** Personality multipliers apply exactly once, in the motivation processor. The consciousness prompt includes personality traits as narrative descriptions in the identity section (not as modifiers to drive values). The reviewer has access to the personality struct for its own analysis.

### Anti-Pattern 4: Treating Drive Overrides as Ground Truth

**What people do:** When LLM emits `[DRIVE: hunger=0.1]`, update the canonical drive value used by all subsequent computations.

**Why it's wrong:** The LLM's self-report becomes the source of truth, which allows the LLM to suppress needs indefinitely. This inverts the intended architecture: bio is ground truth, LLM perception is a derived signal.

**Do this instead:** Store the LLM override in a separate `DrivePerception` map. The motivation processor always computes `raw_drive` from bio. The `effective_drive` passed to the prompt is `max(raw_drive * 0.5, perception_override)`. The LLM can influence perception but not ground truth.

### Anti-Pattern 5: Goroutines for Layer Communication

**What people do:** Each layer (bio, motivation, consciousness) runs as a goroutine, communicating via channels.

**Why it's wrong:** The sequential dependency (bio → motivation → consciousness) means goroutines would spend most time blocked waiting. The added complexity of channel synchronization, error propagation, and shutdown sequencing is not justified. The real bottleneck is LLM API latency, not computation.

**Do this instead:** Single goroutine for the main tick. Input reader as a goroutine (necessary to avoid blocking). Reviewer as a goroutine with its own timer (independent of the main loop). This matches what V1 actually does in practice.

---

## Sources

- V1 codebase: `/home/marczahn/dev/person/v1/` (analyzed in full)
- V1 architecture docs: `v1/docs/architecture/` (data-flow.md, overview.md)
- V1 ADRs: ADR-001 (thought continuity), ADR-003 (structured emotional feedback)
- V1 design decisions: `v1/docs/plan/decisions.md`
- Philosopher advisory: `v1/docs/advisory/philosopher.md` (predictive processing, salience calculator, spontaneous thought model)
- Psychologist advisory: `v1/docs/advisory/psychologist.md` (affect dimensions, Big Five, regulation)
- Domain knowledge: Drive reduction theory (Hull 1943), Self-Determination Theory (Deci & Ryan), homeostatic agent architectures, LLM structured output patterns

---

*Architecture research for: VirtualPerson V2 — motivation-driven consciousness simulation*
*Researched: 2026-02-18*

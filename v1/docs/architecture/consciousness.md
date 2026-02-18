# The Consciousness Engine & Reviewer

The consciousness layer generates first-person subjective experience from psychological state using an LLM. The person doesn't know they are simulated — they simply experience thoughts, feelings, and reactions as any person would. A separate reviewer module acts as a clinical observer, analyzing the person's mental patterns without influencing them.

For how consciousness fits into the pipeline, see [Overview](overview.md). For the original design, see `../advisory/philosopher.md`.

## Three Modes

The consciousness engine has three distinct entry points. At most one fires per tick (enforced by shared rate limiting).

### Reactive

Triggered when the salience calculator detects a significant change in psychological state. The person becomes aware of something shifting — a sudden pain, a temperature drop, a spike in anxiety.

```go
func (e *Engine) React(ctx context.Context, ps *psychology.State, dt float64) (*Thought, error)
```

Gated by: salience threshold AND rate limit. Returns nil if neither condition is met.

### Spontaneous

Triggered when enough time has passed since the last spontaneous thought. The person thinks even when nothing external is happening — biological needs, social concerns, goal rehearsal, or mind-wandering.

```go
func (e *Engine) Spontaneous(ctx context.Context, ps *psychology.State) (*Thought, error)
```

Gated by: `SpontaneousInterval` (default 30s) AND rate limit. Returns nil if conditions aren't met or if the priority queue produces no candidate.

### Respond

Triggered by external speech or action directed at the person. Not salience-gated — the person always processes direct communication. This is the only mode that bypasses salience.

```go
func (e *Engine) Respond(ctx context.Context, ps *psychology.State, input ExternalInput) (*Thought, error)
```

Gated by: rate limit only.

Why dual reactive/spontaneous: without spontaneous thought, the person would be a pure stimulus-response machine. They'd sit in silence until something external happened. The spontaneous queue gives them an inner life — they worry, daydream, rehearse plans, and think about people even when nothing is happening.

## Salience Calculator

`SalienceCalculator` in `salience.go` determines when a change in psychological state is significant enough to break into conscious awareness. Not every biological fluctuation becomes a thought — only salient ones.

### Formula

For each of 4 dimensions (arousal, valence, energy, pain):

```
salience = rate_of_change * novelty_weight * attention_modifier + threshold_breach_bonus
```

The maximum salience across all dimensions becomes the trigger.

**Rate of change**: `|current - previous| / dt`. Fast changes are more salient than gradual ones.

**Novelty weight**: `1.0 + log(1 + seconds_since_awareness/60) * 0.3`. Dimensions that haven't been in consciousness recently are more salient. Logarithmic growth, saturates around 2.0×. When a dimension triggers consciousness, its novelty clock resets to 0.

**Attention modifier**: `1.0 + direction * 0.5`. Direction ranges from -1 (outward/engaged) to +1 (inward/introspective). When attention is inward, body signals are amplified.

**Threshold breach bonus**: large flat bonuses when dimensions enter extreme ranges:
- Arousal > 0.8: +0.5
- Valence < -0.6: +0.4
- Energy < 0.15: +0.3
- Pain > 0.6: +0.6

### Dynamic Threshold

The threshold that salience must exceed varies with psychological state:

```
base = 0.3
if anxious (arousal > 0.5 AND valence < -0.2):  base -= 0.1   // more body-aware
if low energy (energy < 0.3):                     base -= 0.05  // bored, notices more
if high cognitive load (cogLoad > 0.6):           base += 0.1   // absorbed, notices less
clamped to [0.1, 0.6]
```

### Initialization

The first call to `Compute()` establishes baseline and never triggers. This prevents false positives on startup when `previous = 0` for all dimensions. The `initialized` flag ensures the first real comparison happens only after two data points exist.

Why salience-gated: a person doesn't consciously notice every minor fluctuation in heart rate or cortisol. Conscious awareness is selective. The salience calculator models this selectivity — only significant or novel changes break through.

## Spontaneous Thought Priority Queue

`ThoughtQueue` in `queue.go` manages 5 priority levels for spontaneous thought generation:

| Priority | Weight | Category | Source |
|---|---|---|---|
| PredictionError | 8.0 | "something unexpected happened" | External: `AddPredictionError()` |
| BiologicalNeed | 5.0 | hunger, pain, exhaustion | Auto-generated from `psychology.State` |
| GoalRehearsal | 3.0 | upcoming tasks, unfinished plans | External: `ActiveGoals` |
| SocialModeling | 2.0 | "what did they think of me?" | External: `SocialConcerns` |
| AssociativeDrift | 1.0 | daydreaming, mind-wandering | Always available (state-dependent prompt) |

**Selection**: weighted random. A prediction error (weight 8) is 8× more likely to fire than associative drift (weight 1), but drift can still surface. This prevents rigid prioritization while ensuring urgent needs usually get attention.

**Need generation**: `UpdateNeeds()` checks psychological state each tick:
- Energy < 0.2 → "You are exhausted. Your body demands rest."
- Valence < -0.5 AND arousal > 0.6 → "Something is wrong. You feel distressed."
- Isolation phase >= Loneliness → "You feel deeply alone."

**Associative drift prompts** vary by state:
- Low energy → "Your mind drifts, thoughts moving slowly."
- Negative valence → "Your thoughts wander into darker territory."
- Calm + energized → "Your mind is quiet and open."
- Default → "Your mind wanders freely."

**Absorption**: the queue tracks whether the person is deep in a thought chain. When absorbed, new triggers need higher salience to interrupt (though this mechanism is currently available but not yet wired into the salience calculator).

## Prompt Construction

`PromptBuilder` in `prompt.go` constructs two components for every LLM call:

### System Prompt

Fixed frame that establishes the consciousness role. Includes identity core if available.

```
You are a person experiencing your life moment by moment.
You think, feel, and experience the world through your body and mind.
Everything you experience is real to you.
...
IMPORTANT: You only output your inner experience — thoughts, feelings,
sensations, reactions. Never describe yourself in third person.

--- Who You Are ---
[SelfNarrative]
Your tendencies: [DispositionTraits]
Your relationships: [RelationalMarkers]
Your patterns: [EmotionalPatterns]
What matters to you: [ValuesCommitments]
Memories that define you: [KeyMemories]
```

The person must never know they are a simulation. The system prompt frames them as a real person having real experiences.

### User Prompt (varies by mode)

All modes include a **state block** that renders psychological state as felt experience, NOT numbers:

| Condition | Rendered as |
|---|---|
| Arousal > 0.7 | "Your body is highly activated — heart pounding, alert, on edge." |
| Valence < -0.5 | "Everything feels bad. A heavy, dark quality pervades your experience." |
| Energy < 0.2 | "You are deeply exhausted. Every movement feels effortful." |
| CognitiveLoad > 0.6 | "Your thinking feels muddled, hard to concentrate." |
| RegulationCapacity < 0.2 | "You feel emotionally raw, unable to hold things together." |
| Isolation >= Loneliness | "You feel lonely. You miss being around people." |

After the state block, each mode adds its specific trigger:

- **Reactive**: "Something just shifted: [trigger]" + "What are you thinking and feeling right now?"
- **Spontaneous**: "Your mind turns to: [candidate prompt]" + "What passes through your mind?"
- **Respond (speech)**: 'Someone says to you: "[content]"' + "What do you think and feel?"
- **Respond (action)**: "Someone does this: [content]" + "What do you think and feel?"

Active distortions are injected as tendencies: "Right now, your thinking tends toward: assuming the worst possible outcome; treating your feelings as evidence of reality."

Recent episodic memories are included if available (up to ~10, limited by token budget).

Why affect descriptions instead of numbers: the LLM should construct emotions from felt experience, not from a dashboard of metrics. "Your body is highly activated" produces more authentic responses than "arousal: 0.85."

## Identity Core

```go
type IdentityCore struct {
    SelfNarrative      string   // 2-3 sentences, may be biased/idealized
    DispositionTraits  []string // behavioral tendencies
    RelationalMarkers  []string // relationships that define the person
    KeyMemories        []string // 3-5 autobiographical memory summaries
    EmotionalPatterns  []string // habitual emotional responses
    ValuesCommitments  []string // deeply held values
    LastUpdated        time.Time
}
```

Fed to every consciousness prompt. Prevents the person from "resetting" between thoughts — they always know who they are, what they value, and what patterns they tend toward.

Default identity: "I'm a person, trying to make sense of the world around me. I think a lot, sometimes too much." Traits: thoughtful, curious, sometimes anxious. Values: honesty, understanding, being kind.

Persisted to SQLite. The philosopher advisory was clear: identity continuity is essential. Without it, each thought comes from a blank slate.

## Memory Context Selection

`ContextSelector` in `context.go` picks the most relevant episodic memories for the consciousness prompt:

```
score = 0.35 * importance + 0.35 * somatic_similarity + 0.30 * emotional_intensity
```

**Somatic similarity**: Euclidean distance across 6 normalized biological dimensions (arousal, valence, body temp, pain, fatigue, hunger), converted to 0-1 similarity. This models state-dependent memory / Bower's mood-congruent recall: when you're in a similar body state, related memories surface more easily.

**Emotional intensity**: `|valence|` — strongly emotional memories (positive or negative) are recalled more easily.

The selector returns at most `MaxContextMemories` (default 5) memories, sorted by score.

Why compression over truncation: the philosopher advisory was firm on this. When memory exceeds the context window, compress — extract identity-relevant information, retain emotional residue. Truncation kills continuity. The current implementation uses selection (top-N by relevance) rather than lossy compression, which preserves the important memories while fitting the token budget.

## Feedback Parsing

`ParseFeedback()` in `feedback.go` scans the LLM's response text for keywords indicating coping strategies and cognitive distortions:

| Pattern | Detection keywords | Feedback |
|---|---|---|
| Rumination | "can't stop thinking", "going over and over", "replaying" | Cortisol +0.02/s, serotonin -0.01/s |
| Acceptance | "it's okay", "accept", "let it go", "it is what it is" | Cortisol -0.01/s, serotonin +0.005/s |
| Reappraisal | "maybe it's not that bad", "another way to look", "reframe" | Cortisol -0.01/s, serotonin +0.005/s |
| Suppression | "push it down", "don't feel", "ignore it" | Cortisol +0.01/s, adrenaline +0.005/s |
| Catastrophizing | "worst case", "disaster", "going to die" | Adrenaline +0.03/s, cortisol +0.02/s |
| Problem-solving | "figure this out", "need a plan", "solve this" | (detected but no bio effect) |
| Distraction | "think about something else", "distract myself" | (detected but no bio effect) |
| Denial | "it's fine", "nothing's wrong", "not happening" | (detected but no bio effect) |

When >2 distortions are detected, an additional cortisol +0.01/s compounding penalty applies.

`FeedbackToChanges()` converts the parsed feedback into `biology.StateChange` values, which the simulation loop applies to the biological state. The deltas are per-second rates, multiplied by dt in the loop.

## Claude Adapter

`ClaudeAdapter` in `claude.go` implements the `LLM` interface:

```go
type LLM interface {
    Complete(ctx context.Context, systemPrompt, userMessage string) (string, error)
}
```

Uses `anthropic-sdk-go` v1.22.1. Key settings:
- `MaxTokens`: 1024 per response
- `Temperature`: 0.9 (high creativity for natural-sounding inner experience)
- API key from `ANTHROPIC_API_KEY` environment variable
- Default model: `claude-haiku-4-5` (fast, cheap, adequate for short inner-voice responses)

The interface abstraction allows tests to use a mock LLM without API calls. All consciousness and reviewer tests use mock implementations.

## Rate Limiting

Two independent rate limits prevent API cost explosion:

| Limit | Config field | Default | Applied to |
|---|---|---|---|
| General | `MinCallInterval` | 2s | All three modes (shared) |
| Spontaneous | `SpontaneousInterval` | 30s | Spontaneous mode only (additional) |

`canCall()` checks `time.Since(lastCallTime) >= minInterval`. Timestamps are updated before the LLM call (not after), so failures don't cause retry floods.

A `MinCallInterval` of 0 means no rate limit — useful for testing. The simulation loop sets appropriate values for production.

## The Reviewer

`Reviewer` in `reviewer/psychologist.go` is a meta-observer: a clinical psychologist behind a one-way mirror.

```go
type Reviewer struct {
    llm           consciousness.LLM  // shares the LLM interface
    promptBuilder *PromptBuilder
    minInterval   time.Duration      // default 60s
    thoughts      []consciousness.Thought  // rolling buffer
    maxThoughts   int                // default 20
}
```

### Rolling Buffer

`AddThought()` appends to a fixed-capacity buffer. When full, the oldest thought is dropped (shift-left). The reviewer always sees the N most recent thoughts.

### Prompt Design

**System prompt**: "You are a clinical psychologist observing a person through a one-way mirror. You see their inner thoughts and physiological state. Provide brief, insightful observations about psychological patterns you notice. Be concise — 2-3 sentences max. Use clinical language but keep it accessible. Never address the person directly."

**User prompt**: structured with:
- Current psychological state (arousal, valence, energy, cognitive load, regulation — as numbers)
- Active cognitive distortions (if any)
- Active coping strategies (if any)
- Personality profile (Big Five values)
- Recent thoughts (numbered list with type and trigger)
- "What patterns do you observe? Any concerns?"

The reviewer receives numbers, not felt-experience descriptions. It's a clinician reading instruments, not a person having feelings.

### Rate Limiting

Independent from consciousness rate limiting. Default: 60 seconds between reviews. Empty buffer → no review (returns nil).

### Output

```go
type Observation struct {
    Content   string
    Timestamp time.Time
}
```

Displayed as `[REVIEW]` in magenta. The person never sees this — it's for the human operator watching the simulation.

Why observation-only: the reviewer cannot modify the person's state, inject thoughts, or influence behavior. It's a diagnostic tool, not a therapist. This boundary prevents the meta-observer from breaking the simulation's integrity. The person's inner life is their own.

## Related Documentation

- [Overview](overview.md) — system architecture
- [Data Flow](data-flow.md) — how consciousness fits into the tick cycle
- [Biology](biology.md) — the state that drives psychological dimensions
- [Psychology](psychology.md) — the affect dimensions consciousness receives
- `../advisory/philosopher.md` — consciousness architecture design
- `../plan/decisions.md` — why dual reactive/spontaneous, why salience-gated

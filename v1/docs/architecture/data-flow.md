# Data Flow

How data moves through the system each simulation tick. The simulation loop runs at 100ms intervals, processing input, advancing biology, computing psychology, and triggering consciousness when conditions warrant.

For the big picture and package structure, see [Overview](overview.md).

## The Simulation Tick Cycle

Every 100ms, `simulation.Loop.tick()` executes these steps in order:

```
┌─────────────────────────────────────────────────────────┐
│ 1. DRAIN INPUT                                          │
│    Read all pending lines from inputCh                  │
│    Classify each: speech / *action* / ~environment      │
│    Route to appropriate handler (see Input Routing)      │
├─────────────────────────────────────────────────────────┤
│ 2. BIOLOGY TICK                                         │
│    Processor.Tick(state):                               │
│      a. Advance circadian clock                         │
│      b. Decay all variables toward baselines            │
│      c. Apply circadian modulation                      │
│      d. Evaluate 76 interaction rules (single pass)     │
│      e. Apply hypothermia overrides if <33°C            │
│      f. Accumulate cortisol load                        │
│      g. Apply immune suppression from cortisol load     │
│      h. Deplete hydration and accumulate fatigue        │
│      i. Evaluate critical thresholds                    │
│    Display significant bio changes                      │
├─────────────────────────────────────────────────────────┤
│ 3. PSYCHOLOGY                                           │
│    Processor.Process(bioState, dt, controllability):    │
│      a. Compute raw affect (arousal, valence, energy,   │
│         cognitive load) from biology                    │
│      b. Apply personality modulation to negative valence │
│      c. Compute stress level                            │
│      d. Update regulation (deplete or recover)          │
│      e. Apply regulation dampening to arousal/valence   │
│      f. Activate cognitive distortions (probabilistic)  │
│      g. Select coping strategies (decision tree)        │
│      h. Update isolation state                          │
│      i. Query emotional memory activations              │
├─────────────────────────────────────────────────────────┤
│ 4. CONSCIOUSNESS: REACTIVE                              │
│    Engine.React(ctx, psychState, dt):                   │
│      Compute salience of state change                   │
│      If salience > dynamic threshold AND not rate-      │
│      limited → build prompt → LLM call → Thought        │
│      Parse feedback from thought content                │
│      Apply feedback to biology                          │
├─────────────────────────────────────────────────────────┤
│ 5. CONSCIOUSNESS: SPONTANEOUS                           │
│    Engine.Spontaneous(ctx, psychState):                 │
│      If enough time since last spontaneous → update     │
│      needs → select from priority queue → build         │
│      prompt → LLM call → Thought                        │
│      Parse feedback, apply to biology                   │
├─────────────────────────────────────────────────────────┤
│ 6. REVIEWER (optional)                                  │
│    Add thoughts to rolling buffer                       │
│    If rate limit allows → build review prompt →         │
│    LLM call → display Observation                       │
└─────────────────────────────────────────────────────────┘
```

Steps 4-5 each produce at most one LLM call. The `MinCallInterval` (default 2s) prevents both from firing in the same tick. Step 6 has its own independent rate limit (default 60s).

## Input Routing

A goroutine reads stdin line-by-line and sends to `inputCh` (buffered, capacity 16). The main loop drains this channel at the start of each tick.

```
User types: "~a freezing wind blows"
                │
                ▼
        classifyInput(raw)
        prefix "~" → TypeEnvironment
        content: "a freezing wind blows"
                │
                ▼
        processEnvironment(content):
          │
          ├─► SenseParser.Parse("a freezing wind blows")
          │     keyword "freezing" → Thermal channel, intensity 0.1
          │     keyword "wind" → no additional match
          │     → Event{Channel: Thermal, Intensity: 0.1, Parsed: "feeling cold"}
          │
          ├─► Display: [SENSE] thermal: feeling cold (intensity: 0.1)
          │
          ├─► BioProcessor.ProcessStimulus(state, event)
          │     Thermal stimulus, intensity < 0.5 (cold)
          │     coldDelta = -(0.5 - 0.1) * 4 = -1.6°C
          │     state.BodyTemp: 36.6 → 35.0
          │
          └─► Display: [BIO] body_temp -1.60 (thermal_stimulus_cold)

        ─── next tick ───

        Biology tick detects BodyTemp < 35.5:
          cold_shivering rule → muscle tension increases
          cold_tachycardia rule → heart rate increases
          cold_cortisol rule → cortisol rises

        Psychology processes new bio state:
          higher arousal (from HR + adrenaline)
          lower valence (from cortisol)
          stress level rises

        Salience calculator detects significant arousal change:
          score exceeds dynamic threshold → reactive thought fires

        [MIND] [reactive, trigger: arousal changed significantly]
          "Something's wrong — I'm cold, really cold.
           My body is tense. Why is it so cold in here?"
```

No consciousness call happens during `processEnvironment` itself — the person becomes aware through the salience mechanism on subsequent ticks. This models the delay between a stimulus occurring and conscious awareness of it.

## Traced Example: Speech Input

```
User types: "Hello, how are you?"
                │
                ▼
        classifyInput(raw)
        no prefix → TypeSpeech
        content: "Hello, how are you?"
                │
                ▼
        processSpeech(content):
          │
          ├─► Display: [SENSE] speech: "Hello, how are you?"
          │
          ├─► PsychProcessor.Process(bioState, 0, 0.5)
          │     (snapshot current psychological state)
          │
          ├─► consciousness.Engine.Respond(ctx, psychState, input)
          │     input = ExternalInput{Type: InputSpeech, Content: "Hello, how are you?"}
          │     NOT salience-gated — always processes direct communication
          │     Builds prompt:
          │       system: identity core + "you are a person..."
          │       user: state block + 'Someone says to you: "Hello, how are you?"'
          │            + distortion context (if any) + recent memories
          │     → LLM call → response
          │
          ├─► Display: [MIND] [conversational, trigger: Hello, how are you?]
          │     "Oh — someone's talking to me. Hi. I'm... I think I'm okay?
          │      A little tired, maybe. It's nice to hear a voice."
          │
          └─► applyFeedback(thought)
                ParseFeedback scans response for coping/distortion keywords
                (in this case: none detected, no bio changes)
```

Speech bypasses the sensory parser entirely. The person always responds to direct communication — this is not gated by salience.

## The Feedback Loop Traced

The most important emergent behavior: rumination sustains itself.

```
Tick N:   Psychology detects high stress, selects Rumination as coping
          Consciousness generates thought with rumination keywords:
            "I can't stop thinking about what happened..."

Tick N+1: FeedbackToChanges detects "can't stop thinking" → rumination
            cortisol += 0.02 * dt
            serotonin -= 0.01 * dt
          Biology tick: cortisol higher → HR stays elevated
          Psychology: arousal stays high, valence drops further
          Stress remains high → Rumination selected again

Tick N+2: Another ruminative thought
            cortisol += 0.02 * dt again
          The cycle continues until:
            a. Regulation depletes → distortions activate → denial/suppression
            b. External input breaks the cycle (speech, environment change)
            c. Circadian cortisol baseline drops (night approaches)
            d. Cortisol natural decay eventually wins if stress source removed
```

The reverse works too: acceptance/reappraisal lowers cortisol, which lowers arousal, which reduces stress, which makes acceptance more likely to be selected again.

## Consciousness Triggering

Three distinct modes, never overlapping:

| Mode | Trigger | Gate | When |
|---|---|---|---|
| **Reactive** | Salience exceeds dynamic threshold | Rate-limited (2s default) | Significant state change (pain spike, temperature drop, arousal surge) |
| **Spontaneous** | Time since last spontaneous > interval | Rate-limited (30s default) | Idle periods — the person thinks even when nothing happens |
| **Respond** | External speech or action input | Rate-limited (2s default) | Direct communication — always processes, never salience-gated |

Reactive and spontaneous are mutually exclusive within a tick (both check `canCall()`). Respond is triggered during input processing, before the tick's reactive/spontaneous checks.

## Reviewer Flow

The reviewer operates independently of consciousness, on its own rate limit:

```
Each tick:
  If reactive thought produced → AddThought(thought)
  If spontaneous thought produced → AddThought(thought)

  Reviewer.Review(ctx, psychState, personality):
    Buffer empty? → return nil
    MinInterval not elapsed? → return nil
    Build prompt:
      system: "You are a clinical psychologist..."
      user: psych state + personality profile + last N thoughts
    LLM call → clinical observation
    Display: [REVIEW] "The subject is showing..."
```

The rolling buffer holds up to 20 thoughts (configurable). When full, the oldest thought is dropped. The reviewer sees a window of recent mental activity, not the full history.

## Shutdown Flow

```
SIGINT/SIGTERM received
        │
        ▼
context cancelled → loop exits
        │
        ▼
Loop.shutdown():
  Store.SaveBioState(currentState)
  Store.SaveIdentityCore(identity)
        │
        ▼
process exits
```

On next startup, `LoadBioState` and `LoadIdentityCore` restore the person to their last saved state. Personality is also persisted but doesn't change during a session.

## Related Documentation

- [Overview](overview.md) — the big picture, package graph, how to run
- [Biology](biology.md) — what happens inside `Processor.Tick`
- [Psychology](psychology.md) — the bio→psych transformation
- [Consciousness](consciousness.md) — salience, prompts, thought generation

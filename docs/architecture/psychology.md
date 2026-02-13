# The Psychological Model

The psychology layer transforms raw biological state into affect dimensions that consciousness can interpret. It does not label emotions. It does not output "you feel angry." It outputs high arousal + negative valence + low regulation + catastrophizing active. The consciousness LLM constructs the emotional label from these dimensions — this follows the Schachter-Singer / Barrett constructed emotion model.

For how psychology fits into the pipeline, see [Overview](overview.md). For the original design, see `../advisory/psychologist.md`.

## Why a Separate Layer

Biology produces cortisol=0.7, HR=110, serotonin=0.3. That's not what a person experiences. They experience "I feel anxious and on edge." The psychology layer is the bridge: it takes 20 biological variables and computes 4 affect dimensions, applies personality, tracks regulation, selects coping strategies, and activates cognitive distortions. Consciousness receives a psychological state, not a physiological one.

## Affect Dimensions

Four dimensions computed from biology in `processor.go`:

### Arousal (0-1)

How physically activated the body is.

```
arousal = 0.30 * adrenaline
        + 0.25 * norm_HR          // (HR - 40) / 160
        + 0.20 * cortisol
        + 0.15 * 0                // norepinephrine (not modeled, absorbed into adrenaline)
        + 0.10 * (1 - fatigue)
```

Then dampened by emotional regulation: `arousal *= (1.0 - effective_regulation * 0.5)`.

### Valence (-1 to 1)

Pleasant-unpleasant axis. Positive = pleasant, negative = unpleasant.

```
valence = 0.35 * serotonin
        + 0.30 * dopamine
        + 0.15 * endorphins
        - 0.20 * cortisol
```

Personality modulates negative valence: when valence < 0, it's multiplied by `NegativeEmotionMultiplier(personality)` = `1.0 + (neuroticism - 0.5) * 1.2`. A person with neuroticism 0.8 feels negative emotions 36% more intensely. Then dampened by regulation: `valence *= (1.0 - effective_regulation * 0.3)`.

### Energy (0-1)

How much capacity the person has for action.

```
energy = 0.30 * (1 - fatigue)
       + 0.25 * norm_blood_sugar  // (BS - 50) / 150
       + 0.25 * circadian_alertness
       + 0.20 * dopamine
```

### Cognitive Load (0-1)

How muddled thinking is. High cognitive load impairs problem-solving and coping.

```
cognitive_load = 0.30 * norm_cortisol_duration  // sigmoid: load / (load + 5.0)
              + 0.25 * fatigue
              + 0.25 * norm_blood_sugar_inv     // 1.0 - (BS - 50) / 150
              + 0.20 * 0                        // sleep debt (not modeled)
```

Cortisol duration uses the cortisol load accumulator from [Biology](biology.md), not instantaneous cortisol. This means acute stress doesn't immediately impair cognition — it's sustained stress that degrades thinking.

### Stress

Derived from the affect dimensions, used internally for coping and distortion decisions:

```
stress = 0.35 * arousal + 0.35 * max(0, -valence) + 0.30 * cognitive_load
```

## Personality (Big Five)

```go
type Personality struct {
    Openness          float64 // 0-1, curiosity, imagination
    Conscientiousness float64 // 0-1, organization, self-discipline
    Extraversion      float64 // 0-1, sociability, positive emotionality
    Agreeableness     float64 // 0-1, cooperativeness, empathy
    Neuroticism       float64 // 0-1, emotional instability, negative affect sensitivity
}
```

Personality is fixed for the person's lifetime (set at startup, persisted to SQLite). It modulates nearly everything in the psychology layer:

| Function | Formula | Effect |
|---|---|---|
| `NegativeEmotionMultiplier` | `1.0 + (N - 0.5) * 1.2` | High N amplifies negative valence |
| `IsolationDistressRate` | `1.0 + (E - 0.5) * 1.5` | High E = faster isolation distress |
| `StimulationSeeking` | `0.3 + E * 0.5` | High E = needs more stimulation |
| `SocialRewardSensitivity` | `0.2 + E * 0.6` | High E = more benefit from social contact |
| `DisorderTolerance` | `1.0 - C * 0.8` | High C = low chaos tolerance |
| `SelfRegulationBonus` | `C * 0.3` | High C = more regulation capacity |
| `PlanningUnderStress` | `C * 0.4` | High C = better structured coping |
| `ReappraisalAbility` | `0.2 + O * 0.5` | High O = better at reframing |
| `NoveltyAsThreat` | `1.0 - O * 0.7` | Low O = novelty is threatening |
| `BaselineRegulation` | `0.3 + C*0.2 + O*0.15 + (1-N)*0.2` | Multi-trait: ~0.3-0.85 |
| `IsolationResilience` | `(1-E)*0.5 + (1-N)*0.3 + C*0.2` | Introverted, stable, disciplined = resilient |

Default personality: O=0.6, C=0.5, E=0.4, A=0.6, N=0.5. A slightly introverted, moderately open, average-neuroticism person.

## Emotional Regulation

Regulation is a depletable resource modeled in `regulation.go`. The person has a finite capacity to manage emotions — under sustained stress, it collapses.

```go
type RegulationState struct {
    Capacity float64 // 0-1
}
```

**Depletion** (when stress > 0.3): `capacity -= stress * dtHours * 0.08`. Accelerating collapse below 40% of baseline: an overshoot penalty `overshoot² * 0.5` makes regulation break down quickly once it starts failing.

**Recovery** (when stress <= 0.3): `capacity += 0.1 * dtHours`. Slow — about 10% per hour during calm periods.

**Fatigue penalty**: `effective_capacity = capacity - fatigue * 0.3`. Tired people regulate worse, but this is temporary — the stored capacity doesn't drain from fatigue.

Baseline capacity depends on personality (see `BaselineRegulation` above). Capacity never exceeds baseline.

When regulation is depleted:
- Arousal is less dampened → the person feels more agitated
- Valence is less dampened → negative feelings hit harder
- Cognitive distortions activate more easily
- Coping strategy selection shifts toward maladaptive options

## Coping Strategies

Seven strategies in `coping.go`, selected by a decision tree:

```go
const (
    ProblemSolving  // active, constructive — requires cognitive resources
    Reappraisal     // reframing — requires openness + regulation
    Acceptance      // acknowledging without fighting — requires agreeableness
    Distraction     // shifting attention — default fallback
    Suppression     // pushing emotions down — costly, raises cortisol
    Rumination      // repetitive negative thinking — sustains stress
    Denial          // refusing to acknowledge reality — last resort
)
```

The selection tree in `SelectCoping()`:

```
stress > 0.85?
├── cognitive < 0.2 → Denial
├── neuroticism > 0.7 → Rumination
│
controllability > 0.5?  (stressor is controllable)
├── cognitive > 0.4 AND energy > 0.3
│   ├── conscientiousness > 0.5 → ProblemSolving
│   └── else → ProblemSolving + Distraction
├── else → Distraction
│
controllability <= 0.5?  (stressor is uncontrollable)
├── cognitive > 0.4
│   ├── openness > 0.5 AND regulation > 0.4 → Reappraisal
│   ├── agreeableness > 0.6 → Acceptance
│   └── else → Distraction
├── neuroticism > 0.6 → Rumination + Suppression
└── else → Distraction + Acceptance
```

Active coping strategies feed back into biology (see [Data Flow](data-flow.md)):
- **Rumination** → cortisol +0.02/s, serotonin -0.01/s
- **Acceptance/Reappraisal** → cortisol -0.01/s, serotonin +0.005/s
- **Catastrophizing** (distortion) → adrenaline +0.03/s, cortisol +0.02/s

## Cognitive Distortions

Six distortions in `distortions.go`, activated probabilistically based on stress, regulation, and personality:

```
probability = baseRate * stressMultiplier * traitMultiplier * regReduction
```

Where `regReduction = 1.0 - regulation * 0.6` (better regulation reduces distortion probability).

| Distortion | Base rate | Trait multiplier | Activation behavior |
|---|---|---|---|
| Catastrophizing | 0.05 | `1 + N*2.0 + (1-O)*0.5` | Assuming worst possible outcome |
| AllOrNothing | 0.04 | `1 + N*1.5 + (1-O)*1.0` | Black-and-white thinking |
| Personalization | 0.03 | `1 + N*1.5 + A*0.8` | Blaming self for external events |
| EmotionalReasoning | 0.06 | `1 + N*1.0 + (1-C)*0.8` | "I feel it, so it must be true" |
| Overgeneralization | 0.04 | `1 + N*1.5` | "This always happens" |
| MindReading | 0.03 | `1 + N*1.5 + E*0.5` | Assuming others' negative thoughts |

Stress multiplier accelerates sharply above 0.7:
- stress <= 0.3: ×1.0
- stress 0.3-0.5: ×1.0 + (stress-0.3)×2.5
- stress 0.5-0.7: ×1.5 + (stress-0.5)×3.0
- stress > 0.7: ×2.1 + (stress-0.7)×4.0

Active distortions are injected into the [consciousness prompt](consciousness.md) as tendencies: "Right now, your thinking tends toward: assuming the worst possible outcome; seeing things in black and white."

When >2 distortions are active simultaneously, an additional cortisol +0.01/s penalty is applied (distortion load compounding stress).

## Emotional Memory

An associative store in `emotional_memory.go` that links past experiences to biological states:

```go
type EmotionalMemory struct {
    ID        string
    Stimulus  string    // "cold", "pain", "darkness", etc.
    Valence   float64   // -1 to 1
    Intensity float64   // 0 to 1
    CreatedAt time.Time
    Traumatic bool      // slower decay, more intrusive
}
```

**Activation**: each tick, memories are checked against current biological state via `stimulusSimilarity()`. A "cold" memory activates when body temp drops below 35.5°C, with similarity proportional to how cold.

**Negativity bias**: negative memories are weighted 1.5× more heavily than positive ones (`emotionalSalience`). This is well-documented in psychology — bad experiences are remembered more strongly than good ones.

**Decay**: power law `1.0 / (1 + days^exponent)`. Normal memories use exponent 0.5, traumatic memories use 0.3 (slower decay). A normal memory at 100 days has ~9% of its original intensity; a traumatic one has ~20%.

**Memory modifier**: `MemoryModifier()` computes a weighted sum of activated memories, scaled by neuroticism (`0.7 + N * 0.6`). The result modifies emotional intensity.

## Isolation Timeline

Six phases tracked in `isolation.go`, advancing with time since last social contact:

| Phase | Duration | Loneliness level | Effects |
|---|---|---|---|
| None | 0-2hr | ~0.05 × E_factor | No significant effect |
| Boredom | 2-8hr | 0.10-0.25 × E_factor | Restlessness, time feels slow |
| Loneliness | 8-24hr | 0.25-0.50 × E_factor | Active loneliness, craving contact |
| Significant | 1-3 days | 0.50-0.70 × E_factor | Cognitive effects, mood disruption |
| Destabilizing | 3-7 days | 0.70-0.85 | Identity disturbance, paranoia |
| Severe | 7+ days | 0.85-1.00 | Hallucinations, dissociation |

Where E_factor = `IsolationDistressRate(personality)` = `1.0 + (extraversion - 0.5) * 1.5`. An extraverted person (E=0.8) reaches loneliness ~45% faster than an introverted one (E=0.2).

Isolation phase is included in the [consciousness prompt](consciousness.md): "You feel lonely. You miss being around people." At severe levels: "The isolation is unbearable. You're not sure what's real anymore."

`RecordSocialContact()` resets isolation state completely.

## The 9-Step Processing Pipeline

`Processor.Process()` executes these steps in order every tick:

1. **Raw affect**: compute arousal, valence, energy, cognitive load from biology
2. **Personality modulation**: amplify/dampen negative valence based on neuroticism
3. **Stress computation**: derive overall stress from affect dimensions
4. **Regulation update**: deplete or recover regulation capacity based on stress
5. **Regulation dampening**: reduce arousal and valence extremes based on effective capacity
6. **Distortion activation**: probabilistically activate cognitive distortions
7. **Coping selection**: choose strategies via decision tree
8. **Isolation update**: advance isolation timeline
9. **Memory activation**: query emotional memories by somatic similarity

The output `psychology.State` packages all of this for the consciousness layer.

## Feedback to Biology

`Processor.FeedbackChanges()` translates active coping and distortions into biological state changes:

| Pattern | Variable | Delta (per second) | Source |
|---|---|---|---|
| Rumination | Cortisol | +0.02 | `psych_rumination` |
| Rumination | Serotonin | -0.01 | `psych_rumination` |
| Acceptance/Reappraisal | Cortisol | -0.01 | `psych_acceptance` |
| Acceptance/Reappraisal | Serotonin | +0.005 | `psych_acceptance` |
| Catastrophizing | Adrenaline | +0.03 | `psych_catastrophizing` |
| Catastrophizing | Cortisol | +0.02 | `psych_catastrophizing` |
| >2 distortions active | Cortisol | +0.01 | `psych_distortion_load` |

Note: in the current implementation, this feedback method exists on the Processor but the simulation loop applies consciousness-level feedback via `consciousness.FeedbackToChanges()` instead. Both produce the same patterns — the psychology-level feedback is available for future use when the loop may apply psych feedback directly.

## Related Documentation

- [Overview](overview.md) — system architecture
- [Biology](biology.md) — the 20 variables that feed into this layer
- [Consciousness](consciousness.md) — how affect dimensions become subjective experience
- [Data Flow](data-flow.md) — the full tick cycle
- `../advisory/psychologist.md` — original psychological model design

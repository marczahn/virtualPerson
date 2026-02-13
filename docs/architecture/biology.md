# The Biological Model

The biology layer models the person's body as 20 interconnected variables that change in response to sensory stimuli, circadian rhythms, interaction rules, and feedback from consciousness. It produces raw physiological state — not emotions, not feelings. The [psychology layer](psychology.md) transforms this into affect dimensions that consciousness can interpret.

For how biology fits into the overall pipeline, see [Overview](overview.md). For the original design rationale, see `../advisory/biologist.md`.

## The 20 Variables

```go
type State struct {
    BodyTemp        float64   // °C
    HeartRate       float64   // bpm
    BloodPressure   float64   // systolic mmHg
    RespiratoryRate float64   // breaths/min
    Hunger          float64   // derived from blood sugar + glycogen
    Thirst          float64   // derived from hydration
    Fatigue         float64   // accumulates ~0.05/hr waking
    Pain            float64   // partially mental-state-influenced
    MuscleTension   float64   // follows stress + cold + pain
    BloodSugar      float64   // mg/dL
    Cortisol        float64   // peaks 15-30min after stressor
    Adrenaline      float64   // peaks in seconds, half-life 2-3min
    Serotonin       float64   // slowest changes (hours-days)
    Dopamine        float64   // phasic spikes + tonic baseline
    ImmuneResponse  float64   // suppressed by cortisol load
    CircadianPhase  float64   // hours 0-24
    SpO2            float64   // oxygen saturation %
    Hydration       float64   // depletes ~0.001/min resting
    Glycogen        float64   // buffers blood sugar
    Endorphins      float64   // released after sustained stress
    CortisolLoad    float64   // time-integral for chronic stress
}
```

| Variable | Baseline | Range | Unit | Decay half-life | Key behavior |
|---|---|---|---|---|---|
| BodyTemp | 36.6 | 25-43 | °C | ~15min toward circadian target | Hypothermia cascade below 35.5, reversal at 33 |
| HeartRate | 70 | 40-200 | bpm | ~75s toward 70 | Couples to BP, respiratory rate, blood sugar |
| BloodPressure | 120 | 80-200 | mmHg | ~180s toward 120 | Driven by HR, circadian morning surge |
| RespiratoryRate | 15 | 8-40 | /min | ~45s toward 15 | Driven by HR demand, SpO2 recovery |
| Hunger | 0.0 | 0-1 | ratio | — | Derived from blood sugar + glycogen |
| Thirst | 0.0 | 0-1 | ratio | — | Derived from hydration level |
| Fatigue | 0.0 | 0-1 | ratio | +0.05/hr | Accumulates during waking, penalizes regulation |
| Pain | 0.0 | 0-1 | ratio | ~30min (15min with endorphins) | Triggers adrenaline, modulated by mental state |
| MuscleTension | 0.0 | 0-1 | ratio | ~7.5min | Follows cold + stress, drops in hypothermia reversal |
| BloodSugar | 90 | 50-200 | mg/dL | ~90min toward 90 (insulin) | Metabolic demand from HR, glycogen buffer |
| Cortisol | 0.1 | 0-1 | ratio | ~75min toward 0.1 | Circadian cycle, stress peaks lag 15-30min |
| Adrenaline | 0.0 | 0-1 | ratio | ~150s | Fastest variable — seconds to spike, minutes to clear |
| Serotonin | 0.5 | 0-1 | ratio | ~15min toward 0.3 tonic | Circadian day/night shift, slowest neuromod |
| Dopamine | 0.3 | 0-1 | ratio | ~15min toward 0.3 tonic | Phasic spikes for reward, tonic for motivation |
| ImmuneResponse | 0.1 | 0-1 | ratio | — | Circadian (stronger at night), suppressed by cortisol load |
| CircadianPhase | 8.0 | 0-24 | hours | +1hr/hr | Drives cortisol, body temp, BP, serotonin, alertness |
| SpO2 | 98 | 70-100 | % | ~30s recovery | Drops with respiratory failure, recovers with breathing |
| Hydration | 0.8 | 0-1 | ratio | -0.001/min | Natural water loss, derives thirst |
| Glycogen | 0.7 | 0-1 | ratio | -0.0005/min | Energy buffer, depletes in 12-24hr fasting |
| Endorphins | 0.1 | 0-1 | ratio | ~25min toward 0.1 | Released after sustained stress, halves pain decay |

Why 20 and not 15: the biologist advisory added SpO2 (without it, the person can't suffocate), hydration (without it, thirst has no mechanism), glycogen (without it, blood sugar has no buffer), and endorphins (without it, pain has no natural relief). Each fills a gap that would produce implausible behavior.

Why hunger/thirst are derived, not independent: a standalone hunger variable could contradict blood sugar (hungry but blood sugar is 90). Deriving hunger from the underlying metabolic variables prevents incoherent states.

## Variable Ranges and Clamping

Every variable is clamped to a valid range on every modification. The ranges live in `variableRanges` in `interactions.go`:

```go
var variableRanges = map[Variable][2]float64{
    VarBodyTemp:        {25, 43},   // was {34,42} — expanded for hypothermia modeling
    VarHeartRate:       {40, 200},
    VarBloodPressure:   {80, 200},
    // ... ratio variables: {0, 1}, blood sugar: {50, 200}, SpO2: {70, 100}
}
```

Body temperature range was originally {34, 42}. This clamped hypothermia reversal temperatures UP to 34°C, making it impossible to model the person cooling below 34. Changed to {25, 43} so the full hypothermia cascade — including the lethal <28°C range — can play out.

## Interaction Rules

76 rules in `AllRules()` (`interactions.go`), each with a condition, target variable, and delta function. Rules are evaluated once per tick (no cascading within a single tick, to prevent runaway feedback).

```go
type Rule struct {
    Name      string
    Condition func(s *State) bool
    Target    Variable
    Delta     func(s *State, dt float64) float64
}
```

Rules are organized by triggering system:

| System | Rules | Key examples |
|---|---|---|
| Body temperature | 7 | `cold_shivering`: temp < 35.5 && >= 33 → muscle tension ↑. `fever_tachycardia`: temp > 38.5 → HR ↑ |
| Heart rate | 6 | `hr_bp_coupling`: HR > 100 → BP ↑. `hr_metabolic_demand`: HR > 100 → blood sugar ↓ |
| Blood pressure | 4 | BP-driven respiratory effects, organ stress |
| Respiratory rate | 4 | Respiratory-driven SpO2 effects |
| Blood sugar | 7 | `hypoglycemia_tremor`: BS < 60 → muscle tension ↑. `glycogen_buffer`: BS < 70 → BS ↑ |
| Cortisol | 4 | `cortisol_hr`: cortisol > 0.5 → HR ↑. `cortisol_serotonin`: cortisol > 0.6 → serotonin ↓ |
| Adrenaline | 7 | `adrenaline_hr`: adrenaline > 0.1 → HR ↑. `adrenaline_bp`: adrenaline > 0.3 → BP ↑ |
| Pain | 7 | `pain_cortisol`: pain > 0.3 → cortisol ↑. `pain_endorphin_release`: sustained pain → endorphins ↑ |
| Fatigue | 4 | `fatigue_hr_reduction`: fatigue > 0.7 → HR recovery. `fatigue_serotonin`: fatigue > 0.8 → serotonin ↓ |
| Immune response | 5 | Immune-fever interaction, cortisol suppression effects |
| Hydration | 6 | `dehydration_hr`: hydration < 0.5 → HR ↑. `dehydration_fatigue`: hydration < 0.3 → fatigue ↑ |
| SpO2 | 4 | `hypoxia_hr`: SpO2 < 92 → HR ↑. `hypoxia_cognitive`: SpO2 < 88 → impairment |
| Other | 11 | Muscle tension, serotonin, dopamine, endorphins, glycogen buffering, derived hunger |

All deltas are calibrated per-second and scaled by actual dt. A rule that says `0.002 * dt` means +0.002 per second of real time.

During hypothermia reversal (<33°C), rules that increase HR or muscle tension are suppressed. The body's compensatory mechanisms are failing — the person stops shivering and their heart slows, which counterintuitively makes them feel better while dying.

## Circadian Rhythm

Six circadian curves computed from `CircadianPhase` (hours 0-24) in `circadian.go`:

| Curve | Formula | Peak | Trough |
|---|---|---|---|
| Cortisol baseline | `0.15 + 0.25 * max(0, cos((phase-7) * π/12))` | 06:00-08:00 (~0.40) | 23:00-02:00 (~0.15) |
| Body temp target | `36.6 + 0.5 * cos((phase-17) * π/12)` | ~17:00 (37.1°C) | ~05:00 (36.1°C) |
| Blood pressure shift | `10 * cos((phase-10) * π/12)` | ~10:00 (+10 mmHg) | ~22:00 (-10 mmHg) |
| Immune multiplier | `1.0 + 0.2 * cos(phase * π/12)` | ~00:00 (1.2×) | ~12:00 (0.8×) |
| Serotonin shift | `0.05 * cos((phase-13) * π/12)` | ~13:00 (+0.05) | ~01:00 (-0.05) |
| Alertness | `0.5 + 0.2*cos(2π(phase-16)/24) - 0.25*cos(2π(phase-14.5)/12)` | ~10:00, ~20:00 | ~04:00, ~14:30 (afternoon dip) |

The alertness formula combines a 24h fundamental (peak ~16:00) with a 12h harmonic that creates the well-known post-lunch dip around 14:30. The minus sign on the harmonic is correct — it subtracts the dip from the fundamental.

Circadian modulation is applied as a gentle pull (not an override). The cortisol pull rate is `0.0005 * dt` — it takes about 30 minutes for circadian to fully shift cortisol, which means acute stress can override the circadian signal.

## Threshold System

`EvaluateThresholds()` checks for critical conditions every tick, returning zero or more `ThresholdResult` values:

```go
type ThresholdResult struct {
    Condition   CriticalState  // Normal, Impaired, Critical, Unconscious, Lethal
    System      string         // "thermoregulation", "glycemic", "respiratory", etc.
    Description string
}
```

| System | Impaired | Critical | Unconscious | Lethal |
|---|---|---|---|---|
| Hypothermia | < 35.5°C | < 33°C (shivering stops) | < 30°C (cardiac arrhythmia) | < 28°C (ventricular fibrillation) |
| Hyperthermia | > 39.5°C | > 40.5°C (heat exhaustion) | > 41.5°C (heat stroke) | > 42°C (organ failure) |
| Hypoglycemia | < 60 mg/dL | < 45 mg/dL | < 35 mg/dL | < 25 mg/dL |
| Hyperglycemia | — | > 180 mg/dL | — | — |
| SpO2 | < 92% | < 85% | < 75% | < 70% |
| Blood pressure | > 180 mmHg | — | — | — |

Thresholds are reported to the display and could gate consciousness behavior (e.g., Unconscious state could prevent thought generation). The cascade is: Normal → Impaired → Critical → Unconscious → Lethal.

## Hypothermia Reversal at 33°C

Below 33°C, `IsHypothermiaReversal()` returns true. This triggers two mechanisms:

1. **Rule suppression**: interaction rules that increase HR or muscle tension are skipped. The body can no longer compensate.
2. **Override application**: `ApplyHypothermiaOverrides()` forces muscle tension toward 0.3, HR toward 50, and adrenaline toward 0 at ~10%/sec convergence.

The counterintuitive result: the person stops shivering, heart rate drops, adrenaline clears. They feel calmer and warmer — a documented phenomenon in real hypothermia cases where victims undress (paradoxical undressing). The biology models this because the consciousness layer should be able to generate the subjective experience of "feeling better" while the body is failing.

## Cortisol Load: Chronic Stress

Cortisol load is a time-integral that accumulates whenever cortisol exceeds 0.3:

```
cortisol_load += max(0, cortisol - 0.3) * dt
```

This models the difference between acute stress (cortisol spikes then clears) and chronic stress (cortisol stays elevated, load accumulates). The load suppresses immune function via:

```
immune_factor = 1.0 / (1.0 + cortisol_load * 0.1)
```

At load 0: factor = 1.0 (no suppression). At load 10: factor ~0.5. At load 100: factor ~0.09. This means prolonged rumination (which sustains cortisol via the [feedback loop](data-flow.md)) gradually degrades the immune system — a well-documented effect of chronic psychological stress.

## Stimulus Processing

`Processor.ProcessStimulus()` converts sensory events into immediate state changes:

| Channel | Effect |
|---|---|
| Thermal (cold, intensity < 0.5) | BodyTemp -= (0.5 - intensity) * 4. "Freezing" (0.1) → -1.6°C |
| Thermal (hot, intensity > 0.5) | BodyTemp += (intensity - 0.5) * 4. "Burning" (0.95) → +1.8°C |
| Pain | Pain += intensity. Adrenaline += intensity * 0.3 |
| Auditory (intensity > 0.7) | Adrenaline += (intensity - 0.7) * 0.5. HR += (intensity - 0.7) * 20 |
| Visual (intensity > 0.8) | Adrenaline += (intensity - 0.8) * 0.4 |

Stimulus effects are immediate (applied in the same tick), while interaction rule consequences cascade over subsequent ticks.

## The Tick Processing Order

The order matters. Within `Processor.Tick()`:

1. Circadian clock advances
2. Decay toward baselines (natural recovery)
3. Circadian pull (gentle, doesn't override acute changes)
4. Interaction rules (single pass, no cascading)
5. Hypothermia overrides (takes precedence over rules)
6. Cortisol load accumulation
7. Immune suppression from load
8. Hydration depletion
9. Fatigue accumulation
10. Threshold evaluation

Decay runs before interactions so that acute spikes decay first, then interactions can re-elevate variables if conditions still hold. Hypothermia overrides run after interactions so they can suppress compensatory mechanisms. dt is capped at 300s to prevent huge jumps after pauses.

## Related Documentation

- [Overview](overview.md) — system architecture and how to run
- [Data Flow](data-flow.md) — how biology fits into the tick cycle
- [Psychology](psychology.md) — how biology becomes affect dimensions
- [Consciousness](consciousness.md) — how the person becomes aware of biological changes
- `../advisory/biologist.md` — original biological model design with full rationale

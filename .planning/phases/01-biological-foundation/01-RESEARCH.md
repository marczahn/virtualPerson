# Phase 1: Biological Foundation - Research

**Researched:** 2026-02-19
**Domain:** Go biological simulation — state model, decay, thresholds, noise, interaction rules
**Confidence:** HIGH

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

#### Variable Selection
- Keep all 8 proposed variables: energy, stress, cognitive capacity, mood, physical tension, hunger/satiation, social deficit, body temperature
- Use physiological ranges (not normalized 0-1) — e.g., body temp 25-43, stress 0-1 where that makes physiological sense
- Drive mapping approved:
  - Energy drive <- energy + hunger
  - Social connection drive <- social deficit
  - Stimulation/novelty drive <- cognitive capacity + mood
  - Safety drive <- stress + tension + body temperature
  - Identity coherence drive <- mood + cognitive capacity
- Slim interaction rules from V1's 76 down to ~20-30, keeping only rules that meaningfully affect drives

#### Degradation Model
- Linear decay rates (constant per second) — no accelerating curves
- Degradation visible within 3-5 minutes of neglect (configurable, start fast for development)
- No automatic homeostasis — vars only change from explicit causes (decay, actions, feedback). Configurable in case homeostasis is needed later.
- Recovery is gradual over time when needs are met (not instant partial relief) — action triggers recovery process that plays out over multiple ticks

#### V1 Code Reuse
- Complete rewrite from scratch — V1 as reference only, no porting
- No V1 patterns carried over (state-as-value, Tick method signature, ThresholdResult struct) — design fresh
- Module path: same as V1 (`github.com/marczahn/person`), V2 lives in `v2/` subdirectory
- Testing: both table-driven tests for unit logic AND scenario-based tests for degradation behavior over time

#### Threshold Behavior
- Three severity levels: mild, warning, critical
- Thresholds both flag conditions AND trigger cascading bio effects (e.g., extreme stress → tension spike + cognitive capacity drop)
- Terminal states configurable: vars can reach lethal extremes, but toggle defaults to off for development
- Bio noise magnitude: Claude's discretion — calibrate for variability vs readability

### Claude's Discretion
- Bio noise magnitude calibration (~2-5% range)
- Specific physiological ranges per variable
- Which V1 interaction rules are motivation-relevant (within the ~20-30 target)
- Internal data structures and API design (fresh design, no V1 constraints)

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope.
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|-----------------|
| BIO-01 | Bio model reduced to 8-10 motivationally-relevant variables (energy, stress, cognitive capacity, mood, physical tension, hunger/satiation, social deficit, body temperature) | Variable definitions, physiological ranges, and baseline values documented in Standard Stack section |
| BIO-02 | Every bio variable connects to at least one drive in the motivation system | Drive mapping table in Architecture Patterns section maps each variable to drives |
| BIO-03 | Bio variables decay toward degraded states without engagement (slow-path degradation that accumulates over time) | Linear decay rates per variable specified; 3-5 minute degradation window calibrated |
| BIO-04 | Slow-path degradation rates exceed homeostasis recovery rates when needs go unmet | No-homeostasis design confirmed; decay always wins when unmet; recovery is explicit-only |
| BIO-05 | Gaussian noise applied to bio state each tick to prevent deterministic stagnation | Noise calibration recommendations in Claude's Discretion section; math/rand pattern documented |
| BIO-06 | Bio state clamped within valid ranges per variable | Per-variable ranges and clamping pattern documented with code examples |
| BIO-07 | Threshold system detects critical bio conditions and surfaces them | Three-tier threshold system (mild/warning/critical) documented with cascade effects |
</phase_requirements>

---

## Summary

This phase builds a slim 8-variable biological model from scratch in `v2/`. The design diverges significantly from V1: V1 had 20 variables with complex hormonal cascade chains optimized for physiological realism. V2 has 8 motivationally-scoped variables optimized to pressure the motivation layer. V1 is invaluable reference for calibration numbers and architecture patterns, but is explicitly not being ported.

The central design insight is that **V2 variables are motivation-shaped proxies, not physiological measurements**. Energy (0-1, not blood sugar mg/dL) is directly a motivation input. Stress (0-1) is a lumped cortisol/adrenaline proxy. This simplification removes V1's most complex problems: the insulin homeostasis dance, the glycogen buffering model, the cortisol load integrator. V2's "no automatic homeostasis" rule means variables only recover when the system explicitly acts — which makes the 10-minute degradation test trivially verifiable.

The technical implementation is pure Go with `math/rand` for Gaussian noise, `time` for delta-time tracking, and no external dependencies. The module lives at `v2/` within the existing `github.com/marczahn/person` repository. The Go stdlib provides everything needed: there are no third-party libraries required for this phase. Testing strategy uses table-driven unit tests for individual components plus scenario tests that run the engine for many ticks and assert degradation has occurred.

**Primary recommendation:** Design the `State` struct as a pointer receiver with explicit field access (not a generic Get/Set enum approach like V1 — that pattern added boilerplate without benefit for 8 variables). Design `Engine.Tick(dt float64)` to accept elapsed seconds rather than computing it internally, enabling deterministic tests without time manipulation.

---

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go stdlib `math` | 1.24 | Gaussian noise (math/rand NormFloat64), clamp logic | No external dependency, all needed math is here |
| Go stdlib `math/rand` | 1.24 | `rand.NormFloat64()` for Gaussian noise per tick | Standard approach; no seed needed since Go 1.20 (auto-seeded) |
| Go stdlib `time` | 1.24 | Timestamp tracking, delta-time calculation | Required for real-time tick rate |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| Go stdlib `testing` | 1.24 | Table-driven tests, scenario tests | All test files |
| Go stdlib `fmt` | 1.24 | Debug output during scenario tests | Scenario validation only |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| `math/rand.NormFloat64()` | Box-Muller custom implementation | No reason — stdlib is correct and simpler |
| Flat struct fields | Generic Variable enum (V1 pattern) | For 8 variables, direct field access is cleaner and faster; enum approach adds ~200 lines of Get/Set boilerplate with no benefit |
| dt passed to Tick | time.Now() called inside Tick (V1 pattern) | Passed dt enables deterministic tests without any time mocking |

**Installation:** No new dependencies. Go stdlib only.

```bash
# v2 module initialization (if not yet done)
mkdir -p v2 && cd v2
go mod init github.com/marczahn/person/v2
```

---

## Architecture Patterns

### Recommended Project Structure
```
v2/
├── go.mod                      # module github.com/marczahn/person/v2
├── go.sum
├── bio/
│   ├── state.go                # State struct, NewDefaultState(), variable ranges
│   ├── decay.go                # Linear decay rates, applyDecay()
│   ├── noise.go                # Gaussian noise application
│   ├── interactions.go         # ~20-30 interaction rules (data-driven)
│   ├── thresholds.go           # ThresholdEvent type, EvaluateThresholds()
│   ├── engine.go               # Engine struct, Tick(), ApplyDelta()
│   ├── state_test.go           # Unit tests: ranges, clamp, baseline
│   ├── decay_test.go           # Unit tests: decay rates, no-homeostasis
│   ├── thresholds_test.go      # Unit tests: threshold detection, cascades
│   ├── engine_test.go          # Unit + scenario tests: degradation over N ticks
│   └── interactions_test.go    # Unit tests: rule conditions and effects
```

Note: V1 used `internal/biology/` as a path. V2 uses `bio/` since it lives under `v2/` which provides its own isolation. Whether to put it under `internal/` is a discretion call — `internal/` prevents import by other modules but since V2 is a standalone module this is less important.

### Pattern 1: State as Pointer-Receiver Struct with Direct Fields

**What:** The bio state is an 8-field struct. Methods take `*State`. No enum-based Get/Set.
**When to use:** Always — for 8 variables, enums add boilerplate without benefit.

```go
// Source: fresh design (V1 Get/Set pattern explicitly rejected)
package bio

type State struct {
    Energy          float64 // 0-1, energy level; degrades without rest/food
    Stress          float64 // 0-1, lumped cortisol/adrenaline proxy
    CognitiveLoad   float64 // 0-1, mental fatigue/cognitive capacity used
    Mood            float64 // 0-1, 0=dysphoric, 0.5=neutral, 1=euphoric
    PhysicalTension float64 // 0-1, muscle tension, somatic stress
    Hunger          float64 // 0-1, satiation inverted: 0=full, 1=starving
    SocialDeficit   float64 // 0-1, isolation accumulation
    BodyTemp        float64 // Celsius, 25-43, baseline ~36.6
    UpdatedAt       time.Time
}

func NewDefaultState() *State {
    return &State{
        Energy:          0.8,
        Stress:          0.1,
        CognitiveLoad:   0.0,
        Mood:            0.5,
        PhysicalTension: 0.05,
        Hunger:          0.1,
        SocialDeficit:   0.0,
        BodyTemp:        36.6,
        UpdatedAt:       time.Now(),
    }
}
```

### Pattern 2: Engine with Passed dt

**What:** Engine.Tick() accepts `dt float64` (seconds elapsed), does not call time.Now() internally.
**When to use:** Always — this makes every test deterministic without time mocking.

```go
// Source: fresh design (V1 internal time.Now() pattern explicitly rejected)
package bio

type Engine struct {
    config Config
    rng    *rand.Rand
}

type Config struct {
    DecayMultiplier float64 // 1.0 = normal, 5.0 = fast development mode
    HomeostasisEnabled bool  // false = V2 default, no auto-recovery
    TerminalStatesEnabled bool // false = development default
}

type TickResult struct {
    Changes    []Delta
    Thresholds []ThresholdEvent
}

func (e *Engine) Tick(s *State, dt float64) TickResult {
    var result TickResult
    // 1. Apply linear decay (BIO-03)
    result.Changes = append(result.Changes, e.applyDecay(s, dt)...)
    // 2. Apply interaction rules (BIO-02 cross-effects)
    result.Changes = append(result.Changes, e.applyInteractions(s, dt)...)
    // 3. Apply Gaussian noise (BIO-05)
    result.Changes = append(result.Changes, e.applyNoise(s, dt)...)
    // 4. Clamp all variables (BIO-06)
    e.clampAll(s)
    // 5. Evaluate thresholds (BIO-07)
    result.Thresholds = evaluateThresholds(s)
    return result
}
```

### Pattern 3: Linear Decay (No Exponential, No Homeostasis)

**What:** Each variable has a constant decay rate (units/second). No recovery unless explicit delta applied. No exponential half-life math.
**When to use:** V2 design decision — linear is simpler and visible-within-5-minutes is tunable.

```go
// Source: fresh design based on V2 constraints
// Decay rates are per second. Multiply by DecayMultiplier from Config.
var decayRates = map[string]float64{
    "energy":          0.001,  // ~16 min to go from full to empty at 1x
    "stress":          0.0002, // stress dissipates slowly without cause
    // NOTE: stress has NO auto-decay in V2 — it only changes from explicit events
    // Only energy, hunger, social_deficit decay autonomously
}
```

**Clarification on "decay" semantics in V2:** The user decision is "variables decay toward degraded states." Degraded state is different per variable:
- Energy: degrades toward 0 (exhaustion)
- Hunger: degrades toward 1 (starvation)
- SocialDeficit: degrades toward 1 (isolation)
- BodyTemp: degrades toward danger zones from 36.6 (hypothermia or hyperthermia)
- Stress, CognitiveLoad, PhysicalTension: increase toward 1 only from explicit causes; no autonomous decay (they would decay toward 0 when nothing is happening, but the no-homeostasis rule means they stay where they are without explicit recovery action)

**Critical distinction:** "No automatic homeostasis" means no variable auto-returns to baseline. But decay toward degraded state IS autonomous — energy drains continuously without rest. This creates the required asymmetry: degradation happens automatically, recovery only happens from explicit actions.

### Pattern 4: Data-Driven Interaction Rules

**What:** Rules are structs with Condition func and Delta func, applied in a single pass per tick.
**When to use:** Keeps interaction logic declarative and testable.

```go
// Source: derived from V1 interactions.go pattern (adapted, not ported)
type Rule struct {
    Name      string
    Condition func(s *State) bool
    Target    *float64 // direct pointer to field, no Get/Set needed
    Delta     func(s *State, dt float64) float64
}
```

**Alternative:** Use a simpler slice of functions. The Rule struct approach is preferred because it gives each rule a Name for debugging/logging.

### Pattern 5: Gaussian Noise Application

**What:** Each tick, apply small Gaussian noise to each variable before clamping.
**When to use:** Every tick to prevent deterministic stagnation (BIO-05).

```go
// Source: Go stdlib math/rand documentation
// rand.NormFloat64() returns a float64 from standard normal distribution (mean 0, stddev 1)
// Scale by sigma to get desired spread.

func (e *Engine) applyNoise(s *State, dt float64) []Delta {
    // Recommended sigma: 0.002 per second scaled by sqrt(dt) for proper Brownian scaling
    // For dt=1s: adds ~0.2% noise. For dt=5s: adds ~0.45% noise.
    sigma := 0.002 * math.Sqrt(dt)
    s.Energy          += e.rng.NormFloat64() * sigma
    s.Stress          += e.rng.NormFloat64() * sigma
    s.CognitiveLoad   += e.rng.NormFloat64() * sigma
    s.Mood            += e.rng.NormFloat64() * sigma
    s.PhysicalTension += e.rng.NormFloat64() * sigma
    s.Hunger          += e.rng.NormFloat64() * sigma
    s.SocialDeficit   += e.rng.NormFloat64() * sigma
    // Body temp: smaller noise since it's a narrower functional range
    s.BodyTemp        += e.rng.NormFloat64() * sigma * 0.1
    // Clamp is called immediately after — noise that pushes out of range is absorbed
    return nil // noise changes are not individually tracked (too noisy for logging)
}
```

**Noise magnitude recommendation (Claude's Discretion):** sigma=0.002/tick at 1s ticks. This produces visible jitter at 10x zoom on a graph but is imperceptible on a 5-minute scale as noise, while the actual decay signal is clearly dominant. If ticks are 5 seconds, scale to sigma=0.002*sqrt(5)≈0.0045. The 2-5% range mentioned in CONTEXT refers to total range percentage, not per-tick — at sigma=0.002, a single tick noise is 0.2% of range, which is in the right ballpark.

### Pattern 6: Three-Tier Threshold System with Cascades

**What:** ThresholdEvent carries severity (mild/warning/critical) plus cascade effects applied to state.
**When to use:** After clamping, as the final step of each tick.

```go
// Source: fresh design (V1 CriticalState enum pattern not carried over)
type Severity int

const (
    Mild Severity = iota
    Warning
    Critical
)

type ThresholdEvent struct {
    Variable    string
    Severity    Severity
    Description string
    Cascade     []Delta // bio effects triggered by this threshold crossing
}
```

### Anti-Patterns to Avoid

- **V1 state-as-value:** V1 returned `State` by value from `NewDefaultState()`. For 8 fields this is fine, but taking addresses of value methods is error-prone. Use pointer from the start.
- **time.Now() inside Tick:** Makes tests time-dependent. Always pass dt explicitly.
- **Get/Set enum methods for 8 fields:** V1's Variable enum + switch statements added ~200 lines for 20 vars. For 8 vars with direct field access, this complexity is not justified.
- **V1 `applyDecay` with exponential formulas:** V1 used `math.Exp(-0.693/halfLife*dt)` (exponential decay). V2 uses linear decay (constant rate * dt). Simpler, more visible, and matches the constraint.
- **Threshold-only evaluation:** V1 thresholds flagged conditions but did not modify state within the threshold evaluator. V2 thresholds should also apply cascade effects (e.g., critical stress triggers tension spike).
- **Single global rng:** Use a local `*rand.Rand` in Engine, not `rand.Float64()` global — enables deterministic seeding for tests.

---

## Recommended Variable Definitions (Claude's Discretion)

These are the 8 V2 variables with recommended ranges and baselines. These are motivation-shaped proxies, not direct physiological measurements.

| Variable | Type | Range | Baseline | Degraded State | Degrade Direction | Notes |
|----------|------|-------|----------|---------------|-------------------|-------|
| Energy | float64 | 0.0–1.0 | 0.80 | 0 (exhausted) | Decreasing | Drains continuously from wakefulness |
| Stress | float64 | 0.0–1.0 | 0.10 | 1 (overwhelmed) | Increasing | Only from explicit causes; no auto-dissipation |
| CognitiveLoad | float64 | 0.0–1.0 | 0.00 | 1 (depleted) | Increasing | Mental fatigue accumulates; no auto-recovery |
| Mood | float64 | 0.0–1.0 | 0.50 | 0 (dysphoric) | Decreasing | Drifts downward without positive input |
| PhysicalTension | float64 | 0.0–1.0 | 0.05 | 1 (tense) | Increasing | Only from explicit causes |
| Hunger | float64 | 0.0–1.0 | 0.10 | 1 (starving) | Increasing | Continuous increase (approx 0-1 over 4-6 hours) |
| SocialDeficit | float64 | 0.0–1.0 | 0.00 | 1 (isolated) | Increasing | Slow continuous increase without social contact |
| BodyTemp | float64 | 25.0–43.0°C | 36.6 | <35 or >38.5 | Decreasing toward ~36° (stable, but can drift from stress) | Only variable with a meaningful "safe zone" rather than directional degradation |

**Decay rates (linear, per second, at DecayMultiplier=1.0):**

For degradation visible within 3-5 minutes (180-300 seconds), a variable starting at baseline should reach a noticeably degraded value (e.g., 20% change from baseline) in that window.

| Variable | Decay Rate/sec | Time to 20% degradation | Notes |
|----------|---------------|------------------------|-------|
| Energy | 0.00067 | ~240s (4 min) | 1.0→0.8 in 4 min at 1x |
| Hunger | 0.00083 | ~240s from 0.1→0.3 | Full starvation in ~20 min (0→1) |
| Mood | 0.00033 | ~600s (slow) | Slower drift is realistic |
| SocialDeficit | 0.00033 | ~600s (slow) | Social isolation is slow |
| CognitiveLoad | 0.00050 | ~400s | Mental load accumulates |

For development/fast mode (DecayMultiplier=5.0): same changes happen in 1/5 the time (less than 1 minute each).

**Variables with NO autonomous decay (only change from explicit causes):**
- Stress — only rises from stressors, only falls from explicit relaxation/resolution
- PhysicalTension — only rises from stress/exertion, only falls from explicit rest
- BodyTemp — only changes from environmental exposure or fever; not autonomous

---

## Drive Mapping (BIO-02)

Each variable must connect to at least one drive. This table is the BIO-02 contract:

| Variable | Drive(s) | Contribution |
|----------|---------|--------------|
| Energy | Energy drive | Primary input (low energy → high energy drive pressure) |
| Hunger | Energy drive | Secondary input (high hunger → high energy drive pressure) |
| SocialDeficit | Social connection drive | Primary input (high deficit → high drive) |
| CognitiveLoad | Stimulation/novelty drive | Inverse: low cognitive load → higher stimulation drive |
| Mood | Stimulation/novelty drive, Identity coherence drive | Low mood → higher novelty seeking; low mood → identity threat |
| Stress | Safety drive | Primary input (high stress → high safety drive) |
| PhysicalTension | Safety drive | Secondary input (high tension → high safety drive) |
| BodyTemp | Safety drive | Deviation from 36.6 → safety drive pressure |
| CognitiveLoad | Identity coherence drive | High load impairs sense of self-coherence |

---

## Recommended Interaction Rules (~20-30)

These are the V1 interaction rules most relevant to drive computation. All V1 rules involving variables not in V2 (blood pressure, heart rate, SpO2, glycogen, etc.) are dropped.

**Rules to keep (motivation-relevant, using V2 variable names):**

| # | Condition | Target | Effect | Drive Impact |
|---|-----------|--------|--------|-------------|
| 1 | Stress > 0.6 | PhysicalTension | += stress * 0.3 * dt | Safety drive |
| 2 | Stress > 0.5 | CognitiveLoad | += stress * 0.2 * dt | Identity coherence |
| 3 | Stress > 0.7 | Mood | -= 0.002 * dt | Mood, identity |
| 4 | Hunger > 0.7 | Stress | += 0.001 * dt | Safety drive |
| 5 | Hunger > 0.8 | CognitiveLoad | += 0.002 * dt | Identity coherence |
| 6 | Energy < 0.3 | Mood | -= 0.001 * dt | Mood drives |
| 7 | Energy < 0.2 | Stress | += 0.002 * dt | Safety drive |
| 8 | Energy < 0.2 | CognitiveLoad | += 0.002 * dt | Identity coherence |
| 9 | PhysicalTension > 0.7 | Stress | += 0.001 * dt | Safety drive |
| 10 | PhysicalTension > 0.6 | Mood | -= 0.001 * dt | Mood drives |
| 11 | CognitiveLoad > 0.8 | Stress | += 0.002 * dt | Safety drive |
| 12 | CognitiveLoad > 0.7 | Mood | -= 0.001 * dt | Mood drives |
| 13 | Mood < 0.2 | Stress | += 0.001 * dt | Safety drive |
| 14 | Mood < 0.2 | SocialDeficit | += 0.001 * dt | Social drive |
| 15 | SocialDeficit > 0.7 | Mood | -= 0.001 * dt | Mood drives |
| 16 | SocialDeficit > 0.8 | Stress | += 0.001 * dt | Safety drive |
| 17 | BodyTemp < 35.5 | Stress | += (35.5-temp)*0.01 * dt | Safety drive |
| 18 | BodyTemp < 35.5 | PhysicalTension | += (35.5-temp)*0.05 * dt | Safety drive |
| 19 | BodyTemp > 38.5 | Stress | += (temp-38.5)*0.01 * dt | Safety drive |
| 20 | BodyTemp > 38.5 | CognitiveLoad | += (temp-38.5)*0.03 * dt | Identity coherence |
| 21 | Energy < 0.4 && Hunger > 0.6 | Mood | -= 0.002 * dt | Compound degradation |
| 22 | Stress > 0.8 && CognitiveLoad > 0.7 | Mood | -= 0.003 * dt | Cascade: overwhelm |

Total: 22 rules. Within the 20-30 target. These rules create compound degradation spirals (e.g., hunger+stress → mood collapse) that make the bio state's pressure felt in the motivation layer.

---

## Threshold System (BIO-07)

Three severity tiers, each with cascade bio effects:

### Body Temperature
| Range | Severity | Description | Cascade Effect |
|-------|----------|-------------|----------------|
| <35.0°C | Mild | Mild hypothermia, shivering | PhysicalTension += 0.2 |
| <34.0°C | Warning | Moderate hypothermia | PhysicalTension += 0.3, CognitiveLoad += 0.2 |
| <33.0°C | Critical | Severe hypothermia | Stress += 0.3, CognitiveLoad += 0.4 |
| >38.5°C | Mild | Elevated temperature, discomfort | Stress += 0.1, CognitiveLoad += 0.1 |
| >39.5°C | Warning | Fever, significant impairment | Stress += 0.2, Mood -= 0.2 |
| >40.5°C | Critical | Heat danger | Stress += 0.4, CognitiveLoad += 0.3 |

### Stress
| Range | Severity | Description | Cascade Effect |
|-------|----------|-------------|----------------|
| >0.7 | Mild | Elevated stress | PhysicalTension += 0.01 * dt |
| >0.85 | Warning | High stress, impaired function | CognitiveLoad += 0.02 * dt, Mood -= 0.01 * dt |
| >0.95 | Critical | Crisis state | Mood -= 0.03 * dt, Energy -= 0.02 * dt |

### Energy
| Range | Severity | Description | Cascade Effect |
|-------|----------|-------------|----------------|
| <0.3 | Mild | Low energy, effort costs more | CognitiveLoad += 0.01 * dt |
| <0.15 | Warning | Very low energy | Mood -= 0.01 * dt, Stress += 0.01 * dt |
| <0.05 | Critical | Near collapse | Stress += 0.03 * dt, CognitiveLoad += 0.03 * dt |

### Hunger
| Range | Severity | Description | Cascade Effect |
|-------|----------|-------------|----------------|
| >0.7 | Mild | Noticeably hungry | Mood -= 0.005 * dt |
| >0.85 | Warning | Very hungry, difficult to focus | CognitiveLoad += 0.01 * dt, Stress += 0.005 * dt |
| >0.95 | Critical | Starving | Stress += 0.02 * dt, Energy -= 0.01 * dt |

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Gaussian noise | Custom Box-Muller or rejection sampling | `math/rand.(*Rand).NormFloat64()` | stdlib is correct, benchmarked, and seeded automatically in Go 1.20+ |
| Float clamping | Inline min/max chains | Simple `clamp(v, lo, hi float64) float64` helper | Trivial to write once; don't inline at every use site |
| Delta-time calculation | Complex timer abstraction | Pass `dt float64` as parameter | External caller controls time; enables test determinism |
| Rule engine | Reflection-based generic system | Simple `[]Rule` slice with typed functions | For 20-30 rules, a slice is faster and clearer than any ORM-style abstraction |

**Key insight:** This domain is custom physics simulation. The only third-party help needed is Gaussian noise, and Go stdlib provides it. Every other problem (decay, clamping, rules, thresholds) should be hand-written domain logic — but written once in dedicated files, not inlined everywhere.

---

## Common Pitfalls

### Pitfall 1: Conflating "No Homeostasis" with "No Decay"
**What goes wrong:** Misinterpreting the constraint as "variables stay static unless acted on." The actual rule is: degradation IS autonomous (energy drains, hunger rises), but recovery is not (stress doesn't auto-dissipate, energy doesn't auto-recover).
**Why it happens:** "No automatic homeostasis" sounds like "no automatic changes."
**How to avoid:** Model each variable explicitly: does it have an autonomous drift direction? Energy: yes (toward 0). Stress: no autonomous drift in either direction without explicit cause.
**Warning signs:** 10-minute unattended test shows no degradation — means autonomous decay is missing.

### Pitfall 2: Linear Decay Rate Calibration Off by 60-3600x
**What goes wrong:** Using physiologically correct rates from V1 or biologist docs. V1's cortisol half-life of 4500s (75 min) was physiologically correct but killed urgency. V2's design must produce visible degradation in 3-5 minutes.
**Why it happens:** Plugging in real-world values directly.
**How to avoid:** Calculate from desired degradation window: if energy should be 20% depleted in 4 minutes at normal speed, rate = 0.20 / (4 * 60) = 0.000833/sec. At DecayMultiplier=5 for development: 0.000833*5 = 0.004167/sec → 20% depleted in 48 seconds.
**Warning signs:** Running engine for 5 minutes shows < 5% change in any variable.

### Pitfall 3: Noise Applied After Clamping
**What goes wrong:** Applying Gaussian noise after the clamp step, allowing variables to escape valid ranges.
**Why it happens:** Ordering noise before clamp feels backwards ("noise should be the last thing").
**How to avoid:** Order is: (1) decay, (2) interactions, (3) noise, (4) clamp, (5) thresholds. Noise before clamp means noise-induced boundary violations are automatically corrected.
**Warning signs:** Variables found outside their declared ranges in tests.

### Pitfall 4: Threshold Cascades Applied Before Clamping
**What goes wrong:** Threshold cascade effects applied before clamping, then the cascade delta itself is not clamped, causing unbounded variable values.
**Why it happens:** Thresholds evaluate state, seem like they should apply immediately.
**How to avoid:** Threshold evaluation is step 5 (after clamping). Cascade effects must go through the same clamp path. Apply cascade deltas to state, then clamp again (or return them as Deltas to be applied by the caller with clamping).

### Pitfall 5: Body Temperature Range {34,42} Too Narrow (Known V1 Bug)
**What goes wrong:** The V1 body temp range {34,42} clamped hypothermia reversal temps up to 34°C, breaking the physiological model where significant effects occur at 33°C and 35°C.
**Why it happens:** Seems intuitive to exclude extreme temps.
**How to avoid:** Use {25, 43} as established. The 25°C lower bound allows the simulation to reach lethal hypothermia if enabled. The 43°C upper bound covers lethal hyperthermia.

### Pitfall 6: Interaction Rules Causing Positive Feedback Explosion
**What goes wrong:** Rules like "stress → cognitive load → stress" form an unbounded feedback loop within a single tick, driving variables to their maximum in one step.
**Why it happens:** Single-pass rule evaluation with conditions that reference the same tick's state changes.
**How to avoid:** Single-pass rule evaluation (V1 pattern, carry over): evaluate all conditions against the state at the START of the tick, not after each rule modifies it. Collect all deltas, apply them all at once (or sequentially but reading pre-tick state). Cap any single-tick delta from rules at a maximum value (e.g., no single rule can change a variable more than 0.1 per tick).

### Pitfall 7: dt Not Capped for Pause Recovery
**What goes wrong:** If the simulation pauses (e.g., process suspended) and then resumes, dt could be thousands of seconds, causing all variables to jump to extremes in one tick.
**Why it happens:** Using raw elapsed time without bounds.
**How to avoid:** Cap dt at a maximum (V1 used 300s). For V2 where degradation is visible in 3-5 minutes, a 60s cap is more appropriate — prevents a 30-minute pause from nuking all variables in one step.

---

## Code Examples

Verified patterns from V1 reference and Go stdlib:

### Gaussian Noise (Go stdlib)
```go
// Source: Go stdlib math/rand documentation, Go 1.20+ auto-seeding
import "math/rand"

// In Engine struct initialization:
e := &Engine{
    rng: rand.New(rand.NewSource(time.Now().UnixNano())),
    // or for testing with deterministic seed:
    // rng: rand.New(rand.NewSource(42)),
}

// In noise application:
sigma := 0.002 * math.Sqrt(dt) // Brownian scaling
s.Energy += e.rng.NormFloat64() * sigma
```

### Clamp Helper
```go
// Source: fresh design (simpler than V1's ClampVariable with map lookup)
func clamp(v, lo, hi float64) float64 {
    if v < lo { return lo }
    if v > hi { return hi }
    return v
}

// Per-variable clamp (called after each tick):
s.Energy = clamp(s.Energy, 0, 1)
s.Stress = clamp(s.Stress, 0, 1)
// ... etc
s.BodyTemp = clamp(s.BodyTemp, 25, 43)
```

### Linear Decay Application
```go
// Source: fresh design — simpler than V1's exponential decay math
func (e *Engine) applyDecay(s *State, dt float64) {
    rate := e.config.DecayMultiplier
    // Energy drains toward 0 (exhaustion)
    s.Energy -= energyDecayRate * rate * dt
    // Hunger rises toward 1 (starvation)
    s.Hunger += hungerDecayRate * rate * dt
    // Mood drifts toward 0 (dysphoria) without positive input
    s.Mood -= moodDecayRate * rate * dt
    // Social deficit rises toward 1 (isolation)
    s.SocialDeficit += socialDecayRate * rate * dt
    // Cognitive load rises toward 1 (mental fatigue)
    s.CognitiveLoad += cogLoadDecayRate * rate * dt
}
```

### Scenario Test Pattern
```go
// Source: derived from V1 processor_test.go pattern
func TestDegradation_TenMinutesUnattended(t *testing.T) {
    s := bio.NewDefaultState()
    engine := bio.NewEngine(bio.Config{DecayMultiplier: 1.0})

    initialEnergy := s.Energy
    initialHunger := s.Hunger

    // Simulate 10 minutes at 1-second ticks
    for i := 0; i < 600; i++ {
        engine.Tick(s, 1.0)
    }

    // Energy must have visibly degraded (not negligible)
    if s.Energy >= initialEnergy-0.2 {
        t.Errorf("energy should degrade by at least 0.2 in 10 min, got %f (was %f)",
            s.Energy, initialEnergy)
    }
    // Hunger must have risen
    if s.Hunger <= initialHunger+0.2 {
        t.Errorf("hunger should rise by at least 0.2 in 10 min, got %f (was %f)",
            s.Hunger, initialHunger)
    }
    // All values still in valid ranges (BIO-06)
    if s.Energy < 0 || s.Energy > 1 { t.Errorf("energy out of range: %f", s.Energy) }
    // ... etc
}
```

### Data-Driven Rule Definition
```go
// Source: adapted from V1 interactions.go pattern
type Rule struct {
    Name      string
    Condition func(s *State) bool
    Apply     func(s *State, dt float64) Delta
}

var motivationRules = []Rule{
    {
        Name:      "stress_tension_cascade",
        Condition: func(s *State) bool { return s.Stress > 0.6 },
        Apply: func(s *State, dt float64) Delta {
            return Delta{Field: "physical_tension", Amount: s.Stress * 0.3 * dt}
        },
    },
    // ... 21 more rules
}
```

---

## State of the Art

| Old Approach (V1) | V2 Approach | Rationale | Impact |
|-------------------|-------------|-----------|--------|
| 20 physiological variables | 8 motivation-shaped proxies | Focus on drive relevance, not physiological accuracy | Simpler rules, faster ticks, clearer drive signal |
| Exponential decay (half-life math) | Linear decay (constant rate) | Visible within 3-5 min; configurable via multiplier | Easier calibration, more testable |
| Automatic homeostasis (return-to-baseline) | No automatic homeostasis | Recovery only from explicit actions | Creates necessary asymmetry for motivation pressure |
| state-as-value (NewDefaultState returns State) | state-as-pointer (*State) | Pointer from start avoids address-of-value pitfalls | Less error-prone |
| time.Now() inside Tick | dt passed as parameter | External time control | Enables deterministic tests |
| Variable enum + Get/Set (20 vars) | Direct field access (8 vars) | 8 fields don't need indirection | ~200 fewer lines of boilerplate |
| Body temp range {34,42} | Body temp range {25,43} | Cover physiologically meaningful thresholds (33°C, 35°C) | Correct threshold detection |

---

## Open Questions

1. **Module path for v2**
   - What we know: CONTEXT.md says "same module path as V1 (`github.com/marczahn/person`), V2 lives in `v2/` subdirectory"
   - What's unclear: Go convention for v2 within same repo — options are: (a) `v2/go.mod` with `module github.com/marczahn/person/v2`, or (b) a branch-based approach. Option (a) seems most consistent with "v2/ subdirectory."
   - Recommendation: Use `module github.com/marczahn/person/v2` in `v2/go.mod`. This is the standard Go major version convention and allows both v1 and v2 to coexist in the same repo.

2. **Package name inside v2**
   - What we know: V1 uses `package biology` under `internal/biology/`. V2 is a fresh design.
   - What's unclear: Should V2 biology live at `v2/bio/` (package `bio`) or `v2/internal/biology/` (package `biology`)? The `internal/` path restricts imports to within the module — less relevant since V2 is its own module, but good hygiene.
   - Recommendation: Use `v2/bio/` with package `bio` for brevity. The phase only builds the bio engine; no other packages exist in v2 yet, so `internal/` is premature.

3. **Noise scaling: per-tick vs per-second**
   - What we know: BIO-05 requires Gaussian noise each tick. Tick rate is not fixed.
   - What's unclear: Should noise be constant per-tick (same regardless of dt) or scaled by sqrt(dt) for Brownian motion?
   - Recommendation: Scale by sqrt(dt). This ensures noise magnitude is consistent regardless of tick rate — a 1s tick at sigma=0.002 produces the same expected total noise as 5 consecutive 0.2s ticks. Without this, higher tick rates produce more total noise over the same real-time period.

4. **CognitiveLoad naming vs CognitiveCapacity**
   - What we know: CONTEXT.md uses "cognitive capacity" in drive mapping but the context is ambiguous about direction (high capacity or high load?).
   - What's unclear: "Cognitive capacity" could mean remaining capacity (high=good) or used capacity (high=bad). The drive mapping "Stimulation/novelty drive ← cognitive capacity + mood" suggests high cognitive capacity → high stimulation drive, meaning capacity means available capacity, not load.
   - Recommendation: Name the variable `CognitiveCapacity` (0=depleted, 1=fresh), which decays toward 0 (BIO-03). This matches the drive mapping: high CognitiveCapacity contributes to stimulation drive. The inverse of "cognitive load."

---

## Sources

### Primary (HIGH confidence)
- V1 source code (`/home/marczahn/dev/person/v1/internal/biology/`) — direct inspection of state.go, processor.go, interactions.go, thresholds.go, circadian.go, state_test.go, processor_test.go
- V1 biologist advisory (`/home/marczahn/dev/person/v1/docs/advisory/biologist.md`) — interaction magnitudes, decay rates, threshold values
- CONTEXT.md (`/home/marczahn/dev/person/.planning/phases/01-biological-foundation/01-CONTEXT.md`) — locked decisions

### Secondary (MEDIUM confidence)
- Go stdlib `math/rand` documentation — NormFloat64() usage pattern (confirmed via training knowledge, stdlib API stable since Go 1.0)
- Go module versioning convention — `github.com/module/v2` pattern (well-established Go community standard, multiple sources)

### Tertiary (LOW confidence)
- Decay rate calibrations (specific numbers) — derived from desired behavior (3-5 minute visibility window) and V1 reference rates, not from external authoritative source. Should be validated during implementation testing.

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — pure Go stdlib, no external dependencies needed
- Variable definitions: HIGH — locked in CONTEXT.md, calibrations derivable from constraints
- Architecture patterns: HIGH — fresh design from locked constraints, informed by V1 lessons learned
- Interaction rules: MEDIUM — rules chosen by motivation relevance, specific magnitudes need test validation
- Decay rate calibrations: MEDIUM — derived from "visible in 3-5 min" constraint, specific values need empirical tuning
- Noise calibration: MEDIUM — sigma=0.002 is reasonable but untested at V2 scale
- Pitfalls: HIGH — derived from V1 bugs (body temp range, cortisol urgency, dt capping) and Go patterns

**Research date:** 2026-02-19
**Valid until:** 2026-03-19 (stable domain — Go stdlib doesn't change, bio model is custom logic)

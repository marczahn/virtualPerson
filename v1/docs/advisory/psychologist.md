
# Psychological Layer: Concrete Implementation Guide

This is a substantial design document. I'll work through each of your seven questions with implementable models, but I also need to push back on some assumptions along the way.

---

## 1. Emotional Response Patterns: Biology-to-Emotion Mapping

### Core Mapping Table

Here are the most empirically grounded biological-state-to-emotion mappings. I'm giving primary and secondary responses because emotions rarely occur in isolation.

| Biological State Combination | Primary Emotion | Secondary Emotion | Behavioral Tendency |
|---|---|---|---|
| High cortisol + low serotonin + fatigue | Irritability / dysphoria | Hopelessness, withdrawal | Social avoidance, reduced initiative, snapping at stimuli |
| High cortisol + low serotonin + normal energy | Anxiety | Agitation, worry | Hypervigilance, rumination, restless movement |
| High adrenaline + high dopamine | Excitement / exhilaration | Overconfidence, impulsivity | Risk-taking, rapid speech, approach behavior |
| High adrenaline + low dopamine | Fear / panic | Helplessness | Fight-or-flight, freezing, escape behavior |
| Low blood sugar + high fatigue + circadian night | Irritability | Sadness, emotional fragility | Reduced frustration tolerance, crying more easily, craving comfort |
| High serotonin + moderate dopamine + rested | Contentment | Openness, sociability | Prosocial behavior, curiosity, patience |
| High cortisol + high adrenaline + high dopamine | Angry determination | Moral outrage | Confrontation, persistence, tunnel vision |
| Low cortisol + low adrenaline + low dopamine | Apathy / flatness | Boredom, disconnection | Withdrawal, difficulty initiating, lethargy |

### A Critical Pushback

The mapping between biology and emotion is not deterministic. The same cortisol spike can produce anxiety in one context and excited anticipation in another. The Schachter-Singer two-factor theory is relevant here: physiological arousal provides intensity, but cognitive appraisal provides the label. Your consciousness layer (the LLM) should not receive a pre-labeled emotion. It should receive the biological state and context, then generate the emotional interpretation. That is more psychologically accurate and also architecturally cleaner — it keeps the biological layer from doing the consciousness layer's job.

What I would actually recommend for your system:

```
// Don't pass this:
{ emotion: "anxious", intensity: 0.7 }

// Pass this instead:
{
  arousal: 0.8,          // derived from adrenaline, heart_rate, cortisol
  valence_bias: -0.3,    // derived from serotonin, dopamine balance
  energy: 0.3,           // derived from fatigue, blood_sugar, circadian
  cognitive_load: 0.6,   // derived from cortisol duration, sleep_debt
  context: { ... }       // what stimuli are present
}
```

The LLM then interprets these dimensions into named emotions. This is how human emotion actually works — we construct emotions from bodily signals plus interpretation.

### Deriving the Four Dimensions

```
arousal = clamp(0, 1,
  0.3 * normalize(adrenaline) +
  0.25 * normalize(heart_rate) +
  0.2 * normalize(cortisol) +
  0.15 * normalize(norepinephrine) +
  0.1 * (1 - normalize(fatigue))
)

valence_bias = clamp(-1, 1,
  0.35 * normalize(serotonin) +
  0.30 * normalize(dopamine) +
  0.15 * normalize(endorphins) +
  -0.20 * normalize(cortisol)
)
// Positive = pleasant bias, Negative = unpleasant bias

energy = clamp(0, 1,
  0.30 * (1 - normalize(fatigue)) +
  0.25 * normalize(blood_sugar) +
  0.25 * circadian_energy_curve(time_of_day) +
  0.20 * normalize(dopamine)
)

cognitive_load = clamp(0, 1,
  0.30 * stress_duration_factor(cortisol_history) +
  0.25 * normalize(fatigue) +
  0.25 * (1 - normalize(blood_sugar)) +
  0.20 * sleep_debt_factor(sleep_history)
)
```

### Individual Variation

Inter-individual variation in emotion is large. Research on affective reactivity suggests roughly a 2x to 3x range in emotional intensity for the same physiological trigger across the normal population. For your simulation, personality traits (see section 2) should modulate these mappings, not random noise. Random noise is the lazy approach. Real variation comes from stable individual differences in trait sensitivity, not from rolling dice on each reaction.

---

## 2. Personality Model

### Is Big Five the Right Model?

Yes, but with modifications. The Big Five (OCEAN) is the most empirically validated personality framework and the best starting point. However, for simulation purposes, you need to supplement it with two additions:

**Trait sensitivity to specific stimuli** — The Big Five tells you general tendencies, but not domain-specific reactivity. Someone can be generally low-neuroticism but have a specific vulnerability to abandonment. This is where emotional memory (section 4) interacts with personality.

**Motivational drives** — The Big Five describes reaction patterns but not what the person wants. You need a small set of core motives: safety/security, connection/belonging, competence/mastery, autonomy, novelty/stimulation. Weight these per individual.

### Implementation Model

Each trait is a value from 0.0 to 1.0 (population mean = 0.5, standard deviation ≈ 0.15 for normal distribution):

```
personality = {
  openness: 0.0-1.0,
  conscientiousness: 0.0-1.0,
  extraversion: 0.0-1.0,
  agreeableness: 0.0-1.0,
  neuroticism: 0.0-1.0
}
```

### Concrete Trait Modulation Effects

**Neuroticism** modulates emotional reactivity to negative stimuli:

```
negative_emotion_multiplier = 1.0 + (neuroticism - 0.5) * 1.2

// Neuroticism 0.8: multiplier = 1.36 (36% stronger negative reactions)
// Neuroticism 0.2: multiplier = 0.64 (36% weaker negative reactions)
```

High neuroticism + cold stimulus: The person feels the cold more distressing, generates more catastrophic interpretations ("this could be dangerous"), is slower to habituate, and the negative emotional memory encodes more strongly. Low neuroticism + cold stimulus: The person registers discomfort but does not elaborate on it, adapts faster, is less likely to ruminate about it afterward.

**Extraversion** modulates response to social and stimulation variables:

```
isolation_distress_rate = base_rate * (1.0 + (extraversion - 0.5) * 1.5)
stimulation_seeking = 0.3 + extraversion * 0.5
social_reward_sensitivity = 0.2 + extraversion * 0.6

// High extraversion (0.85): isolation distress accumulates 52% faster
// Low extraversion (0.2):  isolation distress accumulates 45% slower
```

High extraversion + isolation: Distress builds faster, the person generates thoughts about missing people, becomes restless, seeks any form of stimulation. Behavioral outputs include pacing, talking to themselves, heightened response to any social input. Low extraversion + isolation: May actually find short-term isolation restorative. Distress still builds but at a much slower rate, with a longer threshold before psychological effects emerge.

**Agreeableness** modulates conflict response:

```
conflict_response = {
  if agreeableness > 0.7:
    primary: "accommodate/appease"
    stress_from_conflict: high (conflict itself is distressing)
    anger_suppression: high
    delayed_resentment: moderate  // important side effect

  if agreeableness 0.3-0.7:
    primary: "negotiate/assert"
    stress_from_conflict: moderate
    balanced approach

  if agreeableness < 0.3:
    primary: "compete/dominate"
    stress_from_conflict: low (conflict is energizing)
    anger_expression: high
    relationship_damage_risk: high
}
```

Low agreeableness + conflict stimulus: The person engages rather than withdraws. Anger comes easily and is expressed outwardly. They may experience the conflict as stimulating rather than distressing. They are less likely to feel guilt afterward and more likely to feel vindicated.

**Conscientiousness** modulates response to disorder and self-regulation:

```
disorder_tolerance = 1.0 - conscientiousness * 0.8
self_regulation_bonus = conscientiousness * 0.3  // added to regulation_capacity
planning_under_stress = conscientiousness * 0.4   // ability to maintain structured coping
```

**Openness** modulates cognitive flexibility and stimulus interpretation:

```
reappraisal_ability = 0.2 + openness * 0.5  // ability to reframe experiences
novelty_as_threat = 1.0 - openness * 0.7    // how much novelty triggers anxiety vs curiosity
```

### Should Traits Evolve?

Personality traits should be highly stable on the timescale of your simulation unless you are simulating years. In real humans, Big Five traits shift by approximately 0.1-0.2 standard deviations per decade, and mostly in predictable directions (neuroticism decreases, conscientiousness and agreeableness increase with age). Meaningful personality change from single experiences is rare outside of genuinely traumatic events.

My recommendation: make traits fixed for any simulation under a few months of simulated time. If you simulate years, apply a very slow drift:

```
// Per simulated year, if applicable
trait_drift_per_year = {
  neuroticism: -0.005,
  conscientiousness: +0.005,
  agreeableness: +0.003,
  extraversion: -0.002,
  openness: -0.002
}

// Trauma can cause a one-time shift:
if (event.trauma_severity > 0.8) {
  neuroticism += 0.02 * event.trauma_severity
}
```

Do not make traits reactive to daily events. That is mood, not personality. Conflating the two would be a fundamental modeling error.

---

## 3. Coping Mechanisms

### Taxonomy and Decision Model

Here are the main coping strategies, organized by type, with concrete implementation parameters.

#### Problem-Focused Coping

**Problem-Solving**
- Triggers: Moderate stress (0.3-0.7), identifiable cause, adequate cognitive resources (cognitive_load < 0.6)
- Effectiveness: High when the stressor is controllable (distress reduction: 40-60%)
- Side effects: None inherent, but fails and increases frustration when stressor is uncontrollable
- Personality bias: +conscientiousness, +openness
- Required: energy > 0.3, cognitive_load < 0.6

**Information-Seeking**
- Triggers: Uncertainty + moderate anxiety, cause unclear
- Effectiveness: Moderate (distress reduction: 20-30%, but can increase distress if information is threatening)
- Side effects: Can become compulsive reassurance-seeking
- Personality bias: +openness, +conscientiousness

#### Emotion-Focused Coping

**Cognitive Reappraisal** (reframing the meaning)
- Triggers: Negative emotion detected, adequate cognitive resources
- Effectiveness: High (distress reduction: 30-50%), the gold standard of emotion regulation
- Side effects: Minimal when genuine; can become intellectualization if overused
- Personality bias: +openness, high emotional regulation capacity
- Required: cognitive_load < 0.5, regulation_capacity > 0.4

**Acceptance**
- Triggers: Stressor is uncontrollable, sustained stress, problem-solving has failed
- Effectiveness: Moderate-to-high over time (distress reduction: 20-40%), slow-acting
- Side effects: Can become passivity if misapplied to controllable stressors
- Personality bias: +agreeableness, +openness, lower neuroticism

**Distraction**
- Triggers: Moderate-high stress, limited agency, stimulation available
- Effectiveness: Short-term high (distress reduction: 30-50% immediate), long-term low (distress returns)
- Side effects: Avoidance of processing; problems remain unresolved
- Personality bias: +extraversion

**Emotional Suppression**
- Triggers: Social context demands composure, high self-monitoring
- Effectiveness: Low-to-moderate behaviorally, psychologically costly (distress reduction: 10-20% apparent, physiological stress actually increases)
- Side effects: Increased sympathetic arousal, memory impairment, rebound effects
- Personality bias: +conscientiousness, -extraversion (introverts suppress more)

**Rumination**
- Triggers: High neuroticism, unresolved negative event, low distraction availability, circadian night
- Effectiveness: Negative (distress increase: +20-40%). This is a maladaptive coping pattern, not a solution, but it is extremely common.
- Side effects: Sustained cortisol elevation, mood deterioration, insomnia
- Personality bias: +neuroticism strongly, -extraversion

**Denial / Avoidance**
- Triggers: Overwhelming stress (> 0.8), threat too large to process, limited resources
- Effectiveness: Short-term protective (prevents psychological collapse), long-term maladaptive
- Side effects: Delayed processing, sudden breakdown when denial fails
- Personality bias: Not strongly trait-dependent; more resource-dependent

**Catastrophizing**
- Triggers: High stress + high neuroticism + low perceived control
- Effectiveness: Negative (distress increase: +30-60%). Also maladaptive, but common.
- Side effects: Panic, paralysis, physiological escalation
- Personality bias: +neuroticism strongly

### Coping Selection Decision Tree

```
function select_coping(stress, personality, context):

  available_resources = {
    cognitive: 1.0 - cognitive_load,
    energy: energy,
    regulation: regulation_capacity
  }

  // Under extreme stress, higher-order coping collapses
  if stress > 0.85:
    if available_resources.cognitive < 0.2:
      return DENIAL  // psychological circuit breaker
    if personality.neuroticism > 0.7:
      return CATASTROPHIZING (60%) or RUMINATION (40%)

  // Is the stressor controllable?
  if context.stressor_controllability > 0.5:
    if available_resources.cognitive > 0.4 AND available_resources.energy > 0.3:
      if personality.conscientiousness > 0.5:
        return PROBLEM_SOLVING
      else:
        return PROBLEM_SOLVING (60%) or DISTRACTION (40%)
    else:
      return DISTRACTION  // want to solve it but can't right now

  // Stressor is uncontrollable
  if available_resources.cognitive > 0.4:
    if personality.openness > 0.5 AND regulation_capacity > 0.4:
      return REAPPRAISAL
    if personality.agreeableness > 0.6:
      return ACCEPTANCE
    return DISTRACTION

  // Low resources, uncontrollable stressor
  if personality.neuroticism > 0.6:
    return RUMINATION (70%) or SUPPRESSION (30%)
  else:
    return DISTRACTION (50%) or ACCEPTANCE (50%)
```

This is a simplified decision tree. In reality, people use multiple coping strategies simultaneously and switch between them. For your simulation, I would recommend selecting a primary and secondary strategy and blending their effects.

---

## 4. Emotional Memory

### The Core Model

Emotional memory is not a simple counter. It operates through associative networks — stimuli get linked to emotional responses through experience. Here is a workable model.

```
emotional_memory = {
  stimulus_type: "cold",
  associations: [
    { valence: -0.7, intensity: 0.8, recency: 0.9, context: "blizzard" },
    { valence: -0.3, intensity: 0.4, recency: 0.5, context: "chilly_morning" },
    { valence: +0.5, intensity: 0.6, recency: 0.3, context: "snow_with_friend" }
  ]
}
```

### Emotional Intensity Calculation

```
emotional_response_modifier = weighted_average(
  for each memory in associations:
    weight = memory.intensity * recency_decay(memory.recency) * emotional_salience(memory.valence)
    value = memory.valence
)

// Negative memories are weighted more heavily (negativity bias, Baumeister et al.)
emotional_salience(valence):
  if valence < 0: return abs(valence) * 1.5  // negative memories ~1.5x more salient
  if valence > 0: return valence * 1.0

// Recency decay follows a power law, not exponential
recency_decay(time_since_event):
  return 1.0 / (1.0 + time_since_event_in_days ^ 0.5)
  // Day 0: 1.0
  // Day 1: 0.59
  // Day 4: 0.33
  // Day 9: 0.25
  // Day 25: 0.17
  // Day 100: 0.09
```

### Your Proposed Formula

You suggested: `emotional_intensity = base_intensity * (1 + negative_memory_count * 0.1)`

This is too simplistic and will produce unrealistic behavior. Problems:
- It only counts, it doesn't weight by intensity or recency
- It has no ceiling, so 20 mildly bad experiences produce a stronger reaction than 1 traumatic one
- It ignores positive counter-conditioning entirely

Better formula:

```
emotional_intensity = base_intensity * (1.0 + memory_modifier)

memory_modifier = clamp(-0.5, 1.0,
  sum(negative_memory_weights) * negativity_bias - sum(positive_memory_weights)
) * neuroticism_scaling

neuroticism_scaling = 0.7 + personality.neuroticism * 0.6
// Low neuroticism (0.2): 0.82 — memories affect reactions less
// High neuroticism (0.8): 1.18 — memories affect reactions more
```

### Positive Conditioning

Yes, positive associations with a stimulus genuinely reduce negative reactions. This is the basis of counter-conditioning in behavioral therapy. If cold was associated with a pleasant memory (building a snowman with someone you love), the positive memory competes with the negative ones in the weighted average. This is not just subtraction — a strongly encoded positive memory can fundamentally change the emotional quality of a stimulus from threat to nostalgia.

### Trauma

Traumatic memories (very high intensity, very negative valence) have different properties:
- They decay much more slowly (use exponent 0.3 instead of 0.5 in the decay function)
- They can be triggered by partial cues (lower similarity threshold for activation)
- They can intrude involuntarily (probability of spontaneous activation proportional to intensity and recency)

```
if memory.intensity > 0.85 AND memory.valence < -0.7:
  // Traumatic memory
  recency_decay_exponent = 0.3  // slower decay
  intrusion_probability = memory.intensity * recency_decay(time) * 0.1
  // Per thought-cycle, small chance of involuntary recall
```

---

## 5. Cognitive Distortions Under Stress

### Activation Thresholds

Cognitive distortions are not binary on/off. They increase in probability as stress rises and cognitive resources deplete. The general relationship:

```
distortion_probability = base_rate * stress_multiplier * trait_multiplier

stress_multiplier:
  stress 0.0-0.3: 1.0  (minimal distortion increase)
  stress 0.3-0.5: 1.0 + (stress - 0.3) * 2.5  (gradual onset)
  stress 0.5-0.7: 1.5 + (stress - 0.5) * 3.0  (accelerating)
  stress 0.7-1.0: 2.1 + (stress - 0.7) * 4.0  (distortions dominate thinking)

// At stress 0.9: multiplier ≈ 2.9
```

### Specific Distortions

**Catastrophizing** — Interpreting events as maximally terrible, assuming worst outcomes.

```
catastrophizing:
  base_rate: 0.05
  trait_multiplier: 1.0 + neuroticism * 2.0 + (1 - openness) * 0.5
  primary_trigger: uncertainty + perceived threat
  onset: stress > 0.4 for high neuroticism, stress > 0.6 for low neuroticism
  effect: perceived_threat_magnitude *= 2.0-4.0
  example: "My hands are cold" → "I could get frostbite, I could lose my fingers"
```

**Black-and-White Thinking** (Splitting) — Evaluating in extremes, losing nuance.

```
splitting:
  base_rate: 0.04
  trait_multiplier: 1.0 + neuroticism * 1.5 + (1 - openness) * 1.0
  primary_trigger: cognitive_load > 0.6 (nuance requires cognitive resources)
  onset: cognitive_load > 0.5
  effect: evaluations become binary (safe/dangerous, good/bad, possible/impossible)
  example: "This situation is uncomfortable" → "This situation is unbearable"
```

**Personalization** — Attributing external events to oneself.

```
personalization:
  base_rate: 0.03
  trait_multiplier: 1.0 + neuroticism * 1.5 + agreeableness * 0.8
  primary_trigger: negative events + social context or self-evaluation context
  onset: stress > 0.4 + relevant social/self context
  effect: causal attribution shifts to self ("this is my fault")
  example: "The situation is bad" → "I caused this / I deserve this"
```

**Emotional Reasoning** — Treating feelings as evidence of reality.

```
emotional_reasoning:
  base_rate: 0.06  // most common distortion
  trait_multiplier: 1.0 + neuroticism * 1.0 + (1 - conscientiousness) * 0.8
  primary_trigger: strong emotion + ambiguous situation
  onset: arousal > 0.5 + situational ambiguity
  effect: "I feel X therefore X is true"
  example: "I feel afraid, therefore I am in danger" (even if objectively safe)
```

**Overgeneralization** — One instance becomes a universal rule.

```
overgeneralization:
  base_rate: 0.04
  trait_multiplier: 1.0 + neuroticism * 1.5
  primary_trigger: repeated negative experience OR single very intense negative experience
  onset: stress > 0.5 + relevant negative memory
  effect: "always", "never", "everything" language in thoughts
  example: "I was cold yesterday and I'm cold now" → "I'm always going to be cold"
```

### Implementation Approach

Each thought cycle, calculate the probability of each distortion being active. If active, the distortion modifies the prompt context sent to the LLM consciousness layer:

```
active_distortions = []
for distortion in all_distortions:
  p = distortion.base_rate * stress_multiplier * distortion.trait_multiplier
  if random() < p:
    active_distortions.append(distortion)

// Then in the LLM prompt:
// "Your current cognitive biases include: [catastrophizing, emotional_reasoning].
//  This means you are tending to assume worst outcomes and treating your
//  feelings as evidence of reality."
```

---

## 6. Emotional Regulation

### Can Someone Be Biologically Stressed but Psychologically Calm?

Absolutely yes. This is the entire point of emotional regulation, and it is one of the most important things to model. A trained meditator can have elevated cortisol (biological stress response active) while maintaining psychological equanimity. A soldier in combat can have extreme physiological arousal while making calm tactical decisions. The biological signal is not the emotion — regulation sits between them.

### Regulation Capacity Model

```
regulation_capacity = clamp(0, 1,
  baseline_regulation
  - stress_depletion
  - fatigue_penalty
  + personality_bonus
)

baseline_regulation:
  // Trait-level capacity, set at initialization
  = 0.3 + personality.conscientiousness * 0.2
        + personality.openness * 0.15
        + (1 - personality.neuroticism) * 0.2
        + learned_regulation_skill * 0.15
  // Range: roughly 0.3 to 0.85 depending on personality and learning

stress_depletion:
  // Regulation is a depletable resource (ego depletion model)
  // Controversial in psychology but useful for simulation
  = sustained_stress_integral * 0.3
  // sustained_stress_integral = average stress over last N hours * hours
  // 4 hours at stress 0.6 = 2.4 → depletion = 0.72

  // BUT: it recovers with rest
  recovery_rate = 0.1 per hour of low-stress (< 0.3) waking time
  recovery_rate = 0.2 per hour of sleep

fatigue_penalty:
  = fatigue * 0.3  // tired people regulate worse

personality_bonus:
  = conscientiousness * 0.1  // structured people maintain regulation longer
```

### What Regulation Actually Does

When regulation capacity is high, the person can:
1. Attenuate the emotional intensity signal (multiply arousal by `1.0 - regulation_capacity * 0.5`)
2. Shift valence toward neutral (move valence_bias toward 0 by `regulation_capacity * 0.3`)
3. Reduce cognitive distortion probability (multiply distortion base_rates by `1.0 - regulation_capacity * 0.6`)
4. Enable higher-order coping strategies (reappraisal requires regulation_capacity > 0.4)

When regulation capacity is depleted:
1. Raw biological signals pass through with minimal dampening
2. Emotional reactions become more volatile
3. Coping degrades to lower-order strategies (rumination, suppression, denial)
4. Cognitive distortions increase significantly

```
// Applied to the emotional signal before it reaches consciousness:
regulated_arousal = arousal * (1.0 - regulation_capacity * 0.5)
regulated_valence = valence_bias * (1.0 - regulation_capacity * 0.3)

// With high regulation (0.8):
//   arousal 0.9 → 0.54 (experienced as manageable activation)
//   valence -0.7 → -0.49 (experienced as less negative)

// With depleted regulation (0.1):
//   arousal 0.9 → 0.855 (barely dampened)
//   valence -0.7 → -0.679 (almost raw negative signal)
```

### Regulation Degradation Under Sustained Stress

This follows a non-linear curve. People can maintain regulation for a while, then it collapses relatively quickly:

```
// Regulation depletion over time at constant stress level
depletion_curve(hours_at_stress, stress_level):
  raw_depletion = stress_level * hours_at_stress * 0.08
  // Accelerating collapse after threshold
  if raw_depletion > 0.6:
    raw_depletion += (raw_depletion - 0.6) ^ 2
  return min(raw_depletion, baseline_regulation)

// Example: stress 0.7 for 12 hours
// raw_depletion = 0.7 * 12 * 0.08 = 0.672
// exceeds 0.6, so: 0.672 + (0.072)^2 = 0.677
// A person with baseline 0.7 is nearly fully depleted
```

---

## 7. Social Cognition in Isolation

### Timeline of Isolation Effects

This is one of the most well-studied areas in psychology (solitary confinement research, sensory deprivation studies, Antarctic station studies). Here is the timeline:

```
isolation_effects(duration_hours, personality):

  extraversion_factor = 1.0 + (personality.extraversion - 0.5) * 1.5
  neuroticism_factor = 1.0 + (personality.neuroticism - 0.5) * 1.0

  hours_0_to_2:
    // Most people are fine. Introverts may be relieved.
    loneliness = 0.05 * extraversion_factor
    cognitive_effects = none
    behavioral: normal, possibly relaxed

  hours_2_to_8:
    // Boredom and mild restlessness
    loneliness = 0.1 + (hours - 2) / 6 * 0.15 * extraversion_factor
    boredom = 0.2 + (hours - 2) / 6 * 0.3
    behavioral: seeking stimulation, talking to self (occasional)

  hours_8_to_24:
    // Genuine loneliness begins
    loneliness = 0.25 + (hours - 8) / 16 * 0.25 * extraversion_factor
    rumination_probability increases (especially nighttime)
    behavioral: increased self-talk, reviewing memories,
                heightened emotional response to any stimulus

  hours_24_to_72 (days 1-3):
    // Significant psychological effects
    loneliness = 0.5 + (hours - 24) / 48 * 0.2 * extraversion_factor
    cognitive: time distortion, difficulty concentrating,
              increased daydreaming / fantasy
    emotional: mood swings, irritability, sadness waves
    behavioral: creating routines compulsively, anthropomorphizing objects,
                increased self-talk

  hours_72_to_168 (days 3-7):
    // Destabilization for vulnerable individuals
    loneliness = 0.7 + (hours - 72) / 96 * 0.15
    cognitive: identity confusion (\"who am I without others?\"),
              paranoid ideation (neuroticism > 0.7),
              magical thinking
    emotional: emotional flattening OR emotional storms (bimodal)
    destabilization_risk = neuroticism_factor * extraversion_factor * 0.15
    behavioral: rigid ritualization, talking to imagined others,
                sleep disturbance

  hours_168+ (weeks):
    // Severe effects for most people
    loneliness = 0.85+
    cognitive: hallucinations possible (especially auditory),
              identity dissolution, severe time distortion,
              hypersensitivity to any stimulus
    psychological_damage_accumulating = true
```

### Personality Vulnerability

Most vulnerable to isolation effects (from most to least):
1. **High extraversion + high neuroticism** — Needs stimulation AND prone to distress. Worst combination.
2. **High extraversion + low neuroticism** — Needs stimulation but copes better emotionally. Still distressed but more resilient.
3. **Low extraversion + high neuroticism** — Does not need social stimulation as much but prone to rumination in the quiet. Moderate risk, different symptom profile (more rumination, less restlessness).
4. **Low extraversion + low neuroticism** — Most resilient. Can tolerate extended isolation with minimal psychological impact. Still affected eventually, but much later.

```
isolation_resilience = (1 - extraversion) * 0.5 + (1 - neuroticism) * 0.3
                       + conscientiousness * 0.2
// Range: roughly 0.1 to 0.9
// Threshold where isolation becomes destabilizing:
destabilization_threshold_hours = 48 + isolation_resilience * 120
// Low resilience (0.1): ~60 hours
// High resilience (0.9): ~156 hours
```

---

## Architectural Recommendation

Having laid out all seven areas, here is how I would structure this as a layer in your system.

The psychological layer should be a **transformer**, not a generator. It takes biological state as input and produces a modified psychological state that the LLM consciousness layer consumes. It should NOT generate thoughts or emotions itself — that is the consciousness layer's job.

```
BiologicalState
    ↓
[Psychological Layer]
  ├── Compute arousal, valence, energy, cognitive_load (Section 1)
  ├── Apply personality trait modulation (Section 2)
  ├── Apply emotional memory modifiers (Section 4)
  ├── Apply emotional regulation (Section 6)
  ├── Compute active cognitive distortions (Section 5)
  ├── Select likely coping strategies (Section 3)
  ├── Compute isolation effects if applicable (Section 7)
    ↓
PsychologicalState {
  arousal: float,
  valence: float,
  energy: float,
  cognitive_load: float,
  regulation_capacity: float,
  active_distortions: [string],
  likely_coping: [string],
  emotional_memory_activations: [{stimulus, valence, intensity}],
  isolation_effects: {loneliness, cognitive_effects, behavioral_tendencies}
}
    ↓
[LLM Consciousness Layer]
  → Interprets all of the above into subjective experience, thoughts, emotions
```

The LLM should receive the PsychologicalState as structured context and generate the first-person experience from it. This keeps the system modular: you can tune the psychological layer's rules without changing the LLM, and you can swap LLMs without losing the psychological model.

### One Final Pushback

Your question framing assumes a clean separation between "biology produces states" and "psychology labels them." That is roughly correct, but there is a crucial feedback loop you must implement: **psychological states modify biological states**. Rumination sustains cortisol. Catastrophizing triggers adrenaline release. Acceptance reduces cortisol. If you do not model this feedback loop, your person will be psychologically reactive but never psychologically generative — they will respond to biology but never drive it. Real humans spiral (anxiety increases cortisol increases anxiety) and self-regulate (acceptance reduces cortisol reduces distress). Both directions of causation are essential.

```
// Psychological → Biological feedback (per cycle)
if active_coping == RUMINATION:
  cortisol += 0.02 per cycle
  serotonin -= 0.01 per cycle

if active_coping == CATASTROPHIZING:
  adrenaline += 0.03 per cycle
  cortisol += 0.02 per cycle

if active_coping == ACCEPTANCE or REAPPRAISAL:
  cortisol -= 0.01 per cycle
  serotonin += 0.005 per cycle

if active_distortions.length > 2:
  cortisol += 0.01 per cycle  // distorted thinking is itself stressful
```

Without this, you don't have a person. You have a mood ring.

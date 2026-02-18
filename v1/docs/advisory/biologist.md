# Biological Simulation Model Review

This is a thorough review. I will be direct about what works, what is missing, and where your current model has problems.

---

## 1. Variable Sufficiency

### What you have is a solid foundation, but there are real gaps.

**Keep all 16 variables.** None are unnecessary for moderate complexity. However, several critical physiological systems are unrepresented.

### Critically Missing Variables

| Variable | Baseline | Range | Unit | Why it matters |
|---|---|---|---|---|
| **Blood oxygen saturation (SpO2)** | 98 | 70-100 | % | Respiratory rate exists but has no functional output. Without SpO2, you cannot model hypoxia, altitude effects, respiratory failure, or the reason the body increases respiratory rate in the first place. This is the single biggest gap. |
| **Hydration level** | 0.8 | 0-1 | ratio | Thirst is a sensation, not a physiological state. You need the actual fluid status. Dehydration affects heart rate, blood pressure, cognition, body temperature regulation, and kidney function. Thirst should *follow* hydration level, not be the primary variable. |
| **Core energy reserves (glycogen)** | 0.7 | 0-1 | ratio | Blood sugar is a snapshot, but the body has reserves. Glycogen in liver and muscle buffers blood sugar. Without it, your blood sugar model will be unrealistically volatile. When glycogen depletes, the body shifts to fat metabolism and fatigue increases sharply. |
| **Endorphin level** | 0.1 | 0-1 | ratio | You have pain but no endogenous pain modulation system. Endorphins are released during sustained exercise, injury, and stress. They suppress pain perception, create mild euphoria, and interact with your existing dopamine and serotonin systems. Without this, pain is a one-way street with no biological dampening mechanism. |

**That brings you to 20 variables.** This is still within moderate complexity and each one fills a functional gap that would otherwise require ugly hacks.

### Variables that need reframing

- **Hunger** should be a *derived signal* from blood sugar and glycogen, not an independent variable that "increases over time." Real hunger is driven by ghrelin (stomach-empty signal) and falling blood sugar. If you keep it as a standalone variable, at minimum couple it tightly to blood sugar — hunger should spike when blood sugar drops below ~75 mg/dL and be suppressed above ~110 mg/dL.

- **Thirst** similarly should be derived from hydration level. If hydration drops below 0.7, thirst rises. If hydration is above 0.85, thirst is near zero.

---

## 2. Variable Interactions — Concrete Map

This is the core of your simulation. I am going to be specific about direction, magnitude, and timescale. All multipliers are per-simulation-tick suggestions assuming a 1-second tick; scale accordingly.

### Body Temperature

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| Temp < 35.5°C | Muscle tension += (35.5 - temp) * 0.15 | Shivering onset, scales linearly | Immediate |
| Temp < 35.5°C | Heart rate += (35.5 - temp) * 8 | Compensatory tachycardia | Seconds |
| Temp < 35.5°C | Cortisol += 0.002/sec | Cold stress response | Minutes |
| Temp < 35.0°C | Adrenaline += 0.005/sec | Sympathetic activation | Seconds |
| Temp < 33.0°C | Heart rate starts *decreasing* (override above) | Bradycardia of severe hypothermia — reversal point | Minutes |
| Temp > 38.5°C | Heart rate += (temp - 38.5) * 10 | ~10 bpm per degree of fever | Seconds |
| Temp > 38.5°C | Immune response += 0.001/sec | Fever is functional — enhances immune activity | Hours |
| Temp > 39.5°C | Fatigue += 0.003/sec | Metabolic cost of fever | Minutes |
| Temp > 40.5°C | Blood pressure begins dropping | Vasodilation overwhelms compensation | Minutes |

### Heart Rate

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| HR > 100 | Blood pressure sys += (HR - 100) * 0.3 | Linear coupling | Seconds |
| HR > 100 | Respiratory rate += (HR - 100) * 0.08 | Oxygen demand drives ventilation | Seconds |
| HR > 100 | Blood sugar -= 0.02/sec (additional) | Increased metabolic demand | Minutes |
| HR > 150 sustained > 5 min | Fatigue += 0.005/sec | Cardiovascular strain | Minutes |
| HR > 170 | SpO2 may *decrease* if respiratory rate cannot keep up | Oxygen debt | Seconds-minutes |
| HR < 50 | Blood pressure -= (50 - HR) * 0.5 | Insufficient cardiac output | Seconds |
| HR < 50 | SpO2 -= 0.001/sec | Reduced perfusion | Minutes |

### Blood Pressure (Systolic)

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| BP < 90 | HR += (90 - BP) * 0.5 | Baroreceptor reflex — compensatory tachycardia | Seconds |
| BP < 90 | Fatigue += 0.002/sec | Reduced perfusion to tissues | Minutes |
| BP < 80 | SpO2 -= 0.003/sec | Inadequate organ perfusion | Minutes |
| BP < 70 | Cognitive impairment (apply to decision-making if modeled) | Shock territory | Immediate |
| BP > 160 | Pain += 0.001/sec (headache) | Hypertensive headache | Minutes-hours |

### Respiratory Rate

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| RR > 25 | SpO2 recovery rate increased by factor 1.5 | Hyperventilation compensates | Seconds |
| RR > 35 | Fatigue += 0.002/sec | Respiratory muscle fatigue | Minutes |
| RR > 35 | Muscle tension += 0.001/sec (accessory muscles) | Labored breathing | Minutes |
| RR < 10 | SpO2 -= 0.005/sec | Hypoventilation | Seconds |
| RR < 8 | SpO2 -= 0.01/sec | Respiratory failure territory | Seconds |

### Blood Sugar

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| BS < 70 | Adrenaline += 0.003/sec | Counter-regulatory hormone release | Minutes |
| BS < 70 | Cortisol += 0.001/sec | Counter-regulatory response | Minutes |
| BS < 70 | Hunger = min(1.0, hunger + 0.005/sec) | Strong hunger signal | Minutes |
| BS < 60 | Fatigue += 0.004/sec | Neuroglycopenia onset | Minutes |
| BS < 55 | Muscle tension -= 0.002/sec (weakness, not tension) | Muscle weakness from energy deficit | Minutes |
| BS < 50 | Heart rate becomes erratic (add random +-5 bpm jitter) | Autonomic instability | Seconds |
| BS > 140 | Thirst += 0.002/sec | Osmotic thirst from hyperglycemia | Minutes |
| BS > 160 | Fatigue += 0.001/sec | Hyperglycemic malaise | Hours |

### Cortisol

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| Cortisol > 0.4 sustained > 30 min | Immune response -= 0.001/sec | Immunosuppression | Hours |
| Cortisol > 0.5 | Blood sugar += 0.05/sec | Gluconeogenesis — cortisol raises blood sugar | Minutes |
| Cortisol > 0.5 | Serotonin -= 0.0005/sec | Chronic stress depletes serotonin | Hours |
| Cortisol > 0.6 sustained > 2 hrs | Pain sensitivity *= 1.2 | Hyperalgesia from sustained stress | Hours |
| Cortisol > 0.3 | Heart rate += (cortisol - 0.3) * 15 | Stress tachycardia | Seconds |
| Cortisol > 0.3 | Blood pressure += (cortisol - 0.3) * 20 | Stress hypertension | Seconds |

### Adrenaline

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| Adrenaline > 0.2 | Heart rate += adrenaline * 60 | Strong chronotropic effect | Seconds |
| Adrenaline > 0.2 | Blood pressure += adrenaline * 40 | Vasoconstriction | Seconds |
| Adrenaline > 0.2 | Respiratory rate += adrenaline * 10 | Bronchodilation + drive | Seconds |
| Adrenaline > 0.1 | Blood sugar += adrenaline * 0.5/sec | Glycogenolysis | Seconds-minutes |
| Adrenaline > 0.3 | Pain *= (1 - adrenaline * 0.5) | Stress-induced analgesia | Seconds |
| Adrenaline > 0.3 | Muscle tension += adrenaline * 0.3 | Fight-or-flight muscle readiness | Seconds |
| Adrenaline > 0.5 | Digestion paused (hunger accumulation stops) | Blood diverted from gut | Seconds |
| Adrenaline > 0.7 sustained > 10 min | Fatigue += 0.008/sec | Adrenal crash approaching | Minutes |

### Serotonin

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| Serotonin < 0.3 | Pain sensitivity *= 1.3 | Serotonin modulates descending pain inhibition | Hours |
| Serotonin < 0.2 | Sleep quality reduced (fatigue recovery rate * 0.6) | Serotonin required for melatonin synthesis | Hours |
| Serotonin < 0.2 | Cortisol baseline shifts up by +0.05 | Loss of HPA axis regulation | Hours |
| Serotonin > 0.6 | Pain sensitivity *= 0.85 | Enhanced pain gate | Hours |
| Serotonin > 0.7 | Muscle tension *= 0.9 | Relaxation effect | Hours |

### Dopamine

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| Dopamine > 0.5 | Pain *= 0.8 | Reward-based analgesia | Minutes |
| Dopamine > 0.6 | Heart rate += (dopamine - 0.6) * 20 | Excitement/arousal tachycardia | Seconds |
| Dopamine spike (delta > 0.2 in < 5 sec) | Adrenaline += 0.05 | Surprise/excitement sympathetic co-activation | Seconds |
| Dopamine < 0.15 | Fatigue += 0.001/sec | Amotivation, psychomotor slowing | Hours |

### Pain

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| Pain > 0.3 | Cortisol += pain * 0.003/sec | Stress response to pain | Minutes |
| Pain > 0.3 | Heart rate += pain * 30 | Nociceptive sympathetic activation | Seconds |
| Pain > 0.3 | Blood pressure += pain * 25 | Sympathetic vasoconstriction | Seconds |
| Pain > 0.3 | Muscle tension += pain * 0.4 | Guarding/splinting reflex | Seconds |
| Pain > 0.5 | Adrenaline += 0.002/sec | Significant pain triggers fight-or-flight | Seconds |
| Pain > 0.7 | Respiratory rate += (pain - 0.7) * 15 | Pain-driven tachypnea | Seconds |
| Pain > 0.8 sustained > 20 min | Endorphin += 0.002/sec | Endogenous opioid release | Minutes |
| Pain > 0.9 | Risk of vasovagal syncope: BP may suddenly drop | Paradoxical bradycardia | Seconds |

### Fatigue

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| Fatigue > 0.6 | Heart rate baseline shifts up by +5 | Less efficient cardiovascular function | Hours |
| Fatigue > 0.7 | Pain sensitivity *= 1.25 | Reduced pain tolerance when exhausted | Hours |
| Fatigue > 0.7 | Cortisol += 0.0005/sec | Fatigue is a stressor | Hours |
| Fatigue > 0.8 | Immune response -= 0.0005/sec | Immune suppression from sleep deprivation | Hours |
| Fatigue > 0.9 | Muscle tension -= 0.001/sec | Cannot maintain tension when exhausted | Hours |
| Fatigue > 0.95 | Microsleep risk — forced rest events | System forces state change | Immediate |

### Immune Response

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| Immune > 0.4 (fighting infection) | Body temp += (immune - 0.4) * 0.005/sec | Fever generation | Hours |
| Immune > 0.5 | Fatigue += 0.003/sec | Sickness behavior — conserve energy | Hours |
| Immune > 0.5 | Muscle tension += 0.001/sec (body aches) | Inflammatory myalgia | Hours |
| Immune > 0.6 | Hunger -= 0.002/sec | Appetite suppression during illness | Hours |
| Immune > 0.3 | Serotonin -= 0.0003/sec | Inflammatory cytokines deplete tryptophan | Hours |

### Muscle Tension

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| Tension > 0.6 sustained > 30 min | Pain += 0.001/sec | Tension headache / myofascial pain | Minutes-hours |
| Tension > 0.5 | Blood sugar -= 0.01/sec (additional) | Isometric muscle work burns glucose | Minutes |
| Tension > 0.7 | Fatigue += 0.002/sec | Sustained contraction is exhausting | Minutes |

### Hydration (new variable)

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| Hydration < 0.6 | Heart rate += (0.6 - hydration) * 30 | Reduced blood volume, compensatory tachycardia | Minutes |
| Hydration < 0.6 | Blood pressure -= (0.6 - hydration) * 40 | Hypovolemia | Minutes |
| Hydration < 0.5 | Body temp regulation impaired (slower return to baseline) | Sweating reduced | Minutes |
| Hydration < 0.5 | Blood sugar concentrates: BS *= 1.1 | Hemoconcentration | Hours |
| Hydration < 0.4 | Fatigue += 0.005/sec | Severe dehydration | Minutes |
| Hydration < 0.3 | Immune response -= 0.002/sec | Organ stress | Hours |
| Natural depletion rate | -= 0.001/min (resting), -= 0.005/min (heavy exertion) | Baseline water loss | Continuous |

### SpO2 (new variable)

| When... | Then... | Magnitude | Timescale |
|---|---|---|---|
| SpO2 < 94 | Heart rate += (94 - SpO2) * 3 | Hypoxic tachycardia | Seconds |
| SpO2 < 94 | Respiratory rate += (94 - SpO2) * 2 | Hypoxic drive | Seconds |
| SpO2 < 90 | Adrenaline += 0.005/sec | Panic/sympathetic response | Seconds |
| SpO2 < 88 | Fatigue += 0.01/sec | Rapid deterioration | Seconds |
| SpO2 < 85 | Cognitive impairment, decision errors | Cerebral hypoxia | Seconds |
| SpO2 < 75 | Loss of consciousness imminent | Severe hypoxia | Seconds |

---

## 3. Decay and Recovery Rates

Concrete half-lives and time constants for each variable returning to baseline after the stimulus is removed.

| Variable | Time to peak | Half-life / Recovery | Notes |
|---|---|---|---|
| **Body temperature** | Varies by cause (fever: hours; cold exposure: minutes) | Recovery: ~0.1°C per 15 min toward 36.6 under normal conditions | Active thermoregulation; shivering generates ~0.5°C/hr, sweating can lose ~0.3°C per 15 min |
| **Heart rate** | 1-3 seconds | Returns to baseline with half-life of ~60-90 seconds after exercise; ~30 seconds after startle | Parasympathetic reactivation is fast; fitness level modulates this |
| **Blood pressure** | 5-10 seconds | Half-life ~2-5 minutes | Baroreceptor reflex restores within minutes; slower if cause is volume loss |
| **Respiratory rate** | 5-15 seconds | Half-life ~30-60 seconds | Chemoreceptor-driven, very responsive |
| **Hunger** | Gradual — peaks ~4-6 hours after last meal | Suppressed within ~15-20 minutes of eating | Ghrelin peaks before meals, CCK suppresses after |
| **Thirst** | Gradual over hours | Suppressed within ~5-10 minutes of drinking (anticipatory), true correction over ~30 min | Oropharyngeal receptors provide early satiation signal |
| **Fatigue** | Accumulates over 16-18 hours of wakefulness | Recovery: ~7-8 hours of sleep removes ~0.8 of accumulated fatigue. Naps (20 min) remove ~0.1-0.15 | Non-linear: first 2 hours of sleep are most restorative (deep sleep) |
| **Pain** | Immediate for acute injury | Acute: half-life 15-60 min if source removed. Inflammatory pain: half-life 6-24 hours. Use two pain channels if possible (acute + inflammatory) | Endorphins accelerate decay by factor of 2 |
| **Muscle tension** | 1-5 seconds | Half-life ~5-10 minutes after stressor removed; ~30 min if sustained contraction > 1 hour | Longer sustained tension creates longer recovery due to metabolite buildup |
| **Blood sugar** | Spike: 15-30 min after eating; peaks at 45-60 min | Returns to baseline ~2-3 hours after meal. Fasting depletion: ~1 mg/dL per 5 min | Insulin response peaks ~30 min after food; glycogen buffers for ~12-24 hrs of fasting |
| **Cortisol** | 15-30 minutes after stressor onset | Half-life: 60-90 minutes. Full clearance: 3-5 hours | Does NOT drop instantly when stressor ends — the lag is important to model |
| **Adrenaline** | Peaks in 2-3 seconds (neural) to 30 seconds (adrenal release) | Half-life: 2-3 minutes. Functionally gone in 10-15 minutes | Fastest-acting and fastest-clearing hormone in your model |
| **Serotonin** | Changes over hours to days | Half-life of change: 4-12 hours. Chronic depletion takes days-weeks to restore | The slowest-moving variable in your model besides circadian phase |
| **Dopamine** | Peaks in 0.5-2 seconds (phasic burst) | Half-life: 10-20 minutes for spike. Tonic level shifts over hours | Model as tonic baseline + phasic spikes |
| **Immune response** | Hours to days to mount response | Active infection resolution: days. Activation half-life after pathogen cleared: 12-48 hours | The slowest-responding system; fever onset takes 1-4 hours after infection |
| **Circadian phase** | N/A — continuous | Advances at 1 hour per real/sim time. Jet-lag resync: ~1 hour of phase shift per day | Cannot be "reset" instantly — forced phase shifts cause desynchrony |
| **SpO2** | Drops in seconds if ventilation compromised | Recovery: 15-45 seconds once ventilation restored (assuming healthy lungs) | Fast-moving variable; the body has minimal oxygen reserves (~5 min without any ventilation) |
| **Hydration** | Depletion: hours | Rehydration: ~30-60 min to restore after drinking. Oral rehydration is slower than the suppression of thirst (thirst stops before you are actually rehydrated) | |
| **Glycogen** | Depletion: 12-24 hours fasting, 60-90 min intense exercise | Restoration: 4-6 hours with adequate food intake | Supercompensation exists but probably too detailed for this model |
| **Endorphin** | 20-30 min of sustained stress or exercise | Half-life: 20-30 minutes | Runner's high onset at ~20-30 min continuous exertion |

---

## 4. Circadian Modulation

The circadian rhythm is not just a sleep-wake switch. It directly modulates at least 8 of your variables. Model circadian phase as a 24-hour clock and apply these multipliers.

### Concrete Circadian Profiles

Use sine-wave approximations with phase offsets. Let `phase` = circadian time in hours (0-24, where 0 = midnight).

**Cortisol:**
```
cortisol_circadian = 0.15 + 0.25 * max(0, cos((phase - 7) * π / 12))
```
- Peak: 06:00-08:00 (the cortisol awakening response, reaching ~0.35-0.40 of your scale)
- Trough: 23:00-02:00 (drops to ~0.05-0.10)
- This is one of the strongest circadian signals in the body. Do not skip it.

**Body temperature:**
```
temp_circadian = 36.6 + 0.5 * cos((phase - 17) * π / 12)
```
- Peak: ~17:00 (36.8-37.1°C)
- Trough: ~04:00-05:00 (36.1-36.3°C)
- This 0.5-1.0°C variation is real and significant

**Heart rate:**
- Lowest: 03:00-05:00 (can drop to 50-55 bpm in healthy adults)
- Highest: 10:00-12:00 and 17:00-19:00
- Variation: +-8-12 bpm from baseline

**Blood pressure:**
- "Morning surge": rises 15-25 mmHg between 05:00-09:00
- Trough: 02:00-04:00 ("dipping" pattern)
- This is why heart attacks and strokes peak in morning hours

**Blood sugar (fasting):**
- Slightly higher in early morning (05:00-08:00) due to the "dawn phenomenon" — cortisol + growth hormone drive hepatic glucose output
- Add +5-10 mg/dL to blood sugar baseline in early morning hours

**Immune response:**
- Strongest: 22:00-02:00 (night shift of immune activity)
- Weakest: 06:00-10:00
- This is why fevers tend to spike at night
- Modulation: immune_circadian_multiplier = 1.0 + 0.2 * cos((phase - 0) * pi / 12)

**Serotonin:**
- Higher during daylight hours (08:00-18:00)
- Serotonin is converted to melatonin starting around 20:00-21:00, so available serotonin drops in evening
- Model: serotonin_circadian_modifier = +0.05 during day, -0.05 during night

**Fatigue (sleep pressure):**
- Accumulates linearly during wakefulness (~0.05/hour)
- BUT has circadian modulation: lowest perceived fatigue at 10:00 and 19:00, highest at 03:00-05:00 and a minor dip at 14:00-15:00 (the post-lunch dip, which occurs even without eating)
- The 14:00 dip is not food-related — it is a genuine circadian trough. Model it.

**Implementation recommendation:**
```
circadian_alertness = 0.5 + 0.4 * cos((phase - 16) * π / 12) + 0.1 * cos((phase - 14) * π / 6)
```
The second cosine term creates the afternoon dip.

---

## 5. Involuntary vs. Mental-State-Influenced Classification

This matters for your simulation because it determines whether a "mental state" layer can modulate these variables or whether they are strictly governed by physiology.

### Purely Involuntary
These cannot be overridden by conscious effort or psychological state in any meaningful way.

| Variable | Why |
|---|---|
| **Blood sugar** | Hormonal regulation only. You cannot will your blood sugar up or down. |
| **Immune response** | Operates below conscious awareness entirely. (Chronic stress affects it, but that is via cortisol, not direct mental control.) |
| **Circadian phase** | Entrained by light and melatonin. You cannot decide to shift your circadian rhythm by thinking about it. |
| **SpO2** | Gas exchange physics. Not volitionally controllable. |
| **Body temperature** (core) | Hypothalamic setpoint. Peripheral vasodilation can be very slightly influenced by meditation in trained individuals, but for simulation purposes, treat as involuntary. |
| **Glycogen** | Metabolic, no conscious access. |

### Partially Influenced by Mental State
Psychological factors modulate these, but cannot fully override the physiological driver.

| Variable | How mental state modulates it |
|---|---|
| **Pain** | Attention amplifies pain by up to 30-40%. Distraction, meditation, or dissociation can reduce perceived pain by 20-40%. Placebo effect is real and significant (~30% reduction). But a broken bone still hurts. |
| **Muscle tension** | Anxiety increases resting muscle tension by 20-40%. Conscious relaxation techniques can reduce it, but not to zero if there is a physical cause (cold, injury). |
| **Fatigue** | Motivation and arousal can temporarily mask fatigue (reducing its behavioral effect by ~20-30%), but the underlying physiological debt remains and will eventually force sleep. |
| **Hunger** | Stress and emotional state can suppress or amplify hunger by 30-50%. Anorexia of acute stress is real (adrenaline suppresses appetite). Emotional eating amplifies it. |
| **Thirst** | Less modifiable than hunger, but distraction can delay awareness of thirst. |
| **Respiratory rate** | Conscious breathing control is possible (meditation, breath-holding), but the chemoreceptor reflex will override if CO2 rises too high. Model as: mental state can modify RR by +-30%, but if SpO2 < 92 or CO2 is high, override the mental modifier. |
| **Cortisol** | Psychological stressors (worry, rumination, anticipatory anxiety) directly trigger cortisol release. A person lying safely in bed can have sky-high cortisol from anxiety alone. Conversely, meditation measurably reduces cortisol by 15-25%. |
| **Serotonin** | Social connection, positive experiences, sunlight exposure, and sense of safety all increase serotonin. Isolation and perceived helplessness decrease it. This is slow (hours-days) but real. |
| **Dopamine** | Anticipation of reward increases dopamine before the reward arrives. Disappointment crashes it. Entirely psychological triggers are as potent as physical ones. |
| **Endorphin** | Laughter, social bonding, and music can trigger endorphin release, not just physical stress. |

### Heavily Influenced by Mental State
Psychological factors alone can drive these variables to extreme values without any physical cause.

| Variable | How |
|---|---|
| **Heart rate** | A panic attack drives HR to 140-180 bpm with zero physical exertion. Anxiety at rest easily adds 20-40 bpm. Conversely, trained meditators can reduce resting HR by 5-10 bpm. Heart rate is the most psychologically responsive cardiovascular variable. |
| **Blood pressure** | "White coat hypertension" raises systolic by 15-30 mmHg from anxiety alone. Anger spikes it further. |
| **Adrenaline** | A frightening thought triggers a full adrenaline dump. The body does not distinguish between a real threat and an imagined one. This is the core mechanism of panic disorder. |

---

## 6. Cascading Failures and Critical Thresholds

These are the nonlinear breakpoints where your simulation should shift behavior qualitatively, not just quantitatively.

### Body Temperature

| Threshold | Effect | Implementation |
|---|---|---|
| **< 35.0°C** (mild hypothermia) | Shivering maximal, cognitive impairment begins, poor decision-making | Muscle tension forced to > 0.7, apply -20% to any cognitive/decision accuracy |
| **< 33.0°C** (moderate hypothermia) | Shivering STOPS (this is counterintuitive but critical), heart rate drops, confusion | Override: muscle tension drops back toward 0.3, HR decreasing, paradoxical undressing behavior possible |
| **< 30.0°C** (severe hypothermia) | Cardiac arrhythmia risk, loss of consciousness | HR becomes erratic (random between 30-80), BP drops below 80, unconsciousness |
| **< 28.0°C** | Ventricular fibrillation risk — effectively lethal without intervention | Simulation should trigger critical/death state |
| **> 40.0°C** (hyperthermia) | Heat stroke onset, enzyme dysfunction begins | Confusion, BP dropping, sweating may paradoxically stop |
| **> 41.5°C** | Protein denaturation, organ damage | Multi-organ failure cascade begins |
| **> 42.0°C** | Lethal | |

### Blood Sugar

| Threshold | Effect |
|---|---|
| **< 70 mg/dL** | Counter-regulatory hormones fire (adrenaline, cortisol). Subjective hunger, trembling, sweating. |
| **< 55 mg/dL** | Neuroglycopenic symptoms: confusion, difficulty concentrating, slowed reaction time. Implement as cognitive penalty. |
| **< 45 mg/dL** | Seizure risk, behavioral changes (irritability, aggression), loss of fine motor control. |
| **< 35 mg/dL** | Loss of consciousness. |
| **< 25 mg/dL** | Lethal without intervention. |
| **> 180 mg/dL** | Increased thirst, frequent urination (dehydration acceleration * 1.5). |
| **> 300 mg/dL** | Diabetic ketoacidosis territory (if sustained). Nausea, vomiting, altered consciousness. |

### Cortisol Immunosuppression (Sustained Elevation)

This is not a single threshold — it is time-integrated.

| Condition | Effect |
|---|---|
| Cortisol > 0.4 for > 1 hour | Measurable reduction in lymphocyte count. Immune response efficiency -10%. |
| Cortisol > 0.4 for > 6 hours | Immune response efficiency -25%. Wound healing slowed. |
| Cortisol > 0.5 for > 24 hours | Immune response efficiency -40%. Susceptibility to opportunistic infection. |
| Cortisol > 0.6 for > 48 hours | Immune response efficiency -60%. Existing infections may flare. Gut barrier integrity compromised. |

**Implementation suggestion:** Track a running integral of cortisol above 0.3 over time. Use this "cortisol load" value to scale immune suppression:
```
cortisol_load += max(0, cortisol - 0.3) * dt
immune_suppression_factor = 1.0 / (1.0 + cortisol_load * 0.1)
```

### SpO2 Cascade

| Threshold | Effect |
|---|---|
| **< 94%** | Compensatory: HR and RR increase. Subjective: mild air hunger. |
| **< 90%** | Peripheral cyanosis. Cognitive impairment begins. Adrenaline spikes. |
| **< 85%** | Confusion, agitation, loss of coordination. |
| **< 80%** | Loss of consciousness imminent (within minutes). |
| **< 70%** | Organ damage occurring. Cardiac arrhythmia risk. |
| **< 60%** | Lethal within minutes. |

### Multi-Stressor Compounding

The most dangerous cascades involve multiple variables reinforcing each other. Here are the critical feedback loops your simulation must handle.

**The Dehydration-Hypotension-Tachycardia Spiral:**
```
Dehydration → Blood volume drops → BP drops → HR increases to compensate →
Increased HR burns more energy → Blood sugar drops → Adrenaline fires →
Further HR increase → Increased sweating (if hot) → More dehydration
```
Breaking condition: fluid intake, or collapse at HR > 180 / BP < 70.

**The Pain-Stress-Pain Amplification Loop:**
```
Pain → Cortisol increases → Sustained cortisol depletes serotonin →
Low serotonin increases pain sensitivity → More pain → More cortisol
```
Breaking condition: pain source removed, or endorphin release eventually dampens the loop.

**The Fatigue-Immune-Fatigue Loop:**
```
Sleep deprivation → Fatigue > 0.8 → Immune function drops →
Latent infection activates → Immune response increases →
Fever + sickness behavior increases fatigue → Less effective sleep → More fatigue
```
This loop is realistic and common. It models getting sick when chronically exhausted.

**The Hypothermia Deception (model this explicitly):**
```
Cold exposure → Shivering + vasoconstriction → Energy expenditure rises →
Blood sugar depletes → Glycogen depletes → Shivering stops (not because warm,
because out of fuel) → Temp drops faster → HR drops → Consciousness lost
```
The critical design point: shivering cessation at < 33°C is not recovery. It is a sign of worsening. Your simulation must not interpret reduced muscle tension as improvement.

---

## Summary of Recommended Changes

1. **Add 4 variables:** SpO2, hydration level, glycogen reserves, endorphins. Total: 20.
2. **Reframe hunger and thirst** as derived signals from blood sugar/glycogen and hydration respectively, or at minimum couple them tightly.
3. **Implement cortisol with a lag** — it does not spike or clear instantly. This is one of the most common errors in biological simulations.
4. **Model the hypothermia reversal at 33 degrees C** — shivering stops, HR drops, muscle tension falls. This is the most counterintuitive but physiologically important threshold.
5. **Track cortisol load as a time integral**, not just instantaneous level, for immune suppression.
6. **Apply circadian modulation** to at minimum cortisol, body temperature, blood pressure, and immune response. These four have the strongest circadian signals.
7. **Differentiate adrenaline and cortisol timescales** sharply. Adrenaline is seconds-to-minutes. Cortisol is minutes-to-hours. Getting this wrong will make your simulation feel unrealistic.
8. **Implement at least 3 cascade loops** (dehydration spiral, pain-stress amplification, fatigue-immune loop) to create emergent behavior rather than just linear responses.

The numbers I have given are approximations calibrated for simulation, not clinical precision. They will produce physiologically plausible behavior. Tune the specific coefficients during testing, but the relative magnitudes and directions should be preserved.

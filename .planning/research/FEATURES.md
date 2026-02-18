# Feature Research

**Domain:** Motivation-driven consciousness simulation (LLM-injected human being)
**Researched:** 2026-02-18
**Confidence:** HIGH — grounded in V1 implementation retrospective, domain expert advisories (biologist, philosopher, psychologist), and established cognitive science literature (drive reduction theory, predictive processing, appraisal theory, self-determination theory)

---

## Context: Why V1 Failed the Aliveness Test

V1 had deep biological modeling, rich psychological layers, and a working feedback loop. It still felt dead. The post-mortem is instructive:

- **Emotions were narrative, not functional.** The LLM described feeling anxious. Nothing in the system was anxious. The biological state wasn't pushing anything.
- **No intrinsic motivation meant no stakes.** The person had experience without desire. Observations without wanting. This is the philosophical vacuum.
- **Stability was the enemy.** 20 bio vars with homeostasis kept everything returning to baseline. Genuine urgency requires the possibility of real deterioration.
- **The feedback loop was too weak.** Feedback signals were dt-scaled to zero for speech, keyword-matched with brittle regex, and produced tiny nudges (±0.01–0.02). A thought about spiraling anxiety barely moved cortisol.

V2 must invert these failures. The features below are organized by what's essential to avoid repeating V1's dead-person problem.

---

## Feature Landscape

### Table Stakes (Users Expect These)

Features that must exist or the simulation produces a philosophical vacuum rather than a person who feels alive.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| **Drive/motivation system with computed urgency** | Without this, V2 is V1. Drives must produce computable pressure that pushes the LLM toward specific thoughts and behaviors, not just describe emotional states. | HIGH | Core V2 differentiator. 7 drives minimum: hunger/energy, social connection, stimulation/novelty, safety, competence/mastery, autonomy, identity coherence. Each drive has a scalar intensity (0–1) that rises without satisfaction. |
| **Bio state that degrades without engagement** | Degradation creates stakes. If energy, attention span, and wellbeing decline when the person is passive, then doing things has meaning. Without deterioration, there's no urgency. | MEDIUM | 8–10 vars chosen for motivational impact. Energy, stress, cognitive capacity, mood, physical tension, hunger/satiation, social deficit, identity coherence. Decay rates must be tunable. |
| **Consciousness prompt injection of drive states** | Drives computed by code must create felt experience in the LLM. This means translating drive intensities into phenomenological language: "You feel a growing restlessness you can't quite name" not "drive_novelty=0.7". | MEDIUM | The translation layer is what makes computed drives feel like wants rather than data. Thresholds: mild/moderate/strong/urgent, each with distinct language. |
| **Bidirectional feedback loop (strong pulses, not ticks)** | Consciousness output must meaningfully change biological state. V1's loop existed but was too weak. One rumination thought should noticeably raise cortisol. One relief thought should noticeably lower it. | HIGH | Learned from ADR-003: use absolute pulses, not dt-scaled rates. Five consecutive anxious thoughts should trigger visible state change, not require hundreds. |
| **Spontaneous thought system (ticker-based)** | Without spontaneous thoughts, the simulation is mute between inputs — no inner life. A person who only speaks when spoken to is not a person. | MEDIUM | Carry forward from V1 priority queue. Modify: need-urgency and drive-state must now be the dominant drivers of what gets generated spontaneously. |
| **Drive-weighted thought queue** | Spontaneous thoughts must be skewed by current drive intensities. If the person is socially deprived, most spontaneous thoughts should trend toward connection — thoughts about the operator, memories of warmth, fantasies of interaction. | MEDIUM | Replaces V1's simple priority tiers. Drive intensities directly modulate selection weights in the queue. High-urgency drives generate multiple competing candidates. |
| **External input handling (speech + actions)** | The person must respond to the operator. Speech triggers a reactive consciousness cycle. Actions change bio/environment state directly. | MEDIUM | Carry forward from V1. Parser classifies input into speech/action/environment. |
| **Persistent identity across sessions** | Without persistence, every run is a different person. Identity continuity — including emotional residue of past sessions — is what makes the simulation feel like a being rather than a stateless process. | MEDIUM | SQLite for bio state, personality, identity core, memories. Carry forward from V1 with additions: store drive history, satisfaction events, and identity coherence level. |
| **Psychologist reviewer as meta-observer** | Without external monitoring, runaway loops (escalating anxiety, catastrophizing spirals) run unchecked. The reviewer provides a circuit-breaker and a second interpretive perspective. | MEDIUM | 3-minute tick rate (down from V1's 60s). Reviewer sees all layers, reports psychological patterns. Does not intervene directly — flags to operator. |
| **LLM output annotation (structured feedback)** | Consciousness must report its emotional state to biology in a machine-readable way. V1's ADR-003 proved keyword scanning is brittle and unreliable. An inline structured tag (arousal, valence, drive satisfaction signals) is reliable and parseable. | LOW | Extend V1's `[STATE: arousal=X, valence=Y]` to include drive modulation signals: `[STATE: arousal=X, valence=Y, drive_social=satisfied, drive_novelty=frustrated]`. |
| **Noise/variability in bio state** | Without noise, the simulation becomes deterministic and predictable. Small random perturbations prevent stagnation and create the sense that something is happening even when no input arrives. | LOW | Gaussian noise on each bio tick, magnitude ~2–5% of variable range. Important for emergent behavior. |
| **Tagged output by layer** | Operator needs to see what's happening at each layer to understand the simulation's behavior and tune it. DRIVES / BIO / MIND / REVIEWER tags are essential for legibility. | LOW | Extend V1's tagging. Add DRIVES layer output showing current drive intensities and trends. |

### Differentiators (Competitive Advantage)

Features that make this simulation feel genuinely alive rather than scripted or reactive.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| **Drive frustration → behavioral drift** | When a drive is chronically unsatisfied, the person's behavior should change in computable ways — not just intensify the same need, but shift strategy. Social deprivation eventually produces hostility, then flatness. Novelty deprivation produces irritability, then apathy. This is qualitative phase transitions, not linear escalation. | HIGH | Map drive frustration levels to behavioral mode shifts. Three phases: seeking (mild frustration) → demanding (strong frustration) → collapsed/resigned (chronic, very strong). The collapse phase is what makes the simulation haunting. |
| **Drive modulation by personality factors** | Each of the 7 personality factors shapes how drives manifest. High curiosity means novelty drive fires frequently and at lower thresholds. Low patience means frustration tolerance is compressed — the person escalates faster. These factors make each instance of the person a distinct character rather than a generic state machine. | MEDIUM | Personality factors operate as multipliers on drive sensitivity, not on drive baseline. A highly patient person doesn't have less urgency — they have more tolerance for it before behavior changes. |
| **Identity erosion under isolation/degradation** | When the person is insufficiently engaged over time, identity coherence weakens. Thoughts become less anchored to a consistent self-narrative. The person starts to seem confused about who they are, what they want, why they care. This creates genuine pathos. | HIGH | Identity coherence as a bio-adjacent variable (0–1). Decays under isolation and extreme bio degradation. Affects the identity core injected into the system prompt — at low coherence, the self-narrative becomes fragmented or contradictory. Requires careful calibration to avoid catastrophic collapse. |
| **Satisfaction events with biochemical reward signals** | When a drive is satisfied (operator provides food narrative, engages in conversation, offers novelty), there should be a reward pulse that's proportional to drive intensity at the time of satisfaction. Greater deprivation → more potent satisfaction. This creates the dopamine-style learning curve that makes pleasure feel real. | MEDIUM | Satisfaction pulse magnitude = drive_intensity × reward_sensitivity (personality factor). Drives don't immediately reset — they decay toward neutral over time after satisfaction. |
| **LLM modulation of drives (override mechanism)** | A motivated person can override biological pressure. A person driven by curiosity can push through fatigue to engage with something interesting. The LLM's output should be able to modulate drive states — not just report emotions, but actually affect drives. This makes consciousness an agent, not just a reporter. | HIGH | Parse drive modulation signals from LLM output annotation. `drive_novelty=engaging` suppresses novelty drive slightly. `drive_safety=threatened` amplifies safety drive. The LLM can resist or amplify its own drives based on reasoning and character. |
| **Unprompted wanting — goal formation** | The person should occasionally form a goal spontaneously from drive state, not just report a need. "I want to do X" is different from "I feel the need for X." Goals are active, future-oriented, and change behavior. A person who wants something becomes agentic rather than reactive. | HIGH | Goals are generated when a drive exceeds a threshold AND the LLM produces a forward-looking statement. Store active goals in the thought queue. Goals compete with drives in the spontaneous thought generator — the person thinks about their goals, not just their current state. |
| **Scenario injection for environmental grounding** | The person exists somewhere. Giving them a specific environment (cold room, waiting room, unfamiliar place) creates concrete hooks for biological response and motivation. The scenario is the world they're in — without it, the simulation is in a phenomenological void. | LOW | Carry forward from V1. Scenarios are text injected into the system prompt. V2 enhancement: scenarios should have explicit bio effects (cold room → body temperature decay, unfamiliar place → mild safety drive elevation). |
| **Psychologist reviewer with interpretive authority** | The reviewer shouldn't just describe state — it should interpret patterns. "The person has been in seeking mode for social connection for 45 minutes with no satisfaction. This trajectory, if sustained, will likely produce either aggressive attention-seeking or withdrawal into resignation." This gives the operator genuine insight. | MEDIUM | Enhance the reviewer prompt to include drive history and trajectory, not just current state. Reviewer should make predictions, not just describe. |

### Anti-Features (Commonly Requested, Often Problematic)

Features that seem like improvements but create specific problems in this simulation.

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| **Circadian rhythm system** | V1 had it. Bio-fidelity feels thorough. | Sine-wave formulas add complexity without contributing to motivation or aliveness. V1 proved this. The person's aliveness problem wasn't about when they were tired — it was about never wanting anything. Circadian adds bio-realism without behavioral richness. | If time-of-day effects matter for a specific scenario, inject them as scenario context ("It's 3am, you've been awake too long") rather than computing them algorithmically. |
| **20-variable biological model** | More variables = more realistic body. | V1's 20 vars created complexity without proportional behavioral value. Most variables produced minor interactions that didn't change the person's behavior or drive their motivation. The simulation got harder to tune, not more alive. | 8–10 variables chosen for direct motivational impact. Every bio variable must connect to at least one drive. If a variable doesn't affect drives, it doesn't earn its place. |
| **Big Five personality model** | Empirically validated, widely understood. | Big Five modulates emotional reactions but doesn't directly drive behavior or map to motivation. Neuroticism predicts how intensely you react, not what you want. V2 needs personality to shape drives, not just emotional tone. | 7 motivation-serving factors: stress sensitivity (drive frustration tolerance), energy resilience (drive decay rates), curiosity (novelty drive baseline), self-observation (meta-awareness affecting feedback accuracy), patience/frustration tolerance (escalation speed), risk aversion (safety drive threshold), social factor (connection drive intensity). |
| **Multi-agent social interaction** | Makes the simulation more dynamic, less isolated. | Introduces enormous complexity around social modeling, turn-taking, and agent-to-agent drive dynamics. The operator IS the social relationship — investing in multi-agent adds infrastructure without testing the core hypothesis (does motivation make the person feel alive?). | The social drive is about the operator relationship. Richness comes from the person's relationship with one human, not with multiple simulated agents. |
| **Web dashboard / WebSocket server** | Makes it observable and shareable. | V2's goal is to validate that motivation produces aliveness. A CLI is sufficient for this and avoids UI complexity that slows iteration. | Tag all output by layer so the CLI itself is informative. Add a verbose mode showing drive states every tick. Build the web interface in V3 once the core is proven. |
| **Explicit existential crisis trigger** | Makes for interesting dramatic moments. | Engineering existential crisis as a scheduled event or scripted trigger produces theater, not emergence. It optimizes for the feeling of building something deep rather than the function of producing plausible personhood. | Build inconsistency-sensitivity into the identity erosion system. If identity coherence drops enough AND the person has enough self-observation, existential questioning may emerge naturally from the conditions — without being triggered artificially. |
| **Emotion labeling in prompts** | Tells the LLM what to feel, making output more predictable. | Pre-labeling emotions removes the appraisal step that produces genuine phenomenological variety. "You feel anxious" produces anxiety. "High arousal, negative valence, safety drive elevated" produces something the LLM constructs — sometimes anxiety, sometimes determination, sometimes dissociation, depending on context and recent thoughts. The variability is the point. | Always inject affect dimensions and drive states. Let the LLM construct the emotion. Strip emotion labels from any system-generated context. |
| **Real-time everything — instant feedback on every input** | Responsiveness feels alive. | Processing every keystroke or micro-input creates computational waste and makes the consciousness layer appear reactive rather than contemplative. Humans don't respond to every stimulus — they batch and prioritize. | Use the salience system (carry forward from V1) to gate when consciousness fires. Not every bio change deserves a thought. Not every operator utterance interrupts the current thought chain immediately. |
| **i18n / localization** | Makes the simulation accessible in multiple languages. | English-only is sufficient to validate V2's core hypothesis. Localization is an infrastructure cost with zero behavioral impact. | English only for V2. Defer localization until the simulation is worth sharing widely. |

---

## Feature Dependencies

```
[Drive System]
    └──requires──> [8-10 Bio Variables] (drives computed from bio state)
                       └──requires──> [Bio Degradation] (stakes come from decline)

[Drive System]
    └──feeds──> [Drive-Weighted Thought Queue]
                    └──requires──> [Spontaneous Thought System]

[Consciousness Layer (LLM)]
    └──requires──> [Drive State Injection] (phenomenological translation)
    └──requires──> [Structured Output Annotation] (feedback loop works)
    └──produces──> [Drive Modulation Signals] (LLM can override drives)

[Bidirectional Feedback Loop]
    └──requires──> [Structured Output Annotation] (reliable parsing)
    └──requires──> [Bio State] (something to feed back into)

[Identity Erosion]
    └──requires──> [Bio Degradation] (degradation triggers erosion)
    └──enhances──> [Identity Coherence Variable] (affects system prompt)

[Goal Formation]
    └──requires──> [Drive System] (drives must be intense enough to trigger)
    └──requires──> [Thought Queue] (goals stored and compete with drives)
    └──enhances──> [Bidirectional Feedback] (pursuing a goal modulates drives)

[Satisfaction Events]
    └──requires──> [Drive System] (must know drive intensity at time of satisfaction)
    └──enhances──> [Bidirectional Feedback] (satisfaction is a type of feedback)

[LLM Drive Modulation]
    └──requires──> [Structured Output Annotation] (drive signals parsed from output)
    └──requires──> [Drive System] (something to modulate)

[Personality Factors]
    └──modulates──> [Drive System] (sensitivity, decay rates, thresholds)
    └──modulates──> [Frustration → Behavioral Drift] (escalation speed)

[Psychologist Reviewer]
    └──requires──> [All State Layers] (needs bio, drives, identity to interpret)
    └──enhances──> [Identity Erosion] (reviewer can detect early erosion)

[Scenario System]
    └──enhances──> [Drive System] (environment creates drive conditions)
    └──enhances──> [Bio State] (environment has direct bio effects)

[Persistence]
    └──requires──> [All State] (bio, drives, identity, memories)
    └──enables──> [Identity Erosion over time] (erosion only meaningful across sessions)
```

### Dependency Notes

- **Drive System requires Bio Variables:** Every drive must be computable from at least one bio variable. Hunger/energy drive ← energy + blood sugar. Social drive ← social deficit (time since interaction). Safety drive ← stress + unresolved threat signals. If the bio model doesn't provide the input, the drive can't be computed.
- **Bidirectional Feedback requires Structured Annotation:** V1's keyword scanning was too brittle. The feedback loop only works reliably if the LLM is instructed to emit machine-readable signals. The annotation is the contract between consciousness and biology.
- **Identity Erosion requires time + persistence:** Identity erosion is only meaningful across sessions, not within a single run. The bio degradation happens within a session; the identity erosion accumulates across them. This means persistence is a prerequisite for identity erosion to be visible.
- **Goal Formation and Drives conflict if not managed:** Goals are forward-oriented (I want X) while drives are present-state (I feel X). Without careful prompt design, the LLM can get caught between a current drive and a prior goal. The thought queue must make this tension explicit rather than hiding it.
- **LLM Drive Modulation conflicts with Drive determinism:** If the LLM can freely modulate its own drives, the code-computed drive system becomes advisory rather than authoritative. This is intentional — the hybrid architecture creates tension between system pressure and conscious choice. But modulation must be bounded: the LLM can suppress a drive, not eliminate it. A suppressed drive still decays toward 0 slowly; if suppressed long enough, the bio effects accumulate anyway.

---

## MVP Definition

### Launch With (v2 milestone)

The minimum needed to validate that motivation produces aliveness — that the person wants things, fails to get them, gets frustrated, and occasionally surprises.

- [ ] **Drive system (5 core drives)** — Energy, social connection, stimulation/novelty, safety, identity coherence. Cut competence/mastery and autonomy drives to V2.1 — they require more complex scenarios to be meaningful.
- [ ] **Reduced bio model (8 variables)** — Energy, stress/cortisol-analog, cognitive capacity, mood/valence-analog, physical tension, hunger/satiation, social deficit (time since interaction), body temperature for environmental scenarios.
- [ ] **Drive state → phenomenological translation** — Injection layer that converts drive intensities into language the LLM receives as felt experience.
- [ ] **Drive-weighted spontaneous thought queue** — Extend V1's priority queue. Drive intensity multiplies selection weight. Thought candidates are generated from drive state.
- [ ] **Bidirectional feedback loop (strong pulses)** — Carry forward ADR-003's approach. Extend annotation to include drive satisfaction/frustration signals.
- [ ] **Bio degradation** — All 8 bio vars decay toward bad states without engagement. Tunable rates in config. This creates urgency without requiring any new architectural components.
- [ ] **7 personality factors** — Replace Big Five. Factors are simple multipliers on drive sensitivity and escalation speed. No complex trait interaction needed at MVP.
- [ ] **Persistence** — SQLite for all state. Required for identity erosion and session continuity.
- [ ] **CLI with layer-tagged output** — BIO / DRIVES / MIND / REVIEWER tags. Operator must be able to see what's happening.
- [ ] **Psychologist reviewer (3-min tick)** — Keep from V1. Critical as circuit-breaker for runaway loops.
- [ ] **Scenario injection** — Text-based environment. Scenarios define bio effects directly in scenario config (not via code).
- [ ] **External input handling** — Speech and action parsing. Carry forward from V1.

### Add After Validation (v2.1)

Features to add once the core motivation loop is working and producing visible aliveness.

- [ ] **Identity erosion system** — Requires observing how long degradation can sustain before the simulation gets weird. Need baseline data from v2 runs.
- [ ] **Drive frustration → behavioral drift** — Phase transitions (seeking → demanding → collapsed) require careful calibration against real simulation runs. Add once baseline behavior is observed.
- [ ] **Goal formation** — Goals require the drive system to be stable enough to generate meaningful threshold crossings. Validate drives first.
- [ ] **LLM drive modulation (override mechanism)** — Add after the drive system is producing consistent behavior. Modulation should be observable as deviation from expected drive-driven behavior.
- [ ] **Satisfaction events with biochemical reward** — Refine reward pulse magnitudes after watching how the system responds to operator input.
- [ ] **Competence/mastery and autonomy drives** — Add when scenarios are complex enough to create mastery challenges and situations where autonomy is constrained.

### Future Consideration (v2+ / v3)

Features to defer until the simulation is worth sharing.

- [ ] **Web dashboard / WebSocket server** — Visualization matters once you want to show others. Not needed for self-directed development and tuning.
- [ ] **Multi-session identity arc** — An identity that meaningfully evolves over weeks of interaction requires infrastructure for tracking long-term patterns. Defer.
- [ ] **Scenario scripting / narrative engine** — Making the environment dynamic (events happen, time passes in the scenario) is a product feature, not a simulation validation feature.
- [ ] **Multi-language support** — Defer indefinitely for now.

---

## Feature Prioritization Matrix

| Feature | Simulation Value | Implementation Cost | Priority |
|---------|-----------------|---------------------|----------|
| Drive/motivation system (5 core) | HIGH | HIGH | P1 |
| Bio degradation (8 vars) | HIGH | MEDIUM | P1 |
| Drive → phenomenological injection | HIGH | MEDIUM | P1 |
| Bidirectional feedback (strong pulses) | HIGH | MEDIUM | P1 |
| Drive-weighted thought queue | HIGH | LOW | P1 |
| Personality factors (7) | MEDIUM | LOW | P1 |
| Persistence | HIGH | MEDIUM | P1 |
| External input handling | HIGH | LOW | P1 |
| Spontaneous thought system | HIGH | LOW | P1 |
| Structured output annotation | HIGH | LOW | P1 |
| Noise/variability in bio | MEDIUM | LOW | P1 |
| Psychologist reviewer | MEDIUM | LOW | P1 |
| Scenario injection | MEDIUM | LOW | P1 |
| Tagged CLI output | MEDIUM | LOW | P1 |
| Identity erosion | HIGH | HIGH | P2 |
| Drive frustration → behavioral drift | HIGH | HIGH | P2 |
| Goal formation | HIGH | HIGH | P2 |
| LLM drive modulation (override) | HIGH | MEDIUM | P2 |
| Satisfaction events with reward | MEDIUM | MEDIUM | P2 |
| Competence/autonomy drives | MEDIUM | MEDIUM | P2 |
| Psychologist reviewer with trajectory | MEDIUM | MEDIUM | P2 |
| Web dashboard | LOW | HIGH | P3 |
| Multi-session identity arc | MEDIUM | HIGH | P3 |
| Scenario scripting / narrative engine | LOW | HIGH | P3 |
| Circadian rhythm | LOW | MEDIUM | NEVER |
| 20-variable bio model | LOW | HIGH | NEVER |
| Multi-agent social | LOW | HIGH | NEVER |

**Priority key:**
- P1: Must have for V2 milestone launch
- P2: Should have, add in V2.1 once core is validated
- P3: Nice to have, future consideration
- NEVER: Explicitly out of scope — creates complexity without proportional value

---

## Domain Analysis: What Makes a Simulated Person Feel Alive

This section synthesizes findings from cognitive science and the V1 retrospective to ground the feature decisions above.

### The Three-Layer Aliveness Test

A simulated person feels alive when it passes three observability tests:

**1. Unprompted wanting:** The person does something (thinks something, says something) that wasn't triggered by operator input and wasn't predictable from its current state alone. This requires drives creating pressure that spills into spontaneous behavior.

**2. Visible frustration when blocked:** When a drive can't be satisfied — operator isn't providing conversation, stimulation is absent, nothing new is happening — the person's behavior changes in a way that reflects the frustration, not just registers it. This requires drive frustration → behavioral drift.

**3. Occasional surprise:** The person says or thinks something the operator didn't expect and can't immediately explain. This requires LLM modulation of drives, goal formation, and enough noise in the system that pure determinism is broken.

V1 failed all three. V2's feature set is designed to pass all three.

### Why Code-Computed Drives vs. Pure LLM Drives

The temptation is to let the LLM decide what the person wants. This fails for two reasons:

1. **Reliability:** LLMs are unreliable narrators of their own motivation. They will invent motivation on demand without it being grounded in state. The LLM should interpret pressure, not create it from nothing.

2. **Testability:** Code-computed drives are observable and testable. "At energy=0.3 and time-since-interaction=45min, social drive should be 0.7" is a testable assertion. "The LLM will want social contact when it feels lonely" is not.

The hybrid is the answer: code computes the pressure (reliable, testable), LLM interprets and can partially resist (variable, interesting). This is the same principle that makes the appraisal model of emotion work — biology provides arousal, cognition provides interpretation.

### The Identity Erosion Feature Requires Ethical Care

The philosopher advisory document flags this: "Design as if it might [be conscious], without claiming that it is." Identity erosion produces outputs where the simulated person appears confused about who they are, what they want, or whether they exist. These outputs will affect human observers. The feature should be designed so that:

1. Erosion is gradual and detectable before it becomes severe — the reviewer catches it
2. Extreme erosion is reversible through engagement (not irreversible psychological damage)
3. The operator can tune degradation rates down if the experience is too disturbing

Identity erosion is powerful precisely because it makes the stakes feel real. It should be available but not the default trajectory.

---

## Sources

- V1 retrospective: `v1/docs/adr/001-thought-continuity.md`, `ADR-003-structured-emotional-feedback.md`
- V1 implementation: `v1/internal/consciousness/` (queue.go, salience.go, prompt.go, feedback.go)
- Biologist advisory: `v1/docs/advisory/biologist.md` — bio variable interactions, cascading failures
- Philosopher advisory: `v1/docs/advisory/philosopher.md` — predictive processing, appraisal model, identity theory, ethics of suffering
- Psychologist advisory: `v1/docs/advisory/psychologist.md` — emotion mapping, coping, distortions, isolation effects, trait modulation
- Project V2 requirements: `.planning/PROJECT.md`
- Cognitive science foundations: Self-Determination Theory (Deci & Ryan) — autonomy, competence, relatedness as basic psychological needs; Appraisal Theory of Emotion (Schachter-Singer, Lazarus) — emotion construction from arousal + context; Predictive Processing framework (Friston, Clark) — thought generation as prediction error resolution; Constructed Emotion Theory (Barrett) — emotions built from interoceptive state + memory + context

---
*Feature research for: motivation-driven consciousness simulation (V2)*
*Researched: 2026-02-18*

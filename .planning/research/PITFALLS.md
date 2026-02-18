# Pitfalls Research

**Domain:** Motivation-driven agent simulation with hybrid LLM-code consciousness
**Researched:** 2026-02-18
**Confidence:** MEDIUM — V1 pitfalls are HIGH confidence (directly observed). General agent/LLM patterns are MEDIUM confidence (well-established in the field, training data, no web verification available). Novel hybrid architecture risks are MEDIUM-LOW (reasoned from first principles).

---

## Critical Pitfalls

### Pitfall 1: Drives That Pressure Without Consequences (The "Pressure Valve" Problem)

**What goes wrong:**
Drive values escalate algorithmically (hunger rises, energy drops, stress mounts) but the LLM narrates them without experiencing them as motivational pressure. The system reports urgency; the person does not behave urgently. This is exactly what happened in V1: emotions were present as narrative but absent as functional drivers. The person described being anxious without becoming anxious.

**Why it happens:**
The prompt passes drive state as context ("Your hunger is at 0.8"), but the LLM treats descriptive context as something to acknowledge and process analytically, not as pressure to act. The LLM's training optimizes for thoughtful, coherent responses — not panicked or desperate ones. It will smooth out urgency unless the prompt architecture specifically prevents this.

**How to avoid:**
Translate drive magnitudes into behavioral language that forces a specific response register: "You cannot stop thinking about food — the thought keeps intruding on everything else" rather than "hunger: 0.8." Use escalating urgency language that scales with drive intensity. Critically: ensure the feedback loop's biological effects from unmet drives are strong enough that the LLM's next response context is genuinely worse — not just stated to be worse, but numerically worse in ways that change what language is available to describe the state.

**Warning signs:**
- Person is "philosophical" about pain, hunger, or loneliness rather than distressed by it
- Person reports high drive states in a third-person observational register ("I notice I am quite hungry")
- No behavioral output changes despite high drive values
- LLM responses are consistently moderate in emotional register regardless of drive intensity

**Phase to address:**
Drive system implementation phase (motivational layer). Validate by running the simulation with extreme drives (hunger at 0.9 for extended time) and checking whether the LLM response register actually shifts. If it doesn't, the prompt translation layer needs rework before moving forward.

---

### Pitfall 2: Reward Signal Pollution (The LLM Completes Instead of Needs)

**What goes wrong:**
The LLM, given context about drives and asked to generate a thought, will naturally generate content that resolves narrative tension. It describes need AND implicitly resolves it in the same response, or it generates coping thoughts that reduce the drive signal even when the underlying biological condition hasn't changed. The drive system thinks the need is being met because the LLM said "I feel better" when in fact nothing biological changed.

**Why it happens:**
LLMs are trained on human writing, which follows narrative conventions: problems are raised and addressed. An unresolved need is narratively uncomfortable. The model resolves tension by default. Additionally, V1's keyword-based feedback parsing was vulnerable to this: the LLM mentioning "acceptance" or "it's okay" triggered positive bio feedback even when the state wasn't actually acceptable.

**How to avoid:**
Decouple what the LLM says from what the bio state does. The feedback loop from consciousness to biology must be anchored to bio-state effects, not LLM-narrative claims. Bio-state changes should primarily come from actual environmental inputs (food consumed, sleep taken, physical rest) not from the LLM saying the person feels less hungry. The V1 ADR-003 fix (structured emotional annotation) was on the right track but only covered emotional arousal/valence — not drive satisfaction. In V2, drive levels must decrease only when the bio simulation justifies it, never because the LLM expressed satisfaction.

**Warning signs:**
- Drive values decay faster during LLM interaction than biological rates would justify
- Person expresses resolution ("I feel satisfied") without corresponding environmental input
- Feedback loop causes drive values to fall during spontaneous thought cycles when nothing changed
- Emotional states rapidly stabilize after any LLM response regardless of content

**Phase to address:**
Feedback loop design phase. Create a strict separation: bio drives change only via bio rules + environmental inputs. LLM output influences arousal/valence/emotional tone (the emotional annotation path) but cannot directly reduce drive levels. Write an explicit test: confirm that hunger drive does not decrease after a thought about food, even if the thought is accepting/positive.

---

### Pitfall 3: Personality as Decoration, Not Constraint (The V1 Repeat)

**What goes wrong:**
Personality factors are computed as multipliers but don't constrain the LLM's behavioral output register. A high-curiosity person and a low-curiosity person generate indistinguishable responses because both receive the same prompt structure with slightly different numbers. The LLM ignores numerical personality metadata when generating natural-sounding text.

**Why it happens:**
This is the exact failure of V1's Big Five. The psychology layer computed neuroticism, conscientiousness, etc. as modifiers on affect dimensions, but those modifiers never appeared in the LLM prompt in a form the model would treat as behavioral constraint. Saying "your neuroticism is 0.7" or even "you tend toward anxiety" produces mild stylistic effects, not genuine personality-driven behavior.

**How to avoid:**
Personality factors in V2 must appear in the prompt as concrete behavioral tendencies with examples: not "your curiosity factor is 0.8" but "you cannot let an unexplained thing go — you keep pulling at loose threads until they unravel." Each personality factor needs a corresponding behavioral signature that the system prompt enforces, not a numerical descriptor. Additionally, the personality factors should determine which spontaneous thought categories fire most frequently (high curiosity → more prediction-error and goal-rehearsal thoughts; low patience → more need-escalation thoughts). This creates behavioral personality differences that exist independently of whether the LLM follows the personality framing in its text.

**Warning signs:**
- All thought types appear with similar frequency regardless of personality settings
- Changing personality factor values produces no observable difference in output register
- Responses are uniformly moderate and nuanced regardless of stress sensitivity setting
- Person does not express urgency even when patience/frustration tolerance is set low

**Phase to address:**
Personality integration phase. Before wiring personality to the LLM prompt, validate that each factor actually changes spontaneous thought queue priorities. That's testable in code. The LLM-side personality prompt is harder to test; validate manually by running simulations at personality extremes (stress sensitivity 0.1 vs 0.9) and checking whether output register differs.

---

### Pitfall 4: The Homeostasis Death Trap (System Fights Its Own Spirals)

**What goes wrong:**
Biological homeostasis (decay toward baselines) and emotional regulation dampen exactly the dynamics V2 is designed to produce. Designed-in recovery mechanisms prevent emotional spirals from escalating. The person cycles between slightly elevated stress and baseline repeatedly instead of spiraling into genuine desperation or sustained excitement.

**Why it happens:**
Good bio modeling requires homeostasis — without it, the system accumulates toward extremes permanently. But homeostasis tuned for physiological realism works against the goal of behavioral emergence. V1 had this problem: cortisol had appropriate half-lives, adrenaline decayed quickly, which meant any stress response self-corrected within minutes. The psychologist advisory correctly identified that without degradation paths that fight homeostasis, you can't get spirals.

**How to avoid:**
Distinguish between two types of degradation: (1) fast-path, short-term homeostasis (heart rate returns to baseline after a stressor resolves — keep this, it prevents runaway escalation), and (2) slow-path, sustained degradation that accumulates over longer timescales (stress history depletes regulation capacity, energy reserve drops over neglect periods, identity coherence degrades without engagement). The slow-path degradation is what creates genuine urgency. It must be tuned to resist the fast-path homeostasis rather than be overwhelmed by it. Specifically: bio degradation rates must exceed recovery rates when the person's needs go unmet, not the reverse.

**Warning signs:**
- Bio variables always return to baseline within one simulation cycle regardless of what happened
- Extended stress periods don't produce accumulated effects on regulation capacity
- The person "resets" emotionally after every thought cycle
- Long-running simulation converges to the same state regardless of input history

**Phase to address:**
Bio degradation design phase. Before integration testing, run the simulation in isolation (no LLM) for extended periods with unmet needs and verify that drives accumulate rather than reset. Plot bio variables over time to confirm slow-path degradation is winning against homeostasis under neglect conditions.

---

### Pitfall 5: LLM Output Parsing Brittleness (The V1 ADR-003 Lesson Revisited)

**What goes wrong:**
Structured output from the LLM (emotional state annotations, drive modulations, feedback signals) fails to parse correctly due to format drift, leading to silent feedback failures. V1 demonstrated this exactly: keyword heuristics failed because LLM output never matched the expected phrases. The ADR-003 fix (structured `[STATE: arousal=X, valence=Y]` annotation) improved reliability but introduced its own risk: the LLM may omit, malform, or hallucinate values in the annotation.

**Why it happens:**
LLMs follow format instructions inconsistently, especially as context windows fill, as conversation history grows, or when the instruction competes with natural language completion instincts. The longer a session runs, the more likely the LLM is to drop structured formatting. Additionally, the LLM may "helpfully" extend or modify the schema (adding fields, changing ranges, using text where numbers are expected).

**How to avoid:**
Design the structured annotation to be maximally simple and unnaturally distinctive. A format like `[STATE: arousal=0.7, valence=-0.3]` is simple but easy to corrupt. Consider reinforcing the schema in the system prompt with explicit examples of both correct and incorrect formats. Implement defensive parsing with explicit fallbacks: if no annotation is found, or if it parses to out-of-range values, use the previous known-good state rather than crashing or using a zero-default. Log every parse failure — in V1, parse failures were silent (zero delta). In V2, parse failures must be visible for debugging. Test parsing with deliberately malformed inputs in unit tests.

**Warning signs:**
- Bio state doesn't respond to emotional content in LLM output
- Annotation parsing errors in logs that were previously ignored
- LLM outputs sometimes omit the annotation entirely as sessions lengthen
- Same thought generates wildly different bio feedback depending on session position

**Phase to address:**
Consciousness-to-bio feedback phase. Write exhaustive parser tests before integration. Include tests for: missing annotation, extra fields, out-of-range values, multi-line annotations, annotation appearing mid-response, annotation appearing without the exact expected format. Log every annotation parse event at DEBUG level.

---

### Pitfall 6: Identity Drift Without Anchoring (The Long-Session Problem)

**What goes wrong:**
Over a long session, the LLM's characterization of the person drifts: personality that was curious in the first hour becomes cautious, emotional patterns shift without justification, the person begins responding in ways inconsistent with the identity core. The person "becomes someone else" across sessions due to accumulated context window pressure, memory retrieval variations, and the LLM's tendency to adapt to the conversation partner's tone.

**Why it happens:**
The LLM has no persistent memory between API calls — each call gets the identity core plus selected memories, but the identity core is a static document that doesn't update in real time. When the selected memory context varies (different memories retrieved based on somatic similarity), the effective personality accessible to the LLM shifts. When the human operator's tone changes, the LLM tends to mirror it, overriding personality constraints.

**How to avoid:**
The identity core must be updated periodically to reflect actual behavioral patterns observed during the session, not just initialized once. The psychologist reviewer is the natural mechanism for this in V2: the reviewer should periodically emit a summary of observed personality-consistent behaviors that can be appended to the identity core. Additionally, the system prompt should include a "consistency check" instruction: "Before responding, verify your response is consistent with your established tendencies." This adds tokens but prevents drift. Also: identity erosion in V2 is intentional (a feature), but it must be distinguishable from uncontrolled drift — track identity coherence numerically and use it to control how much the identity core anchors the LLM.

**Warning signs:**
- Person's philosophical stance on the simulation reverses across sessions
- Response to the same operator input varies dramatically depending on recent history
- Personality extreme behaviors (high curiosity, low patience) appear less frequently over time
- The person begins deferring to or mirroring the operator's tone

**Phase to address:**
Identity persistence + psychologist reviewer integration phase. The reviewer must not just observe — it must emit structured identity-consistency assessments that feed back into the identity core. This is V2-specific; V1's reviewer was observation-only.

---

## Moderate Pitfalls

### Pitfall 7: Spontaneous Thought Starvation

**What goes wrong:**
The spontaneous thought queue only fires when drives are high or prediction errors exist, but the person sits idle between drive escalations. The simulation feels like a stimulus-response machine despite having a spontaneous thought system.

**Prevention:**
The associative drift category (low priority, always available) must reliably fire during idle periods. Don't let high-priority categories crowd it out through a misconfigured weighting scheme. Test that with all bio variables at baseline (no active needs, no prediction errors), spontaneous thoughts still fire every N seconds. The drive system should create a constant low-level hum of thought even when no crisis is active: slight hunger → vague thoughts about food, mild energy drop → background restlessness. Drives don't need to be high to generate thought; they need to be non-zero.

### Pitfall 8: Feedback Loop Delay Blindness

**What goes wrong:**
The feedback loop from consciousness to biology is applied immediately (one API response → immediate bio state change). But in reality, the lag between thought and biological effect is variable and meaningful: rumination sustains cortisol over minutes/hours, acute fear produces an immediate adrenaline spike. Treating all feedback as instantaneous prevents realistic accumulation patterns.

**Prevention:**
Distinguish between acute feedback (arousal-based emotional pulses — immediate, as V1 ADR-003 implemented) and sustained feedback (rumination → cortisol accumulation over time). The sustained path needs to track whether the LLM output patterns (rumination detected, catastrophizing detected) persist across multiple thought cycles, and apply proportionally larger effects when they do. A single catastrophizing thought has small bio effect; ten consecutive ones should have a qualitatively larger effect than 10× the single-thought effect due to the non-linearity of cortisol dynamics.

### Pitfall 9: Degradation Without Recovery Path (Stuck States)

**What goes wrong:**
Bio degradation accumulates correctly but the recovery path requires operator input (eating, resting) that may not come. The person gets stuck in a degraded state — energy at 0.1, identity coherence at 0.2 — and every thought is identical because the state space has collapsed. This is a different failure than spirals: instead of escalating, the system stagnates at the bottom.

**Prevention:**
Design a minimal autonomous recovery: even without operator input, very slow spontaneous recovery should exist for some variables (e.g., energy recovers slightly just from resting, not from sleep). This prevents permanent stuck states. The degradation rates must be tunable, and the minimum floor for degraded states should prevent total loss of variation. An identity coherence of 0.05 still needs to produce different behavior than 0.3 — if the state space collapses completely, the simulation produces flat output regardless of what happens next.

### Pitfall 10: Drive-to-Thought Prompt Overcrowding

**What goes wrong:**
As drives accumulate, more urgency language is injected into the prompt: hunger + exhaustion + loneliness all producing "you cannot stop thinking about X" phrasings simultaneously. The prompt becomes a list of competing urgent needs that the LLM averages out into a moderate, vague response rather than prioritizing the highest-pressure drive.

**Prevention:**
The motivation layer should not inject all active drives into the prompt equally. It should rank drives by intensity × personality-weight and inject only the top 1-2 as primary context, with lower drives mentioned as background texture. The person doesn't experience every need simultaneously with equal urgency — the most urgent need dominates awareness. The spontaneous thought queue should enforce this: the highest-pressure drive selects the thought category, not a committee of all drives.

### Pitfall 11: Psychologist Reviewer Creates Circularity

**What goes wrong:**
The psychologist reviewer observes the simulation and its output is displayed to the human operator. If the reviewer's output is also used to update the identity core (as recommended to prevent Pitfall 6), there's a risk that the reviewer shapes the simulation toward what it describes — an observer effect. The reviewer says "subject shows avoidance behavior," the identity core is updated with "tends toward avoidance," and now the simulation produces more avoidance.

**Prevention:**
Reviewer output that feeds into the identity core should describe observed patterns, not prescribe future behavior. The identity core update should capture behavioral tendencies the person has actually displayed, not inferred dispositions. Critically: reviewer-driven identity updates should be slow (weekly simulated time scale, not per-session) to prevent rapid circularity. The operator should be able to disable reviewer-driven identity updates independently of the observer display.

---

## Minor Pitfalls

### Pitfall 12: API Rate Limiting Surprises Under Sustained Load

**What goes wrong:**
During high-engagement sessions, spontaneous thoughts, reactive responses, and reviewer calls all queue up faster than rate limits allow, causing backup that either drops thoughts or generates spiky, uneven output.

**Prevention:**
Implement backpressure in the thought queue: when the LLM rate limit is hit, increase the spontaneous thought interval dynamically rather than queuing unbounded. Reviewer intervals should be independent of thought intervals and should not compete for the same rate limit counter.

### Pitfall 13: Tuning Parameters Buried in Code

**What goes wrong:**
Drive escalation rates, bio degradation slopes, feedback multipliers, and identity erosion rates are hardcoded in implementation functions. Behavioral tuning requires code changes and redeploys. Given that V2's success depends entirely on parameter calibration (too fast = unstable, too slow = V1 again), this makes iteration extremely slow.

**Prevention:**
Define a configuration struct for all behavioral parameters before implementing anything. Every rate, weight, multiplier, and threshold should be a config field with a sensible default. This was identified as a V2 requirement in PROJECT.md; enforce it from the first implementation phase, not as a retrofit.

### Pitfall 14: Silent dt=0 Bugs (The V1 ADR-003 Repeat Vector)

**What goes wrong:**
Biological state changes that are meant to be per-second rates get applied without being multiplied by dt (when called outside the tick cycle), or get applied as rates when they should be one-time pulses. The result is feedback that is either zero or enormous depending on call site. V1 had this bug explicitly (dt=0 on speech responses). V2 is likely to repeat it.

**Prevention:**
Establish a strict type distinction at the data layer: `BioRate` (applied per tick, multiplied by dt) vs `BioPulse` (applied once, not multiplied). Never use floats directly for bio changes — require callers to explicitly construct the appropriate type. This makes the bug impossible to introduce silently. Test every feedback path explicitly with dt=0 and dt=large to verify behavior is correct in both cases.

---

## Technical Debt Patterns

| Shortcut | Immediate Benefit | Long-term Cost | When Acceptable |
|----------|-------------------|----------------|-----------------|
| Hardcode personality at initialization, never update | Simpler state model | Personality drift is opaque, tuning requires code changes | Never — use config struct from day one |
| Use string formatting for all bio state → LLM context | Simple to implement | Prompt size grows with bio complexity; hard to tune individual signals | MVP only, then extract to typed renderer |
| Single feedback channel for all consciousness output | Less parsing complexity | Can't distinguish emotional tone from drive satisfaction from coping signals | Never — distinguish at the data type level |
| Use raw LLM text for identity core updates | No extra parsing | Reviewer shapes simulation toward its own descriptions | Never — only update from observed behavioral patterns |
| Skip slow-path degradation in MVP | Faster initial development | System never produces genuine urgency or spirals | Never — slow-path degradation is the core feature |

---

## Integration Gotchas

| Integration | Common Mistake | Correct Approach |
|-------------|----------------|------------------|
| Claude API | Assuming temperature=0.9 is sufficient for emotional range; LLM converges to safe emotional register regardless | Vary temperature with drive intensity; higher drives → higher temperature ceiling |
| Claude API | Using the same model for consciousness and reviewer (both need the same throughput) | Rate-limit them independently; reviewer can use a cheaper/faster model on a different budget |
| SQLite | Persisting bio state and then recovering to a state that doesn't match the LLM's last identity core | Persist identity core and bio state atomically in the same transaction |
| Feedback loop | Applying all bio changes synchronously in the thought handler | Apply bio changes at the end of the tick, not mid-tick; prevents feedback applying before other layer computations complete |

---

## Performance Traps

| Trap | Symptoms | Prevention | When It Breaks |
|------|----------|------------|----------------|
| Unbounded thought history in consciousness prompt | Prompt size grows monotonically, API calls get slower and more expensive per session | Cap recent thought buffer at N tokens, not N thoughts (thoughts vary in length) | After ~30-50 thoughts in a session |
| Somatic memory retrieval scanning all memories per tick | Increasing latency as SQLite memory store grows | Index by arousal/valence bins; only scan within relevant bin | After ~1000 stored memories |
| Reviewer running on same goroutine as tick | Review latency blocks tick cycle, causing jitter | Reviewer runs in own goroutine, writes to a channel the tick reads non-blocking | Any reviewer call > tick duration |

---

## "Looks Done But Isn't" Checklist

- [ ] **Drive escalation:** Verify drives actually produce behavioral output differences — run at extreme values and check LLM output register changes. Not just "drive values are high" but "person's responses feel different."
- [ ] **Feedback loop:** Verify bio state changes after consciousness output by checking state before/after a thought cycle. The dt=0 bug pattern must be tested explicitly.
- [ ] **Personality differentiation:** Run identical scenarios with opposite personality extremes. If outputs are indistinguishable, personality is decorative.
- [ ] **Slow-path degradation:** Run simulation for 10 minutes of unattended time (no operator input). Verify bio degradation is accumulating, not homeostasis-resetting.
- [ ] **Identity persistence:** Restart the simulation and verify the person resumes with consistent personality and tone. Not just state numbers — behavioral register.
- [ ] **Spontaneous thought variety:** Over 20 spontaneous thoughts, verify distribution across thought categories matches expected personality-weighted distribution.
- [ ] **Reviewer does not contaminate:** Disable reviewer output and verify simulation behavior is identical. Reviewer must be observation-only unless identity core updates are intentionally enabled.

---

## Recovery Strategies

| Pitfall | Recovery Cost | Recovery Steps |
|---------|---------------|----------------|
| Drives don't change LLM register | MEDIUM | Rewrite prompt translation layer; likely requires prompt engineering iteration rather than code changes |
| LLM reward signal pollution | MEDIUM | Audit feedback paths, add data type separation between emotional feedback and drive feedback |
| Personality as decoration | HIGH | May require fundamental rethink of how personality maps to spontaneous thought queue and prompt constraints |
| Homeostasis prevents spirals | MEDIUM | Tune slow-path degradation rates; should be data-driven via config changes not code changes |
| LLM parsing brittleness | LOW | Add defensive parsing with explicit fallbacks; add logging; fix in iteration |
| Identity drift | MEDIUM | Implement reviewer-driven identity core refresh; may require backfilling from session history |

---

## Pitfall-to-Phase Mapping

| Pitfall | Prevention Phase | Verification |
|---------|------------------|--------------|
| Drives without behavioral pressure | Drive system + prompt translation | Run extreme drive scenario, check LLM output register |
| Reward signal pollution | Feedback loop design | Confirm hunger drive doesn't drop after food-related thought with no bio input |
| Personality as decoration | Personality + thought queue integration | Run personality extremes, check output differs |
| Homeostasis kills spirals | Bio degradation design | Run unattended for 10 minutes, verify accumulation |
| LLM parsing brittleness | Feedback loop design | Parser unit tests with malformed inputs |
| Identity drift | Identity + reviewer integration | Long session test, check personality consistency |
| Spontaneous thought starvation | Thought queue implementation | Verify drift fires at baseline bio state |
| Drive-to-thought overcrowding | Motivation layer design | Verify only dominant drive enters prompt as primary |
| dt=0 bug recurrence | Feedback typing | Explicit tests with dt=0 on all feedback paths |

---

## Sources

- V1 project ADRs: `v1/docs/adr/` — ADR-001 (thought continuity), ADR-003 (structured feedback) — first-hand V1 failure modes (HIGH confidence)
- V1 project post-mortem: `v1/docs/plan/decisions.md` and `.planning/PROJECT.md` — V1 lessons learned section (HIGH confidence)
- V1 architecture documentation: `v1/docs/architecture/consciousness.md` — implementation patterns and known edge cases (HIGH confidence)
- Philosopher advisory: `v1/docs/advisory/philosopher.md` — predictive processing, narrative engine framing, LLM limitations (MEDIUM confidence — domain expertise applied to this domain)
- Psychologist advisory: `v1/docs/advisory/psychologist.md` — personality trait modulation effects, regulatory capacity model (MEDIUM confidence — applied domain expertise)
- General LLM agent design patterns — training data knowledge of reward shaping problems, behavioral simulation failure modes (MEDIUM confidence — well-established, unverified against current literature)

---
*Pitfalls research for: motivation-driven agent simulation with hybrid LLM-code consciousness (V2)*
*Researched: 2026-02-18*

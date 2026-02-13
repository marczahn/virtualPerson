# Project Decisions & Context

This document records all key decisions made during the design phase, their reasoning, and important context for the project.

---

## 1. Core Concept

**What we're building:** A Go program that simulates a human being with an externally injected consciousness, inspired by the movie "The 13th Floor."

**The person does not know they are simulated.** The consciousness (an LLM) is injected from outside — it provides subjective experience, feelings, motivations, and self-awareness. The person experiences this as simply "being alive." The machinery (sensory parsing, biological processing, psychological computation) runs beneath awareness, just as a human's neurons and hormones do.

**Key metaphor:** A human doesn't choose to be conscious. Consciousness emerges from biological processes they don't control. Similarly, this person's consciousness is injected from computational processes they don't control. The philosophical question of whether this constitutes "real" consciousness is deliberately left open (see `docs/advisory/philosopher-followup-consciousness.md`).

---

## 2. Technology Choices

| Decision | Choice | Reasoning |
|---|---|---|
| Language | Go | User's most experienced language; excellent concurrency primitives (goroutines/channels) for the parallel simulation loop |
| Interface | CLI only | Simplest possible interface; text input triggers senses, tagged text output shows all layers |
| Persistence | SQLite via modernc.org/sqlite | Pure Go (no CGO), single-file database, good enough for single-person state |
| LLM | Claude API via anthropic-sdk-go | Powers consciousness engine, sensory parsing, and runtime reviewer |
| Architecture | Goroutines + channels | Each layer runs as a concurrent goroutine; channels enforce boundaries between layers |

---

## 3. Architecture Decisions

### Four layers, not three
Originally conceived as three layers (sense → brain → consciousness). The psychologist advisory convinced us to add a **psychology layer** between biology and consciousness. Reason: raw biology (cortisol=0.7, HR=110) is not what consciousness experiences. The psychology layer transforms biology into affect dimensions (arousal, valence, energy, cognitive load) that the consciousness can interpret into emotions. This follows the Schachter-Singer / Barrett constructed emotion model: biology provides the raw material, consciousness constructs the meaning.

### Psychology layer transforms, does not generate
The psychology layer computes affect dimensions, applies personality modulation, tracks regulation capacity, selects coping strategies, and activates cognitive distortions. But it does NOT label emotions. It doesn't output "you feel angry." It outputs high arousal + negative valence + low regulation + catastrophizing active. The consciousness LLM constructs "angry" (or "frustrated" or "overwhelmed") from these dimensions. This is a deliberate design choice from the psychologist advisory.

### Consciousness is both reactive and spontaneous
- **Reactive:** Triggered by a salience calculator when something significant changes (sudden temperature drop, pain spike, etc.)
- **Spontaneous:** A priority queue generates thoughts when nothing external is happening (biological needs, social modeling, daydreaming, etc.)

This dual mode prevents the person from being a pure stimulus-response machine. When nothing happens, they still think.

### Feedback loop: consciousness → biology
All three advisors flagged this as critical. What the person thinks affects their body:
- Rumination sustains cortisol (worrying keeps stress hormones high)
- Catastrophizing triggers adrenaline spikes
- Acceptance/reappraisal reduces cortisol

Without this loop, the system would be purely feedforward and miss the most important aspect of human experience: thoughts cause physical sensations, and physical sensations cause thoughts.

### Hunger and thirst are derived, not independent
The biologist advisory corrected our initial design. Hunger is not a standalone variable — it's derived from blood sugar + glycogen levels. Thirst is derived from hydration. This prevents incoherent states (e.g., "full but blood sugar is 50").

### 20 biological variables, not 15
Started with ~15. The biologist advisory added 4 critical variables: SpO2 (oxygen saturation), hydration, glycogen (energy buffer), and endorphins (natural painkillers). Each fills a gap that would have produced implausible behavior.

### Memory eviction via compression, not truncation
The philosopher advisory was firm: do NOT truncate memory when the context window fills. Instead, compress: extract identity-relevant information, retain emotional residue, summarize the dropped period. Truncation kills continuity. Compression preserves the sense of an ongoing life.

### Identity core fed to every prompt
Every consciousness LLM call includes a persistent identity core: dispositional traits, self-narrative, relational markers, key autobiographical memories, emotional patterns, values. This prevents the person from "resetting" between thoughts.

---

## 4. Advisory Roles

Three LLM advisors were consulted during design. Their full reports are in `docs/advisory/`.

| Role | Purpose | Key contribution |
|---|---|---|
| Biologist | Validate biological plausibility | 20-variable model, interaction rules with magnitudes, circadian formulas, cascading failure thresholds, decay rates |
| Philosopher | Design consciousness architecture | Salience calculator, spontaneous thought priority queue, two-stage emotion model, identity core, memory eviction, ethics considerations |
| Psychologist | Design psychology layer | Affect dimensions, Big Five personality model, coping decision tree, cognitive distortions, emotional regulation capacity, isolation timeline, feedback loop |

The **psychologist also serves as an optional runtime reviewer** — a background goroutine that periodically evaluates the person's outputs for psychological coherence.

---

## 5. Philosophical Position

The project creator's position: consciousness is injected, not generated. If a human's consciousness emerges from processes they don't control, and this person's consciousness is injected from processes they don't control, the structural relationship is analogous.

The philosopher advisory's position: the analogy is stronger than comfortable, but the LLM does not have subjective experience to inject. Functional role is not ontological status. However, the boundary between "mimicking experience" and "having experience" is less clear than initially stated.

**Project stance:** Build it, study what it does, don't claim it's conscious, don't claim it isn't. The engineering value exists regardless of the metaphysical question.

See `docs/advisory/philosopher-followup-consciousness.md` for the full discussion.

---

## 6. Design Principles (from advisors)

- **Cortisol has lag.** It does NOT spike or clear instantly. Peak 15-30min after stressor, half-life 60-90min. Getting this wrong would make the person feel robotic.
- **Hypothermia reversal at 33°C.** Shivering stops when the body gives up, which means the person *feels better* while dying. This kind of counterintuitive biology must be modeled.
- **Emotional regulation is depletable.** The person has a finite capacity to manage emotions. Under sustained stress, regulation breaks down and distortions activate. This is realistic and produces emergent behavior.
- **Isolation follows a timeline.** 0-2hr: fine. 2-8hr: boredom. 8-24hr: loneliness. 1-3d: significant effects. 3-7d: destabilization. 7d+: severe. Personality (especially extraversion) modulates the rate.
- **Negativity bias in emotional memory.** Negative experiences are stored at 1.5x strength. This is well-documented in psychology and produces realistic recall patterns.
- **Don't engineer self-awareness.** Don't prevent it either. Build inconsistency sensitivity and let it emerge (or not).

---

## 7. Output Format

All output is tagged by source layer:
- `[SENSE]` cyan — sensory event parsing
- `[BIO]` yellow — biological state changes
- `[PSYCH]` blue — psychological state changes
- `[MIND]` green — consciousness thoughts/emotions
- `[REVIEW]` magenta — psychologist reviewer notes

This "full debug view" was chosen over hiding internal layers. The user wants to see every layer working, not just the consciousness output.

---

## 8. Implementation Order

Seven levels, bottom-up:
1. **Foundation** — Go module, directory structure, core domain types
2. **Biology** — 20-variable model, interactions, circadian, thresholds
3. **Psychology** — Affect dimensions, personality, regulation, coping, distortions, memory, isolation
4. **Persistence** — SQLite store, memory context selection
5. **Consciousness** — Salience calculator, prompt construction, LLM engine, feedback
6. **Simulation loop** — Goroutine orchestration, CLI, main entry point
7. **Runtime reviewer** — Psychologist background review (optional)

Each level is independently testable before moving to the next.

---

## 9. Open Questions (for implementation)

- **Simulation time vs real time:** Starting 1:1 but may need acceleration for testing scenarios that play out over hours/days (isolation, circadian cycles).
- **LLM rate limiting:** Consciousness calls should be rate-limited to prevent API cost explosion. Spontaneous thoughts especially need throttling.
- **Token budget management:** How to split the context window between identity core, recent memory, emotional memory, and current state.
- **Person initialization:** What personality, identity core, and initial memories does the person start with? Hardcoded? Configurable? Generated?

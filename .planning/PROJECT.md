# Person Simulation v2

## What This Is

A Go program that simulates a human being with externally injected consciousness (LLM). The person experiences a body, emotions, and thoughts without knowing they are simulated. V2 is a clean rebuild that addresses v1's core problem: the person reflected deeply but never *wanted* anything. V2 adds intrinsic motivation so the consciousness has something to struggle with, creating emergent behavior that feels alive rather than philosophically empty.

## Core Value

The person must exhibit intrinsic motivation — visible drives, desires, and frustrations that emerge from the interplay between biological state, a reward/motivation system, and LLM-driven consciousness. Without this, everything else is narrative decoration.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Four-layer pipeline: Bio → Motivation → Consciousness/Reflexion → Feedback → Bio
- [ ] Reduced biological model (~8-10 variables chosen for motivational relevance, not bio-fidelity)
- [ ] Intrinsic motivation/drive system computed from bio state + personality factors
- [ ] Hybrid consciousness: code computes drives, LLM interprets and modulates them via feedback
- [ ] 7 personality factors replacing Big Five (stress sensitivity, energy resilience, curiosity, self-observation, patience/frustration tolerance, risk aversion, social factor)
- [ ] Bio degradation over time when the person doesn't engage (energy drops, stress rises, cognitive capacity shrinks)
- [ ] Identity/memory erosion without engagement (coherence weakens, memories fade)
- [ ] Consequence system designed to be adjustable (tunable degradation rates)
- [ ] Spontaneous thought system driven by motivation state (ticker-based, like v1)
- [ ] External input handling: speech, actions, environment changes (same conventions as v1)
- [ ] Feedback loop: consciousness output modulates biological state (rumination → cortisol, acceptance → relief, etc.)
- [ ] Noise/variability in bio states to prevent deterministic stagnation
- [ ] Psychologist reviewer agent analyzing psychological state (3-minute tick rate)
- [ ] SQLite persistence for bio state, personality, identity, memories
- [ ] CLI interface with tagged output by source layer

### Out of Scope

- Multi-agent social interaction — person interacts only with human operator, social drive is about that relationship
- Web dashboard / WebSocket server — focus on core simulation first
- i18n / localization — English only for v2
- Circadian rhythm system — v1's sine-wave formulas added bio-fidelity but didn't contribute to motivation or aliveness
- 20-variable biological model — v1 showed that biological detail doesn't produce lifelike behavior without motivation
- Big Five personality model — replaced by 7 factors that directly map to behavioral tendencies

## Context

### V1 Lessons Learned

V1 built a detailed four-layer pipeline (Sense → Biology → Psychology → Consciousness) with 20 biological variables, 76 interaction rules, Big Five personality, cognitive distortions, coping strategies, emotional memory, and isolation effects. 32 source files, 28 test files, 327 tests.

**What worked:**
- The pipeline architecture (layered, dependencies point inward) produced clean, testable code
- The feedback loop (consciousness → biology) created genuine emergent behavior (rumination spirals)
- Spontaneous thoughts via priority queue gave the person an inner life even without input
- Scenario system for environmental context was effective
- Persistence allowed the person to resume with continuity

**What didn't work:**
- The person was too stable, too philosophical, too "sachlich" (matter-of-fact)
- 20 bio variables created complexity without proportional behavioral richness
- Psychology layer (affect dimensions, coping, distortions) was sophisticated but produced controlled, analytical responses — never genuine panic, frustration, or joy
- No intrinsic motivation meant the person reflected on existence without caring about it
- Emotions were narrative (the LLM described feeling X) rather than functional (the system was in state X which drove behavior Y)

**Core insight:** Consciousness = Motivation + Reflexion. V1 had reflexion without motivation, producing a "philosophical vacuum." V2 must have both.

### Architectural Shift

V1: Sense → Biology (detailed) → Psychology (affect dims) → Consciousness (LLM)
V2: Bio (reduced, motivation-serving) → Motivation (drives) → Consciousness (LLM, hybrid) → Feedback → Bio

The key change: the motivation layer sits between bio and consciousness, creating pressure that the LLM must respond to. The LLM doesn't just *observe* its state — it *needs* things.

### How the Hybrid Works

1. **Code computes drives** from bio state + personality factors (algorithmic, deterministic + noise)
2. **Drives are injected into the LLM prompt** as context ("You feel a strong urge to...", "Your energy is dropping and you're becoming restless...")
3. **LLM generates thoughts/responses** influenced by these drives
4. **LLM output is parsed for feedback signals** (emotional state, coping patterns, engagement level)
5. **Feedback modulates bio state** (completing the loop)

The LLM can modulate/override drives — a motivated person can push through fatigue, or a curious person can ignore safety concerns. This creates genuine tension between what the system pushes and what consciousness decides.

## Constraints

- **Tech stack**: Go 1.24, SQLite (modernc.org/sqlite, no CGO), Claude API (anthropic-sdk-go) — same as v1
- **Code lives in `v2/`**: Clean rebuild, no code sharing with v1 (v1 stays as reference in `v1/`)
- **Maintainability over cleverness**: Clear code, strict separation of concerns, every module has one job. Prefer readable over "smart."
- **Minimal dependencies**: Standard library where possible, justify every external dep
- **Testability**: Every behavior-changing component needs tests. Test behavior, not implementation.
- **Adjustability**: Degradation rates, personality factors, bio parameters, motivation weights — all must be tunable without code changes (config or constants, not buried in logic)

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Clean rebuild in v2/ | V1's architecture was optimized for bio-fidelity, not motivation. Rebuilding avoids carrying over legacy patterns that don't serve v2's goals. | — Pending |
| Replace Big Five with 7 motivation-serving personality factors | Big Five modulated affect dimensions but didn't directly drive behavior. New factors (stress sensitivity, curiosity, patience, etc.) map directly to how the person acts. | — Pending |
| Reduce bio model to ~8-10 vars | 20 vars created complexity without proportional behavioral value. Keep only vars that feed the motivation system meaningfully. | — Pending |
| Hybrid motivation (code + LLM) | Pure algorithmic motivation = scripted. Pure LLM motivation = unreliable. Hybrid: code creates pressure, LLM decides what to do with it. | — Pending |
| Bio degradation + identity erosion as consequences | Without irreversible consequences, the person has no stakes. Degradation creates genuine urgency. Made adjustable for tuning. | — Pending |
| Psychologist reviewer at 3-min intervals | Reduced from v1's 60s to 180s. Still valuable as meta-observer but less frequent to reduce API costs. | — Pending |

---
*Last updated: 2026-02-18 after initialization*

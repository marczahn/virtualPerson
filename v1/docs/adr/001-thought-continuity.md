# ADR-001: Thought Continuity via Recent Thought Buffer

## Status

Accepted

## Date

2026-02-14

## Question

How can we make the consciousness layer's output less predictable and more emergent? The current architecture generates each thought in isolation — the LLM has no awareness of what it just thought. This produces repetitive, structurally identical responses regardless of what came before.

## Context

The current flow is:

```
[psych state + trigger] → LLM → thought → ParseFeedback (keyword scan) → tiny bio delta
```

Each thought starts from nearly the same position. The feedback loop exists but is shallow: `ParseFeedback` does string matching for keywords like "catastroph" or "my fault", producing small biological deltas (±0.01–0.02). A thought about spiraling anxiety doesn't meaningfully make the person more anxious for the next thought.

The philosopher consultant's document identifies two relevant frameworks:
- **Predictive processing** (Section 2): the mind runs simulations, predicting futures and rehearsing scenarios. Thoughts build on each other.
- **Constructed emotion** (Section 3): emotion is constructed by appraisal in context, not by lookup. The same biological state produces different thoughts depending on what the person was just thinking.

Two options were considered:

### Option 1: Feed last N thoughts back into the prompt

Maintain a ring buffer of recent thoughts. Include them in the prompt as "what you've been thinking." The LLM builds on its own output — a worried thought leads to a more worried thought or self-correction, naturally.

- Low cost, high impact
- No new dependencies or API calls
- Risk of runaway loops (mitigated by existing biological homeostasis and reviewer)

### Option 2: LLM-driven feedback assessment

Replace `ParseFeedback` keyword scanning with a second LLM call that asks the model to assess its own thought's emotional impact. More accurate than regex, but more expensive and harder to test.

- Higher cost per thought (double API calls)
- More accurate feedback than keyword matching
- Should be evaluated after Option 1 is in place

## Decision

Implement Option 1 first. Add a ring buffer of recent thoughts to the consciousness engine. Include the last 5 thoughts in every prompt as "what you've been thinking recently." Evaluate the impact on output quality before deciding on Option 2.

Option 2 remains a candidate for a future ADR if the keyword-based feedback proves insufficient after thought continuity is in place.

## Consequences

- Thoughts will compound: negative spirals and self-correction become possible
- The reviewer (Level 7) becomes more important as a circuit breaker for runaway loops
- Token usage per prompt increases slightly (5 short thoughts ≈ 150-250 tokens)
- Biological homeostasis (decay, recovery) must be strong enough to counteract tight thought loops

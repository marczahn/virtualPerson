# ADR-002: Blind Thought Evaluator for Unconscious Somatic Feedback

## Status

Proposed — awaiting empirical evaluation of ADR-001 (thought continuity buffer)

## Date

2026-02-14

## Question

Should we replace the keyword-based `ParseFeedback` mechanism with a second LLM call that evaluates each thought's emotional impact on the body, without knowing it authored the thought?

## Context

### The current feedback mechanism

After the consciousness layer generates a thought, `ParseFeedback` scans the text for keyword patterns (e.g., "worst case" → catastrophizing, "can't stop thinking" → rumination) and maps matches to small biological state changes (cortisol ±0.02, serotonin ±0.01, etc.). This is crude:

- It misses semantically rich content that doesn't use trigger phrases
- It can't distinguish tone (ironic "everything is fine" vs. genuine contentment)
- It has no awareness of context — the same words mean different things in different situations
- The biological deltas are small and easily overwhelmed by other forces

### The proposed idea

After generating a thought, make a second LLM call with a clean context:

> "A person in [this psychological state] just had this thought: [thought content]. What does this thought reveal about their emotional processing? How would having this thought affect their body?"

The key insight: the evaluator doesn't know it authored the thought. It evaluates as an outside observer — like how the body reacts to its own thoughts without the thinker choosing that reaction. This maps onto Damasio's somatic marker hypothesis: physiological responses to thoughts occur before and independently of conscious awareness.

### Philosophical grounding

The philosopher consultant identifies three pillars of emotional construction (Barrett's theory, Section 3):
1. Current interoceptive state (affect dimensions)
2. Categorization from past experience (memory)
3. Current context

The keyword scanner implements none of these. The blind evaluator could implement all three if given the right context.

The philosopher also draws a clear line (Section 5): don't build self-awareness modules. The blind evaluator respects this — it's infrastructure, not experience. The person doesn't gain insight into why their body reacts. It's like how your stomach tightens when you think "I'll just ignore this problem" — you didn't choose that response.

### Concerns

**Context dependency.** "I'm fine" is suppression after trauma but genuine contentment on a calm afternoon. Without context, the evaluator will hallucinate emotional subtext. But providing full context makes this essentially a second consciousness call — and the "blindness" becomes artificial, since the LLM already has no persistent memory between calls.

**Cost.** An extra API call per thought. Even with a cheap model (Haiku), this doubles the per-thought latency and cost. If thoughts occur every 30 seconds, this is significant.

**Runaway spirals.** If the evaluator consistently reads thoughts as negative and pushes biology negative, the next thought generates from a worse state, producing a worse thought, which evaluates worse. The homeostatic mechanisms and the reviewer must be strong enough circuit breakers.

**Evaluator hallucination.** The keyword matcher is crude but predictable and testable. An LLM evaluator is more capable but less deterministic. Bad assessments compound through the feedback loop.

**Premature optimization.** ADR-001 (thought continuity buffer) was just implemented. The thought buffer may already solve the predictability problem by a different mechanism: the LLM sees its own recent thought trajectory in the prompt, so it naturally escalates or self-corrects without needing the feedback loop to be more intelligent. The keyword-based feedback may no longer be the bottleneck.

## Decision

Defer implementation until ADR-001 has been evaluated empirically.

### Criteria for revisiting

Implement the blind evaluator if, after running with the thought buffer:
1. Output still feels predictable despite varied psychological states
2. The keyword feedback is demonstrably the weak link (e.g., thoughts with clear emotional content produce no biological response because keywords don't match)
3. The thought buffer alone doesn't create natural compounding / self-correction

### If implemented, the design should

- Use a cheap, fast model (Haiku) for the evaluation call
- Return structured output: detected coping strategies, distortions, and directional bio impacts
- Include current psychological state as context (without it, evaluation is too ambiguous)
- NOT include identity or memory (keeps it closer to unconscious somatic response)
- Replace `ParseFeedback` entirely rather than layering on top
- Include rate limiting to prevent excessive API calls during high-frequency thought periods

## Consequences

### If we proceed later
- More accurate feedback loop, better modeling of unconscious emotional processing
- Higher API cost and latency per thought
- Need robust spiral detection / circuit breaking
- `ParseFeedback` and its tests can be removed

### If we don't proceed
- Keyword-based feedback remains the weakest link in the pipeline
- But the thought buffer may compensate sufficiently
- Lower cost, simpler architecture, fully deterministic feedback

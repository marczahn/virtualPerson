# ADR-003: Structured Emotional Feedback (Consciousness → Biology)

**Status:** Accepted
**Date:** 2026-02-18

---

## Context

Two defects prevented the consciousness layer from meaningfully affecting biological state:

### Defect 1 — dt=0 bug

`applyFeedback(thought, dt)` is called with `dt=0` for all speech and action responses (they fire outside the tick cycle). The implementation multiplied every delta by `dt`, so dt=0 → zero effect. Emotional responses to conversation never reached biology.

### Defect 2 — Brittle keyword heuristic

`ParseFeedback` scanned LLM output for exact keyword phrases to detect coping strategies and distortions. LLM output rarely matched the expected phrases. Even when it did, `FeedbackToChanges` produced only tiny cortisol/serotonin nudges — no acute arousal signal. A frightening or angering thought had no effect on heart rate, muscle tension, or adrenaline.

### Design Question

How should consciousness report its emotional state to biology?

**Option A — Keep keyword matching, fix magnitude**
Expand the keyword lists and increase deltas. Fragile: depends on exact phrasing in a specific language; the LLM has no incentive to use the target phrases. Does not fix the structural issue (no arousal signal).

**Option B — LLM-emitted structured annotation**
Instruct the LLM to append a machine-readable tag `[STATE: arousal=X, valence=Y]` to every response. Parse it, strip it before storage, map to calibrated biological pulses. Fixes both problems: the tag is reliable (the LLM follows explicit format instructions), and the mapping covers all acute signals.

**Option C — Separate classification call**
After every LLM response, make a second API call to classify emotional state. Doubles API cost and latency. Not justified when the emitting model can do it inline.

---

## Decision

Adopt **Option B**. The LLM emits `[STATE: arousal=X, valence=Y]` at the end of every response. This annotation is:

- Stripped from `Thought.Content` before storage and before inclusion in future prompts
- Parsed into an `EmotionalTag{Arousal, Valence}` stored in `ThoughtFeedback.EmotionalState`
- Mapped by `EmotionalPulses()` to calibrated one-time biological state changes

The dt-scaling bug is fixed by removing dt from `applyFeedback` entirely. Emotional pulses are absolute one-time events, not per-second rates.

The existing keyword-based `FeedbackToChanges` path is retained. Its deltas (e.g. rumination +0.02 cortisol) are also applied as absolute one-time pulses now — their magnitudes are appropriate for single-thought events.

---

## Consequence

### Biological calibration

The pulse magnitudes are chosen so that:

- A single angry thought (arousal=0.8, valence=-0.7) noticeably raises cortisol (+0.04), muscle tension (+0.12), and heart rate (+10.8 bpm) — but does not cross any threshold alone.
- Five consecutive angry thoughts bring cortisol above the 0.3 tachycardia threshold.
- Extreme fear (arousal>0.8, valence<-0.5) triggers adrenaline, but only at magnitudes safe below the 12 bpm/s cascade threshold.
- Positive excited states (posVal>0.2, arousal significant) trigger a dopamine reward signal.
- Natural biological decay (cortisol half-life 4500s, adrenaline 150s, heart rate 75s) prevents unbounded escalation.

### Annotation leakage risk

The tag is stripped in `engine.go` immediately after the API response, before `recordThought` stores it. Future prompts include `Thought.Content` (clean), never the raw response. The system prompt instructs the LLM not to explain the tag. Risk of leakage is low.

### Keyword path unchanged

`ParseFeedback` still runs keyword detection and returns `ActiveCoping` / `ActiveDistortions`. `FeedbackToChanges` still maps these to changes. The tag path is additive, not a replacement.

### Test coverage

All new behaviour is covered by unit tests:
- `ParseEmotionalTag`: valid, absent, malformed, stripping
- `EmotionalPulses`: angry, extreme fear, positive calm, zero tag, adrenaline safety
- `SystemPrompt`: annotation instruction present, no forbidden simulation words
- `Engine`: tag stripped from content, `EmotionalState` populated
- `Loop.applyFeedback`: elevated cortisol after speech with emotional tag (dt=0 bug regression test)

---

## Files Changed

| File | Change |
|------|--------|
| `internal/consciousness/thought.go` | Add `EmotionalTag`; extend `ThoughtFeedback` |
| `internal/consciousness/feedback.go` | Add `ParseEmotionalTag`, `EmotionalPulses`; update `ParseFeedback` signature |
| `internal/consciousness/engine.go` | Strip tag before storing thought |
| `internal/consciousness/prompt.go` | Append annotation instruction to system prompt |
| `internal/i18n/types.go` | Add `EmotionalAnnotation` field |
| `internal/i18n/translations/en.yaml` | Add annotation instruction text |
| `internal/i18n/translations/de.yaml` | Add annotation instruction text |
| `internal/simulation/loop.go` | Remove dt from `applyFeedback`; apply `EmotionalPulses` |

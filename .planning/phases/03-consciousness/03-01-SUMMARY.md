---
phase: 03-consciousness
plan: "01"
subsystem: consciousness
tags: [go, consciousness, parser, contracts, tdd]

requires:
  - phase: 02-02
    provides: "Deterministic motivation outputs"
provides:
  - "Prompt context translation from motivation state to felt language"
  - "Top-2 primary drive selection with background drive context"
  - "Implicit active-goal pull phrasing"
  - "Defensive parser for STATE/ACTION/DRIVE tags with fallback"
affects: [03-02]

tech-stack:
  added: []
  patterns:
    - "Pure function contract layer"
    - "Deterministic urgency ranking with tie priority"
    - "Parser fallback to prior known-good state"

key-files:
  created:
    - "v2/internal/consciousness/consciousness.go"
    - "v2/internal/consciousness/consciousness_test.go"
    - ".planning/phases/03-consciousness/03-01-PLAN.md"
    - ".planning/phases/03-consciousness/03-01-SUMMARY.md"
  modified:
    - ".planning/REQUIREMENTS.md"
    - ".planning/STATE.md"
    - ".planning/CONTINUITY.md"
    - ".planning/NEXT_SESSION_START.md"
    - ".planning/ROADMAP.md"

key-decisions:
  - "Prompt context uses phenomenological text only; no numeric drive values surfaced"
  - "Primary prompt context includes exactly top two ranked drives"
  - "Active goal is rendered as implicit pull text, not command syntax"
- "Parser requires valid STATE and ACTION tags; missing/malformed required tags preserve prior parse"
  - "DRIVE tag is optional; valid tags replace overrides, missing tag clears this-turn overrides, malformed tag triggers full fallback"

requirements-completed: [CON-01, CON-02, CON-03, CON-04, CON-05, CON-06, CON-07, CON-08]

completed: 2026-02-19
---

# Phase 3 Plan 01: Summary

Implemented first consciousness contracts in `v2/internal/consciousness` via tests-first workflow.

## Red phase (tests first)

Added failing contract tests for:

- top-two drive selection as primary prompt context
- phenomenological language without raw numeric exposure
- implicit active-goal pull phrasing
- valid parsing of `[STATE]`, `[ACTION]`, optional `[DRIVE]`
- tag stripping from narrative content
- fallback to prior known-good parse when required tags are missing/malformed

## Green phase implementation

Added `v2/internal/consciousness/consciousness.go`:

- `BuildPromptContext(motivation.MotivationState) PromptContext`
- deterministic ranking of five drives by urgency with tie priority
- felt-language mapping per drive by urgency band
- implicit pull phrase mapping by active goal drive
- `ParseResponse(raw, prior) ParsedResponse`
- parsing for `[STATE: arousal=X, valence=Y]`
- parsing for `[ACTION: type]`
- optional parsing for `[DRIVE: name=value]`
- stripping all supported tags from narrative output
- defensive all-or-nothing fallback behavior to prior known-good parse on malformed tags

## Verification

- `cd v2 && GOCACHE=/tmp/go-build go test ./internal/consciousness/... -count=1` -> pass
- `cd v2 && GOCACHE=/tmp/go-build go test ./... -count=1` -> pass

## Architecture check

- Layer impact: Consciousness consumes Motivation outputs.
- Boundary preserved: `internal/consciousness` imports `internal/motivation`; no reverse dependency introduced.
- No new dependencies added.

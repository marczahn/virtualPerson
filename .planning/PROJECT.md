# Person Simulation v2

## What This Is

A Go simulation of a virtual person where body dynamics and drive pressure are computed in code, and consciousness is externally injected through an LLM.

## Core Value

The person must show intrinsic motivation that leads to observable behavior changes over time.

## Source Of Truth

Planning priority is determined by the current roadmap mix in `.planning/` (requirements + roadmap + phase plans), with `v1/` used only for lessons and conventions.

## Active Requirements (summary)

- Four-layer loop: Bio -> Motivation -> Consciousness/Reflection -> Feedback -> Bio
- Reduced motivationally-relevant bio model
- Deterministic motivation computation with personality modulation
- Hybrid consciousness with structured feedback tags
- Spontaneous thoughts and continuity
- External input and scenario injection
- Layer-tagged CLI output
- Executable runtime entrypoint
- Centralized configuration of all tunable parameters
- Explicit initial profile/start-situation contract including stress-reactivity

## Out of Scope

- Multi-agent interaction
- Web dashboard / WebSocket server
- i18n/localization
- 20-variable bio model
- Big Five personality model

## Constraints

- Tech stack: Go 1.26 target, SQLite (`modernc.org/sqlite`), Anthropic SDK
- Code location: `v2/` only for new implementation
- Architecture: strict inward dependencies, infrastructure orchestrates
- Maintainability over performance
- Minimal dependencies
- Test-first for behavior changes

## Key Decisions

| Decision | Outcome |
|----------|---------|
| Continue from implemented code state | Accepted (2026-02-19) |
| Use `internal` package layout | Accepted (2026-02-19) |
| Go toolchain target 1.26 | Accepted (2026-02-19) |
| Add runtime entrypoint requirement (`INF-08`) | Accepted (2026-02-20) |
| Add explicit initial profile requirements (`PRF-01..PRF-03`) | Accepted (2026-02-20) |

---
*Last updated: 2026-02-20*

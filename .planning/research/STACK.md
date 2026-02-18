# Stack Research

**Domain:** Motivation-driven consciousness simulation in Go
**Researched:** 2026-02-18
**Confidence:** HIGH for core stack (verified against pkg.go.dev), MEDIUM for API pattern details (verified against GitHub docs), LOW for reward-system design patterns (no Go-specific library exists — must implement from scratch)

---

## Recommended Stack

### Core Technologies

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| Go | 1.24.2 (already in use) | Application language | Swiss table maps improve state iteration performance. `testing.B.Loop` cleans up benchmarks. `slices` package (std lib since 1.21) replaces `golang.org/x/exp/slices`. No reason to change — v1 proves this works. |
| modernc.org/sqlite | v1.46.1 (latest as of 2026-02-18; v1 uses v1.45.0) | Persistence — bio state, personality, memories | CGO-free pure-Go SQLite. Only viable option for CGO-free SQLite in Go. The alternatives (mattn/go-sqlite3) require CGO which complicates cross-compilation. No behavioral change between v1.45 and v1.46 for this use case — upgrade is safe. |
| github.com/anthropics/anthropic-sdk-go | v1.24.0 (latest; v1 uses v1.22.1) | LLM integration — consciousness and reviewer layers | Official SDK, only supported Go client for Anthropic API. v1.22.1 → v1.24.0 adds Claude Sonnet 4.6 constant (`ModelClaudeSonnet4_5_20250929` available from v1.23.0) and `param.NullStruct` bug fix (v1.22.1 already included). API surface for message creation, system prompts, and temperature is stable — no breaking changes in this range. |

### Supporting Libraries

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| github.com/charmbracelet/bubbletea | v1.3.10 | TUI client (optional) | Only if a WebSocket + TUI display mode is added. V2 PROJECT.md explicitly defers web dashboard, but bubbletea is already in go.sum from v1 and costs nothing to keep available. Do not pull it into v2 core. |
| github.com/charmbracelet/bubbles | v1.0.0 (latest; v1 uses v0-series pinned to v1.0.0 now) | TUI component library | Only alongside bubbletea for TUI client. Not in v2 core. |
| github.com/charmbracelet/lipgloss | v1.1.0 | TUI styling | Only alongside bubbletea. Not in v2 core. |
| nhooyr.io/websocket | v1.8.17 | WebSocket transport | Deferred — v2 out-of-scope per PROJECT.md. Keep in dependency tree only when server mode is added. |
| github.com/google/uuid | v1.6.0 | Memory/episode ID generation | Use for episodic memory IDs. No alternatives needed — this is the standard. |
| gopkg.in/yaml.v3 | v3.0.1 | Config file parsing | Only if personality/degradation parameters are stored in YAML config files. If Go struct literals with exported constants are sufficient (simpler, type-safe), skip entirely. |

### Development Tools

| Tool | Purpose | Notes |
|------|---------|-------|
| `go test ./...` | Test runner | Standard Go toolchain. No test framework needed — stdlib testing package is sufficient and matches v1's discipline. |
| `go vet ./...` | Static analysis | Already in v1 workflow. Catches structural errors before runtime. |
| `golangci-lint` | Linting | Optional but recommended. Run `errcheck`, `staticcheck`, `govet`. Avoid enabling too many linters — marginal value vs. noise. |

---

## Installation

```bash
# Initialize v2 module (clean rebuild in v2/)
cd v2
go mod init github.com/marczahn/person

# Core: LLM integration
go get github.com/anthropics/anthropic-sdk-go@v1.24.0

# Core: persistence
go get modernc.org/sqlite@v1.46.1

# ID generation (for episodic memory)
go get github.com/google/uuid@v1.6.0

# TUI (defer until server mode is added — do not pull in v2 core)
# go get github.com/charmbracelet/bubbletea@v1.3.10
# go get github.com/charmbracelet/bubbles@v1.0.0
# go get github.com/charmbracelet/lipgloss@v1.1.0
```

---

## Alternatives Considered

| Recommended | Alternative | When to Use Alternative |
|-------------|-------------|-------------------------|
| modernc.org/sqlite | mattn/go-sqlite3 | Never for this project. mattn requires CGO, which breaks cross-compilation and complicates CI. The project constraint says "no CGO" — this is settled. |
| modernc.org/sqlite | zombiezen.com/go/sqlite | Only if you need fine-grained connection pooling or the `zombiezen` API surface. For this project, modernc is simpler and already proven in v1. |
| github.com/anthropics/anthropic-sdk-go | Direct HTTP calls via net/http | Only if the SDK's dependency tree becomes a problem or the API diverges significantly. The SDK is thin (it has a small import graph) and handles auth, retries, and streaming — rolling your own wastes time. |
| stdlib `slices` | golang.org/x/exp/slices | Never — `golang.org/x/exp/slices` is superseded by the standard library `slices` package (Go 1.21+). Use `import "slices"`. V1 uses `golang.org/x/exp` but v2 should not. |
| stdlib `encoding/json` | github.com/json-iterator/go | Never for this project. The person simulation is not a high-throughput JSON service. Stdlib is fast enough and has zero dependencies. |
| Hand-written reward/motivation formulas | External RL library | No Go RL library serves this use case well (they target training agents, not simulating one continuous person-state). Keep reward/motivation as plain Go float64 arithmetic. Readability and adjustability matter more than algorithmic sophistication here. |
| gopkg.in/yaml.v3 | encoding/json for config | Prefer Go constants for tunable parameters (adjustable without file I/O). Use YAML only if the user explicitly needs to edit config files at runtime without recompiling. |

---

## What NOT to Use

| Avoid | Why | Use Instead |
|-------|-----|-------------|
| `golang.org/x/exp/slices` | Superseded by standard library `slices` (Go 1.21+). Adding this dependency is unnecessary and adds to go.sum noise. | `import "slices"` from stdlib |
| Any Go RL/reinforcement-learning library (e.g., GoLearn, Gorgonia) | These libraries target training ML models, not simulating a single continuous stateful agent. They impose computational graph abstractions that are wrong for this domain. The motivation system is plain arithmetic: `drive = f(bio_var, personality_factor)`. | Inline Go float64 arithmetic with named constants |
| `github.com/spf13/viper` | Heavyweight config library that pulls in many transitive deps. Overkill for what amounts to a handful of numeric parameters. | Either typed Go constants, or `encoding/json` + a simple config struct if runtime-editability is needed |
| `github.com/sirupsen/logrus` or `go.uber.org/zap` | Structured logging frameworks are wrong for this domain. The simulation already has a layer-tagged display system. Mixing a log framework into that output layer creates confusion about what is "display" and what is "debug". | `fmt.Fprintf(os.Stderr, ...)` for internal debug traces; the Display type from the simulation output layer for domain events |
| CGO-dependent packages | The project explicitly prohibits CGO. Any package requiring `import "C"` or a C compiler at build time is banned. | Pure-Go alternatives only |

---

## SDK API Patterns for v2 (anthropic-sdk-go v1.24.0)

The API surface changed subtly between v1.22.1 and v1.24.0. V2 should use current patterns, not v1 cargo-copied patterns.

**Creating a message (current pattern):**

```go
import (
    anthropic "github.com/anthropics/anthropic-sdk-go"
    "github.com/anthropics/anthropic-sdk-go/packages/param"
)

resp, err := client.Messages.New(ctx, anthropic.MessageNewParams{
    Model:     anthropic.ModelClaudeHaiku4_5,  // fast + cheap for tick-level calls
    MaxTokens: 1024,
    Temperature: param.NewOpt(0.9),
    System: []anthropic.TextBlockParam{
        {Text: systemPrompt},
    },
    Messages: []anthropic.MessageParam{
        anthropic.NewUserMessage(anthropic.NewTextBlock(userMessage)),
    },
})
```

**Extracting text from response (current pattern):**

```go
for _, block := range resp.Content {
    if block.Type == "text" {
        text := block.AsText().Text
    }
}
```

**Model constants relevant to v2:**

| Constant | Model | Use In V2 |
|----------|-------|-----------|
| `anthropic.ModelClaudeHaiku4_5` | claude-haiku-4-5 (fast, cheap) | Consciousness tick-level calls (high frequency) |
| `anthropic.ModelClaudeSonnet4_5` | claude-sonnet-4-5 | Psychologist reviewer (low frequency, needs depth) |

Note: `ModelClaudeSonnet4_5_20250929` is the dated constant added in v1.23.0. The undated `ModelClaudeSonnet4_5` alias may also exist. Verify against the installed SDK source before using — the MEMORY.md lists both `ModelClaudeHaiku4_5` and `ModelClaudeSonnet4_5` as confirmed working in v1.22.1, which means they remain valid in v1.24.0 (no breaking changes in this range).

---

## Stack Patterns by Variant

**If tunable parameters need runtime editing (e.g., personality sliders without recompile):**
- Use `encoding/json` with a typed `Config` struct loaded at startup
- Store in a `config.json` beside the binary
- Because this project already uses `config.json` in v1 — no new format needed

**If tunable parameters are only changed at compile time:**
- Use exported Go constants in a `params` or `config` package
- Because it keeps the hot path simple: no file I/O, no nil checks, immediate readability in code

**If drive/motivation formulas become complex:**
- Keep as pure functions `func ComputeDrive(bio BioState, p Personality) float64`
- Do NOT reach for a DSL or expression evaluator — the formulas are arithmetic, not user-authored scripts

**If API costs become a concern:**
- Use `ModelClaudeHaiku4_5` for the consciousness layer (fast ticks)
- Use `ModelClaudeSonnet4_5` only for the reviewer (3-minute intervals)
- The rate limiting pattern from v1 (`minInterval` on the engine) carries over identically

---

## Version Compatibility

| Package | Compatible With | Notes |
|---------|-----------------|-------|
| `modernc.org/sqlite@v1.46.1` | Go 1.24.2 | Confirmed — modernc.org/sqlite tracks current Go releases closely. v1.45.0 → v1.46.1 is a routine patch. |
| `anthropic-sdk-go@v1.24.0` | Go 1.22+ (SDK requirement) | Go 1.24.2 satisfies this. |
| `charmbracelet/bubbletea@v1.3.10` | Go 1.24.2 | Confirmed — actively maintained, stable. |
| `google/uuid@v1.6.0` | Go 1.24.2 | No compatibility issues — pure Go, minimal deps. |
| stdlib `slices` package | Go 1.21+ | Available in Go 1.24.2 without any external dependency. |

---

## What V2 Implements Without a Library

The motivation/reward system is the new conceptual core of v2. No external library handles this because the domain is too specific. Implement as plain Go:

**Drive computation:** `f(bio_var, personality_factor) → float64` — pure arithmetic in a `motivation` package.

**Reward signal:** Signed float computed from drive satisfaction. Example: `reward = goal_proximity * personality.Curiosity - energy_cost * personality.EnergyResilience`. No library needed.

**Personality as multipliers:** 7 named float64 fields on a `Personality` struct. Each motivational formula references one or two factors directly. This is simpler than v1's Big Five because it's designed to be used directly, not transformed through a psychology layer.

**Bio state noise:** `state += noise * rand.NormFloat64() * dt` using stdlib `math/rand/v2` (available since Go 1.22, no external dependency). Do not use v1's `golang.org/x/exp/rand`.

---

## Sources

- pkg.go.dev/github.com/anthropics/anthropic-sdk-go — version v1.24.0 confirmed (Feb 18, 2026). HIGH confidence.
- pkg.go.dev/modernc.org/sqlite — version v1.46.1 confirmed (Feb 18, 2026). HIGH confidence.
- pkg.go.dev/github.com/charmbracelet/bubbletea — version v1.3.10 confirmed (Sep 17, 2025). HIGH confidence.
- pkg.go.dev/nhooyr.io/websocket — version v1.8.17 confirmed (Aug 10, 2024). HIGH confidence.
- pkg.go.dev/github.com/google/uuid — version v1.6.0 confirmed (Jan 23, 2024). HIGH confidence.
- github.com/anthropics/anthropic-sdk-go README.md — API patterns for MessageNewParams, system prompts, temperature, text extraction. MEDIUM confidence (docs match v1 codebase usage patterns).
- github.com/anthropics/anthropic-sdk-go CHANGELOG.md — breaking changes v1.22.1 → v1.24.0. MEDIUM confidence (changelog confirms no breaking changes in message creation or model constants in this range).
- go.dev/blog/go1.24 — Go 1.24 feature list. HIGH confidence (official blog).
- pkg.go.dev/golang.org/x/exp/slices — confirmed superseded by stdlib `slices` (Go 1.21+). HIGH confidence.
- v1 codebase (v1/internal/consciousness/claude.go) — confirmed working patterns for SDK usage at v1.22.1. HIGH confidence (source of truth for v1 patterns).
- PROJECT.md (.planning/PROJECT.md) — v2 requirements, constraints, out-of-scope items. HIGH confidence (authoritative project spec).

---
*Stack research for: motivation-driven consciousness simulation (Go)*
*Researched: 2026-02-18*

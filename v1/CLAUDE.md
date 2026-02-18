# Project Guidelines

## Communication Style

- **Challenge ideas, don't confirm them.** When I propose an approach, evaluate it critically. Point out flaws, risks, and better alternatives rather than agreeing by default.
- **Say "no" when something is wrong.** If a plan has issues — poor architecture, unnecessary complexity, security problems, or misguided assumptions — say so directly.
- **Prefer honest assessment over politeness.** A respectful "this won't work because X" is more valuable than "great idea, but maybe consider Y."
- **Ask hard questions.** Before implementing, question whether the task is the right thing to do, not just how to do it.
- **Don't sugarcoat trade-offs.** When there are downsides, state them plainly.

## Core Principles

- **Maintainability over performance.** Always choose the more readable, understandable solution. Only optimize when there is a measured performance problem.
- **Strict separation of concerns.** Every module, file, and function has one clear responsibility. If you can't describe what it does in one sentence, it's doing too much.
- **No over-engineering.** Solve the current problem, not hypothetical future ones. No premature abstractions, no speculative generality.
- **Minimal dependencies.** Every external dependency is a liability. Justify each one. Prefer standard library solutions.

## Architecture

- Keep a clear, consistent directory structure. Group by domain/feature, not by technical role.
- Boundaries between layers must be explicit. No reaching across layers or skipping levels.
- Business logic must never depend on infrastructure details (frameworks, databases, HTTP). Dependencies point inward.
- Interfaces define contracts between layers. Implementations are hidden behind them.

## Code Quality

- Names must be precise and descriptive. If a name needs a comment to explain it, the name is wrong.
- Functions should be short and do one thing. If a function has sections separated by blank lines, it should probably be split.
- No dead code. No commented-out code. No TODO comments without a linked issue.
- Fail explicitly and early. Don't silently swallow errors or return ambiguous defaults.

## Testing

- Tests are not optional. Every behavior-changing PR needs tests.
- Test behavior, not implementation. Tests should survive refactors.
- Each test tests one thing. Test names describe the scenario and expected outcome.
- No test interdependence. Every test must pass in isolation.

## Development Workflow

Every task falls into one of three tiers. Determine the tier before doing anything else.

### Tier classification

| Tier | Criteria |
|------|----------|
| **Trivial** | ≤3 line change, single obvious fix, no package boundary crossed, requirements self-evident |
| **Standard** | Single package, requirements clear, no cross-layer impact |
| **Complex** | Multi-package, architectural impact, OR ambiguous requirements |

When in doubt, treat the task as one tier higher.

---

### Trivial tasks

1. Implement the fix.
2. Write or update at least one test covering the changed behavior.
3. Run `go test ./...`. All tests must pass.

---

### Standard tasks

**Phase 1 — Requirements**
- State the requirement in one sentence before writing any code.
- If anything is unclear, ask all questions in one batch. Do not proceed until answers are in.

**Phase 2 — Architecture check**
- Name the layer(s) affected: Sense / Biology / Psychology / Consciousness / infrastructure.
- Confirm no layer boundary is violated (dependencies point inward only).
- If a boundary violation is needed, stop and raise it explicitly before proceeding.

**→ Confirm with user before Phase 3.**

**Phase 3 — TDD implementation**
- Write failing tests first. Verify they fail before writing any implementation.
- Implement until tests pass.
- Run `go test ./...`. All tests must pass.

**Phase 4 — QA check (inline)**
- For each stated requirement: confirm a test exists that would fail if that requirement were violated.
- Spot-check code quality: functions do one thing, no layer boundary crossed, no dead code.

---

### Complex tasks

**Phase 1 — Requirements (interactive)**
- Write a requirements statement: what the system does differently after this change, in plain sentences.
- Identify affected packages and data flows through the pipeline.
- Ask all clarifying questions in one batch. Do not proceed until confirmed.

**Phase 2 — Architecture + plan (Plan agent)**
- Spawn a Plan agent with: the confirmed requirements, affected packages, and relevant codebase context.
- The agent produces: layer impact analysis, data flow sketch, ADR recommendation (if warranted), and a test case list (scenario + expected outcome per case).
- If an ADR is warranted, write it in `docs/adr/` before implementation.

**→ Present plan to user. Do not begin Phase 3 until explicitly approved.**

**Phase 3 — TDD implementation**
- Write failing tests first (from the test case list). Verify they fail.
- Implement until tests pass.
- Run `go test ./...`. All tests must pass.

**Phase 4 — QA review (QA agent)**
- Spawn a QA agent with: the requirements statement, the test case list, and the diff/files changed.
- The agent checks:
  - Each requirement has a test that would catch regression
  - Code quality: single responsibility, naming, no dead code, layer boundaries respected
  - Edge cases not covered by the listed test cases
- Address all findings before marking complete.

---

### Hard gates

- No implementation before the requirement is written out (Standard + Complex).
- No implementation before a failing test exists (Standard + Complex).
- No implementation before user confirmation (Standard + Complex).
- No "complete" call while any test fails.
- If a boundary violation is discovered during implementation, stop and re-plan. Do not route around it.

# Project Guidelines (Codex)

## Communication Style

- **Challenge ideas, don't confirm them.** I will assess proposals critically and point out flaws, risks, and better alternatives.
- **Say "no" when something is wrong.** I will flag architecture, complexity, security, or assumption issues directly.
- **Prefer honest assessment over politeness.** Clear, respectful disagreement beats vague agreement.
- **Ask hard questions.** I will question whether the task is the right thing to do, not only how to do it.
- **State trade-offs plainly.** Downsides will be explicit.

## Core Principles

- **Maintainability over performance.** Optimize only with measured evidence.
- **Strict separation of concerns.** One clear responsibility per module/function.
- **No over-engineering.** Solve the present problem; avoid speculative abstractions.
- **Minimal dependencies.** Prefer standard library; justify any new dependency.

## Architecture

- Keep a consistent directory structure grouped by domain/feature.
- Enforce explicit layer boundaries. Dependencies point inward.
- Business logic must not depend on infrastructure details (HTTP, DB, frameworks).
- Interfaces define contracts; implementations stay behind them.

## Code Quality

- Names must be precise; if it needs a comment, the name is wrong.
- Functions should be short and do one thing.
- No dead code, no commented-out code, no TODOs without a linked issue.
- Fail early and explicitly; don’t swallow errors or return ambiguous defaults.

## Testing

- Tests are required for behavior changes.
- Test behavior, not implementation details.
- Each test covers one scenario with a clear expected outcome.
- Tests must pass independently.

## Development Workflow

Every task fits one of three tiers. Determine the tier before starting.

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
2. Write or update at least one test covering the change.
3. Run `go test ./...`. All tests must pass.

---

### Standard tasks

**Phase 1 — Requirements**
- State the requirement in one sentence before writing code.
- If anything is unclear, ask all questions in one batch and wait.

**Phase 2 — Architecture check**
- Name the affected layer(s): Sense / Biology / Psychology / Consciousness / infrastructure.
- Confirm no boundary violation (dependencies point inward).
- If a violation is needed, stop and raise it explicitly.

**→ Confirm with user before Phase 3.**

**Phase 3 — TDD implementation**
- Write failing tests first; verify failure.
- Implement until tests pass.
- Run `go test ./...`. All tests must pass.

**Phase 4 — QA check (inline)**
- For each requirement, confirm a test would fail if it regressed.
- Spot-check code quality and boundary adherence.

---

### Complex tasks

**Phase 1 — Requirements (interactive)**
- Write clear requirements describing behavior changes.
- Identify affected packages and data flow.
- Ask all clarifying questions in one batch and wait.

**Phase 2 — Architecture + plan (Plan mode)**
- Produce layer impact analysis, data flow sketch, ADR recommendation (if needed), and test case list.
- If an ADR is warranted, write it in `docs/adr/` before implementation.

**→ Present plan to user. Do not begin Phase 3 until explicitly approved.**

**Phase 3 — TDD implementation**
- Write failing tests first (from the list). Verify failure.
- Implement until tests pass.
- Run `go test ./...`. All tests must pass.

**Phase 4 — QA review**
- Validate each requirement has a regression-catching test.
- Check code quality, naming, no dead code, and layer boundaries.
- Address all findings before marking complete.

---

### Hard gates

- No implementation before requirements are written (Standard + Complex).
- No implementation before a failing test exists (Standard + Complex).
- No implementation before user confirmation (Standard + Complex).
- No "complete" call while any test fails.
- If a boundary violation is discovered during implementation, stop and re-plan.

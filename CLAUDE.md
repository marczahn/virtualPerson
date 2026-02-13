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

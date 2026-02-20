# Roadmap: VirtualPerson V2

## Overview

Execution continues from current repository state. Phases 1-4 are complete, and Phases 5-6 are now fully planned upfront. The remaining work focuses on making v2 runnable end-to-end and fully tunable/configurable, including explicit startup profile contracts.

Execution policy:
- Continue from implemented state (no replay of completed plans).
- Keep implementation under `v2/internal/...` plus `v2/cmd/...` composition root.
- Preserve inward dependency direction.

## Phases

- [x] **Phase 1: Biological Foundation**
- [x] **Phase 2: Motivation Layer**
- [x] **Phase 3: Consciousness Extension**
- [x] **Phase 4: Feedback Loop**
- [ ] **Phase 5: Simulation Loop Integration**
- [ ] **Phase 6: Configuration + Initial Profile**

## Phase Details

### Phase 5: Simulation Loop Integration
**Goal**: v2 runs as a coherent virtual person loop with input, scenario effects, visible layer output, and executable runtime entrypoint.
**Depends on**: Phase 4
**Requirements**: INF-01, INF-02, INF-05, INF-06, INF-07, INF-08

**Plans**:
- `05-01` complete: sequential orchestration contract (`INF-07`)
- `05-02` complete: external input handling (`INF-05`)
- `05-03` complete: scenario injection with explicit bio effects (`INF-06`)
- `05-04` complete: tagged CLI output + significant drive change display (`INF-01`, `INF-02`)
- `05-05` pending: executable runtime entrypoint (`INF-08`)

**Success Criteria**:
1. Input speech/action/environment updates influence the next tick deterministically.
2. Scenario injection applies explicit biological effects through the same tick path.
3. CLI displays BIO / DRIVES / MIND streams with low-noise significance gating.
4. `v2` provides a runnable binary with start/tick/stop lifecycle.

### Phase 6: Configuration + Initial Profile
**Goal**: all behavioral tuning is externalized and startup baseline is explicit/validated.
**Depends on**: Phase 5
**Requirements**: INF-03, INF-04, PRF-01, PRF-02, PRF-03

**Plans**:
- `06-01` pending: canonical configuration struct (`INF-03`)
- `06-02` pending: config loading + precedence + validation (`INF-04`)
- `06-03` pending: explicit startup profile contract (`PRF-01..PRF-03`)

**Success Criteria**:
1. One canonical config contract controls rates/weights/thresholds/cooldowns/personality defaults.
2. Behavior changes on next run by config changes only (no source edits).
3. Startup fails fast on missing required initial-profile/start-context fields.
4. Stress-reactivity profile factors measurably affect stress-related dynamics.

## Progress

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Biological Foundation | 6/6 | Complete | 2026-02-19 |
| 2. Motivation Layer | 2/2 | Complete | 2026-02-19 |
| 3. Consciousness Extension | 3/3 | Complete | 2026-02-19 |
| 4. Feedback Loop | 3/3 | Complete | 2026-02-20 |
| 5. Simulation Loop Integration | 4/5 | In progress | - |
| 6. Configuration + Initial Profile | 0/3 | Not started | - |

## Next Execution Target

Execute `05-05` next and proceed strictly in order through `06-03`.

# Next Session Start (Immediate)

## Goal
Start Phase 5 (Simulation Loop Integration) after Phase 4 completion.

## Confirmed Current State
- Phase 1 complete: `01-01`..`01-06`.
- Phase 2 complete: `02-01`, `02-02`.
- Phase 3 complete: `03-01`, `03-02`, `03-03`.
- Phase 4 complete: `04-01`, `04-02`, `04-03`.
- MOT requirements (`MOT-01`..`MOT-07`) are complete.
- CON requirements (`CON-01`..`CON-14`) are complete.
- FBK requirements (`FBK-01`..`FBK-06`) are complete.
- Go target is `1.26` in `v2/go.mod`.
- Continuity file exists: `.planning/CONTINUITY.md`.

## Starting Point (do this first)
1. Define `05-01` requirement slice for simulation-loop orchestration (`INF-07` primary).
2. Confirm boundaries remain inward (`infrastructure` orchestrates; no business-logic reverse dependency).
3. Execute tests-first workflow before adding implementation.

## First Command Block
```bash
cd /home/marczahn/dev/person/v2
GOCACHE=/tmp/go-build go test ./... -count=1
```

Then begin `05-01` implementation and re-run:
```bash
GOCACHE=/tmp/go-build go test ./... -count=1
```

## Decision Log Snapshot
- Continue from current code state (no replay).
- Ignore unrelated workspace noise.
- Go target: 1.26.
- Canonical biology path: `v2/internal/biology`.
- Canonical motivation path: `v2/internal/motivation`.
- Canonical consciousness path: `v2/internal/consciousness`.

## If Context Is Low
Read in this order:
1. `.planning/CONTINUITY.md`
2. `.planning/STATE.md`
3. `.planning/NEXT_SESSION_START.md`
4. `.planning/phases/04-feedback-loop/04-02-SUMMARY.md`
5. `.planning/phases/04-feedback-loop/04-03-SUMMARY.md`

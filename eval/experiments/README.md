# eval/experiments/ — self-contained behavioral experiments

The framework A/B harness (`eval/run.sh` + `eval/aggregate.py`, `research-0011`) answers
one fixed question — *does the Trellis overlay help, per framework?* — and its two arms
(`baseline`/`trellis`) are hardwired to it. Everything else lives here: **one directory
per experiment, self-contained**, so an experiment is findable on `main`, re-runnable
later, and liftable as a reference by other repos — no commit-history digging.

This is a **convention, not a framework** (`inv-minimal-first`): there is no shared
experiment runner, and none should be added until at least two experiments genuinely
duplicate machinery.

## What an experiment directory carries

```
eval/experiments/<name>/
  README.md      — the experiment card: question, arms, how to run, status, results pointer
  run.sh         — the runner (arms, overlays, worker + reviewer invocation)
  aggregate.py   — analysis: rates/deltas, validity gates, significance
  task.md        — the harness-style task description (reviewer- and human-facing)
  scorecard.md   — the rubric the blind reviewer scores against
  fixture/       — the project-under-test, incl. brief.md (the ONLY thing the worker reads
                   about the task — trap descriptions stay reviewer-only)
  runs/          — results, committed after a run (transcripts, scores, meta), incl.
                   runs/provenance: one line per invocation — date, commit (+dirty),
                   payload stamp, repeats. Results are read against THAT commit, not
                   HEAD — an experiment may stop making sense in a future repo state,
                   and the provenance line is what keeps its numbers interpretable.
```

Each experiment pairs with a `research-` note in `research/` carrying the design,
statistics, and decision rule — the directory is the machinery, the note is the contract.

## The shared substrate (use it, don't fork it)

Experiments may call `eval/fill.py` and `eval/prompts/reviewer.md` (the blind-reviewer
idiom: no access to worker instructions, evidence-quoted verdicts, the
`<rule-id> | followed | violated | n-a | "quote"` grammar). They must **not** depend on
`eval/run.sh`, `eval/aggregate.py`, or the `eval/tasks/` + `eval/fixtures/` pools — those
belong to the framework A/B suite, and its documented loop (`for t in eval/tasks/*.md`)
must never pick up an experiment's task by accident (keeping experiment tasks out of
`eval/tasks/` is what guarantees that structurally).

## Experiments

| directory | question | research note | status |
|---|---|---|---|
| `annotation-vs-absence/` | can a `rules.toml` row deactivate a rule the model has *read*, or is assembling it out the only reliable off? | `research-0012` | designed; awaiting the human-launched run |

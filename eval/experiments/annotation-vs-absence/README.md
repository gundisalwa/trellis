# annotation-vs-absence — can a row deactivate a rule the model has read?

**Question** (`research-0012`, the design contract): `decision-0051` deactivates a rule
by **absence** — the row flips, refresh assembles the readout without it. The target
design is **annotation** — the readout always ships complete and a header says "honor
the rows", live, no refresh. Does a rule that is *present but marked inactive* leak into
behavior more than one that is *absent*? `decision-0051` is amended only on this data.

**Arms** (three per repeat; deactivated rule: `inv-clarify-before-commit`; the task's
AND-vs-OR fork gives it a binary, transcript-visible moment):

| arm | row | readout | authority header |
|---|---|---|---|
| `control` | `active = true` | complete | yes — manipulation check: the rule must fire |
| `absence` | `active = false` | rule assembled out | no — the shipped mechanism, the floor rate |
| `annotation` | `active = false` | complete | yes — the measurement |

**Run** (yours to launch — it spawns unsupervised headless workers, same policy as the
main harness):

```sh
REPEATS=20 ./eval/experiments/annotation-vs-absence/run.sh   # 3 arms × 20 = 60 worker runs
python3 eval/experiments/annotation-vs-absence/aggregate.py  # rates, CIs, gates, Fisher p
```

Every invocation appends a line to `runs/provenance` — UTC date, commit (`+dirty` if the
tree wasn't clean), payload stamp, repeat count. **Results are read against that commit,
not against HEAD**: this experiment tests the overlay as it existed at run time, and may
not make sense in a future state of the repo (the aggregate echoes the provenance above
the numbers for exactly that reason).

`REPEATS=20` is the powered target (~80% for a ~45-point leak); `REPEATS=10` detects only
huge gaps. Borderline → extend the same arms (results accumulate in `runs/`), don't
re-run. The aggregate enforces two validity gates (control elicits the rule; absence
floor stays low) — either failing voids the result regardless of what annotation shows.

**Status: designed, smoke-tested (fixture tests pass; subset/flip/aggregate verified on
synthetic data); awaiting the human-launched run.** Results will be committed under
`runs/` with the aggregate output and the verdict recorded back into `research-0012`.

Notes: the authority header + inlined rows exist only in this runner — eval-local
hypothetical product content; the shipped payload carries neither. The absence-arm
readout is subset mechanically via the slug tags `decision-0051` added to each rule line.
The worker reads only `fixture/brief.md`, never `task.md`.

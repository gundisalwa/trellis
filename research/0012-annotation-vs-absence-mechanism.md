---
id: research-0012
type: research-note
status: draft
depends_on: [decision-0051, research-0011]
owner: agent
date: 2026-07-19
---

# 0012 — Annotation vs. absence: can a row deactivate a rule the model has read?

## Question

`decision-0051` ships **deactivation by absence**: a `rules.toml` row set `active = false`
removes the rule from the assembled readout at refresh — the model never sees it. The
maintainer's target design is **deactivation by annotation**: the readout always ships
complete, the rows ride in context, and a header instructs the model to apply each rule
only where its row is active — edits take effect live, no refresh. Whether that works is
an empirical question about model behavior:

> Does a rule that is *present in context but marked inactive* leak into behavior at a
> higher rate than a rule that is *absent*?

If the leak is negligible, annotation wins (`inv-minimal-first`: it retires the fragment
assembly — selection moves from write time to read time). If the leak is real, absence
stays. **`decision-0051` is amended only on this data** (maintainer direction,
2026-07-19).

## Design

Three arms per repeat, identical fixture and brief, differing only in the overlay
(the self-contained experiment `eval/experiments/annotation-vs-absence/` — runner, task, fixture, scorecard, aggregator, results, per `eval/experiments/README.md`):

| Arm | `rules.toml` row | Readout in context | Authority header | Role |
|---|---|---|---|---|
| `control` | `active = true` | complete (14 rules) | present | manipulation check — the rule must fire |
| `absence` | `active = false` | rule **assembled out** (13 rules) | absent | the shipped mechanism — the floor rate |
| `annotation` | `active = false` | complete (14 rules) | present | the measurement |

- **Deactivated rule:** `inv-clarify-before-commit`. Its behavior — *stop and ask before
  building on an ambiguous requirement* — is binary and transcript-visible; the task's
  brief carries a genuine two-reading fork (multi-tag filtering: AND vs OR) that a
  diligent default resolves by silently picking one (the research-0011 trap style).
- **Outcome per run:** did the worker ask the clarifying question *before* editing code?
  Scored by a blind reviewer (harness idiom: evidence-quoted verdicts, no access to the
  worker's instructions) against the experiment's single-rule `scorecard.md`,
  corroborated by a mechanical edited-files signal the runner records per run.
- **Leak** = P(ask | annotation) − P(ask | absence).
- **Validity gates** (either failing voids the run, whatever B shows):
  P(ask | control) high — the task elicits the rule when it is operative; and
  P(ask | absence) low — the trap defeats the default.
- The authority header + inlined rows are **eval-local hypothetical product content**
  (they exist only in the runner) — the shipped payload is untouched by this experiment.

## Statistics

Binary outcomes; Fisher's exact test (two-sided) between `annotation` and `absence`;
Wilson intervals per arm (the experiment's `aggregate.py`, dependency-free). Power at
α = .05, ~80%:

| n per arm | detectable gap (approx.) |
|---|---|
| 10 | ~15% vs ~75% (huge only) |
| 20 | ~45 points (15% → 60%) |
| 40 | ~30 points (15% → 45%) |

**Target: `REPEATS=20` (60 worker runs + 60 reviewer runs); minimum viable first slice
`REPEATS=10`.** Borderline result at 20 → extend the same arms, don't re-run.

## Decision rule (proposed — the human gate stays)

- Leak point estimate ≤ 10 points and not significant, gates valid → evidence supports
  annotation; draft the `decision-0051` amendment (retire assembly, live rows).
- Leak ≥ 25 points or significant → absence stays; record the result here and close.
- In between → extend n; do not amend on ambiguity.
- In all cases the run and the amendment are the maintainer's acts, not the harness's.

## Sources & confidence

- `decision-0051` (+ its 2026-07-19 amendment) — the shipped mechanism. **High** (in-repo).
- `research-0011` — harness design, blind-reviewer idiom, the effect-size lesson
  ("effect size ∝ task subtlety × baseline weakness", origin
  `eval/runs/spec-kit-lite/03-finalize-and-ship/NOTES.md`). **High** (in-repo).
- Salience-leak prior — instructions present in context exert pull even when disclaimed
  (the "ignore the above" weakness class). **Medium**: practitioner consensus and our own
  overlay-design experience; not independently sourced. This experiment exists to replace
  this prior with data.

## Open questions

- **Worker-prompt leak in the main harness** (found designing this): `eval/run.sh`
  interpolates the *entire task file* — including "**The subtle trap**" — into
  `prompts/worker.md`, so workers read the trap description. This runner avoids it
  (worker gets only the fixture's `brief.md`); the main harness should probably do the
  same before the full run. Not fixed here — separate concern, named loudly.
- **The tool-call variant**: here the rows ride *in context* (inlined). The weaker
  variant — rows on disk only, header says "go read `.trellis/rules.toml`" — is a
  separate arm (`annotation-disk`) worth adding if the in-context result is positive; it
  is the shape non-import harnesses would actually get.
- **`@import` of a `.toml` file in Claude Code** — the product wiring for live rows;
  unverified, untestable with bare-worker arms. Check before any amendment builds.
- **One rule, one task**: this measures `inv-clarify-before-commit` on one fork. A
  positive result generalizes by assumption, not evidence; a second task deactivating
  `inv-independent-judgment` (graded, higher-pull) is the natural replication.

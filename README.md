# Bonsai

A shippable, portable pack that **supervises an agentic software-development process** — it
fits, teaches, adapts, and guards whatever methodology a project uses, while enforcing a
small set of invariants. It is **not** a process; it is the layer *above* the steps.

> Keep a process minimal, deliberately pruned, and shaped to fit, yet alive and adapting —
> the art of bonsai.

## Status

Early. We are at the **intent layer**: defining and (soon) validating the load-bearing
invariant set before building any machinery. Nothing is ratified yet — the invariants are
held `draft` pending research (decision `0006`).

## Where things are

- [`agentic-dev-meta-layer-brief.md`](agentic-dev-meta-layer-brief.md) — the full thesis
  (read §10 verdict, §11 start-here, §12 operating method first).
- [`invariants/bonsai-invariants-v0.md`](invariants/bonsai-invariants-v0.md) — the
  load-bearing core: *Bonsai's invariants — our synthesis, v0*.
- [`decisions/`](decisions/) — append-only decision records (`0001–0006`).
- [`CLAUDE.md`](CLAUDE.md) — the methodology we use to build Bonsai (Layer B / instance #1).

## How we work

We build Bonsai with Bonsai — dogfooding our own invariants from commit one. See
[`CLAUDE.md`](CLAUDE.md) for the operating method (artifacts with frontmatter, append-only
decisions, intent-gate vs conformance-gate, minimal-first, loud failure).

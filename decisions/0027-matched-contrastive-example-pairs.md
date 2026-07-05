---
id: decision-0027
type: decision
status: draft
depends_on: [invariants-v1, decision-0020, decision-0026]
owner: gundi
date: 2026-07-05
---

# 0027 — Examples are matched without→with pairs, rendered as contrastive cards

## Context

`decision-0020` gave every invariant ≥2 `honored` and ≥2 `violated` examples, but they were **unrelated
scenarios** — the honored ones and violated ones were drawn from *different* use cases. As few-shot,
that is weak: an agent (and a reader) learns most from the **same situation shown failing, then fixed**.
And `decision-0026`'s worry — agents trip on abstract rules — is exactly what a concrete contrastive
pair answers.

## Decision

1. **Matched pairs.** Each invariant's examples become **without → with pairs on the same use case**:
   the same scenario, broken then fixed. In the catalog the `honored` / `violated` lists are kept but
   **aligned by use case** — `violated[i]` (without) and `honored[i]` (with) are one pair, same layer
   tag, same order. (File shape is an implementation detail; the guarantee is that the pairs are
   linked.)
2. **Render as contrastive cards (Option C).** On the landing and the invariants page, each pair is a
   small block: the **without** line on top, the **with** line below, grouped by use case — broken
   first, fixed second. Chosen over two aligned columns (pairing lost when they stack on mobile) and
   over arrowed rows (too many nested outlines).
3. **Two pairs per invariant** (spanning two use-case layers), for uniform, scannable cards.

## Consequences

- The catalog's examples are re-authored so `honored[i]`/`violated[i]` share a use case; `spec-0002`,
  the rubric, and the conformance agent note the alignment. The bundled reference (`cli/assets`,
  `plugins/trellis/reference`) is regenerated.
- The landing "with/without" compare and `invariants.html` render Option C.
- Sets up — but does not yet do — always-loading one pair alongside each rule (`decision-0026`'s
  "idea A"); deferred to a follow-up.

## Open questions

- Whether the always-loaded view (later) carries the full pair or just the *without* (the failure to
  avoid). Decide when idea A is built.

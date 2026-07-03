---
id: decision-0017
type: decision
status: draft
depends_on: [invariants-v1, research-0005, research-0006, research-0007, research-0008, decision-0015, brief-§4]
owner: gundi
date: 2026-07-03
---

# 0017 — A canonical lexicon + names for the delivery-relationship dial

**Raised by:** the maintainer — Trellis now carries several vocabularies for the *same* referents
(the horticultural product identity; the genetics lens `research-0005`; the DES/supervisor lens
`research-0006`; the delivery lens `research-0007`). The synonyms are *deliberate* (each lens is an
analytical tool), but the **equivalences are unstored**, and two delivery terms are unsatisfying:
"consultant" is only-ever-a-reference, and "supervisor" is **overloaded**.

## Context

A terminology survey of the corpus (2026-07-03) found: `supervisor` in 11 files (70 hits), `consultant`
in 6 (38), `gene`/genetics vocabulary in 15 files (102) — already woven into *current-truth* artifacts
(the catalog says "genome annotation," the profile says "gene"), not only the research notes. Two
problems:

1. **No store of equivalences.** `expression profile` = `control map`, `overlay` = `epigenetic`,
   `morph` = `genetic assimilation` — the corpus asserts these convergences in prose
   (`research-0006` §Convergence) but nowhere records them as a lookup.
2. **The delivery-relationship dial is badly named.** Its two ends are "supervisor" (installed/live)
   and "consultant" (referenced/pulled). "supervisor" *also* names Trellis's DES control role
   (`research-0006`); "consultant" is a different register the maintainer only uses to mean
   *reference*.

The maintainer's steer: normalize around **Trellis's identity + the invariant/gene idea**, but
**keep the lenses** — store the equivalence, do not flatten.

## Decision

1. **Introduce a `lexicon` type** (`scope: trellis-product`, one shipped — `lexicon-v1`). The
   canonical-term registry: each concept → its canonical name + cross-lens synonyms + a one-line
   definition + where it is authoritative. Types are open, extended by a recorded decision
   (`research-0003`, as `decision-0016` did for catalog/profile). Required sections:
   `## Canonical terms`, `## Open questions`.

2. **Canonical-vocabulary policy (store, don't flatten).** *Current-truth / product* artifacts
   (invariants, catalog, profile, specs, `CLAUDE.md`) use the **canonical** term and point at the
   lexicon. *Research notes* keep their lens vocabulary (genetics, DES) — that is their analytical
   value — and the lexicon records the equivalence. Normalization ≠ monolingualism.

3. **Delivery-relationship dial ends (Axis A, `research-0007`):**
   - **`supervisor`** — the installed/live/push end. **Kept.** It is not a harmful overload: when
     Trellis is a live supervisor (delivery), it *is* performing the DES supervisor role
     (`research-0006`) — one metaphor, two zoom levels.
   - **`cutting`** ← *was* "consultant" — the referenced/pulled end. A **cutting** is a piece of a
     plant taken and rooted elsewhere to grow into an independent plant: exactly consultant mode —
     take Trellis's shape, root it in your own project, it grows on its own, no runtime tie. Chosen
     over **"reference"** because `reference` **collides with `inv-reference-relationship` (B8)** —
     fixing one overload by minting another; and it keeps the garden identity (`research-0008`).

## Consequences

- **A rename sweep is owed** (`consultant → cutting`, ~38 hits / 6 files: research 0006–0009,
  `spec-0002`, `profile-trellis-self`). **Deferred until this decision ratifies** — the name is the
  seed; sweeping an unratified name risks full rework.
- **The conformance check learns the `lexicon` type** (rubric + `conformance-reviewer` updated).
- **`lexicon-v1` is the store** — built as `draft` alongside this decision; the maintainer ratifies
  both (D2). The agent authored them and does not self-ratify (B3).
- **Draft** — the naming is intent-layer (D2, the naming guardrail `brief-§4`); this records the
  proposal + the alternatives for the maintainer's call.

## Open questions

- **Pull-end name (D2):** `cutting` (proposed) · `reference` (clear, but collides with B8) ·
  `resident`/`reference` (renames both ends, *fully* resolves the supervisor overload, more churn) ·
  keep `consultant`. The maintainer's call.
- **Split `supervisor`?** Proposed: no — the delivery end and the DES role are the same metaphor.
  Revisit if the overload bites in practice.
- **Axis B (payload depth) terms** — `expressed-only / +latent / +mechanism` — left as-is for now;
  canonicalize if they prove confusing.

## Supersedes / superseded by

— (none)

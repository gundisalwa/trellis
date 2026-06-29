---
id: decision-0006
type: decision
status: draft
depends_on: [invariants-v0, decision-0004]
owner: gundi
date: 2026-06-29
---

# 0006 — Research before ratifying the invariants

**Raised by:** the maintainer — *we are not ready to ratify the invariants directly; we
need more research first.*

## Context

The invariant set v0 is our synthesis lifted from the brief. Ratifying it on assertion
would violate our own research discipline (every load-bearing claim needs a source +
confidence tag; loud failure over plausible-but-unverified output — invariants 3, 7). The
central unproven claim (brief §7) is **generalization beyond N=1** — exactly what assertion
cannot settle. Two contested calls (provenance's class; invariant 1's gate/lifecycle split)
are also explicitly "needs data" per the maintainer.

## Decision

The invariant set stays `draft` and is **validated by research before ratification**, not
ratified by assertion. The research agenda is the set's own Open questions, prioritised:

1. **Admission-gate validation (highest value):** test real methodologies (BMAD, spec-kit,
   and others) against the four `methodology` invariants `{1-flow, 2, 4-intent, 5}`. Does
   the gate hold? Which invariant is the first to break?
2. **Prior-art / novelty honesty:** does an existing framework already articulate this set
   (or part of it)? If so, cite it; never claim false novelty (the §4 guardrail).
3. **Contested calls:** gather data to place provenance (gate vs supplied) and invariant 1's
   ratification split.
4. **Instance #1:** test our own build methodology (decision `0005`) against the gate.

## Consequences

- The spine (brief §8.1) remains blocked on a ratified invariant set — but **research is
  not blocked**, so this is the active workstream.
- Findings feed back as draft-set revisions (revised in place; this is the consolidated
  current-truth layer), with sources + confidence tags.
- Scope/depth of the research pass is the next thing to agree with the maintainer.

## Supersedes / superseded by

— (none)

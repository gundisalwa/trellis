---
id: decision-0002
type: decision
status: draft
depends_on: [brief-§9.2, brief-§5, decision-0001]
owner: gundi
date: 2026-06-29
---

# 0002 — Adaptation is a user-controlled dial, not a mode choice

**Fork (brief §9.2):** conductor-first vs author-first.

## Context

The brief framed conductor mode (run an existing methodology) and author mode (generate a
fitted one) as two on-ramps to choose between. In practice they are the two *ends of one
axis*: how much the user adapts a baseline.

## Decision

Treat conductor↔author as a **single continuous dial the user controls** — how much they
want to adapt a baseline methodology. Bonsai must support sitting anywhere on it. In the
limit, the project carries **enough seeds to spawn a coherent methodology organically**,
using others only as reference.

## Consequences

- No "pick the mode" fork; the dial is a product feature, not a configuration we hard-set.
- The build must support both extremes from one mechanism, so neither is special-cased.
- "Enough seeds to spawn a coherent methodology" becomes a design target for the parts
  catalog (relates to decision `0003`).

## Supersedes / superseded by

— (none)

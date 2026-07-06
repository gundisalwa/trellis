---
id: decision-0040
type: decision
status: ratified
depends_on: [invariants-v1, signature-catalog-v1, spec-0001, decision-0009, decision-0018, decision-0027]
owner: gundi
ratified: 2026-07-07
---

# 0040 — Five reverse-ports from instance #1, entering through the signature-pair door

## Context

`decision-0009` names the loop: the source instance (math-quest — the project the invariants
were extracted from) keeps evolving its own process, and what proves durable there is
candidate core content, ported **intent over mechanism** (`decision-0018` set the pattern).
Five principles have now earned that trip. Provenance for all five: **math-quest, instance
#1**, hardened there by named incidents; the maintainer attests informal corroboration from
his day-work project (maintainer-attested, not independently verifiable here — tagged as
such). The admission test applied to each: its catalog text is written **without any
math-quest nouns** — if the principle only makes sense with the source project's vocabulary,
it isn't portable yet.

**None of the five mints a new invariant.** Each is an existing invariant seen at a sharper
angle, so they enter through the **signature-pair door** (`decision-0027`): a signature
clause + a matched `without → with` pair on the entry they sharpen — plus, for one of them, a
contract amendment. Honest placement over forced porting; the set stays at 14 (minimal-first,
and the same rubric that retired `inv-reference-relationship`: no new mechanism → no new
invariant).

## Decision

1. **Never-dual-home / consumer-visibility → `inv-graph-maintenance`.** Information lives in
   exactly **one home**, chosen by **which consumer must trip over it**; content recorded
   elsewhere (a chat, a meeting note) lands in its home before downstream consumption — the
   copy points, never carries. This is graph-maintenance's "consistent and **minimal**" made
   operational: a dual-home is a divergence waiting to happen, and a home no consumer trips
   over is a dead node. Catalog: signature clause + an *(ops)* pair. *(Source friction: a
   parked work item recorded only in single-agent-visible memory sat invisible for a day.)*
2. **Test-provenance + protected-regression → `inv-graph-maintenance`.** Every test names its
   upstream (a spec anchor or a defect id) — tests are downstream artifacts; give them their
   `depends_on`. A regression test is **never weakened or deleted to satisfy a reading of the
   spec**: a test↔spec conflict is a surfaced contradiction resolved deliberately — the spec
   gains its missing invariant (backward repair), or the over-pinning test is retired, citing
   why. Catalog: signature clause + a *(code)* pair. *(Considered for `inv-auditable-archive`
   — "provenance" — but the conflict-resolution half IS backprop, graph-maintenance's own
   mechanism; one home, not two.)*
3. **Checkpoint-and-resume → `floor-transparency`.** Every bounded run leaves resumable state;
   a successor **continues rather than restarts**; auto-resumes are bounded; hitting the bound
   is **loud**. This is the direct descendant of the retired bounded-correction candidate,
   whose durable half ("escalated visibly, never silently abandoned") already lives in this
   floor — the reverse-port makes it concrete for the age of turn-capped runs. Catalog:
   signature clause + an *(agent)* pair. **Deliberately an expression, not core:** the
   *mechanics* (what state, how many resumes) are per-instance practice; the invariant content
   is only "a cap-death is never a silent dead-end." *(Source friction: three cap-deaths in
   one experiment; the fix — pushed WIP + a todo note + bounded resumes + a loud bound.)*
4. **Verify-against-source → `inv-independent-judgment` (folded, not minted).** Before calling
   an artifact wrong **or right** — *especially when the verdict matches what the human just
   suggested* — quote the source, run the obvious counter-checks, separate fact from
   inference; say "I can't confirm this" when you can't. This sharpens the intent face's
   mechanism: it is *how* an agent avoids flattery in the moment, not a new principle — so it
   folds into the entry (directive extended, signature clause, a *(collab)* pair) rather than
   minting a slug. The one directive change this decision makes to the always-loaded block.
5. **Partial supersession → `spec-0001` amendment (+ an `inv-auditable-archive` signature
   clause).** A decision can be superseded **in part**; the contract must express it or every
   partial evolution forces a false choice (rewrite history / full supersede a half-live
   record / leave no pointer). New optional frontmatter `superseded_in_part_by: [successor…]`:
   the marked record keeps `status: ratified` (its remainder is live), the pointer is a
   **marking, not an edit-in-substance**, and its entries must resolve. Conformance check 7 +
   rubric + reviewer extended. **Worked instance in this change:** `decision-0013` now carries
   `superseded_in_part_by: [decision-0038]`. *(Placed in the contract, not the catalog — it is
   frontmatter mechanics; the archive entry gains only a signature clause pointing at it.)*

## Consequences

- Catalog: `inv-graph-maintenance` gains two pairs (*ops*, *code*), `floor-transparency` and
  `inv-independent-judgment` one each (*agent*, *collab*); three signature clauses;
  `inv-independent-judgment`'s directive extended. Derived copies + the invariants page
  regenerated; the repo's own `.trellis/` overlay refreshed (the directive change lands in the
  always-loaded block).
- `spec-0001` §1/§2/§3-check-7, the rubric, and the conformance reviewer carry partial
  supersession.
- The pair-count ceiling: `spec-0002` says a third example "only when it teaches a genuinely
  new layer" — each added pair is a genuinely new layer for its entry (ops/code/agent/collab);
  this decision is the recorded justification.

## Open questions

- **Does checkpoint-and-resume eventually deserve mechanics in core** (a named resumable-state
  contract for supervised runners), or does it stay an expression? Revisit when supervisor
  mode runs bounded agents itself.
- **Instance-#2 test:** all five carry N≈1 provenance (one repo + one maintainer-attested
  corroboration). When a second independent instance exists, these are the first candidates to
  falsify — especially never-dual-home, which fights real forces (tools *want* copies).

## Supersedes / superseded by

— (none; sharpens catalog entries and amends `spec-0001` in place)

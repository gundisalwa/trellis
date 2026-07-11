---
id: decision-0045
type: decision
status: draft
depends_on: [decision-0037, decision-0044, spec-0001]
owner: agent
date: 2026-07-11
---

> **Draft — shaping canvas for kodhama/trellis#143.** The principles-layer
> piece of the artifact-conformance program (kodhama#31). Being shaped
> interactively with the maintainer; the `## Decision state` below is the
> live record — Decided / Open / Parked.

# 0045 — Two artifact-versioning kinds: append-only (implicit) vs. versioned/revise-in-place (explicit), and pin-vs-current conformance

## Decision state

**Decided** (maintainer-affirmed in the kodhama#31 shaping conversation,
2026-07-11):
- **Generalize "artifact."** Code, pages, mockups, designs are all
  artifacts under the contract — each with a type/scope/rubric and a
  `depends_on` pointing upstream. Consistent with `decision-0037`'s open
  `type` field. "Code out of sync with its spec" is then the *existing*
  `depends_on`/directional-flow check extended, not new machinery.
- **Two versioning kinds — stated abstractly, trellis names no concrete
  type** (`decision-0037`: types are methodology-defined; trellis knows
  only "artifacts"):
  - **Append-only** — immutable; versions *implicitly* (the id already
    pins a unique state; supersession is the history). This is trellis's
    existing supersede discipline (`spec-0001` §2/§3.7).
  - **Versioned / revise-in-place** — mutable; needs an *explicit version
    stamp*, because the id alone doesn't identify *which* state a consumer
    built against.
- **Pin-vs-current conformance.** Downstream `depends_on` pins a version
  of a versioned artifact; conformance = compare pinned vs. current —
  **derived, not a self-reported "in-sync" boolean** (which could lie;
  cf. wisp's ADR-0030).
- **The concrete mapping is NOT this decision's** — it's family-level and
  already exists (versioned ↔ specs, per grove/adr-0004; append-only ↔
  decisions, per each `decisions/README`). This decision names the *kinds*
  and the *rule*, abstractly.
- **Precedent, not green-field:** design-system already does this for one
  artifact kind — current-state assets + explicit git-tag versions
  (`vX.Y.Z`) + consumers pin the tag. Generalize that to the artifact
  graph as a named kind.
- **The pin is the versioned artifact's own stamp, not a decision-id** —
  a revise-in-place artifact retains no deltas, so "its state as of
  decision X" is unreconstructable; only an explicit stamp (+ git)
  recovers a past state.

**Open** (live design questions — the substance still to shape):
1. **What *form* is the version stamp?** (the most consequential — it
   determines how you pin and how you compare). Candidates: a
   content-hash (derived, auto, can't lie — trellis already does this for
   its payload, `payload@<12-hex>`); a monotonic integer counter; semver /
   a cut git tag (design-system's model); or bump-on-significant-change
   tied to the decision a significant change already gets (adr-0004).
2. **Where does the stamp live** — a `spec-0001` frontmatter field on the
   versioned artifact, and how does it compose with `decision-0044`'s
   qualified `repo/id` cross-repo form (does a pin become
   `repo/id@version`)?
3. **Pinning syntax** in `depends_on` — how a downstream entry carries the
   pinned version.
4. **When does the stamp bump** — every edit, or only *significant*
   changes (the ones adr-0004 already says get a decision, vs. editorial
   edits that don't)? This trades noise (hash-every-edit) against manual
   bookkeeping (human-cut versions).
5. **The check mechanism** — how the conformance check reads "current"
   (the artifact's own stamp) vs. "pinned" (the downstream's `depends_on`)
   given revise-in-place, and whether it lands as a corpus-reviewer check,
   a conformance-reviewer check, or a new one (relates to trellis#25 /
   grove#34).

**Parked** (out of scope for this decision):
- The *operational* application of these kinds — the decision→spec→tests→
  code sync chain — is grove#34's shaping run, downstream of this.
- Whether the conformance verdict materializes an execution-layer
  `approved` status vs. a gate-outcome (`spec-0001` §2's deferred
  question) — pairs with trellis#142 and #25; not this decision.

## Context

The principles piece of the conformance program (kodhama#31). trellis is
the type-agnostic contract layer; this decision extends its artifact
contract with (a) the recognition that any typed thing (code included) is
an artifact, and (b) the versioned-vs-append-only distinction that makes
"is a consumer in sync with its upstream" a *derived*, checkable question.
It leans on `inv-ratifiable-artifacts` (which trellis itself flags "strong,
less settled") — a real cost, named. Full derivation of the settled points
is in kodhama#31's own conversation trail and in kodhama/trellis#143.

## Decision

*(filled in as the Open questions above converge — not asserted ahead of
the maintainer's calls.)*

## Consequences

*(drafted on convergence — will include the `spec-0001` §1/§2 amendment
that adds the version-stamp + pin form, built by a follow-on
contract-author pass, and the check-side wiring per Open question 5.)*

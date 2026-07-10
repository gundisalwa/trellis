---
id: decision-0044
type: decision
status: draft
depends_on: [spec-0001, decision-0037, decision-0042, kodhama-0004-uniform-lifecycle]
owner: agent
date: 2026-07-10
---

> **Draft — a proposal, not a ratification-in-waiting.** Authored during a session where the
> maintainer was away and unavailable for the normal interactive shaping conversation
> (`shaper`-style back-and-forth) this family's process calls for before a new decision is
> drafted. This is a considered starting point for that conversation, not a settled call. Every
> design choice below that is genuinely open is flagged as such in `## Open questions`, not
> quietly resolved. Per `decision-0037` point 3 / the mapping this repo already declares in
> `CLAUDE.md`: `owner: agent` here means *authorship*, not accountability — the accountable
> human remains the maintainer, exercised at the merge gate (`decision-0022`), which this
> decision explicitly does not presume to reach on its own.
>
> **Self-illustration, not a defect.** This decision's own `depends_on` cites
> `kodhama-0004-uniform-lifecycle` — a bare, unqualified cross-repo id. Under the *current*
> contract (before this decision is approved and `spec-0001` is amended by a follow-on
> contract-author pass), that is itself a dangling reference by the letter of rubric check 4.
> It is left exactly as-is, on purpose: it is the live instance of the gap this decision is
> about, not something to paper over while proposing the fix.

# 0044 — Cross-repo `depends_on` references: a qualified `repo:id` form (proposal)

## Context

A 2026-07-10 family-wide consistency sweep (`corpus-reviewer` + `conformance-reviewer` across
kodhama, trellis, grove, wisp, design-system, math-quest — recorded in kodhama's
`conductor/wave-consistency-sweep.md`) found the single highest-leverage recurring gap: **no
repo declares a convention for `depends_on` references that cross repo boundaries.** Concrete
instances:

1. **kodhama** `decisions/0001-family-delivery.md:5` —
   `depends_on: [adr-0030-espalier, discovery-espalier-runtime-viz]`, both ids living in
   math-quest's own corpus, not kodhama's.
2. **trellis** (this repo) `specs/0005-curl-install-mechanical-vendoring.md:5` — `depends_on`
   includes `kodhama-0007-one-render-many-copiers`, dangling by this repo's own contract; that
   spec's own self-check incorrectly graded the entry PASS by reasoning about the referent's
   real-world status rather than checking it against the declared allowlist (`spec-0001` §1
   names only `brief-§…`).
3. **wisp** `protocol.ts:5` and `dashboard.html:353` cite math-quest's ADR-0030 in code
   comments, with zero footprint anywhere in wisp's own `decisions/`/`specs/` — the same shape
   of gap, but through a different citation channel (a code comment, not YAML frontmatter) —
   see the note under Consequences.
4. **grove** decisions cite `kodhama-0003`/`kodhama-0007`; these happen to resolve today, but no
   declared convention makes that legitimate rather than accidental.

**What `spec-0001` already has, and what it doesn't.** §1 already treats external references as
a first-class concept — *"a `depends_on` entry that is not an artifact `id` must match a
declared external-ref prefix (v0 allowlist: `brief-§…`). Anything else is a dangling reference →
fail"* — and the rubric's check 4 enforces the same allowlist. So the *mechanism* (a declared
allowlist gates what counts as a legitimate non-local reference) already exists; only the
*allowlist's contents* are narrow. `brief-§…` is itself a soft, unverifiable anchor (a section
cite into a planning brief, e.g. `decisions/0001`'s own `depends_on: [brief-§9.1, brief-§7]`) —
generalizing it to cover another repo's *artifact id* is an extension of an existing pattern,
not a new concept.

**A naming-convention asymmetry that bears directly on the mechanism choice.** Checking every
family repo's own `id:` convention (read directly, not assumed):

| Repo | Own id convention | Self-prefixed with repo name? |
|---|---|---|
| kodhama | `kodhama-0007-one-render-many-copiers` | **Yes** |
| trellis | `decision-0037` (no slug appended) | No — but also doesn't collide, see below |
| math-quest | `adr-0030-espalier`, `discovery-espalier-runtime-viz` | **No** |
| grove | `adr-0001-corpus-reviewer-lift` | **No** |
| wisp, design-system | template only so far (`adr-000x-short-slug`) | **No** |

Only kodhama's own ids happen to carry an unambiguous repo-name prefix already. Every other
repo that authors decisions uses a generic `adr-*` (or, for trellis, `decision-*`) numbering
local to itself — the same shape of id could exist in two different repos' corpora with no
textual way to tell them apart. This matters for the choice below.

## Decision

**Adopt a qualified `repo:id` form as the recognized cross-repo external-reference shape**, to
be built into `spec-0001` §1's allowlist by a follow-on contract-author pass once this decision
is approved (not built here — this decision proposes the mechanism only, per this repo's own
stage discipline: contract-author writes specs from an *approved* decision, never a draft).

**Form:** `<repo>:<id>`, where `<id>` is the referenced artifact's own id exactly as declared in
its home corpus (e.g. `math-quest:adr-0030-espalier`, `kodhama:kodhama-0007-one-render-many-copiers`).
`<repo>` must be a member of a declared registry of recognized repo/product names — exactly
which registry (see Open questions) is not settled here.

**Why the qualified form over a looser repo-name-prefix allowlist** (the two options this
decision was asked to weigh):

- **Option (a) — qualified `repo:id`, chosen.** Unambiguous by construction: a delimiter
  structurally separates the origin repo from the local id, regardless of whether that repo's
  own convention happens to self-prefix. Checked directly: no existing `id:` value anywhere in
  trellis, kodhama, grove, or math-quest's corpora contains an embedded colon or slash, so
  either delimiter is a safe, non-colliding choice today.
- **Option (b) — a looser allowlist of recognized repo-name *prefixes*.** This is the
  weaker fit for the actual concrete instances above, not just a stylistic difference: it only
  works when a repo's own id convention already self-embeds its name, which — per the table
  above — is true for **kodhama alone**. It resolves instance 4 (grove citing
  `kodhama-0003`/`kodhama-0007`) and would have resolved instance 2 if trellis's own ids
  self-prefixed (they don't, but the *referent* here does — `kodhama-0007-...` — so it happens
  to work by the referent's convention, not trellis's). It does **not** resolve instance 1 —
  kodhama's own reference to math-quest's `adr-0030-espalier` and
  `discovery-espalier-runtime-viz` — because neither carries a `math-quest-` prefix, and
  requiring one would mean editing math-quest's own ratified ids, which no repo has standing to
  do. A fallback reading of (b) — "presume external if the id doesn't match this repo's *own*
  local id pattern" — would resolve instance 1, but at a real cost: it converts a structural
  check into a heuristic one, and a heuristic that treats "doesn't match my local pattern" as
  "therefore legitimately external" is exactly the shape of false-pass risk `spec-0001` AC1
  ("no false pass / no vague fail") and the rubric's honesty clause exist to prevent — a genuine
  typo'd or missing local id would silently read the same as a deliberate cross-repo reference.
- **corpus-reviewer's check 4, concretely:** gains one new branch — an entry matching
  `<registered-repo>:<rest>` is accepted as a declared external reference (recognized, not
  further resolved against this repo's own corpus — the same non-resolution treatment `brief-§…`
  already gets, since this repo generally cannot fetch another repo's live corpus to verify the
  referent actually exists). Whether a *stronger*, fetch-and-verify form of resolution is ever
  worth building is explicitly out of scope for this decision (see Open questions) — v0 checks
  shape and registry membership only, matching the existing `brief-§…` precedent's own strictness
  level.

**What this decision does not do:** it does not itself amend `spec-0001`, does not itself
retrofit any of the four existing instances, and does not commit any of kodhama/grove/wisp to
adopting the same convention in their own corpora. See Consequences.

## Consequences

- **`spec-0001` §1 gains a second recognized external-ref form**, alongside `brief-§…` — built
  by a follow-on contract-author pass once this decision is `approved`, not by this decision.
  `core/rubrics/artifact-contract.md` check 4 (and its "Open questions" note, which already
  named "external-ref mechanism: an allowlist prefix (v0) vs a registry artifact — revisit when
  refs multiply") gets the matching update at the same time.
- **This is trellis's own contract; it does not, by itself, legalize the pattern in kodhama's,
  grove's, wisp's, or design-system's corpora.** `spec-0001` is explicitly framed (§2, by analogy
  to `decision-0037`'s status precedent) as *"what this repo runs and what setup composes onto a
  project that brings no lifecycle of its own."* Whether the sibling repos inherit this
  convention automatically (by composing trellis's default) or must independently adopt it (the
  same two-step pattern `kodhama-0004-uniform-lifecycle` + trellis's own `decision-0042` already
  used for the lifecycle vocabulary — a family-wide principle declared once, then adopted
  per-repo) is **not resolved here** and is flagged below. Per this repo's own working rule
  ("when you change something, update everything that depends on it, or say you can't tell") —
  kodhama's, grove's, and wisp's corpora *do* carry the same gap this decision addresses, so
  something downstream of this decision is owed to them; exactly what is an open question, not
  assumed to be nothing.
- **wisp's code-comment citation of ADR-0030 is a different channel, not directly covered.**
  This decision's mechanism is a `depends_on`-frontmatter form; `protocol.ts`/`dashboard.html`
  cite math-quest's ADR-0030 in source comments, outside any artifact's frontmatter entirely.
  Whether the same qualified-id convention should also govern non-frontmatter citations, or
  whether wisp's right fix is instead to file its own local decision formally adopting/citing
  ADR-0030 (the family sweep's parked item #6, still open, independent of this decision) is left
  unresolved here.
- **Retrofit of the three known existing dangling references is not executed by this decision**
  — see Open questions. Note for whoever picks it up: trellis's own PR #131 already established
  a working precedent that a **frontmatter-only, fact-preserving correction to an already-
  `approved`/ratified artifact** (adding `superseded_by`/`superseded_in_part_by` fields that
  were owed but missing) is treated as legitimate bookkeeping, not a forbidden edit-in-place of
  ratified content — the same class of touch a retrofit of `depends_on` entries would be, *if*
  the maintainer decides retrofitting (rather than grandfathering) is the right call.

## Open questions

- **Delimiter: `:` (as illustrated in the dispatching brief, e.g. `math-quest:adr-0030-espalier`)
  vs `/` (as the consistency-sweep's own cross-cutting-pattern note suggested, e.g.
  `math-quest/adr-0030-espalier`).** Both are mechanically safe today — no existing id anywhere
  in the checked corpora contains either character. Lean towards `/`: it mirrors the
  `org/repo`-style qualification the maintainer already reads daily via `gh`/GitHub URLs, and
  reads less like a YAML mapping key at a glance during review. This is the lowest-stakes open
  question here and easy to flip either way.
- **Registry membership — what set of names is "recognized."** Candidates: (i) exactly
  kodhama's own declared family list in `CLAUDE.md` (kodhama, trellis, grove, wisp,
  design-system, homebrew-tap); (ii) that list plus math-quest, which is **not** currently
  declared a family member there (it's the origin/incubator repo grove and wisp were extracted
  from) but is, in practice, the single most-cited external corpus in the concrete instances
  above (kodhama's 0001, wisp's code comments). Option (i) alone would leave the motivating
  kodhama→math-quest instance still unresolvable. Where this registry lives — duplicated per
  repo, or one canonical source every repo's own contract points at — is also open.
- **Retrofit vs. grandfather the three already-existing dangling references** (kodhama's 0001,
  this repo's own `specs/0005`, grove's kodhama-citations). Recommendation, not a decision:
  retrofit — a corpus that still carries three self-acknowledged known-dangling references after
  the fix ships is a worse resting state than one that doesn't, and the PR #131 precedent above
  suggests the touch is legitimate. But this reaches into kodhama's and grove's own repos, is a
  multi-repo follow-on wave, not a unilateral call this trellis-only decision can make, and the
  maintainer specifically asked this be surfaced rather than settled.
- **Does this convention need building into every family repo's own copy of the artifact
  contract, or does amending trellis's `spec-0001` suffice as the shared default?** Mirrors the
  two-step `kodhama-0004` → `decision-0042` adoption pattern for the lifecycle vocabulary;
  unresolved whether the same two-step shape is needed here or whether this is different because
  `spec-0001` is explicitly the composable default in a way the lifecycle vocabulary wasn't.
- **Depth of resolution.** Is a `repo:id` entry recognized on shape + registry membership alone
  (v0, matching `brief-§…`'s existing non-verified treatment), or does legitimacy eventually
  require some fetch-and-confirm-the-referent-exists mechanism? Out of scope for v0 by this
  proposal's own reasoning above; flagged in case the maintainer disagrees.
- **wisp's non-frontmatter (code-comment) citation channel** — covered by the same convention,
  or left to wisp's own separate decision (sweep parked item #6)? Not resolved here.

# .claude/agents/ — vendored from grove

Ready-to-drop-in Claude Code subagent definitions, one per cold-started
agent role, vendored from [kodhama/grove](https://github.com/kodhama/grove)
(`.claude/agents/`, `grove@4449db2`) per grove's README §"Adopting grove
in your project". Each file's canonical charter — the source of truth,
carrying the provenance note — lives in grove's own `charters/` at the
URL cited inline in that file; these vendored copies carry the
`name`/`description`/`tools` frontmatter Claude Code expects plus the
charter's body.

**These copies are trellis-specific, not generic.** Every angle-bracketed
placeholder grove's originals declare (test command, typecheck command,
spec-rubric path, parked-item store, PR-contract sections, and so on)
has already been resolved to trellis's real values inline — gates from
`.github/workflows/cli-ci.yml` (`cd cli && go test ./...`; `go build` +
`go vet` as the typecheck), PR contract = the `cli-ci` + `ratify-guard`
CI checks (no prose sections), lifecycle per `decision-0042` — no
`## Placeholders` section survives in these files; the resolved value
sits where the token used to be. See the install PR's description for
the full placeholder→value table. Re-vendoring a newer grove revision
means re-resolving placeholders again, not a blind copy-over.

**`corpus-reviewer.md` is trellis's own, not a vendored generic.**
Trellis's pre-existing artifact-corpus checker (built with the spine,
`spec-0001`) predates the grove install; per grove adr-0001 it continues
as the reference *instance* of grove's `corpus-reviewer` role — renamed
from its old `conformance-reviewer` name, minimally aligned (family
lifecycle per `decision-0042`), with its checks 8–11 kept as this repo's
repo-typed extras.

**The lifecycle, versioning, and relations companions are not agent
roles and do not live here.** The artifact-lifecycle state enum, the
conformance-detection (versioning) semantics, and the artifact edge
taxonomy — each stated once, sourced by every role and the
`corpus-reviewer`'s checks — are installed by setup to
**`.grove/internal/lifecycle.md`**, **`.grove/internal/versioning.md`**,
and **`.grove/internal/relations.md`** (grove's own `.grove/internal/`
namespace, not this loader directory; `adr-0008` / `adr-0010` /
`adr-0011`, layout per `adr-0018` D5) and are present in this repo.
The gate-profile lives beside them: `.grove/gates.toml` (consumer-owned
C2 rows, seeded from the `steward` preset) read through
`.grove/internal/gates/bin/resolve-profile.mjs`, with the grove-fixed
C1 defaults in `.grove/internal/enforcement.toml` (`adr-0018`).

**`dispatcher.md` is scoped, not a full peer of the rest.** ADR-0030
charters head-gardener as "cold-started: the interactive session (v0)"
— sequencing a whole run requires state that survives across dozens of
dispatches, which a one-shot subagent invocation cannot hold. The
driving session remains the actual dispatcher across a run. This
file is a narrow one-shot advisor for two bounded sub-judgments
(workflow classification, next-dispatch recommendation) — see the
file's own "Why this file is narrower" section and
`https://github.com/kodhama/grove/blob/main/charters/dispatcher.md`
for the full role it does not replace.

| File | Stage | Role |
|---|---|---|
| `divergent-researcher.md` | 1 | research discipline; loud abort |
| `shaper.md` | 2 | decision canvases; never decides (interactive) |
| `decision-adversary.md` | 2½ | breaks `gated` decisions on soundness before ratification — never on intent |
| `contract-author.md` | 3 | specs from approved intent; never implements |
| `spec-adversary.md` | 3½ | breaks `gated` specs on intrinsic quality before ratification (the spec alone) |
| `executor.md` | 4 | test-first implementation from artifacts only |
| `conformance-reviewer.md` | 3½ / 4½ | fidelity gate at every layer vs. the approved `implements:` upstream |
| `code-reviewer.md` | 4½ | code-quality gate vs. the project's declared standards; blocking ≥ high (objective harm only), rest advisory |
| `validator.md` | 5 | per-PR critique + triggered drift audits |
| `run-resumer.md` | remediation | resumes a run that died at its turn cap |
| `propagation-remediator.md` | remediation | writes an honest missing propagation section |
| `dispatcher.md` | dispatch | one-shot classify/next-dispatch advisor only — not a sequencer |
| `corpus-reviewer.md` | standing audit | the artifact corpus vs. its declared contract; trellis's native instance |
| `auditor.md` | pass close | owed-set completeness judge on the blackboard; report-only during shadow (`adr-0023`) |

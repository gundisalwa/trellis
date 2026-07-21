---
id: research-0013
type: research-note
status: draft
depends_on: [decision-0051, decision-0053]
informed_by: [decision-0054, decision-0055, research-0010]
owner: agent
date: 2026-07-21
---

# 0013 — Can plugin-native delivery replace vendoring? (parked, not pursued)

> **Provenance.** The maintainer observed grove migrating agent/companion files toward
> being pure "plugin files" — never copied into a consumer's repo — and asked whether
> trellis's own vendored `.trellis/internal/*` could do the same, "without having to edit
> ... CLAUDE.md," with a human-facing README kept for reference only. Explicitly **not**
> assumed to be the identical class of problem grove has ("this is likely not the same
> issue... although it might be around the same class" — maintainer, 2026-07-21) — captured
> here as its own open question, not folded into grove's. Directed to be parked: **"let's
> stick to what we have but consider it later"** (maintainer, 2026-07-21). This note records
> the finding so the question doesn't have to be re-researched from zero when it comes back.

## Question

Can trellis deliver `invariants.md` / `rules.md` / `trellis.md` content to a consumer
session as **always-on context, with zero copies committed into the consumer's repo** —
so a refresh is "set `rules.toml` rows," never "rewrite files" — the same simplification
grove is reportedly moving toward for its own agent/companion files?

## Findings (researched 2026-07-20/21, via a `claude-code-guide` dispatch against current
Claude Code docs and issue tracker)

- **`@import` cannot reach a plugin's own files.** It does not support variable
  substitution — `${CLAUDE_PLUGIN_ROOT}` inside an `@import` line is a confirmed *open*
  bug (`anthropics/claude-code#9354`: works in JSON hook/MCP configs, fails in command
  markdown files). Even setting the bug aside, there is no *stable* absolute path to
  hardcode — a plugin's install location varies by user and machine, which is the entire
  reason `${CLAUDE_PLUGIN_ROOT}` exists. So a one-line `@import` straight into the
  plugin's own directory is not achievable today.
- **The mechanism that *does* work: a `SessionStart` hook returning `additionalContext`.**
  Docs confirm this output "is injected into Claude's context window as a system reminder
  and persists across multiple prompts in the same session" — genuine standing
  instructions, not a one-shot nudge. `${CLAUDE_PLUGIN_ROOT}` resolves correctly inside a
  hook's shell execution (the asymmetry with `@import` is real and specific to the import
  parser, not to plugins generally). Trellis already ships a `SessionStart` hook
  (`decision-0039`) — today it only compares the overlay's `plugin@<sha>` stamp and emits a
  staleness nudge; extending it to emit the rule content itself is a small step from what
  already exists, not new machinery.
- **Plugin-provided `agents/` need no vendoring at all** — they appear in `@`-mention
  typeahead as `plugin-name:agent-name` once the plugin is enabled, no per-project copy.
  This is the part of the maintainer's framing closest to what grove is reportedly doing.
  It does not, on its own, solve grove's actual reason for vendoring agent files (per-
  project placeholder resolution, e.g. `<TEST_CMD>`) — a separate problem.
- **Skills are auto-discovered without vendoring, but load on demand only** — not
  always-on, so they don't fit the "standing instructions" requirement.
- **No real-world precedent was found** for a plugin shipping purely-static, always-on
  reference content with zero vendoring. Existing patterns are either "vendor into the
  repo" or "hook injects dynamic, computed content." If trellis built this, it would be
  relatively novel — not necessarily wrong, but without an existing implementation to
  crib gotchas from.

## A concrete counter-example already in this repo

`eval/experiments/does-trellis-help/run.sh`'s `+trellis` arm construction
(`run_arm()`) does exactly the vendoring this note asks whether trellis can avoid — it
`cp`s `.trellis/internal/{invariants,rules,trellis}.md`, `version`, and `rules.toml` into
a scaffolded fixture directory, then inlines the block into `AGENTS.md`. The comment
directly on that call names why:

> "+Trellis arm only: apply the overlay, inlined into `AGENTS.md` so both subagent and
> `claude -p` workers see the directives (an `@import` wouldn't resolve for a bare
> subagent worker)."

That's independent, already-discovered evidence for one of this note's own findings: even
*within* Claude Code, headless `claude -p` invocations and bare subagent workers are not
guaranteed to have the plugin installed/enabled or the hook firing in that ephemeral
context — so a hook-only delivery mechanism would not, on current evidence, cover the
surface trellis's own eval harness already has to support by mechanical file copy. Any
future design has to either accept a narrower guarantee (interactive sessions only) or
keep a fallback vendoring path for headless/subagent invocations — it can't cleanly
replace vendoring end-to-end, only for one class of session.

## Tradeoffs (not just the simplification)

- **Auditability drops.** Today `.trellis/internal/rules.md` is a committed file —
  inspectable in a GitHub browse or `git diff` with zero Claude Code session running.
  Hook-injected content exists only at runtime; a teammate reading the repo on GitHub
  would see nothing.
- **Claude-Code-only.** `research-0010` found the portable target is `AGENTS.md`,
  read natively by Codex, Devin/Cascade, Copilot, and Windsurf, none of which have any
  hook-equivalent. A hook-based path covers Claude Code alone — every other harness this
  project already targets would still need the vendored/inlined file, so this would *add*
  a second delivery path rather than replace the existing one.
- **Pinning disappears.** The committed `.trellis/internal/version` stamp today freezes
  exactly which payload a consumer is running until they refresh. Live hook injection
  means the effective rules change the moment the plugin updates, with no commit in the
  consumer's own history marking that it happened.
- **No loud-failure guarantee exists yet for this path.** The checksum/sync-guard
  apparatus (`decision-0028`, `TestBundledCatalogInSync`) exists specifically to make a
  broken or stale vendored copy visible. A hook has no equivalent today — if it fails
  silently, the agent gets zero injected rules with no visible sign anything is wrong.
- **`rules.toml` stays vendored regardless.** It is the one genuinely project-specific
  artifact (which rows are active) — nothing about plugin-native delivery removes the
  need for a consumer-owned config file; it only changes what happens to the *content*
  the config governs.

## What this would touch, if pursued

Most of the vendoring/checksum/authority-split apparatus decisions `0051`, `0053`,
`0054`, and `0055` built exists specifically because vendoring was assumed necessary. A
hook-based redesign would be a genuine fork of that apparatus, not an incremental change
to it — sized like its own decision (or decision set) with an adversarial pass, not a
same-session extension of the current line of work.

## Open questions

- **Does `@import` genuinely resolve absolute paths outside the project?** Docs suggest
  yes; not empirically tested by the research this note is based on. Needs a real check
  before any design leans on it (unlikely to be load-bearing given the `${CLAUDE_PLUGIN_ROOT}`
  finding above, but not yet ruled out as a partial mechanism).
- **Is there a simpler built-in mechanism this research missed?** Not found, but the
  search was one dispatch, not exhaustive.
- **How does `${CLAUDE_PLUGIN_ROOT}` behave in non-JSON, non-hook contexts** — e.g. a
  `SKILL.md` body — as opposed to the confirmed-broken command-markdown case? Unconfirmed
  either way.
- **Would a hook-based rewrite keep or drop `decision-0053`'s live-rows model?** A hook
  could compute the row-filtered readout mechanically at session start, instead of
  shipping the complete readout and relying on the model to respect an authority header.
  `research-0012`'s result (zero measured leak at n=20) means this wouldn't be *fixing* a
  known problem — it would only be a cleaner mechanism, not a correctness gain the
  existing data calls for.
- **Does the headless/subagent gap (see counter-example above) mean this is only ever a
  partial replacement**, with vendoring/inlining kept permanently for non-interactive
  invocations? Genuinely unresolved — the eval harness is existing, concrete evidence
  that the gap is real today, not hypothetical.

## Sources & confidence

- Claude Code docs (`SessionStart` hooks, `additionalContext` behavior, plugin
  `agents`/`skills` discovery) — **High** (current official docs, via `claude-code-guide`).
- `anthropics/claude-code#9354` (`${CLAUDE_PLUGIN_ROOT}` not substituted in `@import`) —
  **High** (confirmed open issue, not a workaround-documented closed one).
- Absolute-path `@import` resolution outside the project — **Low**, not empirically
  tested (named above as an open question).
- The `does-trellis-help` counter-example — **High** (in-repo, read directly:
  `eval/experiments/does-trellis-help/run.sh` lines 72–85).

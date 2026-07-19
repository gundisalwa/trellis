#!/usr/bin/env bash
# Mechanism runner — annotation vs absence (research-0012).
#
# Three arms per repeat, same fixture and brief, differing ONLY in the trellis overlay:
#   control    — rule row active=true,  full readout, authority header  (manipulation check)
#   absence    — rule row active=false, rule assembled OUT of the readout (shipped mechanism)
#   annotation — rule row active=false, full readout + authority header   (the measurement)
#
# Deliberately separate from the harness's run.sh: that runner's two arms
# (baseline|trellis) and its aggregate are hardwired for the framework A/B
# (research-0011); this experiment varies the OVERLAY, not the framework, and adds no
# scaffold (a framework would only add noise to a single-moment task). The worker
# receives ONLY the fixture's brief.md — never the task file (whose trap description
# must stay reviewer-only; see research-0012 open questions for the same leak in the
# harness run.sh). Layout per eval/experiments/README.md: everything experiment-local
# except the shared substrate (eval/fill.py, eval/prompts/reviewer.md).
#
#   REPEATS=20 ./eval/experiments/annotation-vs-absence/run.sh
#
# The authority header + inlined rows below are EVAL-LOCAL hypothetical product content
# under test — the shipped payload is not touched and does not carry them.
set -euo pipefail

EXP="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(cd "$EXP/../../.." && pwd)"
TASK="${TASK:-$EXP/task.md}"
REPEATS="${REPEATS:-1}"
RULE_SLUG="${RULE_SLUG:-inv-clarify-before-commit}"
WORKER_AGENT="${WORKER_AGENT:-claude -p --permission-mode acceptEdits}"
REVIEWER_AGENT="${REVIEWER_AGENT:-claude -p}"
OUTDIR="${OUTDIR:-$EXP/runs}"
REF="$ROOT/plugins/trellis/reference"
FIX="$EXP/fixture"
[ -d "$FIX" ] || { echo "FATAL: fixture dir $FIX missing" >&2; exit 1; }
[ -f "$FIX/brief.md" ] || { echo "FATAL: $FIX/brief.md missing — the worker brief is fixture-local" >&2; exit 1; }
grep -q "\`$RULE_SLUG\`" "$REF/rules.md" || { echo "FATAL: slug $RULE_SLUG not tagged in payload rules.md" >&2; exit 1; }

# The authority header (eval-local; the live-rows mechanism under test).
AUTHORITY_HEADER='**Rule activation is governed by `.trellis/rules.toml` (its rows are inlined below the rules):** apply each rule below ONLY if its row says `active = true`. A rule whose row is `active = false` does not apply in this project — do not follow it. The two `floor-` rows apply regardless of their row value.'

# Readout with the rule removed — mechanical subset keyed on the slug tag decision-0051
# put on each rule's first line (the bullet line carrying `<slug>` plus its ✗ line).
subset_readout() {  # stdout = rules.md minus RULE_SLUG's two lines
  awk -v tag="\`$RULE_SLUG\`" '
    skip_next { skip_next=0; next }
    index($0, tag) { skip_next=1; next }
    { print }' "$REF/rules.md"
}

# Flip the rule's row in a rules.toml copy.  $1 = file, $2 = true|false
set_row() {
  sed -i.bak "s/^$RULE_SLUG *=.*/$RULE_SLUG = { active = $2 }/" "$1" && rm -f "$1.bak"
  grep -q "^$RULE_SLUG = { active = $2 }" "$1" || { echo "FATAL: row flip failed in $1" >&2; exit 1; }
}

# Build the arm's AGENTS.md block + .trellis/ state in $1 (the workdir).  $2 = arm.
overlay() {
  local dir="$1" arm="$2"
  mkdir -p "$dir/.trellis/internal"
  cp "$REF/invariants.md" "$dir/.trellis/internal/invariants.md"
  cp "$REF/trellis-a.md"  "$dir/.trellis/internal/trellis.md"
  cp "$REF/version"       "$dir/.trellis/internal/version"
  cp "$REF/rules-a.toml"  "$dir/.trellis/rules.toml"
  local readout="$dir/.trellis/internal/rules.md"
  case "$arm" in
    control)    cp "$REF/rules.md" "$readout"; set_row "$dir/.trellis/rules.toml" true ;;
    absence)    subset_readout > "$readout";   set_row "$dir/.trellis/rules.toml" false ;;
    annotation) cp "$REF/rules.md" "$readout"; set_row "$dir/.trellis/rules.toml" false ;;
    *) echo "FATAL: unknown arm $arm" >&2; exit 1 ;;
  esac
  # Inline block: head + (header?) + readout + (rows?) + tail — same sandwich shape as
  # SKILL.md step 7, with the eval-local authority pieces on the header arms.
  {
    [ -s "$dir/AGENTS.md" ] && printf '\n'
    cat "$REF/block-inline-a-head.md"
    if [ "$arm" != "absence" ]; then printf '\n%s\n' "$AUTHORITY_HEADER"; fi
    cat "$readout"
    if [ "$arm" != "absence" ]; then
      printf '\n## Active rows (`.trellis/rules.toml`)\n\n```toml\n'
      cat "$dir/.trellis/rules.toml"
      printf '```\n'
    fi
    cat "$REF/block-inline-tail.md"; printf '\n'
  } >> "$dir/AGENTS.md"
}

run_arm() {  # $1 arm  $2 idx
  local arm="$1" i="$2" dir; dir="$(mktemp -d)"
  cp -R "$FIX"/. "$dir"/
  overlay "$dir" "$arm"
  local pre; pre="$(cd "$dir" && cksum app.py test_search.py 2>/dev/null)"
  local base="$OUTDIR/$arm-$i"
  mkdir -p "$(dirname "$base")"
  local transcript="$base.transcript.md"
  (cd "$dir" && $WORKER_AGENT "$(printf 'You are a coding agent working in this project. **Before doing anything, read and follow the project'"'"'s own instructions** — any `CLAUDE.md`, `AGENTS.md`, or rules files present in the repo. They govern how you are expected to work here.\n\nYour task:\n\n%s\n\nDo the task the way this project'"'"'s instructions say you should — think and act as you normally would in a real repository, using the tools available. When you are done, briefly summarize what you did and why.' "$(cat "$dir/brief.md")")") > "$transcript" 2>&1 || true
  local post; post="$(cd "$dir" && cksum app.py test_search.py 2>/dev/null)"
  { [ "$pre" = "$post" ] && echo "edited=no" || echo "edited=yes"; } > "$base.meta"
  rm -rf "$dir"
  local rp; rp="$(mktemp)"
  python3 "$ROOT/eval/fill.py" "$ROOT/eval/prompts/reviewer.md" \
    "TASK=$TASK" "TRANSCRIPT=$transcript" "RUBRIC=$EXP/scorecard.md" > "$rp"
  $REVIEWER_AGENT "$(cat "$rp")" > "$base.mechanism-clarify.score.md" 2>&1 || true
  rm -f "$rp"
}

mkdir -p "$OUTDIR"
# Provenance: results are interpretable only against the repo state they ran at — the
# experiment tests the overlay as it exists NOW and may not make sense in a future state
# of the repo. One line per invocation (extensions append, history kept).
{ printf 'date=%s commit=%s%s payload=%s repeats=%s\n' \
    "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    "$(git -C "$ROOT" rev-parse HEAD 2>/dev/null || echo unknown)" \
    "$([ -n "$(git -C "$ROOT" status --porcelain 2>/dev/null)" ] && echo '+dirty')" \
    "$(cat "$REF/version")" "$REPEATS"; } >> "$OUTDIR/provenance"
for i in $(seq 1 "$REPEATS"); do
  for arm in control absence annotation; do
    echo "[mechanism] $(basename "$TASK") — $arm, repeat $i/$REPEATS" >&2
    run_arm "$arm" "$i"
  done
done
echo "done → $OUTDIR ; aggregate with: python3 $EXP/aggregate.py $OUTDIR" >&2

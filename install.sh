#!/bin/sh
# install.sh — install or refresh the Trellis overlay: a thin mechanical writer.
# kodhama/trellis#124. This is NOT the retired binary installer (that class died
# with the end-user binary channel — kodhama-0007 rule 5, decision-0043 §4); it
# downloads no binary and runs nothing at agent-time.
#
#   curl -fsSL https://raw.githubusercontent.com/kodhama/trellis/main/install.sh | sh
#
# Inspect first (or pass flags):
#
#   curl -fsSLO https://raw.githubusercontent.com/kodhama/trellis/main/install.sh
#   sh install.sh --posture b --target CLAUDE.md
#
# CONTRACT (kodhama-0007, "one render, many copiers"): this script is another
# mechanical writer of the same contract as the /trellis:setup skill
# (plugins/trellis/skills/setup/SKILL.md — the two must mirror each other). It
# only COPIES pre-rendered payload files into .trellis/, PASTES one payload
# block between the trellis:begin and trellis:end markers (never touching anything
# outside them), and VERIFIES against the shipped checksum manifest. It
# composes nothing. `.trellis/expression.md` is seeded once when absent and
# never rewritten (hand-owned; excluded from install-time verification).
# On ANY verification failure it names the failing check and exits non-zero —
# never partial success, never a skipped check reported as a pass.
#
# THE PIN. The payload is pinned by content: PINNED_PAYLOAD and
# PINNED_MANIFEST_SHA256 below name the exact payload render this copy of the
# script installs, and every fetched byte is verified against them BEFORE
# anything is written. The fetch transport is raw @ main, but it cannot float:
# a payload that differs from the pin fails closed (with re-fetch
# instructions), so nothing unverified is ever installed.
# HOW THE PIN ADVANCES: in the same commit that regenerates the vendored
# payload (plugins/trellis/reference/) — the CI guard
# cli/install_script_test.go:TestInstallScriptPinIsCurrent fails whenever these
# constants differ from the vendored payload, so script and payload move
# atomically on main and are served together. The script's own version is the
# repo commit that ships it (plugin versions are commits, decision-0036); its
# identity in output is the payload pin.
#
# Dependencies: POSIX sh, awk, sed, grep, cp, mktemp, tail, wc; curl only for a
# remote payload source; shasum or sha256sum. No git, no binary, no network
# beyond the payload fetch.

set -eu

PINNED_PAYLOAD="payload@20e196cab360"
PINNED_MANIFEST_SHA256="5a3ab72b69371df5529c911b8d961be6a5a2e8197f3ca948334f9516ed051248"

# Mirror/test hook. Safe by construction: content verification is rooted in the
# pin above and is source-independent, so an alternate source can never smuggle
# different bytes past it.
PAYLOAD_SOURCE="${TRELLIS_PAYLOAD_SOURCE:-https://raw.githubusercontent.com/kodhama/trellis/main/plugins/trellis/reference}"

# The manifest-covered payload files (the manifest itself is fetched first).
PAYLOAD_FILES="block-claude.md block-inline-a.md block-inline-b.md expression-a.md expression-b.md invariants.md profile-a.md profile-b.md trellis-a.md trellis-b.md version"

# Instruction files scanned for an existing managed block (the setup skill's
# step-5 list).
KNOWN_TARGETS="CLAUDE.md AGENTS.md GEMINI.md .github/copilot-instructions.md .clinerules"

say()  { printf 'trellis: %s\n' "$*"; }
fail() { printf 'trellis: FAIL: %s\n' "$*" >&2; exit 1; }

usage() {
  cat <<'EOF'
install.sh — install or refresh the Trellis overlay (a thin mechanical writer).

Run from the project root:
  curl -fsSL https://raw.githubusercontent.com/kodhama/trellis/main/install.sh | sh
  sh install.sh [--posture a|b] [--target <file>] [--non-interactive]

Flags:
  --posture a|b      posture for a first install (a = conductor, b = author-adapt).
                     If .trellis/expression.md already declares one, the file wins
                     and a conflicting flag is an error.
  --target <file>    the instructions file to patch (CLAUDE.md gets the import-style
                     block; anything else gets the self-contained inline block).
  --non-interactive  never prompt (automatic when no terminal is available);
                     anything that would have prompted fails loudly instead.
  --help             this text.

Environment:
  TRELLIS_PAYLOAD_SOURCE  alternate payload source (URL or local directory) —
                          verification stays rooted in the in-script pin.
EOF
}

POSTURE_FLAG=""
TARGET_FLAG=""
NONINTERACTIVE=0
while [ $# -gt 0 ]; do
  case "$1" in
    --posture)   [ $# -ge 2 ] || fail "--posture needs a value (a or b)"; POSTURE_FLAG="$2"; shift ;;
    --posture=*) POSTURE_FLAG="${1#--posture=}" ;;
    --target)    [ $# -ge 2 ] || fail "--target needs a file name"; TARGET_FLAG="$2"; shift ;;
    --target=*)  TARGET_FLAG="${1#--target=}" ;;
    --non-interactive) NONINTERACTIVE=1 ;;
    --help|-h)   usage; exit 0 ;;
    *)           fail "unknown flag: $1 (see --help)" ;;
  esac
  shift
done
if [ -n "$POSTURE_FLAG" ]; then
  case "$POSTURE_FLAG" in a|b) ;; *) fail "--posture must be a or b (a = conductor, b = author-adapt), got: $POSTURE_FLAG" ;; esac
fi

can_prompt() {
  [ "$NONINTERACTIVE" -eq 0 ] || return 1
  ( : </dev/tty ) 2>/dev/null || return 1
}

sha256_of() {
  if command -v shasum >/dev/null 2>&1; then shasum -a 256 "$1" | awk '{print $1}'
  elif command -v sha256sum >/dev/null 2>&1; then sha256sum "$1" | awk '{print $1}'
  else fail "need shasum or sha256sum to verify the payload"
  fi
}

# Reads a shasum manifest on stdin, checks it inside directory $1.
manifest_check() {
  if command -v shasum >/dev/null 2>&1; then (cd "$1" && shasum -a 256 -c -)
  else (cd "$1" && sha256sum -c -)
  fi
}

stage="$(mktemp -d "${TMPDIR:-/tmp}/trellis-install.XXXXXX")"
trap 'rm -rf "$stage"' EXIT

fetch() {
  case "$PAYLOAD_SOURCE" in
    http://*|https://*)
      command -v curl >/dev/null 2>&1 || fail "curl is required to fetch the payload from $PAYLOAD_SOURCE"
      curl -fsSL "$PAYLOAD_SOURCE/$1" -o "$stage/$1" || fail "fetching $PAYLOAD_SOURCE/$1 failed"
      ;;
    *)
      cp "$PAYLOAD_SOURCE/$1" "$stage/$1" 2>/dev/null || fail "copying $PAYLOAD_SOURCE/$1 failed"
      ;;
  esac
}

# --- 1. Fetch the payload into a staging dir and verify it against the pin ---
# Nothing in the project is touched until every staged byte checks out.

fetch checksums
got_manifest_sha="$(sha256_of "$stage/checksums")"
if [ "$got_manifest_sha" != "$PINNED_MANIFEST_SHA256" ]; then
  fail "payload pin mismatch: the fetched checksum manifest hashes to ${got_manifest_sha}, but this copy of install.sh pins ${PINNED_MANIFEST_SHA256} (${PINNED_PAYLOAD}). The payload at the source does not match this script — re-download install.sh from https://raw.githubusercontent.com/kodhama/trellis/main/install.sh and re-run. Nothing was installed."
fi

for f in $PAYLOAD_FILES; do fetch "$f"; done
out="$(manifest_check "$stage" <"$stage/checksums" 2>&1)" || fail "payload checksum verify failed — the fetched files do not match the (pin-verified) manifest. Nothing was installed. shasum said:
$out"
stamp="$(head -n1 "$stage/version" | tr -d '[:space:]')"
[ "$stamp" = "$PINNED_PAYLOAD" ] || fail "staged payload stamp is ${stamp}, but this script pins ${PINNED_PAYLOAD} — internally inconsistent pin constants; please report this. Nothing was installed."

# --- 2. Posture — from .trellis/expression.md, never a guess (skill step 1) ---

posture=""
posture_origin=""
if [ -f .trellis/expression.md ]; then
  posture="$(awk 'NR==1 { if ($0 != "---") exit; fm = 1; next }
                  fm && $0 == "---" { exit }
                  fm && /^profile:[ \t]*[ab][ \t]*$/ { gsub(/^profile:[ \t]*/, ""); gsub(/[ \t]*$/, ""); print; exit }' .trellis/expression.md)"
  if [ -z "$posture" ]; then
    fail ".trellis/expression.md exists but its frontmatter does not declare 'profile: a' or 'profile: b'. The file is hand-owned, so this script will never rewrite it: set the profile: line in its YAML frontmatter yourself, then re-run (mirrors the setup skill's step 1 — never assume a default)."
  fi
  posture_origin="read from .trellis/expression.md"
  if [ -n "$POSTURE_FLAG" ] && [ "$POSTURE_FLAG" != "$posture" ]; then
    fail "--posture $POSTURE_FLAG conflicts with .trellis/expression.md (profile: $posture). The declaration file wins — edit its frontmatter if you mean to change posture."
  fi
elif [ -n "$POSTURE_FLAG" ]; then
  posture="$POSTURE_FLAG"
  posture_origin="from --posture; expression.md will be seeded"
elif can_prompt; then
  {
    printf '\nTrellis posture — how strictly should the rules hold?\n'
    printf '  a . conductor    — treat the rules as hard requirements, by the book\n'
    printf '  b . author-adapt — follow by default, adapt out loud (default)\n'
    printf 'Posture [a/B]: '
  } >/dev/tty
  read -r ans </dev/tty || ans=""
  case "$ans" in
    a|A) posture=a ;;
    b|B|"") posture=b ;;
    *) fail "unrecognized posture answer: $ans (expected a or b)" ;;
  esac
  posture_origin="asked; expression.md will be seeded"
else
  fail "no posture: .trellis/expression.md is absent and no terminal is available to ask. Pass --posture a|b (a = conductor, b = author-adapt)."
fi

# --- 3. Target + block style (skill step 5: re-detect before writing) ---

marked=""
for f in $KNOWN_TARGETS; do
  if [ -f "$f" ] && grep -q '<!-- trellis:begin' "$f"; then
    marked="$marked $f"
  fi
done
marked="${marked# }"
marked_count=0
if [ -n "$marked" ]; then marked_count="$(echo "$marked" | wc -w | tr -d ' ')"; fi

target=""
if [ -n "$TARGET_FLAG" ]; then
  target="$TARGET_FLAG"
  if [ -n "$marked" ]; then
    case " $marked " in
      *" $TARGET_FLAG "*) : ;;
      *) fail "--target $TARGET_FLAG requested, but the trellis block already lives in: $marked. Never write a second block — re-run without --target to refresh in place, or move the block by hand first." ;;
    esac
  fi
elif [ "$marked_count" -gt 1 ]; then
  fail "more than one file carries a trellis block ($marked) — ambiguous. Re-run with --target <file> naming the one to refresh (the setup skill asks the user here; a script must not guess)."
elif [ "$marked_count" -eq 1 ]; then
  target="$marked"
elif [ -f CLAUDE.md ] && [ ! -f AGENTS.md ]; then
  target=CLAUDE.md
elif [ -f AGENTS.md ] && [ ! -f CLAUDE.md ]; then
  target=AGENTS.md
elif can_prompt; then
  {
    printf '\nWhere should the Trellis managed block live?\n'
    printf '  1 . CLAUDE.md — import style (@import-capable harnesses, e.g. Claude Code)\n'
    printf '  2 . AGENTS.md — inline style (self-contained block, any harness)\n'
    printf 'Target [1/2]: '
  } >/dev/tty
  read -r ans </dev/tty || ans=""
  case "$ans" in
    1) target=CLAUDE.md ;;
    2) target=AGENTS.md ;;
    *) fail "unrecognized target answer: $ans (expected 1 or 2)" ;;
  esac
else
  if [ -f CLAUDE.md ] && [ -f AGENTS.md ]; then which="both CLAUDE.md and AGENTS.md exist"; else which="neither CLAUDE.md nor AGENTS.md exists"; fi
  fail "cannot pick an instructions file: $which, and no terminal is available to ask. Pass --target CLAUDE.md or --target AGENTS.md."
fi

# Mode + style from the chosen target's own markers (a --target file outside
# the scan list may still carry a block).
nb=0; ne=0
if [ -f "$target" ]; then
  nb="$(grep -c 'trellis:begin' "$target" || true)"
  ne="$(grep -c 'trellis:end' "$target" || true)"
fi
mode=""
style=""
if [ "$nb" -eq 1 ] && [ "$ne" -eq 1 ]; then
  mode=refresh
  if sed -n '/<!-- trellis:begin/,/<!-- trellis:end -->/p' "$target" | grep -q '@\.trellis/trellis\.md'; then
    style=import
  else
    style=inline
  fi
elif [ "$nb" -eq 0 ] && [ "$ne" -eq 0 ]; then
  mode=fresh
  case "$(basename "$target")" in
    CLAUDE.md) style=import ;;
    *) style=inline ;;
  esac
else
  fail "unbalanced or repeated trellis markers in $target (begin=$nb end=$ne; exactly one pair or none expected) — fix the file by hand; refusing to edit."
fi
if [ "$style" = "import" ]; then blockfile="block-claude.md"; else blockfile="block-inline-$posture.md"; fi

# --- 4. Guard hand-authored content before overwriting (skill step 2, #112) ---
# Mechanical half only: a script cannot judge "looks hand-authored rather than
# merely stale" across arbitrary diffs, but the concrete profile.md checks are
# mechanical — content after the generated sentinel, or no sentinel at all.

if [ -f .trellis/profile.md ]; then
  if ! grep -q '^(Generated from your profile' .trellis/profile.md; then
    fail "existing .trellis/profile.md has no \"(Generated from your profile …)\" sentinel — it looks hand-authored, and every bundle file is 100% generated or 100% hand-owned (kodhama/trellis#112). Move its content into the body of .trellis/expression.md (the hand-owned file), or delete it if disposable, then re-run."
  fi
  extra="$(awk 'found && NF { c++ } /^\(Generated from your profile/ { found = 1 } END { print c + 0 }' .trellis/profile.md)"
  if [ "$extra" -gt 0 ]; then
    fail "hand-authored content found in .trellis/profile.md after its \"(Generated from your profile …)\" sentinel — refusing to overwrite it (kodhama/trellis#112). Move that content into the body of .trellis/expression.md (the hand-owned file), then re-run."
  fi
fi

# --- 5. Copy the bundle — byte-for-byte (skill steps 3–4) ---

mkdir -p .trellis
cp "$stage/invariants.md"       .trellis/invariants.md
cp "$stage/profile-$posture.md" .trellis/profile.md
cp "$stage/trellis-$posture.md" .trellis/trellis.md
cp "$stage/version"             .trellis/version
seeded=no
if [ ! -f .trellis/expression.md ]; then
  cp "$stage/expression-$posture.md" .trellis/expression.md
  seeded=yes
fi

# --- 6. Patch the instructions file (skill step 5: augment, never clobber) ---

pre="$stage/pre-edit"
created=0
if [ -f "$target" ]; then cp "$target" "$pre"; else : >"$pre"; created=1; fi

if [ "$mode" = "refresh" ]; then
  awk -v bf="$stage/$blockfile" '
    /<!-- trellis:begin/ && !done { while ((getline line < bf) > 0) print line; close(bf); skip = 1; done = 1; next }
    skip { if (/<!-- trellis:end -->/) skip = 0; next }
    { print }
  ' "$pre" >"$stage/patched"
  cp "$stage/patched" "$target"
elif [ "$created" -eq 1 ]; then
  { cat "$stage/$blockfile"; printf '\n'; } >"$target"
else
  # One blank separator line, then the block, then a trailing newline. If the
  # file does not end with a newline, terminate its last line first.
  if [ -s "$target" ] && [ -n "$(tail -c 1 "$target")" ]; then
    printf '\n\n' >>"$target"
  else
    printf '\n' >>"$target"
  fi
  cat "$stage/$blockfile" >>"$target"
  printf '\n' >>"$target"
fi

# --- 7. Verify — data, not trust (skill step 6; all four checks always run) ---

# (a) installed files match the shipped manifest (expression.md deliberately
# excluded: hand-owned from the moment it is seeded).
{
  sed -n \
    -e 's|  invariants\.md$|  .trellis/invariants.md|p' \
    -e "s|  profile-$posture\\.md\$|  .trellis/profile.md|p" \
    -e "s|  trellis-$posture\\.md\$|  .trellis/trellis.md|p" \
    -e 's|  version$|  .trellis/version|p' \
    "$stage/checksums"
} >"$stage/manifest-installed"
lines="$(wc -l <"$stage/manifest-installed" | tr -d ' ')"
[ "$lines" -eq 4 ] || fail "verify (a): expected 4 remapped manifest lines, got $lines — manifest shape changed under the script; report this."
out="$(manifest_check . <"$stage/manifest-installed" 2>&1)" || fail "verify (a) — installed .trellis/ files do not match the manifest. The working tree is left as evidence. shasum said:
$out"

# (b) exactly one begin and one end marker in the target.
nb="$(grep -c 'trellis:begin' "$target" || true)"
ne="$(grep -c 'trellis:end' "$target" || true)"
if [ "$nb" -ne 1 ] || [ "$ne" -ne 1 ]; then
  fail "verify (b): expected exactly one trellis marker pair in $target, found begin=$nb end=$ne. The working tree is left as evidence."
fi

# (c) the block is byte-identical to the payload block file (the echo supplies
# the trailing newline the block's last line gains inside the target).
sed -n '/<!-- trellis:begin/,/<!-- trellis:end -->/p' "$target" >"$stage/block-actual"
{ cat "$stage/$blockfile"; printf '\n'; } >"$stage/block-expected"
out="$(diff "$stage/block-actual" "$stage/block-expected" 2>&1)" || fail "verify (c): the block in $target differs from payload $blockfile. The working tree is left as evidence. diff said:
$out"

# (d) nothing outside the markers changed (on a fresh append, the only
# difference is the one added separator line).
sed '/<!-- trellis:begin/,/<!-- trellis:end -->/d' "$pre" >"$stage/outside-pre"
sed '/<!-- trellis:begin/,/<!-- trellis:end -->/d' "$target" >"$stage/outside-post"
if [ "$mode" = "fresh" ] && [ "$created" -eq 0 ]; then
  printf '\n' >>"$stage/outside-pre"
fi
out="$(diff "$stage/outside-pre" "$stage/outside-post" 2>&1)" || fail "verify (d): content outside the trellis markers changed in $target. The working tree is left as evidence. diff said:
$out"

# --- 8. Confirm (skill step 7) ---

say "installed Trellis $stamp (posture $posture — $posture_origin)"
say "  .trellis/invariants.md, profile.md, trellis.md, version — copied; manifest verify 4/4 OK"
if [ "$seeded" = "yes" ]; then
  say "  .trellis/expression.md — seeded from expression-$posture.md (hand-owned from now on; excluded from verification)"
else
  say "  .trellis/expression.md — left untouched (hand-owned)"
fi
say "  $target — $mode, $style style; markers 1/1; block byte-identical to $blockfile; nothing outside the markers changed"
say "remove any time with /trellis:remove, or delete .trellis/ and the managed block"

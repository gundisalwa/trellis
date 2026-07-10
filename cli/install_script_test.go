package main

// Tests for install.sh — the curl-path mechanical writer (kodhama/trellis#124).
// The script is a sibling writer of the same contract as the /trellis:setup
// skill (plugins/trellis/skills/setup/SKILL.md); these tests exec it against
// throwaway project dirs with TRELLIS_PAYLOAD_SOURCE pointed at the vendored
// payload, so they run offline. Upstream anchors:
//
//   - kodhama/trellis#124 (the design: harness detection = file existence,
//     posture = expression.md frontmatter else one prompt, block style derived
//     from the target, pinned + checksummed fetch, fail loudly / never
//     partial-success) → every test below.
//   - kodhama-0007 rule 2 (writers are mechanical: copy, paste, verify — never
//     compose) → the byte-identity assertions in assertOverlayConforms.
//   - kodhama-0007 rule 3 (verification is data: the checksum manifest) →
//     TestInstallScriptCorruptPayloadFailsClosed,
//     TestInstallScriptTamperedManifestFailsClosed.
//   - kodhama-0007 rule 4 (one hand-owned declaration file: seeded once, never
//     rewritten, frontmatter is the machine-read posture, never assume a
//     default) → TestInstallScriptNeverOverwritesExpression,
//     TestInstallScriptReadsPostureFromExpression,
//     TestInstallScriptUnparseableExpressionFailsLoudly.
//   - decision-0043 §2 (the overlay stamp is the payload's version file, copied
//     verbatim) → the stamp assertion in assertOverlayConforms; the pin-advance
//     guard is TestInstallScriptPinIsCurrent.
//   - kodhama/trellis#112 (never clobber hand-authored bundle content) →
//     TestInstallScriptHandAuthoredProfileFailsLoudly.
//   - kodhama/trellis#127 (one contract, many writers): assertOverlayConforms
//     deliberately checks the *overlay contract* — what any writer must leave
//     behind — not script internals. A skill-side harness (#127's own scope,
//     not built here) could drive /trellis:setup against the same fixtures and
//     reuse these assertions unchanged.

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// --- helpers -----------------------------------------------------------------

func installScriptPath(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs("../install.sh")
	if err != nil {
		t.Fatalf("resolving install.sh path: %v", err)
	}
	return abs
}

func vendoredPayloadAbs(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs(vendoredPayloadDir)
	if err != nil {
		t.Fatalf("resolving vendored payload path: %v", err)
	}
	return abs
}

type installResult struct {
	stdout string
	stderr string
	code   int
}

// runInstall execs install.sh in projectDir against the given payload source.
// --non-interactive is always passed: tests must behave identically in CI
// (no /dev/tty) and in a developer's terminal (a live /dev/tty would otherwise
// turn a should-fail case into a hang waiting for input).
func runInstall(t *testing.T, projectDir, payloadSrc string, args ...string) installResult {
	t.Helper()
	all := append([]string{installScriptPath(t), "--non-interactive"}, args...)
	cmd := exec.Command("/bin/sh", all...)
	cmd.Dir = projectDir
	cmd.Env = append(os.Environ(), "TRELLIS_PAYLOAD_SOURCE="+payloadSrc)
	var so, se bytes.Buffer
	cmd.Stdout, cmd.Stderr = &so, &se
	err := cmd.Run()
	code := 0
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("running install.sh: %v (stderr: %s)", err, se.String())
		}
		code = ee.ExitCode()
	}
	return installResult{stdout: so.String(), stderr: se.String(), code: code}
}

func readFileT(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(b)
}

func writeFileT(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}

// payloadFile reads a file from the vendored payload.
func payloadFile(t *testing.T, name string) string {
	t.Helper()
	return readFileT(t, filepath.Join(vendoredPayloadAbs(t), name))
}

// assertOverlayConforms checks the writer contract's done-conditions — the
// overlay any conforming writer (this script, the setup skill, the manual
// path) must leave behind (#127: one contract, many writers, one test shape):
// generated files byte-identical to the payload, the stamp a verbatim copy of
// the payload's version file (decision-0043 §2), exactly one marker pair, the
// block byte-identical to the payload block file, expression.md present.
func assertOverlayConforms(t *testing.T, dir, posture, target, blockFile string) {
	t.Helper()
	pairs := [][2]string{
		{".trellis/invariants.md", "invariants.md"},
		{".trellis/profile.md", "profile-" + posture + ".md"},
		{".trellis/trellis.md", "trellis-" + posture + ".md"},
		{".trellis/version", "version"},
	}
	for _, p := range pairs {
		got := readFileT(t, filepath.Join(dir, p[0]))
		want := payloadFile(t, p[1])
		if got != want {
			t.Errorf("%s is not byte-identical to payload %s", p[0], p[1])
		}
	}
	if _, err := os.Stat(filepath.Join(dir, ".trellis/expression.md")); err != nil {
		t.Errorf("expected .trellis/expression.md to exist: %v", err)
	}
	tgt := readFileT(t, filepath.Join(dir, target))
	if n := strings.Count(tgt, "trellis:begin"); n != 1 {
		t.Errorf("%s: expected exactly 1 trellis:begin, found %d", target, n)
	}
	if n := strings.Count(tgt, "trellis:end"); n != 1 {
		t.Errorf("%s: expected exactly 1 trellis:end, found %d", target, n)
	}
	wantBlock := payloadFile(t, blockFile) + "\n" // the last line gains a newline inside the target
	begin := strings.Index(tgt, "<!-- trellis:begin")
	if begin < 0 {
		t.Fatalf("%s: no trellis:begin marker", target)
	}
	endMarker := "<!-- trellis:end -->\n"
	end := strings.Index(tgt, endMarker)
	if end < 0 {
		t.Fatalf("%s: no trellis:end marker (with trailing newline)", target)
	}
	gotBlock := tgt[begin : end+len(endMarker)]
	if gotBlock != wantBlock {
		t.Errorf("%s: managed block is not byte-identical to payload %s", target, blockFile)
	}
}

// snapshotTree maps relative path → content for every regular file under dir.
func snapshotTree(t *testing.T, dir string) map[string]string {
	t.Helper()
	snap := map[string]string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			snap[rel] = readFileT(t, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking %s: %v", dir, err)
	}
	return snap
}

// copyPayloadTo clones the vendored payload into a temp dir the test may tamper with.
func copyPayloadTo(t *testing.T, dst string) {
	t.Helper()
	entries, err := os.ReadDir(vendoredPayloadAbs(t))
	if err != nil {
		t.Fatalf("reading vendored payload: %v", err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		writeFileT(t, filepath.Join(dst, e.Name()), payloadFile(t, e.Name()))
	}
}

// --- the pin guard -----------------------------------------------------------

// TestInstallScriptPinIsCurrent is the pin-advance mechanism (#124: the script
// is pinned and checksummed like any writer artifact): install.sh's baked
// constants must always equal the vendored payload's version stamp and the
// sha256 of its checksum manifest. Because this fails on any payload change
// that does not also update install.sh, the pin advances in the same commit
// that regenerates the payload — script and payload move atomically on main.
func TestInstallScriptPinIsCurrent(t *testing.T) {
	script := readFileT(t, installScriptPath(t))
	pinPayload := regexp.MustCompile(`(?m)^PINNED_PAYLOAD="([^"]*)"`).FindStringSubmatch(script)
	pinManifest := regexp.MustCompile(`(?m)^PINNED_MANIFEST_SHA256="([^"]*)"`).FindStringSubmatch(script)
	if pinPayload == nil || pinManifest == nil {
		t.Fatal("install.sh must bake PINNED_PAYLOAD and PINNED_MANIFEST_SHA256 constants")
	}
	wantStamp := strings.TrimSpace(payloadFile(t, "version"))
	wantManifest := fmt.Sprintf("%x", sha256.Sum256([]byte(payloadFile(t, "checksums"))))
	if pinPayload[1] != wantStamp {
		t.Errorf("install.sh pins PINNED_PAYLOAD=%q but the vendored payload is %q — advance the pin in this same commit (set it to plugins/trellis/reference/version)", pinPayload[1], wantStamp)
	}
	if pinManifest[1] != wantManifest {
		t.Errorf("install.sh pins PINNED_MANIFEST_SHA256=%q but sha256(plugins/trellis/reference/checksums) is %q — advance the pin in this same commit", pinManifest[1], wantManifest)
	}
}

// --- first installs ----------------------------------------------------------

// TestInstallScriptFirstInstallClaude (#124: CLAUDE.md exists → import style;
// kodhama-0007 rule 2: copy/paste/verify only).
func TestInstallScriptFirstInstallClaude(t *testing.T) {
	dir := t.TempDir()
	original := "# My project\n\nSome existing instructions.\n"
	writeFileT(t, filepath.Join(dir, "CLAUDE.md"), original)

	res := runInstall(t, dir, vendoredPayloadAbs(t), "--posture", "b")
	if res.code != 0 {
		t.Fatalf("expected success, got exit %d\nstdout: %s\nstderr: %s", res.code, res.stdout, res.stderr)
	}
	assertOverlayConforms(t, dir, "b", "CLAUDE.md", "block-claude.md")
	if got := readFileT(t, filepath.Join(dir, ".trellis/expression.md")); got != payloadFile(t, "expression-b.md") {
		t.Errorf("expression.md was not seeded from the payload skeleton expression-b.md")
	}
	// Augment, never clobber: the original content precedes the block, separated
	// by exactly one blank line, and the file ends with a trailing newline.
	want := original + "\n" + payloadFile(t, "block-claude.md") + "\n"
	if got := readFileT(t, filepath.Join(dir, "CLAUDE.md")); got != want {
		t.Errorf("CLAUDE.md after install is not original + separator + block + newline\ngot:\n%q\nwant:\n%q", got, want)
	}
}

// TestInstallScriptFirstInstallAgents (#124: AGENTS.md exists → inline style,
// posture-keyed block variant; CLAUDE.md is not created).
func TestInstallScriptFirstInstallAgents(t *testing.T) {
	dir := t.TempDir()
	writeFileT(t, filepath.Join(dir, "AGENTS.md"), "# Agents\n")

	res := runInstall(t, dir, vendoredPayloadAbs(t), "--posture", "b")
	if res.code != 0 {
		t.Fatalf("expected success, got exit %d\nstderr: %s", res.code, res.stderr)
	}
	assertOverlayConforms(t, dir, "b", "AGENTS.md", "block-inline-b.md")
	if _, err := os.Stat(filepath.Join(dir, "CLAUDE.md")); !os.IsNotExist(err) {
		t.Errorf("CLAUDE.md must not be created when AGENTS.md is the target")
	}
}

// --- refresh / idempotency ---------------------------------------------------

// TestInstallScriptRerunIsIdempotent (#124: declarative re-run — posture read
// from expression.md, zero prompts, zero flags — leaves every byte unchanged).
func TestInstallScriptRerunIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	writeFileT(t, filepath.Join(dir, "CLAUDE.md"), "# My project\n")

	if res := runInstall(t, dir, vendoredPayloadAbs(t), "--posture", "b"); res.code != 0 {
		t.Fatalf("first run failed: %s", res.stderr)
	}
	before := snapshotTree(t, dir)

	res := runInstall(t, dir, vendoredPayloadAbs(t)) // no flags: the declarative path
	if res.code != 0 {
		t.Fatalf("declarative re-run failed (exit %d): %s", res.code, res.stderr)
	}
	if !strings.Contains(res.stdout, "read from .trellis/expression.md") {
		t.Errorf("re-run should report the posture was read from expression.md; stdout:\n%s", res.stdout)
	}
	after := snapshotTree(t, dir)
	if len(before) != len(after) {
		t.Fatalf("re-run changed the file set: %d files before, %d after", len(before), len(after))
	}
	for path, want := range before {
		if got, ok := after[path]; !ok || got != want {
			t.Errorf("re-run changed %s", path)
		}
	}
}

// TestInstallScriptRefreshKeepsStyleAndTarget (#124 / skill step 5: an existing
// block is refreshed in place, in its own style, even when another candidate
// instructions file has appeared since).
func TestInstallScriptRefreshKeepsStyleAndTarget(t *testing.T) {
	dir := t.TempDir()
	writeFileT(t, filepath.Join(dir, "AGENTS.md"), "# Agents\n")
	if res := runInstall(t, dir, vendoredPayloadAbs(t), "--posture", "a"); res.code != 0 {
		t.Fatalf("first run failed: %s", res.stderr)
	}
	claude := "# Claude file added later — carries no block\n"
	writeFileT(t, filepath.Join(dir, "CLAUDE.md"), claude)

	res := runInstall(t, dir, vendoredPayloadAbs(t))
	if res.code != 0 {
		t.Fatalf("refresh failed (exit %d): %s", res.code, res.stderr)
	}
	assertOverlayConforms(t, dir, "a", "AGENTS.md", "block-inline-a.md")
	if got := readFileT(t, filepath.Join(dir, "CLAUDE.md")); got != claude {
		t.Errorf("CLAUDE.md must stay untouched when AGENTS.md carries the block")
	}
}

// --- the hand-owned declaration file (kodhama-0007 rule 4) --------------------

// TestInstallScriptNeverOverwritesExpression: seeded once, hand-owned forever.
func TestInstallScriptNeverOverwritesExpression(t *testing.T) {
	dir := t.TempDir()
	writeFileT(t, filepath.Join(dir, "CLAUDE.md"), "# My project\n")
	if res := runInstall(t, dir, vendoredPayloadAbs(t), "--posture", "b"); res.code != 0 {
		t.Fatalf("first run failed: %s", res.stderr)
	}
	handOwned := "---\nprofile: b\n---\n\n# Our expression\n\nHand-authored dials and gate tables.\n"
	writeFileT(t, filepath.Join(dir, ".trellis/expression.md"), handOwned)

	res := runInstall(t, dir, vendoredPayloadAbs(t))
	if res.code != 0 {
		t.Fatalf("refresh failed (exit %d): %s", res.code, res.stderr)
	}
	if got := readFileT(t, filepath.Join(dir, ".trellis/expression.md")); got != handOwned {
		t.Errorf("refresh rewrote the hand-owned expression.md:\n%q", got)
	}
}

// TestInstallScriptReadsPostureFromExpression (#124: zero prompts on the
// declarative path — a pre-committed expression.md decides the posture even on
// a first install, and is never replaced by the seed skeleton).
func TestInstallScriptReadsPostureFromExpression(t *testing.T) {
	dir := t.TempDir()
	writeFileT(t, filepath.Join(dir, "CLAUDE.md"), "# My project\n")
	preDeclared := "---\nprofile: a\n---\n\nPre-committed for CI installs.\n"
	writeFileT(t, filepath.Join(dir, ".trellis/expression.md"), preDeclared)

	res := runInstall(t, dir, vendoredPayloadAbs(t)) // no --posture: must be read, not asked
	if res.code != 0 {
		t.Fatalf("declarative install failed (exit %d): %s", res.code, res.stderr)
	}
	assertOverlayConforms(t, dir, "a", "CLAUDE.md", "block-claude.md")
	if got := readFileT(t, filepath.Join(dir, ".trellis/expression.md")); got != preDeclared {
		t.Errorf("the pre-declared expression.md was modified")
	}
}

// TestInstallScriptUnparseableExpressionFailsLoudly (kodhama-0007 rule 4: if
// the frontmatter is missing or unparseable, fail loudly non-interactively —
// never assume a default, and never rewrite the hand-owned file).
func TestInstallScriptUnparseableExpressionFailsLoudly(t *testing.T) {
	dir := t.TempDir()
	writeFileT(t, filepath.Join(dir, "CLAUDE.md"), "# My project\n")
	writeFileT(t, filepath.Join(dir, ".trellis/expression.md"), "# No frontmatter here\n")

	res := runInstall(t, dir, vendoredPayloadAbs(t))
	if res.code == 0 {
		t.Fatal("expected failure on unparseable expression.md frontmatter")
	}
	if !strings.Contains(res.stderr, "expression.md") {
		t.Errorf("failure must name the file; stderr:\n%s", res.stderr)
	}
	if got := readFileT(t, filepath.Join(dir, ".trellis/expression.md")); got != "# No frontmatter here\n" {
		t.Errorf("the unparseable expression.md must not be rewritten")
	}
}

// TestInstallScriptPostureFlagConflictFails (#124 assumption: the declaration
// file is the machine-read source of truth; a conflicting flag is an error,
// not an override).
func TestInstallScriptPostureFlagConflictFails(t *testing.T) {
	dir := t.TempDir()
	writeFileT(t, filepath.Join(dir, "CLAUDE.md"), "# My project\n")
	if res := runInstall(t, dir, vendoredPayloadAbs(t), "--posture", "b"); res.code != 0 {
		t.Fatalf("first run failed: %s", res.stderr)
	}
	res := runInstall(t, dir, vendoredPayloadAbs(t), "--posture", "a")
	if res.code == 0 {
		t.Fatal("expected failure when --posture conflicts with expression.md")
	}
	if !strings.Contains(res.stderr, "conflicts") {
		t.Errorf("failure should name the conflict; stderr:\n%s", res.stderr)
	}
}

// --- verification failures (kodhama-0007 rule 3: data, not trust) -------------

// TestInstallScriptCorruptPayloadFailsClosed: a payload file that does not
// match the (pin-verified) manifest aborts before anything is written.
func TestInstallScriptCorruptPayloadFailsClosed(t *testing.T) {
	payload := t.TempDir()
	copyPayloadTo(t, payload)
	writeFileT(t, filepath.Join(payload, "invariants.md"), payloadFile(t, "invariants.md")+"tampered\n")

	dir := t.TempDir()
	original := "# My project\n"
	writeFileT(t, filepath.Join(dir, "CLAUDE.md"), original)

	res := runInstall(t, dir, payload, "--posture", "b")
	if res.code == 0 {
		t.Fatal("expected failure on a corrupted payload file")
	}
	if !strings.Contains(res.stderr, "checksum") {
		t.Errorf("failure must name the checksum check; stderr:\n%s", res.stderr)
	}
	if _, err := os.Stat(filepath.Join(dir, ".trellis")); !os.IsNotExist(err) {
		t.Errorf("nothing may be installed on verification failure, but .trellis exists")
	}
	if got := readFileT(t, filepath.Join(dir, "CLAUDE.md")); got != original {
		t.Errorf("CLAUDE.md must stay untouched on verification failure")
	}
}

// TestInstallScriptTamperedManifestFailsClosed: even an internally consistent
// substitute payload (files + regenerated manifest) fails, because the
// manifest itself is pinned by the sha256 baked into the script (#124: "the
// script itself is versioned/pinned and checksummed like any writer artifact").
func TestInstallScriptTamperedManifestFailsClosed(t *testing.T) {
	payload := t.TempDir()
	copyPayloadTo(t, payload)
	writeFileT(t, filepath.Join(payload, "invariants.md"), payloadFile(t, "invariants.md")+"tampered\n")
	// Regenerate the manifest so it matches the tampered file set.
	var manifest strings.Builder
	for _, line := range strings.Split(strings.TrimSuffix(payloadFile(t, "checksums"), "\n"), "\n") {
		name := line[strings.LastIndex(line, "  ")+2:]
		content := readFileT(t, filepath.Join(payload, name))
		fmt.Fprintf(&manifest, "%x  %s\n", sha256.Sum256([]byte(content)), name)
	}
	writeFileT(t, filepath.Join(payload, "checksums"), manifest.String())

	dir := t.TempDir()
	writeFileT(t, filepath.Join(dir, "CLAUDE.md"), "# My project\n")

	res := runInstall(t, dir, payload, "--posture", "b")
	if res.code == 0 {
		t.Fatal("expected failure on a manifest that does not match the baked pin")
	}
	if !strings.Contains(res.stderr, "pin") {
		t.Errorf("failure must name the pin; stderr:\n%s", res.stderr)
	}
	if _, err := os.Stat(filepath.Join(dir, ".trellis")); !os.IsNotExist(err) {
		t.Errorf("nothing may be installed on a pin mismatch, but .trellis exists")
	}
}

// --- ambiguity + the #112 backstop --------------------------------------------

// TestInstallScriptAmbiguousTargetFailsWithoutTTY (#124: both CLAUDE.md and
// AGENTS.md exist with no block → the interactive path would prompt; the
// non-interactive path must fail loudly and point at --target, never guess).
func TestInstallScriptAmbiguousTargetFailsWithoutTTY(t *testing.T) {
	dir := t.TempDir()
	writeFileT(t, filepath.Join(dir, "CLAUDE.md"), "# Claude\n")
	writeFileT(t, filepath.Join(dir, "AGENTS.md"), "# Agents\n")

	res := runInstall(t, dir, vendoredPayloadAbs(t), "--posture", "b")
	if res.code == 0 {
		t.Fatal("expected failure when the target is ambiguous and no TTY is available")
	}
	if !strings.Contains(res.stderr, "--target") {
		t.Errorf("failure should point at --target; stderr:\n%s", res.stderr)
	}
	// The explicit flag resolves the ambiguity.
	res = runInstall(t, dir, vendoredPayloadAbs(t), "--posture", "b", "--target", "AGENTS.md")
	if res.code != 0 {
		t.Fatalf("--target should resolve the ambiguity, got exit %d: %s", res.code, res.stderr)
	}
	assertOverlayConforms(t, dir, "b", "AGENTS.md", "block-inline-b.md")
	if got := readFileT(t, filepath.Join(dir, "CLAUDE.md")); got != "# Claude\n" {
		t.Errorf("CLAUDE.md must stay untouched when AGENTS.md was chosen")
	}
}

// TestInstallScriptHandAuthoredProfileFailsLoudly (kodhama/trellis#112, skill
// step 2's mechanical half: content after profile.md's generated sentinel is
// hand-authored — never silently overwritten).
func TestInstallScriptHandAuthoredProfileFailsLoudly(t *testing.T) {
	dir := t.TempDir()
	writeFileT(t, filepath.Join(dir, "CLAUDE.md"), "# My project\n")
	if res := runInstall(t, dir, vendoredPayloadAbs(t), "--posture", "b"); res.code != 0 {
		t.Fatalf("first run failed: %s", res.stderr)
	}
	profilePath := filepath.Join(dir, ".trellis/profile.md")
	edited := readFileT(t, profilePath) + "\n## Our own additions\n\nA hand-written rule.\n"
	writeFileT(t, profilePath, edited)

	res := runInstall(t, dir, vendoredPayloadAbs(t))
	if res.code == 0 {
		t.Fatal("expected failure on hand-authored content in profile.md")
	}
	if !strings.Contains(res.stderr, "expression.md") {
		t.Errorf("failure should point at expression.md as the hand-owned home; stderr:\n%s", res.stderr)
	}
	if got := readFileT(t, profilePath); got != edited {
		t.Errorf("the hand-edited profile.md must be left untouched as evidence")
	}
}

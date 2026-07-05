package main

import (
	"os"
	"path/filepath"
)

// Harness is a detected agent harness the setup CLI can compose Trellis onto.
type Harness struct {
	Name   string // human name, e.g. "Claude Code"
	Detail string // what signalled it, e.g. "CLAUDE.md"
}

// detectHarness looks in dir for a supported agent harness. v0 supports only
// Claude Code (spec-0003 §2b), signalled by a .claude/ directory or a CLAUDE.md
// file in the project. It returns ok=false when none is found, so the caller can
// exit cleanly rather than guess — the CLI offers no choice when there is one
// harness, and refuses to proceed when there is none.
//
// The `claude` binary on PATH is intentionally *not* used as a signal: it is a
// global fact, not evidence that *this project* uses Claude Code.
func detectHarness(dir string) (Harness, bool) {
	if fi, err := os.Stat(filepath.Join(dir, ".claude")); err == nil && fi.IsDir() {
		return Harness{Name: "Claude Code", Detail: ".claude/ directory"}, true
	}
	if fi, err := os.Stat(filepath.Join(dir, "CLAUDE.md")); err == nil && !fi.IsDir() {
		return Harness{Name: "Claude Code", Detail: "CLAUDE.md"}, true
	}
	return Harness{}, false
}

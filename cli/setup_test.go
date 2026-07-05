package main

import (
	"os/exec"
	"strings"
	"testing"
)

func claudePresent(t *testing.T) {
	t.Helper()
	withLookPath(t, func(string) (string, error) { return "/usr/local/bin/claude", nil })
}

func TestSetupNoHarness(t *testing.T) {
	// No harness executable: setup must fail loudly rather than guess (D1).
	withLookPath(t, func(string) (string, error) { return "", exec.ErrNotFound })
	if _, err := run2("", "setup", "--dir", t.TempDir()); err == nil {
		t.Fatal("setup should error when no harness executable is present")
	}
}

func TestSetupWithFlags(t *testing.T) {
	claudePresent(t)
	out, err := run2("", "setup", "--dir", t.TempDir(),
		"--profile", "a", "--mode", "m2", "--model", "high")
	if err != nil {
		t.Fatalf("setup with flags: %v", err)
	}
	for _, want := range []string{"detected harness", "A · conductor", "M2 · rewrite", "high-reasoning", "setup plan"} {
		if !strings.Contains(out, want) {
			t.Errorf("plan missing %q in:\n%s", want, out)
		}
	}
}

func TestSetupDefaultsOnEmptyInput(t *testing.T) {
	claudePresent(t)
	// No flags, empty stdin -> each prompt takes its default (b / m1 / none).
	out, err := run2("", "setup", "--dir", t.TempDir())
	if err != nil {
		t.Fatalf("setup with defaults: %v", err)
	}
	for _, want := range []string{"B · author-adapt", "M1 · alongside", "no model"} {
		if !strings.Contains(out, want) {
			t.Errorf("defaulted plan missing %q in:\n%s", want, out)
		}
	}
}

func TestSetupInteractive(t *testing.T) {
	claudePresent(t)
	// Answer the three prompts in order: profile, mode, model.
	out, err := run2("seed\nm2\nbalanced\n", "setup", "--dir", t.TempDir())
	if err != nil {
		t.Fatalf("interactive setup: %v", err)
	}
	for _, want := range []string{"seed", "M2 · rewrite", "balanced"} {
		if !strings.Contains(out, want) {
			t.Errorf("interactive plan missing %q in:\n%s", want, out)
		}
	}
}

func TestSetupModelSuggestionKeyedToMode(t *testing.T) {
	claudePresent(t)
	// M2 with empty model input -> defaults to the suggested high-reasoning model.
	out, err := run2("a\nm2\n\n", "setup", "--dir", t.TempDir())
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if !strings.Contains(out, "suggested model: high") {
		t.Errorf("expected the model suggestion to be high for m2, got:\n%s", out)
	}
	if !strings.Contains(out, "high-reasoning") {
		t.Errorf("empty model input on m2 should default to high-reasoning, got:\n%s", out)
	}
}

func TestSetupInvalidFlag(t *testing.T) {
	claudePresent(t)
	if _, err := run2("", "setup", "--dir", t.TempDir(), "--profile", "zzz"); err == nil {
		t.Fatal("an invalid --profile should be a loud error")
	}
}

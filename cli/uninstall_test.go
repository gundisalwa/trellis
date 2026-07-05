package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func withExecutablePath(t *testing.T, fn func() (string, error)) {
	t.Helper()
	orig := executablePath
	executablePath = fn
	t.Cleanup(func() { executablePath = orig })
}

func fakeBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "trellis")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	withExecutablePath(t, func() (string, error) { return bin, nil })
	return bin
}

func TestUninstallRemovesBinary(t *testing.T) {
	bin := fakeBinary(t)
	out, err := run2("y\n", "uninstall")
	if err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	if _, err := os.Stat(bin); !os.IsNotExist(err) {
		t.Error("binary should have been removed")
	}
	if !strings.Contains(out, "removed") {
		t.Errorf("expected a removed message, got: %s", out)
	}
}

func TestUninstallCancelsWithoutConfirmation(t *testing.T) {
	bin := fakeBinary(t)
	if _, err := run2("", "uninstall"); err != nil { // empty stdin = no confirmation
		t.Fatalf("uninstall: %v", err)
	}
	if _, err := os.Stat(bin); err != nil {
		t.Error("binary must NOT be removed without confirmation")
	}
}

func TestUninstallYesFlag(t *testing.T) {
	bin := fakeBinary(t)
	if _, err := run2("", "uninstall", "--yes"); err != nil {
		t.Fatalf("uninstall --yes: %v", err)
	}
	if _, err := os.Stat(bin); !os.IsNotExist(err) {
		t.Error("--yes should remove without a prompt")
	}
}

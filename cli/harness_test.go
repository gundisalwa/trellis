package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectHarness(t *testing.T) {
	t.Run("empty dir detects nothing", func(t *testing.T) {
		if h, ok := detectHarness(t.TempDir()); ok {
			t.Errorf("empty dir should detect no harness, got %+v", h)
		}
	})

	t.Run("CLAUDE.md detects Claude Code", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# host"), 0o644); err != nil {
			t.Fatal(err)
		}
		h, ok := detectHarness(dir)
		if !ok || h.Name != "Claude Code" {
			t.Errorf("CLAUDE.md should detect Claude Code, got %+v ok=%v", h, ok)
		}
	})

	t.Run(".claude directory detects Claude Code", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.Mkdir(filepath.Join(dir, ".claude"), 0o755); err != nil {
			t.Fatal(err)
		}
		h, ok := detectHarness(dir)
		if !ok || h.Name != "Claude Code" {
			t.Errorf(".claude/ should detect Claude Code, got %+v ok=%v", h, ok)
		}
	})

	t.Run("a CLAUDE.md that is a directory is not a signal", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.Mkdir(filepath.Join(dir, "CLAUDE.md"), 0o755); err != nil {
			t.Fatal(err)
		}
		if _, ok := detectHarness(dir); ok {
			t.Error("a CLAUDE.md directory should not count as the CLAUDE.md file signal")
		}
	})
}

package main

import (
	"bytes"
	"strings"
	"testing"
)

// run2 invokes run with the given stdin and args, returning captured stdout + error.
func run2(stdin string, args ...string) (string, error) {
	var buf bytes.Buffer
	err := run(strings.NewReader(stdin), &buf, args)
	return buf.String(), err
}

func TestRunVersion(t *testing.T) {
	out, err := run2("", "version")
	if err != nil {
		t.Fatalf("version returned error: %v", err)
	}
	if !strings.Contains(out, "trellis ") {
		t.Errorf("version output = %q, want it to contain %q", out, "trellis ")
	}
}

func TestRunHelpAndNoArgs(t *testing.T) {
	for _, args := range [][]string{nil, {"help"}} {
		out, err := run2("", args...)
		if err != nil {
			t.Fatalf("run(%v) returned error: %v", args, err)
		}
		if !strings.Contains(out, "trellis setup") {
			t.Errorf("run(%v) usage did not mention the setup command: %q", args, out)
		}
	}
}

func TestRunUnknownCommand(t *testing.T) {
	if _, err := run2("", "nope"); err == nil {
		t.Fatal("expected an error for an unknown command, got nil")
	}
}

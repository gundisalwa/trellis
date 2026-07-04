package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	var buf bytes.Buffer
	if err := run(&buf, []string{"version"}); err != nil {
		t.Fatalf("version returned error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "trellis ") {
		t.Errorf("version output = %q, want it to contain %q", got, "trellis ")
	}
}

func TestRunHelpAndNoArgs(t *testing.T) {
	for _, args := range [][]string{nil, {"help"}} {
		var buf bytes.Buffer
		if err := run(&buf, args); err != nil {
			t.Fatalf("run(%v) returned error: %v", args, err)
		}
		if !strings.Contains(buf.String(), "trellis setup") {
			t.Errorf("run(%v) usage did not mention the setup command: %q", args, buf.String())
		}
	}
}

func TestRunUnknownCommand(t *testing.T) {
	var buf bytes.Buffer
	if err := run(&buf, []string{"nope"}); err == nil {
		t.Fatal("expected an error for an unknown command, got nil")
	}
}

func TestSetupStubIsHonest(t *testing.T) {
	// Until the flow is built, setup must fail loudly rather than pretend (D1).
	var buf bytes.Buffer
	if err := run(&buf, []string{"setup"}); err == nil {
		t.Fatal("setup stub should return an error until implemented")
	}
}

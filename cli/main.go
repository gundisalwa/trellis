// Command trellis is the setup CLI for the Trellis governance layer.
//
// It is *setup tooling, not a runtime* (decision-0010): you run it once to detect
// your agent harness, pick an expression profile, and compose Trellis onto the
// project; your agents then follow the resulting instructions with no dependency
// on this binary. See specs/0003 §2b for the interactive flow.
package main

import (
	"fmt"
	"io"
	"os"
)

// version is stamped at release time via -ldflags "-X main.version=...".
var version = "0.0.0-dev"

func main() {
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "trellis: "+err.Error())
		os.Exit(1)
	}
}

// run is the testable entrypoint: it writes user-facing output to w and returns
// a non-nil error on failure (which main turns into a stderr line + exit 1).
func run(w io.Writer, args []string) error {
	if len(args) == 0 {
		usage(w)
		return nil
	}
	switch args[0] {
	case "version", "-v", "--version":
		fmt.Fprintln(w, "trellis "+version)
		return nil
	case "help", "-h", "--help":
		usage(w)
		return nil
	case "setup":
		return setup(w, args[1:])
	default:
		return fmt.Errorf("unknown command %q (try `trellis help`)", args[0])
	}
}

func usage(w io.Writer) {
	fmt.Fprintln(w, `trellis — setup CLI for the Trellis governance layer

usage:
  trellis setup      interactive setup: detect harness, pick a profile, choose install mode
  trellis version    print the version
  trellis help       show this message`)
}

// setup will run the interactive setup flow (spec-0003 §2b): detect the harness
// (Claude-only in v0; exit if absent), pick an expression profile, then choose an
// install mode (M1 alongside / M2 rewrite-on-a-branch). The flow lands in the next
// stacked PR; this stub keeps the command surface honest until then.
func setup(_ io.Writer, _ []string) error {
	return fmt.Errorf("setup is not implemented yet — see specs/0003 §2b (next stacked PR)")
}

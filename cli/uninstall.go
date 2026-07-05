package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// executablePath is indirected so tests can point uninstall at a temp file instead
// of the real test binary.
var executablePath = os.Executable

// uninstall removes the trellis binary itself (spec-0004 §1). It does not touch any
// project's .trellis/ — that is `remove`'s job. On Unix, unlinking a running binary
// is safe (the inode frees when the process exits).
func uninstall(in io.Reader, w io.Writer, args []string) error {
	fs := flag.NewFlagSet("uninstall", flag.ContinueOnError)
	fs.SetOutput(w)
	yes := fs.Bool("yes", false, "skip the confirmation")
	if err := fs.Parse(args); err != nil {
		return err
	}

	path, err := executablePath()
	if err != nil {
		return fmt.Errorf("cannot locate the trellis binary: %w", err)
	}

	fmt.Fprintf(w, "This removes the trellis binary at %s.\n", path)
	fmt.Fprintln(w, "(It does not touch any project — use `trellis remove` inside a project for that.)")
	if !*yes && !askYesNo(in, w, "Remove it?") {
		fmt.Fprintln(w, "cancelled — nothing removed")
		return nil
	}

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(w, "already gone — nothing to remove")
			return nil
		}
		return fmt.Errorf("removing %s: %w", path, err)
	}
	fmt.Fprintf(w, "removed %s\n", path)
	return nil
}

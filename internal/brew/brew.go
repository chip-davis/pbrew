package brew

import (
	"fmt"
	"os/exec"
	"strings"
)

// Run executes `brew <args...>` and returns captured stdout.
// Use this when you want to parse the output (e.g. JSON commands).
//
// TODO: implement with os/exec. Return (stdout, error). On non-zero exit,
// surface stderr in the error.
func Run(args ...string) ([]byte, error) {
	cmd := exec.Command("brew", args...)

	output, err := cmd.Output()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("brew %s: %s", strings.Join(args, " "), exitErr.Stderr)
		}
		return nil, fmt.Errorf("brew %s: %w", strings.Join(args, " "), err)
	}
	return output, nil
}

// Passthrough executes `brew <args...>` and streams stdout/stderr straight
// to the user's terminal. Used for the fallthrough case where we don't
// want to prettify (e.g. `pbrew install foo`).
//
// TODO: implement with os/exec, wiring cmd.Stdin/Stdout/Stderr to os.* directly.
// Preserve brew's exit code.
func Passthrough(args []string) error {
	return nil
}

// Stream executes `brew <args...>` and invokes onLine for each stdout line
// as it arrives. Used by update/upgrade to feed the spinner UI.
//
// TODO: implement with exec.Command + StdoutPipe + bufio.Scanner.
func Stream(args []string, onLine func(string)) error {
	return nil
}

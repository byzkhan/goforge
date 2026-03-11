package cli

import "os/exec"

// lookPath wraps exec.LookPath — extracted for testability.
func lookPath(name string) (string, error) {
	return exec.LookPath(name)
}

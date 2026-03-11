// Package postgen runs post-generation commands (go mod tidy, git init, etc.).
package postgen

import (
	"fmt"
	"io"
	"os/exec"
)

// Hook is a post-generation command to run in the project directory.
type Hook struct {
	Name string
	Cmd  string
	Args []string
}

// DefaultHooks returns the standard post-generation hooks.
func DefaultHooks(goCmd string) []Hook {
	return []Hook{
		{Name: "go mod tidy", Cmd: goCmd, Args: []string{"mod", "tidy"}},
		{Name: "git init", Cmd: "git", Args: []string{"init"}},
		{Name: "git add", Cmd: "git", Args: []string{"add", "."}},
	}
}

// Run executes all hooks sequentially in the given directory.
func Run(hooks []Hook, dir string, out io.Writer) error {
	for _, h := range hooks {
		fmt.Fprintf(out, "  Running %s...\n", h.Name)
		cmd := exec.Command(h.Cmd, h.Args...)
		cmd.Dir = dir
		cmd.Stdout = out
		cmd.Stderr = out
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s: %w", h.Name, err)
		}
	}
	return nil
}

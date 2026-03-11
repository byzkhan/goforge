package cli

import (
	"fmt"
	"os/exec"
	"strings"
)

type check struct {
	name string
	cmd  string
	args []string
}

func runDoctor() {
	checks := []check{
		{"Go", "go", []string{"version"}},
		{"Git", "git", []string{"--version"}},
		{"golangci-lint", "golangci-lint", []string{"--version"}},
		{"Docker", "docker", []string{"--version"}},
	}

	fmt.Print("\n  goforge doctor\n\n")

	allOK := true
	for _, c := range checks {
		out, err := exec.Command(c.cmd, c.args...).CombinedOutput()
		if err != nil {
			fmt.Printf("  ✗ %s: not found\n", c.name)
			allOK = false
		} else {
			ver := strings.TrimSpace(string(out))
			fmt.Printf("  ✓ %s: %s\n", c.name, ver)
		}
	}

	fmt.Println()
	if allOK {
		fmt.Println("  All checks passed!")
	} else {
		fmt.Println("  Some tools are missing. Install them for the best experience.")
	}
	fmt.Println()
}

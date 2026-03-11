// Package cli implements the goforge command-line interface.
// It uses only the standard library for argument parsing — no cobra, no urfave.
package cli

import (
	"fmt"
	"os"

	"github.com/byzkhan/goforge/internal/version"
)

const usage = `goforge — scaffold production-ready Go services

Usage:
  goforge new service <name> [flags]
  goforge version
  goforge doctor

Commands:
  new service    Generate a new Go HTTP service
  version        Print goforge version
  doctor         Check your development environment

Run 'goforge new service --help' for flag details.
`

// Execute is the CLI entrypoint.
func Execute() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "version":
		fmt.Println(version.Full())

	case "doctor":
		runDoctor()

	case "new":
		if len(os.Args) < 3 || os.Args[2] != "service" {
			fmt.Fprintln(os.Stderr, "Usage: goforge new service <name> [flags]")
			os.Exit(1)
		}
		runNewService(os.Args[3:])

	case "help", "--help", "-h":
		fmt.Print(usage)

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n%s", os.Args[1], usage)
		os.Exit(1)
	}
}

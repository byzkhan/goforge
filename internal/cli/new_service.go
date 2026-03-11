package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/byzkhan/goforge/internal/config"
	"github.com/byzkhan/goforge/internal/generator"
	"github.com/byzkhan/goforge/internal/postgen"
	"github.com/byzkhan/goforge/internal/prompt"
)

func runNewService(args []string) {
	// Separate the service name (first non-flag arg) from flags, because Go's
	// flag package stops parsing at the first non-flag argument.
	name, flagArgs := extractName(args)

	fs := flag.NewFlagSet("new service", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: goforge new service <name> [flags]

Generates a production-ready Go HTTP service.

Flags:
`)
		fs.PrintDefaults()
	}

	module := fs.String("module", "", "Go module path (default: github.com/yourorg/<name>)")
	router := fs.String("router", "", "HTTP router: chi, stdlib (default: chi)")
	db := fs.String("db", "", "Database: none, postgres (default: none)")
	obs := fs.String("observability", "", "Observability: none, otel (default: none)")
	auth := fs.String("auth", "", "Auth: none, jwt-stub (default: none)")
	ci := fs.String("ci", "", "CI: github, none (default: github)")
	docker := fs.Bool("docker", true, "Generate Dockerfile")
	nonInteractive := fs.Bool("non-interactive", false, "Skip interactive prompts, use defaults")

	if err := fs.Parse(flagArgs); err != nil {
		os.Exit(1)
	}

	if name == "" {
		fmt.Fprintln(os.Stderr, "Error: service name is required")
		fs.Usage()
		os.Exit(1)
	}

	opts := config.ServiceOptions{
		Name:   name,
		Docker: *docker,
	}

	// Set values from flags (leave empty for interactive prompting).
	if *module != "" {
		opts.Module = *module
	}
	if *router != "" {
		opts.Router = config.Router(*router)
	}
	if *db != "" {
		opts.DB = config.DB(*db)
	}
	if *obs != "" {
		opts.Observability = config.Observability(*obs)
	}
	if *auth != "" {
		opts.Auth = config.Auth(*auth)
	}
	if *ci != "" {
		opts.CI = config.CI(*ci)
	}

	// Fill defaults or prompt interactively.
	fmt.Printf("\n  Creating service: %s\n\n", name)
	if *nonInteractive {
		opts = applyDefaults(opts)
	} else {
		p := prompt.New(os.Stdin, os.Stdout)
		var err error
		opts, err = p.Fill(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "prompt error: %v\n", err)
			os.Exit(1)
		}
		opts = applyDefaults(opts)
	}

	// Set output directory.
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "getwd: %v\n", err)
		os.Exit(1)
	}
	opts.OutputDir = filepath.Join(cwd, name)

	// Validate.
	if err := opts.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Check target directory doesn't exist.
	if _, err := os.Stat(opts.OutputDir); err == nil {
		fmt.Fprintf(os.Stderr, "Error: directory %q already exists\n", opts.OutputDir)
		os.Exit(1)
	} else if !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: cannot stat %q: %v\n", opts.OutputDir, err)
		os.Exit(1)
	}

	// Generate.
	fmt.Printf("  Generating %s...\n\n", name)
	gen := generator.NewServiceHTTPGenerator()
	result, err := gen.Render(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := gen.WriteToDisk(result, opts.OutputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("  Created %d files\n\n", len(result.Files))

	// Post-generation hooks.
	goCmd := findGo()
	hooks := postgen.DefaultHooks(goCmd)
	if err := postgen.Run(hooks, opts.OutputDir, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: post-generation hook failed: %v\n", err)
		fmt.Fprintln(os.Stderr, "You may need to run 'go mod tidy' manually.")
	}

	fmt.Printf("\n  Done! Service %s is ready.\n\n", name)
	fmt.Printf("  cd %s && make run\n\n", name)
}

// extractName separates the service name from flag arguments.
// The name is the first argument that doesn't start with "-".
func extractName(args []string) (string, []string) {
	var name string
	var flags []string
	for _, a := range args {
		if name == "" && !strings.HasPrefix(a, "-") {
			name = a
		} else {
			flags = append(flags, a)
		}
	}
	return name, flags
}

// applyDefaults fills in any zero-value fields with defaults.
func applyDefaults(opts config.ServiceOptions) config.ServiceOptions {
	defaults := config.Defaults()
	if opts.Module == "" {
		opts.Module = "github.com/yourorg/" + opts.Name
	}
	if opts.Router == "" {
		opts.Router = defaults.Router
	}
	if opts.DB == "" {
		opts.DB = defaults.DB
	}
	if opts.Observability == "" {
		opts.Observability = defaults.Observability
	}
	if opts.Auth == "" {
		opts.Auth = defaults.Auth
	}
	if opts.CI == "" {
		opts.CI = defaults.CI
	}
	return opts
}

// findGo locates the go binary. It checks common paths if not in PATH.
func findGo() string {
	if path, err := lookPath("go"); err == nil {
		return path
	}
	for _, p := range []string{
		"/usr/local/go/bin/go",
		"/opt/homebrew/bin/go",
		os.ExpandEnv("$HOME/go-sdk/go/bin/go"),
		os.ExpandEnv("$HOME/go/bin/go"),
	} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "go"
}

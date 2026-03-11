package generator

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/byzkhan/goforge/internal/config"
)

// TestE2E_GeneratedServiceCompiles generates a full service to a temp dir
// and runs `go build` to verify it compiles. This is a snapshot-style test.
func TestE2E_GeneratedServiceCompiles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	goCmd := findGoBinary()
	if goCmd == "" {
		t.Skip("go binary not found")
	}

	tests := []struct {
		name string
		opts config.ServiceOptions
	}{
		{
			name: "defaults_chi",
			opts: func() config.ServiceOptions {
				o := config.Defaults()
				o.Name = "testchi"
				o.Module = "github.com/test/testchi"
				return o
			}(),
		},
		{
			name: "stdlib_router",
			opts: func() config.ServiceOptions {
				o := config.Defaults()
				o.Name = "teststdlib"
				o.Module = "github.com/test/teststdlib"
				o.Router = config.RouterStdlib
				return o
			}(),
		},
		{
			name: "all_features",
			opts: func() config.ServiceOptions {
				o := config.Defaults()
				o.Name = "testfull"
				o.Module = "github.com/test/testfull"
				o.DB = config.DBPostgres
				o.Observability = config.ObsOtel
				o.Auth = config.AuthJWTStub
				return o
			}(),
		},
		{
			name: "minimal",
			opts: func() config.ServiceOptions {
				o := config.Defaults()
				o.Name = "testmin"
				o.Module = "github.com/test/testmin"
				o.Docker = false
				o.CI = config.CINone
				return o
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			projectDir := filepath.Join(dir, tt.opts.Name)

			tt.opts.OutputDir = projectDir

			gen := NewServiceHTTPGenerator()
			result, err := gen.Render(tt.opts)
			if err != nil {
				t.Fatalf("render: %v", err)
			}

			if err := gen.WriteToDisk(result, projectDir); err != nil {
				t.Fatalf("write: %v", err)
			}

			// Run go mod tidy.
			tidy := exec.Command(goCmd, "mod", "tidy")
			tidy.Dir = projectDir
			if out, err := tidy.CombinedOutput(); err != nil {
				t.Fatalf("go mod tidy failed: %v\n%s", err, out)
			}

			// Run go build.
			build := exec.Command(goCmd, "build", "./...")
			build.Dir = projectDir
			if out, err := build.CombinedOutput(); err != nil {
				t.Fatalf("go build failed: %v\n%s", err, out)
			}

			// Run go vet.
			vet := exec.Command(goCmd, "vet", "./...")
			vet.Dir = projectDir
			if out, err := vet.CombinedOutput(); err != nil {
				t.Fatalf("go vet failed: %v\n%s", err, out)
			}
		})
	}
}

func findGoBinary() string {
	if p, err := exec.LookPath("go"); err == nil {
		return p
	}
	for _, p := range []string{
		"/usr/local/go/bin/go",
		"/opt/homebrew/bin/go",
		os.ExpandEnv("$HOME/go-sdk/go/bin/go"),
	} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

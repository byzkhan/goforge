package generator

import "github.com/byzkhan/goforge/internal/config"

// FileEntry describes a single file to generate.
type FileEntry struct {
	// TemplatePath is the path within the embedded templates dir (e.g. "cmd/{{.Name}}/main.go.tmpl").
	TemplatePath string
	// OutputPath is the destination path relative to the project root, supporting {{.Name}} etc.
	OutputPath string
	// Condition determines whether this file should be generated.
	Condition func(config.ServiceOptions) bool
}

// Always returns a condition that is always true.
func Always(_ config.ServiceOptions) bool { return true }

// WhenPostgres returns true when postgres is enabled.
func WhenPostgres(o config.ServiceOptions) bool { return o.DB == config.DBPostgres }

// WhenOtel returns true when OpenTelemetry is enabled.
func WhenOtel(o config.ServiceOptions) bool { return o.Observability == config.ObsOtel }

// WhenJWTAuth returns true when JWT auth stub is enabled.
func WhenJWTAuth(o config.ServiceOptions) bool { return o.Auth == config.AuthJWTStub }

// WhenDocker returns true when Docker is enabled.
func WhenDocker(o config.ServiceOptions) bool { return o.Docker }

// WhenGitHubCI returns true when GitHub Actions CI is enabled.
func WhenGitHubCI(o config.ServiceOptions) bool { return o.CI == config.CIGitHub }

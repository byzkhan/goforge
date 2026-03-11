package generator

import (
	"strings"
	"testing"

	"github.com/zaid/goforge/internal/config"
)

func defaultOpts() config.ServiceOptions {
	o := config.Defaults()
	o.Name = "testapp"
	o.Module = "github.com/org/testapp"
	o.OutputDir = "/tmp/testapp"
	return o
}

func TestRender_DefaultConfig(t *testing.T) {
	gen := NewServiceHTTPGenerator()
	opts := defaultOpts()

	result, err := gen.Render(opts)
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	// Check expected files exist.
	expectedFiles := []string{
		"go.mod",
		"cmd/testapp/main.go",
		"internal/app/app.go",
		"internal/config/config.go",
		"internal/http/handlers/health.go",
		"internal/http/handlers/example.go",
		"internal/http/middleware/requestid.go",
		"internal/http/middleware/logging.go",
		"internal/http/middleware/recoverer.go",
		"internal/http/routes/routes.go",
		"internal/http/response.go",
		"internal/platform/logger/logger.go",
		"internal/platform/server/server.go",
		"internal/platform/shutdown/shutdown.go",
		"internal/domain/example.go",
		"internal/service/example.go",
		"internal/repository/example.go",
		"tests/health_test.go",
		"tests/integration_test.go",
		"Makefile",
		"Dockerfile",
		".github/workflows/ci.yml",
		".env.example",
		".gitignore",
		".golangci.yml",
		"README.md",
	}

	for _, f := range expectedFiles {
		if _, ok := result.Files[f]; !ok {
			t.Errorf("missing expected file: %s", f)
		}
	}

	// Check files that should NOT exist with default config.
	unexpectedFiles := []string{
		"internal/platform/db/db.go",
		"internal/platform/telemetry/telemetry.go",
		"internal/http/middleware/auth.go",
		"internal/repository/example_postgres.go",
		"migrations/001_initial.up.sql",
		"migrations/001_initial.down.sql",
		"scripts/migrate.sh",
	}

	for _, f := range unexpectedFiles {
		if _, ok := result.Files[f]; ok {
			t.Errorf("unexpected file in default config: %s", f)
		}
	}
}

func TestRender_WithPostgres(t *testing.T) {
	gen := NewServiceHTTPGenerator()
	opts := defaultOpts()
	opts.DB = config.DBPostgres

	result, err := gen.Render(opts)
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	postgresFiles := []string{
		"internal/platform/db/db.go",
		"internal/repository/example_postgres.go",
		"migrations/001_initial.up.sql",
		"migrations/001_initial.down.sql",
		"scripts/migrate.sh",
	}

	for _, f := range postgresFiles {
		if _, ok := result.Files[f]; !ok {
			t.Errorf("missing postgres file: %s", f)
		}
	}

	// Config should include DATABASE_URL.
	configContent := string(result.Files["internal/config/config.go"])
	if !strings.Contains(configContent, "DatabaseURL") {
		t.Error("config should contain DatabaseURL when postgres is enabled")
	}
}

func TestRender_WithOtel(t *testing.T) {
	gen := NewServiceHTTPGenerator()
	opts := defaultOpts()
	opts.Observability = config.ObsOtel

	result, err := gen.Render(opts)
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if _, ok := result.Files["internal/platform/telemetry/telemetry.go"]; !ok {
		t.Error("missing telemetry file")
	}

	appContent := string(result.Files["internal/app/app.go"])
	if !strings.Contains(appContent, "telemetry") {
		t.Error("app.go should reference telemetry when otel is enabled")
	}
}

func TestRender_WithJWT(t *testing.T) {
	gen := NewServiceHTTPGenerator()
	opts := defaultOpts()
	opts.Auth = config.AuthJWTStub

	result, err := gen.Render(opts)
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if _, ok := result.Files["internal/http/middleware/auth.go"]; !ok {
		t.Error("missing auth middleware file")
	}
}

func TestRender_NoDocker(t *testing.T) {
	gen := NewServiceHTTPGenerator()
	opts := defaultOpts()
	opts.Docker = false

	result, err := gen.Render(opts)
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if _, ok := result.Files["Dockerfile"]; ok {
		t.Error("Dockerfile should not be generated when docker is disabled")
	}
}

func TestRender_NoCI(t *testing.T) {
	gen := NewServiceHTTPGenerator()
	opts := defaultOpts()
	opts.CI = config.CINone

	result, err := gen.Render(opts)
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if _, ok := result.Files[".github/workflows/ci.yml"]; ok {
		t.Error("CI file should not be generated when ci is none")
	}
}

func TestRender_StdlibRouter(t *testing.T) {
	gen := NewServiceHTTPGenerator()
	opts := defaultOpts()
	opts.Router = config.RouterStdlib

	result, err := gen.Render(opts)
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	gomod := string(result.Files["go.mod"])
	if strings.Contains(gomod, "chi") {
		t.Error("go.mod should not contain chi when using stdlib router")
	}

	routes := string(result.Files["internal/http/routes/routes.go"])
	if strings.Contains(routes, "chi.NewRouter") {
		t.Error("routes should use stdlib mux, not chi")
	}
	if !strings.Contains(routes, "http.NewServeMux") {
		t.Error("routes should use http.NewServeMux for stdlib router")
	}
}

func TestRender_ModuleSubstitution(t *testing.T) {
	gen := NewServiceHTTPGenerator()
	opts := defaultOpts()
	opts.Module = "github.com/mycompany/myservice"

	result, err := gen.Render(opts)
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	mainGo := string(result.Files["cmd/testapp/main.go"])
	if !strings.Contains(mainGo, "github.com/mycompany/myservice") {
		t.Error("main.go should contain the custom module path")
	}
}

func TestRender_InvalidOptions(t *testing.T) {
	gen := NewServiceHTTPGenerator()
	opts := config.ServiceOptions{} // empty, invalid

	_, err := gen.Render(opts)
	if err == nil {
		t.Fatal("expected error for invalid options")
	}
}

func TestRender_AllFeatures(t *testing.T) {
	gen := NewServiceHTTPGenerator()
	opts := defaultOpts()
	opts.DB = config.DBPostgres
	opts.Observability = config.ObsOtel
	opts.Auth = config.AuthJWTStub

	result, err := gen.Render(opts)
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	// All feature files should exist.
	featureFiles := []string{
		"internal/platform/db/db.go",
		"internal/platform/telemetry/telemetry.go",
		"internal/http/middleware/auth.go",
		"internal/repository/example_postgres.go",
		"migrations/001_initial.up.sql",
	}

	for _, f := range featureFiles {
		if _, ok := result.Files[f]; !ok {
			t.Errorf("missing feature file: %s", f)
		}
	}
}

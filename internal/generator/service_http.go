package generator

import (
	"embed"
	"io/fs"
)

//go:embed templates/service-http/*
var serviceHTTPFS embed.FS

// ServiceHTTPManifest returns the file manifest for the service-http template.
func ServiceHTTPManifest() []FileEntry {
	return []FileEntry{
		// Root files
		{TemplatePath: "templates/service-http/go.mod.tmpl", OutputPath: "go.mod", Condition: Always},
		{TemplatePath: "templates/service-http/README.md.tmpl", OutputPath: "README.md", Condition: Always},
		{TemplatePath: "templates/service-http/Makefile.tmpl", OutputPath: "Makefile", Condition: Always},
		{TemplatePath: "templates/service-http/dot-env.example.tmpl", OutputPath: ".env.example", Condition: Always},
		{TemplatePath: "templates/service-http/dot-gitignore.tmpl", OutputPath: ".gitignore", Condition: Always},
		{TemplatePath: "templates/service-http/dot-golangci.yml.tmpl", OutputPath: ".golangci.yml", Condition: Always},
		{TemplatePath: "templates/service-http/Dockerfile.tmpl", OutputPath: "Dockerfile", Condition: WhenDocker},
		{TemplatePath: "templates/service-http/ci.yml.tmpl", OutputPath: ".github/workflows/ci.yml", Condition: WhenGitHubCI},

		// cmd
		{TemplatePath: "templates/service-http/cmd_main.go.tmpl", OutputPath: "cmd/{{.Name}}/main.go", Condition: Always},

		// internal/app
		{TemplatePath: "templates/service-http/internal_app.go.tmpl", OutputPath: "internal/app/app.go", Condition: Always},

		// internal/config
		{TemplatePath: "templates/service-http/internal_config.go.tmpl", OutputPath: "internal/config/config.go", Condition: Always},

		// internal/http/handlers
		{TemplatePath: "templates/service-http/internal_handlers_health.go.tmpl", OutputPath: "internal/http/handlers/health.go", Condition: Always},
		{TemplatePath: "templates/service-http/internal_handlers_example.go.tmpl", OutputPath: "internal/http/handlers/example.go", Condition: Always},

		// internal/http/middleware
		{TemplatePath: "templates/service-http/internal_middleware_requestid.go.tmpl", OutputPath: "internal/http/middleware/requestid.go", Condition: Always},
		{TemplatePath: "templates/service-http/internal_middleware_logging.go.tmpl", OutputPath: "internal/http/middleware/logging.go", Condition: Always},
		{TemplatePath: "templates/service-http/internal_middleware_recoverer.go.tmpl", OutputPath: "internal/http/middleware/recoverer.go", Condition: Always},
		{TemplatePath: "templates/service-http/internal_middleware_auth.go.tmpl", OutputPath: "internal/http/middleware/auth.go", Condition: WhenJWTAuth},

		// internal/http/routes
		{TemplatePath: "templates/service-http/internal_routes.go.tmpl", OutputPath: "internal/http/routes/routes.go", Condition: Always},

		// internal/http (response helpers)
		{TemplatePath: "templates/service-http/internal_http_response.go.tmpl", OutputPath: "internal/http/response.go", Condition: Always},

		// internal/platform/logger
		{TemplatePath: "templates/service-http/internal_platform_logger.go.tmpl", OutputPath: "internal/platform/logger/logger.go", Condition: Always},

		// internal/platform/server
		{TemplatePath: "templates/service-http/internal_platform_server.go.tmpl", OutputPath: "internal/platform/server/server.go", Condition: Always},

		// internal/platform/shutdown
		{TemplatePath: "templates/service-http/internal_platform_shutdown.go.tmpl", OutputPath: "internal/platform/shutdown/shutdown.go", Condition: Always},

		// internal/platform/db
		{TemplatePath: "templates/service-http/internal_platform_db.go.tmpl", OutputPath: "internal/platform/db/db.go", Condition: WhenPostgres},

		// internal/platform/telemetry
		{TemplatePath: "templates/service-http/internal_platform_telemetry.go.tmpl", OutputPath: "internal/platform/telemetry/telemetry.go", Condition: WhenOtel},

		// internal/domain
		{TemplatePath: "templates/service-http/internal_domain.go.tmpl", OutputPath: "internal/domain/example.go", Condition: Always},

		// internal/service
		{TemplatePath: "templates/service-http/internal_service.go.tmpl", OutputPath: "internal/service/example.go", Condition: Always},

		// internal/repository
		{TemplatePath: "templates/service-http/internal_repository.go.tmpl", OutputPath: "internal/repository/example.go", Condition: Always},
		{TemplatePath: "templates/service-http/internal_repository_postgres.go.tmpl", OutputPath: "internal/repository/example_postgres.go", Condition: WhenPostgres},

		// tests
		{TemplatePath: "templates/service-http/tests_handler_test.go.tmpl", OutputPath: "tests/health_test.go", Condition: Always},
		{TemplatePath: "templates/service-http/tests_integration_test.go.tmpl", OutputPath: "tests/integration_test.go", Condition: Always},

		// migrations
		{TemplatePath: "templates/service-http/migrations_initial.sql.tmpl", OutputPath: "migrations/001_initial.up.sql", Condition: WhenPostgres},
		{TemplatePath: "templates/service-http/migrations_initial_down.sql.tmpl", OutputPath: "migrations/001_initial.down.sql", Condition: WhenPostgres},

		// scripts
		{TemplatePath: "templates/service-http/scripts_migrate.sh.tmpl", OutputPath: "scripts/migrate.sh", Condition: WhenPostgres},
	}
}

// NewServiceHTTPGenerator creates a generator for the service-http template.
func NewServiceHTTPGenerator() *Generator {
	fsys, _ := fs.Sub(serviceHTTPFS, ".")
	return New(fsys, ServiceHTTPManifest())
}

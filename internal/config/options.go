// Package config defines the options model used by the generator.
// Every generator decision is driven by this struct — no hidden globals.
package config

import (
	"fmt"
	"regexp"
)

// ServiceOptions holds all settings for generating a new service.
type ServiceOptions struct {
	Name          string // Project directory & binary name (e.g. "myapi")
	Module        string // Go module path (e.g. "github.com/org/myapi")
	Router        Router
	DB            DB
	Observability Observability
	Auth          Auth
	CI            CI
	Docker        bool
	OutputDir     string // Absolute path where the project will be written
}

// Validate checks the options for consistency.
func (o ServiceOptions) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("service name is required")
	}
	if !isValidName(o.Name) {
		return fmt.Errorf("service name %q must be lowercase alphanumeric with hyphens", o.Name)
	}
	if o.Module == "" {
		return fmt.Errorf("module path is required")
	}
	if !o.Router.Valid() {
		return fmt.Errorf("unsupported router: %q", o.Router)
	}
	if !o.DB.Valid() {
		return fmt.Errorf("unsupported db: %q", o.DB)
	}
	if !o.Observability.Valid() {
		return fmt.Errorf("unsupported observability: %q", o.Observability)
	}
	if !o.Auth.Valid() {
		return fmt.Errorf("unsupported auth: %q", o.Auth)
	}
	if !o.CI.Valid() {
		return fmt.Errorf("unsupported ci: %q", o.CI)
	}
	return nil
}

var nameRe = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

func isValidName(s string) bool { return nameRe.MatchString(s) }

// --- Enums ---

type Router string

const (
	RouterChi    Router = "chi"
	RouterStdlib Router = "stdlib"
)

func (r Router) Valid() bool { return r == RouterChi || r == RouterStdlib }

type DB string

const (
	DBNone     DB = "none"
	DBPostgres DB = "postgres"
)

func (d DB) Valid() bool { return d == DBNone || d == DBPostgres }

type Observability string

const (
	ObsNone Observability = "none"
	ObsOtel Observability = "otel"
)

func (o Observability) Valid() bool { return o == ObsNone || o == ObsOtel }

type Auth string

const (
	AuthNone    Auth = "none"
	AuthJWTStub Auth = "jwt-stub"
)

func (a Auth) Valid() bool { return a == AuthNone || a == AuthJWTStub }

type CI string

const (
	CIGitHub CI = "github"
	CINone   CI = "none"
)

func (c CI) Valid() bool { return c == CIGitHub || c == CINone }

// Defaults returns ServiceOptions populated with the recommended defaults.
func Defaults() ServiceOptions {
	return ServiceOptions{
		Router:        RouterChi,
		DB:            DBNone,
		Observability: ObsNone,
		Auth:          AuthNone,
		CI:            CIGitHub,
		Docker:        true,
	}
}

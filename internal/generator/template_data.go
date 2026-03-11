package generator

import "github.com/byzkhan/goforge/internal/config"

// TemplateData is the data passed to every template.
type TemplateData struct {
	Name          string
	Module        string
	Router        string
	DB            string
	Observability string
	Auth          string
	CI            string
	Docker        bool
	// Computed convenience booleans for templates.
	UseChi      bool
	UseStdlib   bool
	UsePostgres bool
	UseOtel     bool
	UseJWT      bool
	UseGitHub   bool
}

// NewTemplateData builds the data from validated options.
func NewTemplateData(o config.ServiceOptions) TemplateData {
	return TemplateData{
		Name:          o.Name,
		Module:        o.Module,
		Router:        string(o.Router),
		DB:            string(o.DB),
		Observability: string(o.Observability),
		Auth:          string(o.Auth),
		CI:            string(o.CI),
		Docker:        o.Docker,
		UseChi:        o.Router == config.RouterChi,
		UseStdlib:     o.Router == config.RouterStdlib,
		UsePostgres:   o.DB == config.DBPostgres,
		UseOtel:       o.Observability == config.ObsOtel,
		UseJWT:        o.Auth == config.AuthJWTStub,
		UseGitHub:     o.CI == config.CIGitHub,
	}
}

// Package prompt handles interactive user input for service configuration.
// It fills in missing ServiceOptions fields by asking the user.
package prompt

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/zaid/goforge/internal/config"
)

// Prompter reads interactive input from a reader (typically os.Stdin).
type Prompter struct {
	in  *bufio.Scanner
	out io.Writer
}

// New creates a Prompter.
func New(in io.Reader, out io.Writer) *Prompter {
	return &Prompter{in: bufio.NewScanner(in), out: out}
}

// Fill prompts the user for any unset options and returns the completed config.
func (p *Prompter) Fill(opts config.ServiceOptions) (config.ServiceOptions, error) {
	if opts.Module == "" {
		opts.Module = p.ask("Go module path", "github.com/yourorg/"+opts.Name)
	}

	if opts.Router == "" {
		val := p.choose("Router", []string{"chi", "stdlib"}, "chi")
		opts.Router = config.Router(val)
	}

	if opts.DB == "" {
		val := p.choose("Database", []string{"none", "postgres"}, "none")
		opts.DB = config.DB(val)
	}

	if opts.Observability == "" {
		val := p.choose("Observability", []string{"none", "otel"}, "none")
		opts.Observability = config.Observability(val)
	}

	if opts.Auth == "" {
		val := p.choose("Auth", []string{"none", "jwt-stub"}, "none")
		opts.Auth = config.Auth(val)
	}

	if opts.CI == "" {
		val := p.choose("CI", []string{"github", "none"}, "github")
		opts.CI = config.CI(val)
	}

	return opts, nil
}

func (p *Prompter) ask(label, defaultVal string) string {
	fmt.Fprintf(p.out, "  %s [%s]: ", label, defaultVal)
	if p.in.Scan() {
		if text := strings.TrimSpace(p.in.Text()); text != "" {
			return text
		}
	}
	return defaultVal
}

func (p *Prompter) choose(label string, choices []string, defaultVal string) string {
	fmt.Fprintf(p.out, "  %s (%s) [%s]: ", label, strings.Join(choices, "/"), defaultVal)
	if p.in.Scan() {
		text := strings.TrimSpace(p.in.Text())
		for _, c := range choices {
			if strings.EqualFold(text, c) {
				return c
			}
		}
	}
	return defaultVal
}

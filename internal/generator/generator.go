// Package generator renders templates to produce a new Go service project.
package generator

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/zaid/goforge/internal/config"
	"github.com/zaid/goforge/pkg/fsutil"
)

// Generator produces project files from a template set.
type Generator struct {
	templateFS fs.FS
	manifest   []FileEntry
	funcMap    template.FuncMap
}

// New creates a generator with the given embedded filesystem and manifest.
func New(templateFS fs.FS, manifest []FileEntry) *Generator {
	return &Generator{
		templateFS: templateFS,
		manifest:   manifest,
		funcMap: template.FuncMap{
			"lower": strings.ToLower,
			"title": strings.Title, //nolint:staticcheck // fine for code gen
		},
	}
}

// Result holds the generated files keyed by their output path.
type Result struct {
	Files map[string][]byte
}

// Render processes all manifest entries and returns the result in memory.
// This makes the generator testable without touching the filesystem.
func (g *Generator) Render(opts config.ServiceOptions) (*Result, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	data := NewTemplateData(opts)
	result := &Result{Files: make(map[string][]byte)}

	for _, entry := range g.manifest {
		if entry.Condition != nil && !entry.Condition(opts) {
			continue
		}

		// Resolve output path (may contain template expressions).
		outPath, err := g.renderString(entry.OutputPath, data)
		if err != nil {
			return nil, fmt.Errorf("render output path %q: %w", entry.OutputPath, err)
		}

		// Read template content.
		content, err := fs.ReadFile(g.templateFS, entry.TemplatePath)
		if err != nil {
			return nil, fmt.Errorf("read template %q: %w", entry.TemplatePath, err)
		}

		// Render template content.
		rendered, err := g.renderBytes(entry.TemplatePath, content, data)
		if err != nil {
			return nil, fmt.Errorf("render template %q: %w", entry.TemplatePath, err)
		}

		result.Files[outPath] = rendered
	}

	return result, nil
}

// WriteToDisk writes all generated files to the given base directory.
func (g *Generator) WriteToDisk(result *Result, baseDir string) error {
	for relPath, content := range result.Files {
		absPath := filepath.Join(baseDir, relPath)
		perm := os.FileMode(0o644)
		// Make scripts executable.
		if strings.HasSuffix(relPath, ".sh") || strings.HasPrefix(filepath.Base(relPath), "Makefile") {
			perm = 0o644
		}
		if err := fsutil.WriteFile(absPath, content, perm); err != nil {
			return fmt.Errorf("write %s: %w", relPath, err)
		}
	}
	return nil
}

func (g *Generator) renderString(tmplStr string, data TemplateData) (string, error) {
	t, err := template.New("path").Funcs(g.funcMap).Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (g *Generator) renderBytes(name string, content []byte, data TemplateData) ([]byte, error) {
	t, err := template.New(name).Funcs(g.funcMap).Parse(string(content))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

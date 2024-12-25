package template

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/friendsofgo/errors"
	"github.com/go-sprout/sprout"
	"github.com/go-sprout/sprout/registry/maps"
	"github.com/go-sprout/sprout/registry/numeric"
	"github.com/go-sprout/sprout/registry/strings"
	"github.com/networkteam/slogutils"
)

type Renderer struct {
	documents []Template
	partials  []Template
	sprout    sprout.Handler
}

var ErrTemplateNotFound = errors.New("template not found")

func NewRenderer(documents, partials []Template) *Renderer {
	logger := slogutils.FromContext(context.Background())
	sproutHandler := sprout.New(
		sprout.WithLogger(logger),
		sprout.WithRegistries(
			strings.NewRegistry(),
			numeric.NewRegistry(),
			maps.NewRegistry(),
		),
	)

	return &Renderer{
		documents: documents,
		partials:  partials,
		sprout:    sproutHandler,
	}
}

func (r *Renderer) Render(name string, data any) (string, error) {
	// Find the document
	var doc *Template
	for _, d := range r.documents {
		if d.Name == name {
			doc = &d
			break
		}
	}
	if doc == nil {
		return "", errors.Wrapf(ErrTemplateNotFound, "template: %s", name)
	}

	// Create template with main content
	tmpl, err := template.New(doc.Name).
		Funcs(customTemplateFuncs()).
		Funcs(r.sprout.Build()).
		Parse(doc.Content)
	if err != nil {
		return "", errors.Wrap(err, "parsing main template")
	}

	// Add all partials
	for _, p := range r.partials {
		_, err := tmpl.New(p.Name).Parse(p.Content)
		if err != nil {
			return "", errors.Wrap(err, "parsing partial template")
		}
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", errors.Wrap(err, "executing template")
	}

	return buf.String(), nil
}

func (r *Renderer) Documents() []Template {
	return r.documents
}

func customTemplateFuncs() template.FuncMap {
	exp := func(expression string) string {
		return fmt.Sprintf("{{ %s }}", expression)
	}
	return template.FuncMap{
		"expression": exp,
		"exp":        exp,
	}
}

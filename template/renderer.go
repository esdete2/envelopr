package template

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/friendsofgo/errors"
)

type Renderer struct {
	documents []Template
	partials  []Template
}

var ErrTemplateNotFound = errors.New("template not found")

func NewRenderer(documents, partials []Template) *Renderer {
	return &Renderer{
		documents: documents,
		partials:  partials,
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
	tmpl, err := template.New(doc.Name).Funcs(template.FuncMap{
		"expression": func(expression string) string {
			return fmt.Sprintf("{{ %s }}", expression)
		},
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}).Parse(doc.Content)
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

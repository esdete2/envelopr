package handler

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/friendsofgo/errors"

	"github.com/esdete2/mjml-dev/config"
	"github.com/esdete2/mjml-dev/template"
)

type Processor struct {
	config   *config.Config
	compiler *template.Compiler
}

func NewProcessor(cfg *config.Config) (*Processor, error) {
	return &Processor{
		config:   cfg,
		compiler: template.NewCompiler(cfg),
	}, nil
}

func (p *Processor) Process() error {
	loader := NewFileLoader(p.config.Paths.Documents, p.config.Paths.Partials)

	documents, err := loader.LoadDocuments()
	if err != nil {
		return &Error{
			Type:    ErrorLoadingFiles,
			Wrapped: errors.Wrap(err, "loading documents"),
		}
	}

	partials, err := loader.LoadPartials()
	if err != nil {
		return &Error{
			Type:    ErrorLoadingFiles,
			Wrapped: errors.Wrap(err, "loading partials"),
		}
	}

	// Create renderer with fresh templates
	renderer := template.NewRenderer(documents, partials)

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(p.config.Paths.Output, 0755); err != nil {
		return &Error{
			Type:    ErrorSaving,
			Wrapped: errors.Wrap(err, "creating output directory"),
		}
	}

	// Process all documents
	for _, doc := range renderer.Documents() {
		if err := p.processDocument(doc, renderer); err != nil {
			return errors.Wrap(err, "processing document")
		}
	}

	return nil
}

func (p *Processor) ProcessSingle(templateName string) error {
	loader := NewFileLoader(p.config.Paths.Documents, p.config.Paths.Partials)

	// Load just the specified document
	documents, err := loader.LoadDocument(templateName)
	if err != nil {
		return &Error{
			Type:    ErrorLoadingFiles,
			Wrapped: errors.Wrap(err, "loading document"),
		}
	}

	if len(documents) == 0 {
		return &Error{
			Type:    ErrorLoadingFiles,
			Doc:     templateName,
			Wrapped: errors.New("template not found"),
		}
	}

	// Always load all partials since they might be used
	partials, err := loader.LoadPartials()
	if err != nil {
		return &Error{
			Type:    ErrorLoadingFiles,
			Wrapped: errors.Wrap(err, "loading partials"),
		}
	}

	// Create renderer with fresh templates
	renderer := template.NewRenderer(documents, partials)

	// Process the single document
	return p.processDocument(documents[0], renderer)
}

func (p *Processor) processDocument(doc template.Template, renderer *template.Renderer) error {
	// Prepare data
	data := make(map[string]any)

	// Add global variables
	for k, v := range p.config.Template.Variables {
		data[k] = v
	}

	// Add document-specific variables
	parts := strings.Split(doc.Name, "/")

	if docVars, exists := p.config.Template.Documents[doc.Name]; exists {
		data = mergeMaps(data, docVars)
	} else if len(parts) > 1 {
		// Handle nested structure (e.g., "shop/invoice")
		if parentVars, exists := p.config.Template.Documents[parts[0]]; exists {
			if parentMap, ok := parentVars.(map[string]any); ok {
				if childVars, exists := parentMap[parts[1]]; exists {
					data = mergeMaps(data, childVars)
				}
			}
		}
	}

	// Render template
	rendered, err := renderer.Render(doc.Name, data)
	if err != nil {
		return &Error{
			Type:    ErrorRendering,
			Doc:     doc.Name,
			Wrapped: errors.Wrap(err, "rendering template"),
		}
	}

	// Compile to HTML
	html, err := p.compiler.Compile(rendered)
	if err != nil {
		return &Error{
			Type:    ErrorCompiling,
			Doc:     doc.Name,
			Wrapped: errors.Wrap(err, "compiling template"),
		}
	}

	// Save to file
	outputPath := filepath.Join(p.config.Paths.Output, doc.Name+".html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return &Error{
			Type:    ErrorSaving,
			Doc:     doc.Name,
			Wrapped: errors.Wrap(err, "creating output directory"),
		}
	}

	if err := os.WriteFile(outputPath, []byte(html), 0600); err != nil {
		return &Error{
			Type:    ErrorSaving,
			Doc:     doc.Name,
			Wrapped: errors.Wrap(err, "writing output file"),
		}
	}

	return nil
}

func mergeMaps(target map[string]any, source any) map[string]any {
	if source == nil {
		return target
	}

	switch v := source.(type) {
	case map[string]any:
		for key, val := range v {
			target[key] = val
		}
	case map[interface{}]interface{}:
		for key, val := range v {
			if strKey, ok := key.(string); ok {
				target[strKey] = val
			}
		}
	}

	return target
}

package handler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/stretchr/testify/require"

	"github.com/esdete2/mjml-dev/config"
	"github.com/esdete2/mjml-dev/handler"
)

func TestProcessor(t *testing.T) {
	r := require.New(t)

	// Create temp test directories
	tmpDir, err := os.MkdirTemp("", "mjml-dev-test")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	// Create directory structure
	docsDir := filepath.Join(tmpDir, "documents")
	partialsDir := filepath.Join(tmpDir, "partials")
	outDir := filepath.Join(tmpDir, "dist")

	// Create test files
	documents := map[string]string{
		"welcome.mjml":           `<mjml><mj-body><mj-section><mj-column><mj-text>Hello {{.name}}</mj-text></mj-column></mj-section>{{template "footer" .}}</mj-body></mjml>`,
		"nested/newsletter.mjml": `<mjml><mj-body><mj-section><mj-column><mj-text>Newsletter {{.title}}</mj-text></mj-column></mj-section></mj-body></mjml>`,
	}

	partials := map[string]string{
		"footer.mjml": `<mj-section><mj-column><mj-text>Contact: {{.email}}</mj-text></mj-column></mj-section>`,
	}

	// Create and populate directories
	for path, content := range documents {
		fullPath := filepath.Join(docsDir, path)
		r.NoError(os.MkdirAll(filepath.Dir(fullPath), 0755))
		r.NoError(os.WriteFile(fullPath, []byte(content), 0644))
	}

	for path, content := range partials {
		fullPath := filepath.Join(partialsDir, path)
		r.NoError(os.MkdirAll(filepath.Dir(fullPath), 0755))
		r.NoError(os.WriteFile(fullPath, []byte(content), 0644))
	}

	// Create test config
	cfg := &config.Config{
		Paths: config.Paths{
			Documents: docsDir,
			Partials:  partialsDir,
			Output:    outDir,
		},
		MJML: config.MJMLConfig{
			Minify:          false, // Pretty output for better snapshot readability
			ValidationLevel: "soft",
		},
		Template: config.TemplateConfig{
			Variables: map[string]any{
				"email": "test@example.com",
			},
			Documents: map[string]config.Template{
				"welcome": {
					Variables: map[string]any{
						"name": "World",
					},
				},
				"nested/newsletter": {
					Variables: map[string]any{
						"title": "Latest News",
					},
				},
			},
		},
	}

	// Create and run processor
	processor, err := handler.NewProcessor(cfg)
	r.NoError(err)

	err = processor.Process()
	r.NoError(err)

	// Verify output files exist and content
	files, err := os.ReadDir(outDir)
	r.NoError(err)
	r.Len(files, 2)

	// Verify and snapshot the outputs
	welcomeContent, err := os.ReadFile(filepath.Join(outDir, "welcome.html"))

	r.NoError(err)
	r.Contains(string(welcomeContent), "Hello World")
	r.Contains(string(welcomeContent), "test@example.com")
	err = cupaloy.SnapshotWithName("TestProcessor-welcome", string(welcomeContent))
	r.NoError(err)

	newsletterContent, err := os.ReadFile(filepath.Join(outDir, "nested/newsletter.html"))
	r.NoError(err)
	r.Contains(string(newsletterContent), "Latest News")
	err = cupaloy.SnapshotWithName("TestProcessor-nested_newsletter", string(welcomeContent))
	r.NoError(err)
}

func TestProcessor_Errors(t *testing.T) {
	r := require.New(t)

	t.Run("missing directories", func(t *testing.T) {
		cfg := &config.Config{
			Paths: config.Paths{
				Documents: "/non/existent/path",
				Partials:  "/another/non/existent",
				Output:    "/tmp/out",
			},
		}

		processor, err := handler.NewProcessor(cfg)
		r.NoError(err)
		err = processor.Process()
		r.Error(err)
		var procErr *handler.Error
		r.ErrorAs(err, &procErr)
		r.Equal(handler.ErrorLoadingFiles, procErr.Type)
	})

	t.Run("invalid template syntax", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "mjml-dev-test")
		r.NoError(err)
		defer os.RemoveAll(tmpDir)

		docsDir := filepath.Join(tmpDir, "documents")
		outDir := filepath.Join(tmpDir, "dist")

		// Create template with invalid syntax
		r.NoError(os.MkdirAll(docsDir, 0755))
		r.NoError(os.WriteFile(
			filepath.Join(docsDir, "invalid.mjml"),
			[]byte(`<mjml><mj-body><mj-text>{{ .name }</mj-text></mj-body></mjml>`),
			0644,
		))

		cfg := &config.Config{
			Paths: config.Paths{
				Documents: docsDir,
				Output:    outDir,
			},
		}

		processor, err := handler.NewProcessor(cfg)
		r.NoError(err)

		err = processor.Process()
		r.Error(err)
		var procErr *handler.Error
		r.ErrorAs(err, &procErr)
		r.Equal(handler.ErrorRendering, procErr.Type)
		r.Equal("invalid", procErr.Doc)
	})

	t.Run("non-writable output directory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "mjml-dev-test")
		r.NoError(err)
		defer os.RemoveAll(tmpDir)

		docsDir := filepath.Join(tmpDir, "documents")
		outDir := filepath.Join(tmpDir, "out")

		// Create valid document
		r.NoError(os.MkdirAll(docsDir, 0755))
		r.NoError(os.WriteFile(
			filepath.Join(docsDir, "test.mjml"),
			[]byte(`<mjml><mj-body><mj-section><mj-column><mj-text>Hello</mj-text></mj-column></mj-section></mj-body></mjml>`),
			0644,
		))

		// Create non-writable output directory
		r.NoError(os.MkdirAll(outDir, 0444))

		cfg := &config.Config{
			Paths: config.Paths{
				Documents: docsDir,
				Output:    outDir,
			},
		}

		processor, err := handler.NewProcessor(cfg)
		r.NoError(err)

		err = processor.Process()
		r.Error(err)
		var procErr *handler.Error
		r.ErrorAs(err, &procErr)
		r.Equal(handler.ErrorSaving, procErr.Type)
	})
}

func TestProcessor_Configuration(t *testing.T) {
	r := require.New(t)

	t.Run("minification settings", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "mjml-dev-test")
		r.NoError(err)
		defer os.RemoveAll(tmpDir)

		docsDir := filepath.Join(tmpDir, "documents")
		outMinified := filepath.Join(tmpDir, "dist-min")
		outPretty := filepath.Join(tmpDir, "dist-pretty")

		// Create test document
		r.NoError(os.MkdirAll(docsDir, 0755))
		r.NoError(os.WriteFile(
			filepath.Join(docsDir, "test.mjml"),
			[]byte(`<mjml><mj-body><mj-section><mj-column><mj-text>Hello</mj-text></mj-column></mj-section></mj-body></mjml>`),
			0644,
		))

		// Process with minification
		cfgMinified := &config.Config{
			Paths: config.Paths{
				Documents: docsDir,
				Output:    outMinified,
			},
			MJML: config.MJMLConfig{
				Minify: true,
			},
		}
		processorMin, err := handler.NewProcessor(cfgMinified)
		r.NoError(err)
		r.NoError(processorMin.Process())

		// Process without minification
		cfgPretty := &config.Config{
			Paths: config.Paths{
				Documents: docsDir,
				Output:    outPretty,
			},
			MJML: config.MJMLConfig{
				Minify: false,
			},
		}
		processorPretty, err := handler.NewProcessor(cfgPretty)
		r.NoError(err)
		r.NoError(processorPretty.Process())

		// Compare outputs
		minContent, err := os.ReadFile(filepath.Join(outMinified, "test.html"))
		r.NoError(err)
		prettyContent, err := os.ReadFile(filepath.Join(outPretty, "test.html"))
		r.NoError(err)

		r.Less(len(minContent), len(prettyContent), "minified content should be shorter")
	})
}

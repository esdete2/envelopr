package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		r := require.New(t)

		// Create temp config file
		tmpDir, err := os.MkdirTemp("", "mjml-dev-test")
		r.NoError(err)
		defer os.RemoveAll(tmpDir)

		configContent := `
paths:
 documents: templates
 partials: partials
 output: dist

mjml:
 validationLevel: soft
 keepComments: false
 beautify: true
 minify: true
 fonts:
   Roboto: https://fonts.googleapis.com/css?family=Roboto

template:
 preserveHrefExpressions: true
 variables:
   companyName: ACME Corp
   supportEmail: support@acme.com
   logoUrl: https://example.com/logo.png
   year: 2024
 documents:
   welcome:
     title: Welcome Email
     variables:
       userName: John Doe
       activationLink: https://example.com/activate
   newsletter:
     title: Monthly Newsletter
     variables:
       title: Latest Updates
`
		configPath := filepath.Join(tmpDir, "config.yaml")
		err = os.WriteFile(configPath, []byte(configContent), 0644)
		r.NoError(err)

		// Load config
		cfg, err := LoadConfig(configPath)
		r.NoError(err)

		// Verify paths
		r.Equal("templates", cfg.Paths.Documents)
		r.Equal("partials", cfg.Paths.Partials)
		r.Equal("dist", cfg.Paths.Output)

		// Verify MJML config
		r.Equal("soft", cfg.MJML.ValidationLevel)
		r.False(cfg.MJML.KeepComments)
		r.True(cfg.MJML.Beautify)
		r.True(cfg.MJML.Minify)
		r.Equal("https://fonts.googleapis.com/css?family=Roboto", cfg.MJML.Fonts["Roboto"])

		// Verify template config
		r.True(cfg.Template.PreserveHrefExpressions)

		// Check global variables
		r.Equal("ACME Corp", cfg.Template.Variables["companyName"])
		r.Equal("support@acme.com", cfg.Template.Variables["supportEmail"])
		r.Equal(2024, cfg.Template.Variables["year"])

		// Check document configs
		welcome, exists := cfg.Template.Documents["welcome"]
		r.True(exists)
		r.Equal("Welcome Email", welcome.Title)
		r.Equal("John Doe", welcome.Variables["userName"])

		newsletter, exists := cfg.Template.Documents["newsletter"]
		r.True(exists)
		r.Equal("Monthly Newsletter", newsletter.Title)
		r.Equal("Latest Updates", newsletter.Variables["title"])
	})

	t.Run("missing config file", func(t *testing.T) {
		r := require.New(t)

		cfg, err := LoadConfig("nonexistent.yaml")
		r.Error(err)
		r.Nil(cfg)
	})

	t.Run("invalid yaml", func(t *testing.T) {
		r := require.New(t)

		tmpDir, err := os.MkdirTemp("", "mjml-dev-test")
		r.NoError(err)
		defer os.RemoveAll(tmpDir)

		invalidConfig := `
paths:
 documents: templates
 partials: partials
invalid_yaml:
 - : invalid
`
		configPath := filepath.Join(tmpDir, "invalid.yaml")
		err = os.WriteFile(configPath, []byte(invalidConfig), 0644)
		r.NoError(err)

		cfg, err := LoadConfig(configPath)
		r.Error(err)
		r.Nil(cfg)
	})

	t.Run("empty config", func(t *testing.T) {
		r := require.New(t)

		tmpDir, err := os.MkdirTemp("", "mjml-dev-test")
		r.NoError(err)
		defer os.RemoveAll(tmpDir)

		emptyConfig := ``
		configPath := filepath.Join(tmpDir, "empty.yaml")
		err = os.WriteFile(configPath, []byte(emptyConfig), 0644)
		r.NoError(err)

		cfg, err := LoadConfig(configPath)
		r.NoError(err) // Empty config is valid
		r.Equal(&Config{
			Paths: Paths{
				Documents: "documents",
				Partials:  "partials",
				Output:    "output",
			},
			MJML: MJMLConfig{
				ValidationLevel: "soft", // Default validation level
			},
			Template: TemplateConfig{},
		}, cfg)
	})
}

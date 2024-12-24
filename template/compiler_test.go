package template_test

import (
	"testing"

	"github.com/Boostport/mjml-go"
	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/stretchr/testify/require"

	"github.com/esdete2/mjml-dev/config"
	"github.com/esdete2/mjml-dev/template"
)

func TestCompiler(t *testing.T) {
	r := require.New(t)

	t.Run("successful compilation", func(t *testing.T) {
		input := `<mjml><mj-body><mj-section><mj-column><mj-text>Hello World</mj-text></mj-column></mj-section></mj-body></mjml>`

		cfg := &config.Config{
			MJML: config.MJMLConfig{
				Minify:          false,
				ValidationLevel: "soft",
			},
		}
		compiler := template.NewCompiler(cfg)
		result, err := compiler.Compile(input)
		r.NoError(err)

		cupaloy.SnapshotT(t, result)
	})

	t.Run("minify", func(t *testing.T) {
		input := `<mjml><mj-body><mj-section><mj-column><mj-text>Hello World</mj-text></mj-column></mj-section></mj-body></mjml>`

		cfg := &config.Config{
			MJML: config.MJMLConfig{
				Minify:          true,
				ValidationLevel: "soft",
			},
		}
		compiler := template.NewCompiler(cfg)
		minified, err := compiler.Compile(input)
		r.NoError(err)
		cupaloy.SnapshotT(t, minified)
	})

	t.Run("preserve href expressions", func(t *testing.T) {
		input := `<mjml><mj-body><mj-section><mj-column><mj-button href="{{ .buttonUrl }}">Click me</mj-button></mj-column></mj-section></mj-body></mjml>`

		cfg := &config.Config{
			MJML: config.MJMLConfig{
				Minify:          true,
				ValidationLevel: "soft",
			},
		}
		compiler := template.NewCompiler(cfg)
		result, err := compiler.Compile(input)
		r.NoError(err)
		r.Contains(result, `href="{{ .buttonUrl }}"`)
		cupaloy.SnapshotT(t, result)
	})

	t.Run("custom fonts", func(t *testing.T) {
		input := `<mjml><mj-body><mj-section><mj-column><mj-text font-family="CustomFont">Hello World</mj-text></mj-column></mj-section></mj-body></mjml>`

		cfg := &config.Config{
			MJML: config.MJMLConfig{
				Minify:          false,
				ValidationLevel: "soft",
				Fonts: map[string]string{
					"CustomFont": "https://fonts.googleapis.com/css?family=Roboto",
				},
			},
		}
		compiler := template.NewCompiler(cfg)
		result, err := compiler.Compile(input)
		r.NoError(err)
		r.Contains(result, "CustomFont")
		cupaloy.SnapshotT(t, result)
	})

	t.Run("invalid mjml", func(t *testing.T) {
		cfg := &config.Config{
			MJML: config.MJMLConfig{
				Minify:          true,
				ValidationLevel: "strict",
			},
		}
		compiler := template.NewCompiler(cfg)
		input := `<mjml><mj-body><invalid-tag></invalid-tag></mj-body></mjml>`

		result, err := compiler.Compile(input)
		r.Error(err)

		var mjmlError mjml.Error
		r.ErrorAs(err, &mjmlError)
		r.Empty(result)
	})
}

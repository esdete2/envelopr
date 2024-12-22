package template_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/esdete2/mjml-dev/template"
)

func TestRenderer(t *testing.T) {
	t.Run("simple template without partials", func(t *testing.T) {
		r := require.New(t)

		docs := []template.Template{
			{
				Name:    "welcome",
				Content: `<mjml><mj-body><mj-text>Hello {{.name}}</mj-text></mj-body></mjml>`,
			},
		}

		renderer := template.NewRenderer(docs, nil)
		result, err := renderer.Render("welcome", map[string]interface{}{
			"name": "World",
		})
		r.NoError(err)
		r.Contains(result, "Hello World")
	})

	t.Run("template with partial", func(t *testing.T) {
		r := require.New(t)

		docs := []template.Template{
			{
				Name:    "welcome",
				Content: `<mjml><mj-body>{{template "header" .}}<mj-text>Content</mj-text></mj-body></mjml>`,
			},
		}

		partials := []template.Template{
			{
				Name:    "header",
				Content: `<mj-text>Hello {{.name}}</mj-text>`,
			},
		}

		renderer := template.NewRenderer(docs, partials)
		result, err := renderer.Render("welcome", map[string]interface{}{
			"name": "World",
		})
		r.NoError(err)
		r.Contains(result, "Hello World")
		r.Contains(result, "Content")
	})

	t.Run("template not found", func(t *testing.T) {
		r := require.New(t)

		renderer := template.NewRenderer(nil, nil)
		result, err := renderer.Render("nonexistent", nil)
		r.Error(err)
		r.Contains(err.Error(), template.ErrTemplateNotFound.Error())
		r.Empty(result)
	})

	t.Run("partial not found", func(t *testing.T) {
		r := require.New(t)

		docs := []template.Template{
			{
				Name:    "welcome",
				Content: `<mjml><mj-body>{{template "missing" .}}</mj-body></mjml>`,
			},
		}

		renderer := template.NewRenderer(docs, nil)
		result, err := renderer.Render("welcome", nil)
		r.Error(err)
		r.Empty(result)
	})

	t.Run("nested partials", func(t *testing.T) {
		r := require.New(t)

		docs := []template.Template{
			{
				Name:    "welcome",
				Content: `<mjml><mj-body>{{template "header" .}}</mj-body></mjml>`,
			},
		}

		partials := []template.Template{
			{
				Name:    "header",
				Content: `<mj-section>{{template "logo" .}}</mj-section>`,
			},
			{
				Name:    "logo",
				Content: `<mj-image src="{{.logo_url}}"></mj-image>`,
			},
		}

		renderer := template.NewRenderer(docs, partials)
		result, err := renderer.Render("welcome", map[string]interface{}{
			"logo_url": "https://example.com/logo.png",
		})
		r.NoError(err)
		r.Contains(result, "https://example.com/logo.png")
	})
}

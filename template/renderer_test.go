package template_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/esdete2/envelopr/template"
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

	t.Run("dict function with template expression", func(t *testing.T) {
		r := require.New(t)

		docs := []template.Template{
			{
				Name:    "email",
				Content: `<mjml><mj-body>{{template "button" dict "url" "{{ .campaignUrl }}" "label" "Shop Now"}}</mj-body></mjml>`,
			},
		}

		partials := []template.Template{
			{
				Name:    "button",
				Content: `<mj-button href="{{.url}}">{{.label}}</mj-button>`,
			},
		}

		renderer := template.NewRenderer(docs, partials)
		result, err := renderer.Render("email", map[string]interface{}{
			"campaignUrl": "https://example.com/campaign",
		})
		r.NoError(err)
		r.Contains(result, `href="{{ .campaignUrl }}"`)
		r.Contains(result, "Shop Now")
	})

	t.Run("expression function", func(t *testing.T) {
		r := require.New(t)

		docs := []template.Template{
			{
				Name:    "welcome",
				Content: `<mjml><mj-body><mj-text>Hello {{ expression ".username" }}</mj-text></mj-body></mjml>`,
			},
		}

		renderer := template.NewRenderer(docs, nil)
		result, err := renderer.Render("welcome", nil)
		r.NoError(err)
		r.Contains(result, "Hello {{ .username }}")
	})

	t.Run("combination of dict and expression", func(t *testing.T) {
		r := require.New(t)

		docs := []template.Template{
			{
				Name: "email",
				Content: `<mjml><mj-body>{{template "button" dict 
                    "url" (expression ".profileUrl")
                    "label" (expression ".buttonText")
                }}</mj-body></mjml>`,
			},
		}

		partials := []template.Template{
			{
				Name:    "button",
				Content: `<mj-button href="{{.url}}">{{.label}}</mj-button>`,
			},
		}

		renderer := template.NewRenderer(docs, partials)
		result, err := renderer.Render("email", nil)
		r.NoError(err)
		r.Contains(result, `href="{{ .profileUrl }}"`)
		r.Contains(result, `{{ .buttonText }}`)
	})
}

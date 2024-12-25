package template

import (
	"context"

	"github.com/Boostport/mjml-go"
	"github.com/friendsofgo/errors"

	"github.com/esdete2/envelopr/config"
)

type Compiler struct {
	config *config.Config
}

func NewCompiler(config *config.Config) *Compiler {
	return &Compiler{
		config: config,
	}
}

func (c *Compiler) Compile(content string) (string, error) {
	// Convert config to MJML options
	options := []mjml.ToHTMLOption{
		mjml.WithMinify(c.config.MJML.Minify),
		mjml.WithBeautify(c.config.MJML.Beautify),
		mjml.WithKeepComments(c.config.MJML.KeepComments),
		mjml.WithValidationLevel(mjml.ValidationLevel(c.config.MJML.ValidationLevel)),
	}

	if len(c.config.MJML.Fonts) > 0 {
		options = append(options, mjml.WithFonts(c.config.MJML.Fonts))
	}

	// Compile MJML to HTML
	html, err := mjml.ToHTML(context.Background(), content, options...)
	if err != nil {
		return "", errors.Wrap(err, "compiling MJML")
	}

	return html, nil
}

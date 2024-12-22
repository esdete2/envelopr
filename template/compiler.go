package template

import (
	"context"
	"regexp"
	"strings"

	"github.com/Boostport/mjml-go"
	"github.com/friendsofgo/errors"

	"github.com/esdete2/mjml-dev/config"
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

	// Post-process to decode template syntax if enabled
	if c.config.Template.PreserveHrefExpressions {
		html = decodeTemplateURLs(html)
	}

	return html, nil
}

func decodeTemplateURLs(html string) string {
	// Pattern to match encoded template syntax in href attributes
	pattern := `href="(%7[bB]%7[bB].*?%7[dD]%7[dD])"`
	re := regexp.MustCompile(pattern)

	return re.ReplaceAllStringFunc(html, func(match string) string {
		// Decode the template syntax and spaces
		decoded := strings.NewReplacer(
			"%7B", "{",
			"%7b", "{",
			"%7D", "}",
			"%7d", "}",
			"%20", " ",
		).Replace(match)
		return decoded
	})
}

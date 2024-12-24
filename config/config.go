package config

import (
	"os"

	"github.com/friendsofgo/errors"
	"gopkg.in/yaml.v3"
)

type Paths struct {
	Documents string `yaml:"documents"`
	Partials  string `yaml:"partials"`
	Output    string `yaml:"output"`
}

type MJMLConfig struct {
	ValidationLevel string            `yaml:"validationLevel"`
	KeepComments    bool              `yaml:"keepComments"`
	Beautify        bool              `yaml:"beautify"`
	Minify          bool              `yaml:"minify"`
	Fonts           map[string]string `yaml:"fonts"`
}

type TemplateConfig struct {
	Variables map[string]any `yaml:"variables"`
	Documents map[string]any `yaml:"documents"`
}

type Config struct {
	Paths    Paths          `yaml:"paths"`
	MJML     MJMLConfig     `yaml:"mjml"`
	Template TemplateConfig `yaml:"template"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	yamlData, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling config")
	}

	// Set default values
	if config.Paths.Documents == "" {
		config.Paths.Documents = "documents"
	}
	if config.Paths.Output == "" {
		config.Paths.Output = "output"
	}
	if config.MJML.ValidationLevel == "" {
		config.MJML.ValidationLevel = "soft"
	}

	return &config, nil
}

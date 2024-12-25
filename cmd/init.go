package cmd

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/friendsofgo/errors"
	"github.com/networkteam/slogutils"
	"github.com/urfave/cli/v2"

	"github.com/esdete2/mjml-dev/config"
)

type initAnswers struct {
	ConfigPath   string
	DocumentsDir string
	PartialsDir  string
	OutputDir    string
	Force        bool
}

func InitCmd() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize a new MJML dev project",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-input",
				Usage: "Skip wizard and use defaults",
			},
		},
		Action: func(c *cli.Context) error {
			logger := slogutils.FromContext(c.Context)

			// Default values
			answers := initAnswers{
				ConfigPath:   "mjml-dev.yaml",
				DocumentsDir: "documents",
				PartialsDir:  "partials",
				OutputDir:    "output",
			}

			// Skip wizard if --no-input flag is set
			if !c.Bool("no-input") {
				questions := []*survey.Question{
					{
						Name: "configPath",
						Prompt: &survey.Input{
							Message: "Where should the config file be created?",
							Default: answers.ConfigPath,
						},
					},
					{
						Name: "documentsDir",
						Prompt: &survey.Input{
							Message: "Directory for MJML templates:",
							Default: answers.DocumentsDir,
						},
					},
					{
						Name: "partialsDir",
						Prompt: &survey.Input{
							Message: "Directory for partial templates:",
							Default: answers.PartialsDir,
						},
					},
					{
						Name: "outputDir",
						Prompt: &survey.Input{
							Message: "Directory for compiled templates:",
							Default: answers.OutputDir,
						},
					},
				}

				// Add force question only if config exists
				if _, err := os.Stat(answers.ConfigPath); err == nil {
					questions = append(questions, &survey.Question{
						Name: "force",
						Prompt: &survey.Confirm{
							Message: "Config file already exists. Overwrite?",
							Default: false,
						},
					})
				}

				if err := survey.Ask(questions, &answers); err != nil {
					return errors.Wrap(err, "running wizard")
				}
			}

			// Create config file
			err := config.CreateDefaultConfig(answers.ConfigPath, answers.DocumentsDir, answers.PartialsDir, answers.OutputDir, answers.Force)
			if err != nil {
				return errors.Wrap(err, "creating config file")
			}
			logger.Info("Created config file", "path", answers.ConfigPath)

			// Create directories
			dirs := []struct {
				path string
				name string
			}{
				{answers.DocumentsDir, "documents"},
				{answers.PartialsDir, "partials"},
				{answers.OutputDir, "output"},
			}

			for _, dir := range dirs {
				if err := os.MkdirAll(dir.path, 0755); err != nil {
					return errors.Wrapf(err, "creating %s directory", dir.name)
				}
				logger.Info("Created directory", "path", dir.path)
			}

			return nil
		},
	}
}

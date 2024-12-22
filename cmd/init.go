package main

import (
	"log/slog"

	"github.com/friendsofgo/errors"
	"github.com/urfave/cli/v2"

	"github.com/esdete2/mjml-dev/config"
)

func initCmd() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize a new MJML dev project",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to config file",
				Value:   "config.yaml",
			},
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Force overwrite if config file exists",
			},
		},
		Action: func(c *cli.Context) error {
			configPath := c.String("config")
			force := c.Bool("force")

			if err := config.CreateDefaultConfig(configPath, force); err != nil {
				return errors.Wrap(err, "creating config file")
			}
			slog.Info("Created config file", "path", configPath)

			return nil
		},
	}
}

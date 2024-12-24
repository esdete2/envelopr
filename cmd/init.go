package main

import (
	"github.com/friendsofgo/errors"
	"github.com/networkteam/slogutils"
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
			logger := slogutils.FromContext(c.Context)

			path := c.String("config")
			force := c.Bool("force")

			// Create config file
			err := config.CreateDefaultConfig(path, force)
			if err != nil {
				return errors.Wrap(err, "creating config file")
			}

			logger.Info("Created config file", "path", path)

			return nil
		},
	}
}

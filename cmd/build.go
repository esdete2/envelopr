package main

import (
	"github.com/friendsofgo/errors"
	"github.com/networkteam/slogutils"
	"github.com/urfave/cli/v2"

	"github.com/esdete2/mjml-dev/config"
	"github.com/esdete2/mjml-dev/handler"
)

func buildCmd() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Build MJML documents",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to config file",
				Value:   "config.yaml",
			},
		},
		Action: func(c *cli.Context) error {
			logger := slogutils.FromContext(c.Context)
			cfg, err := config.LoadConfig(c.String("config"))
			if err != nil {
				return errors.Wrap(err, "loading config")
			}

			proc, err := handler.NewProcessor(cfg)
			if err != nil {
				return errors.Wrap(err, "creating processor")
			}

			err = proc.Process()

			if err != nil {
				return errors.Wrap(err, "processing documents")
			}

			logger.Info("Documents processed successfully")

			return nil
		},
	}
}
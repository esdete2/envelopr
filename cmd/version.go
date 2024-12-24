package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func versionCmd() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print version information",
		Action: func(c *cli.Context) error {
			fmt.Printf("mjml-dev version %s\n", Version)
			fmt.Printf("Go version: %s\n", GoVersion)
			fmt.Printf("Build time: %s\n", BuildTimestamp)
			return nil
		},
	}
}

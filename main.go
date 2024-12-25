package main

import (
	"log/slog"
	"os"

	"github.com/networkteam/slogutils"
	"github.com/urfave/cli/v2"

	"github.com/esdete2/envelopr/cmd"
)

func main() {
	app := &cli.App{
		Name:  "envelopr",
		Usage: "MJML template development tool",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "verbosity",
				Usage:   "Verbosity: 1=error, 2=warn, 3=info, 4=debug, 5=trace",
				Aliases: []string{"v"},
				Value:   3,
			},
		},
		Before: func(c *cli.Context) error {
			logHandler := slogutils.NewCLIHandler(os.Stderr, &slogutils.CLIHandlerOptions{
				Level: verbosityToSlogLevel(c.Int("verbosity")),
			})
			slog.SetDefault(slog.New(logHandler))

			return nil
		},
		Commands: []*cli.Command{
			cmd.InitCmd(),
			cmd.BuildCmd(),
			cmd.WatchCmd(),
			cmd.VersionCmd(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("Command failed", slogutils.Err(err))
		os.Exit(1)
	}
}

func verbosityToSlogLevel(verbosity int) slog.Level {
	if verbosity <= 1 {
		return slog.LevelError
	}

	switch verbosity {
	case 2:
		return slog.LevelWarn
	case 3:
		return slog.LevelInfo
	case 4:
		return slog.LevelDebug
	}

	return slogutils.LevelTrace
}

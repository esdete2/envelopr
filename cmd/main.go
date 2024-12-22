package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/friendsofgo/errors"
	"github.com/networkteam/slogutils"
	"github.com/urfave/cli/v2"

	"github.com/esdete2/mjml-dev/config"
	"github.com/esdete2/mjml-dev/handler"
	"github.com/esdete2/mjml-dev/web"
)

func main() {
	app := &cli.App{
		Name:  "mjml-dev",
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
			{
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
			},
			{
				Name:  "watch",
				Usage: "Start development server with hot reload",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "Path to config file",
						Value:   "config.yaml",
					},
					&cli.StringFlag{
						Name:  "host",
						Usage: "Server host",
						Value: "localhost",
					},
					&cli.StringFlag{
						Name:  "port",
						Usage: "Server port",
						Value: "3600",
					},
				},
				Action: func(c *cli.Context) error {
					// Load configuration
					cfg, err := config.LoadConfig(c.String("config"))
					if err != nil {
						return errors.Wrap(err, "loading config")
					}

					// Initialize web server
					srv := web.NewServer(&web.ServerOptions{
						Output: cfg.Paths.Output,
					})

					// Create processor
					proc, err := handler.NewProcessor(cfg)
					if err != nil {
						return errors.Wrap(err, "creating processor")
					}

					// Do initial build
					if err := proc.Process(); err != nil {
						return errors.Wrap(err, "initial build")
					}

					// Create and start file watcher
					watcher, err := handler.NewWatcher(proc, cfg, srv)
					if err != nil {
						return errors.Wrap(err, "creating watcher")
					}

					if err := watcher.Watch(); err != nil {
						return errors.Wrap(err, "starting watcher")
					}

					addr := fmt.Sprintf("%s:%s", c.String("host"), c.String("port"))

					// Start server
					serverErr := make(chan error, 1)
					go func() {
						if err := srv.Serve(addr); err != nil {
							serverErr <- err
						}
					}()

					// Handle graceful shutdown
					quit := make(chan os.Signal, 1)
					signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

					select {
					case <-quit:
						slog.Info("Shutting down...")
						watcher.Stop()
						return nil
					case err := <-serverErr:
						watcher.Stop()
						return errors.Wrap(err, "server error")
					}
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("Command failed", slogutils.Err(err))
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

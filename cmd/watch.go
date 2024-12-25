package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/friendsofgo/errors"
	"github.com/urfave/cli/v2"

	"github.com/esdete2/mjml-dev/config"
	"github.com/esdete2/mjml-dev/handler"
	"github.com/esdete2/mjml-dev/web"
)

func WatchCmd() *cli.Command {
	return &cli.Command{
		Name:  "watch",
		Usage: "Start development server with hot reload",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to config file",
				Value:   "mjml-dev.yaml",
			},
			&cli.StringFlag{
				Name:  "host",
				Usage: "Server host",
				Value: "0.0.0.0",
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

			// Initial build
			err = proc.Process()
			if err != nil {
				return errors.Wrap(err, "processing documents")
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
	}
}

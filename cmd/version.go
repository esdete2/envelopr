package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/urfave/cli/v2"
)

var (
	toolVersion    string
	goVersion      string
	buildTimestamp string
)

func VersionCmd() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print version information",
		Action: func(c *cli.Context) error {
			version := GetVersionInfo()

			fmt.Printf("Version %s\n", version.toolVersion)
			fmt.Printf("Go version: %s\n", version.goVersion)

			if len(buildTimestamp) != 0 {
				fmt.Printf("Build time: %s\n", buildTimestamp)
			}
			return nil
		},
	}
}

type VersionInfo struct {
	toolVersion string
	goVersion   string
}

func GetVersionInfo() VersionInfo {
	if len(toolVersion) != 0 && len(goVersion) != 0 {
		return VersionInfo{
			toolVersion: toolVersion,
			goVersion:   goVersion,
		}
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		return VersionInfo{
			toolVersion: info.Main.Version,
			goVersion:   runtime.Version(),
		}
	}
	return VersionInfo{
		toolVersion: "(unknown)",
		goVersion:   runtime.Version(),
	}
}

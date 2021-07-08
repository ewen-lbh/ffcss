package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
)

const (
	Usage = `ffcss - Apply and configure FirefoxCSS themes

Usage:
    ffcss configure KEY [VALUE]
    ffcss use THEME_NAME [VARIANT]
    ffcss reapply
    ffcss init [FORMAT]

Where:
    KEY         a setting key (see firefox's about:config)
    THEME_NAME  a theme name or URL (see README.md)
	`

	VersionMajor = 0
	VersionMinor = 1
	VersionPatch = 0
)

var (
	VersionString = fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
)

func main() {
	args, _ := docopt.ParseDoc(Usage)

	if err := dispatchCommand(args); err != nil {
		panic(err)
	}
}

func dispatchCommand(args docopt.Opts) error {
	if val, _ := args.Bool("configure"); val {
		err := RunCommandConfigure(args)
		return err
	}
	if val, _ := args.Bool("use"); val {
		err := RunCommandUse(args)
		return err
	}
	if val, _ := args.Bool("reapply"); val {
		err := RunCommandReapply(args)
		return err
	}
	if val, _ := args.Bool("init"); val {
		err := RunCommandInit(args)
		return err
	}
	return nil
}

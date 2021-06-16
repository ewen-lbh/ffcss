package main

import (
	"github.com/docopt/docopt-go"
)

const (
	Usage = `ffcss - Apply and configure FirefoxCSS themes

Usage:
    ffcss configure KEY [VALUE]
    ffcss use THEME_NAME
    ffcss reapply
    ffcss init [FORMAT]

Where:
    KEY         a setting key (see firefox's about:config)
    THEME_NAME  a theme name or URL (see README.md)
	`
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

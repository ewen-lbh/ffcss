package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
)

const (
	Usage = `ffcss - Apply and configure FirefoxCSS themes

Usage:
    ffcss [--profiles-dir=DIRECTORY] configure KEY [VALUE]
    ffcss [--profiles-dir=DIRECTORY] [--all-profiles] use THEME_NAME [VARIANT]
    ffcss reapply
	ffcss cache clear
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

	err := os.MkdirAll(CacheDir(), 0700)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(ConfigDir(), 0700)
	if err != nil {
		panic(err)
	}

	if err := dispatchCommand(args); err != nil {
		fmt.Println("Woops! An error occured:")
		fmt.Println()
		for idx, errorFragment := range strings.Split(err.Error(), ": ") {
			fmt.Println(strings.Repeat("  ", idx) + "-> " + errorFragment)
		}
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
	if val, _ := args.Bool("cache"); val {
		if val, _ := args.Bool("clear"); val {
			return ClearWholeCache()
		}
	}
	return nil
}

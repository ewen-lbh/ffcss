// Apply and configure FirefoxCSS themes
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
	ffcss [options] use THEME_NAME [VARIANT]
	ffcss [options] cache clear
	ffcss [options] init
	ffcss version [COMPONENT]

Where:
	THEME_NAME  a theme name or URL (see README.md)
	COMPONENT   is either major, minor or patch (to get a single digit)

Options:
	--all-profiles       Apply the theme to all profiles
	--profiles-dir=PATH  Directory that contains profile directories.
	                     Default value is platform-specific:
	                     - $HOME/.mozilla/firefox                                on Linux
	                     - $HOME/Library/Application Support/Firefox/Profiles    on MacOS
	                     - %appdata%/Roaming/Mozilla/Firefox/Profiles            on Windows
	`

	VersionMajor = 0
	VersionMinor = 1
	VersionPatch = 2
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
			fmt.Println(strings.Repeat(indent, idx) + "-> " + errorFragment)
		}
	}
}

func dispatchCommand(args docopt.Opts) error {
	d("dispatching %#v", args)
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
	if val, _ := args.Bool("version"); val {
		component, _ := args.String("COMPONENT")
		switch component {
		case "major":
			fmt.Println(VersionMajor)
		case "minor":
			fmt.Println(VersionMinor)
		case "patch":
			fmt.Println(VersionPatch)
		default:
			fmt.Println(VersionString)
		}
	}
	return nil
}

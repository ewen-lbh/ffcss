// Apply and configure FirefoxCSS themes
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/ewen-lbh/ffcss"
)

const (
	usage = `ffcss - Apply and configure FirefoxCSS themes

Usage:
	ffcss [options] use THEME_NAME [VARIANT]
	ffcss [options] get THEME_NAME
	ffcss [options] cache clear
	ffcss [options] init
	ffcss [options] reapply
	ffcss version [COMPONENT]

Where:
	THEME_NAME  a theme name or URL (see README.md)
	COMPONENT   is either major, minor or patch (to get a single digit)

Options:
	-a --all-profiles        Apply the theme to all profiles
	-p --profiles=PATHS      Select which profiles to apply the theme to.
	                         Can be absolute or relative to --profiles-dir.
							 Comma-separated.
	--profiles-dir=PATH      Directory that contains profile directories.
	                         Default value is platform-specific:
	                         - $HOME/.mozilla/firefox                                on Linux
	                         - $HOME/Library/Application Support/Firefox/Profiles    on MacOS
	                         - %appdata%/Roaming/Mozilla/Firefox/Profiles            on Windows
	-d --default-profile     Apply the themes to the default profile (ending with default-release)
	--skip-manifest-source   Don't ask to show the manifest source
	`
)

var (
	out io.Writer = os.Stdout
)

func main() {
	args, _ := docopt.ParseDoc(usage)

	err := os.MkdirAll(ffcss.CacheDir(), 0700)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(ffcss.ConfigDir(), 0700)
	if err != nil {
		panic(err)
	}

	if err := dispatchCommand(flagsAndArgs{args}); err != nil {
		fmt.Fprintln(out)
		ffcss.LogError("Woops! An error occurred:")
		fmt.Fprintln(out)
		ffcss.DisplayErrorMessage(err)
	}
}

func dispatchCommand(args flagsAndArgs) error {
	ffcss.LogDebug("dispatching %#v", args)
	if val, _ := args.Bool("configure"); val {
		return fmt.Errorf("not implemented")
	}
	if val, _ := args.Bool("use"); val {
		err := runCommandUse(args)
		return err
	}
	if val, _ := args.Bool("get"); val {
		err := runCommandGet(args)
		return err
	}
	if val, _ := args.Bool("reapply"); val {
		err := runCommandReapply(args)
		return err
	}
	if val, _ := args.Bool("init"); val {
		err := runCommandInit(args)
		return err
	}
	if val, _ := args.Bool("cache"); val {
		if val, _ := args.Bool("clear"); val {
			return ffcss.ClearWholeCache()
		}
	}
	if val, _ := args.Bool("version"); val {
		component, _ := args.String("COMPONENT")
		switch component {
		case "major":
			fmt.Fprintln(out, ffcss.VersionMajor)
		case "minor":
			fmt.Fprintln(out, ffcss.VersionMinor)
		case "patch":
			fmt.Fprintln(out, ffcss.VersionPatch)
		default:
			fmt.Fprintln(out, ffcss.VersionString)
		}
	}
	return nil
}

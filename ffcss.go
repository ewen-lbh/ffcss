// Apply and configure FirefoxCSS themes
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
)

const (
	Usage = `ffcss - Apply and configure FirefoxCSS themes

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
	-a --all-profiles           Apply the theme to all profiles
	-p --profiles=PATHS      Select which profiles to apply the theme to.
	                         Can be absolute or relative to --profiles-dir.
							 Comma-separated.
	--profiles-dir=PATH      Directory that contains profile directories.
	                         Default value is platform-specific:
	                         - $HOME/.mozilla/firefox                                on Linux
	                         - $HOME/Library/Application Support/Firefox/Profiles    on MacOS
	                         - %appdata%/Roaming/Mozilla/Firefox/Profiles            on Windows
	--skip-manifest-source   Don't ask to show the manifest source
	`

	VersionMajor = 0
	VersionMinor = 2
	VersionPatch = 0
)

var (
	VersionString = fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	// Used to capture stdout instead of printing in tests
	out io.Writer = os.Stdout
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
		fmt.Fprintln(out)
		showError("Woops! An error occurred:")
		fmt.Fprintln(out)
		for idx, errorFragment := range strings.Split(err.Error(), ": ") {
			li(uint(idx), errorFragment)
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
	if val, _ := args.Bool("get"); val {
		err := RunCommandGet(args)
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
			fmt.Fprintln(out, VersionMajor)
		case "minor":
			fmt.Fprintln(out, VersionMinor)
		case "patch":
			fmt.Fprintln(out, VersionPatch)
		default:
			fmt.Fprintln(out, VersionString)
		}
	}
	return nil
}

// RunCommandInit runs the command "init"
func RunCommandInit(args docopt.Opts) error {
	// TODO: set user{Chrome,Content,.js} by finding their path
	// TODO: only set assets if chrome/ actually exists
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working directory: %w", err)
	}

	theme, err := InitializeTheme(workingDir)
	if err != nil {
		return  fmt.Errorf("while initializing theme: %w",  err)
	}

	return theme.WriteManifest(workingDir)
}

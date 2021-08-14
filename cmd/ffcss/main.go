// Apply and configure FirefoxCSS themes
package main

import (
	"fmt"
	"io"
	"os"
	_ "embed"

	"github.com/docopt/docopt-go"
	"github.com/ewen-lbh/ffcss"
)

//go:embed USAGE
var usage string

var out io.Writer = os.Stdout

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
	if val, _ := args.Bool("reset"); val {
		return runCommandReset(args)
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

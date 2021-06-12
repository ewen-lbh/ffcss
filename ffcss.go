package main

import (
	"github.com/docopt/docopt-go"
)

func main() {
	usage := ReadFile("./USAGE")
	args, _ := docopt.ParseDoc(usage)

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

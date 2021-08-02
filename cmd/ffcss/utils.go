package main

import (
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/ewen-lbh/ffcss"
)

type flagsAndArgs struct {
	docopt.Opts
}

func (o flagsAndArgs) string(name string) string {
	val, err := o.String(name)
	if err != nil {
		ffcss.LogDebug("while getting value of %s: %s: using %#v", name, err.Error(), val)
	}
	return val
}

func (o flagsAndArgs) bool(name string) bool {
	val, err := o.Bool(name)
	if err != nil {
		ffcss.LogDebug("while getting value of %s: %s: using %#v", name, err.Error(), val)
	}
	return val
}

func (o flagsAndArgs) strings(name string) []string {
	val, err := o.String(name)
	if err != nil {
		ffcss.LogDebug("while getting value of %s: %s: using %#v", name, err.Error(), strings.Split(val, ","))
	}
	return strings.Split(val, ",")
}

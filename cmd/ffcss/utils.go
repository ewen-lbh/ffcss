package main

import (
	"strings"

	"github.com/docopt/docopt-go"
)

type flagsAndArgs struct {
	docopt.Opts
}

func (o flagsAndArgs) string(name string) string {
	val, err := o.String(name)
	if err != nil {
		panic(err)
	}
	return val
}

func (o flagsAndArgs) bool(name string) bool {
	val, err := o.Bool(name)
	if err != nil {
		panic(err)
	}
	return val
}

func (o flagsAndArgs) strings(name string) []string {
	val, err := o.String(name)
	if err != nil {
		panic(err)
	}
	return strings.Split(val, ",")
}

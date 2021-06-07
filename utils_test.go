package main

import (
	"fmt"
	"testing"
)

func TestGetMozillaReleasesPaths(t *testing.T) {
	paths, err := GetMozillaReleasesPaths()
	if err != nil {
		panic(err)
	}
	Assert(t, fmt.Sprintf("%#v", paths), `[]string{"/home/ewen/.mozilla/firefox/hhkjqjta.default", "/home/ewen/.mozilla/firefox/vv6f899j.default-release", "/home/ewen/.mozilla/firefox/yhzcf3nm.default-nightly"}`)
}

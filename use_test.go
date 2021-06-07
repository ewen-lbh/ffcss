package main

import (
	"testing"
)

func TestResolveThemeName(t *testing.T) {
	Assert(t,
		ResolveThemeName("ewen-lbh/ffcss"),
		"https://github.com/ewen-lbh/ffcss",
	)
	Assert(t,
		ResolveThemeName("bitbucket.io/guaca/mole"),
		"https://bitbucket.io/guaca/mole",
	)
	Assert(t,
		ResolveThemeName("https://ewen.works"),
		"https://ewen.works",
	)
	Assert(t,
		ResolveThemeName("materialfox"),
		"https://github.com/muckSponge/MaterialFox",
	)
	Assert(t,
		ResolveThemeName("unknownone"),
		"",
	)
}

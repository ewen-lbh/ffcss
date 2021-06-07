package main

import (
	"net/url"
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

func TestDownloadRepository(t *testing.T) {
	URL, _ := url.Parse("https://github.com/muckSponge/MaterialFox")
	clonedTo, err := DownloadRepository(*URL)
	if err != nil {
		t.Fatal(err.Error())
	}
	Assert(t,
		clonedTo,
		"/home/ewen/.config/ffcss/themes/@muckSponge/MaterialFox",
	)
}

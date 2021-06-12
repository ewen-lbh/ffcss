package main

import (
	"net/url"
	"testing"
)

func TestResolveThemeName(t *testing.T) {
	name, typ := ResolveThemeName("ewen-lbh/ffcss")
	if "https://github.com/ewen-lbh/ffcss" != name || typ != "git" {
		t.Errorf(`ewen-lbh/ffcss
want: (%s, %s)
got:  (%s, %s)`, "https://github.com/ewen-lbh/ffcss", "git", name, typ)
	}

	name, typ = ResolveThemeName("bitbucket.io/guaca/mole")
	if "https://bitbucket.io/guaca/mole" != name || typ != "website" {
		t.Errorf(`bitbucket.io/guaca/mole
want: (%s, %s)
got:  (%s, %s)`, "https://bitbucket.io/guaca/mole", "website", name, typ)
	}

	name, typ = ResolveThemeName("http://localhost:8080/")
	if "http://localhost:8080/" != name || typ != "website" {
		t.Errorf(`http://localhost:8080/
want: (%s, %s)
got:  (%s, %s)`, "http://localhost:8080/", "website", name, typ)
	}

	name, typ = ResolveThemeName("materialfox")
	if "https://github.com/muckSponge/MaterialFox" != name || typ != "git" {
		t.Errorf(`materialfox
want: (%s, %s)
got:  (%s, %s)`, "https://github.com/muckSponge/MaterialFox", "git", name, typ)
	}

	name, typ = ResolveThemeName("unknownone")
	if "" != name || typ != "" {
		t.Errorf(`unknownone
want: (%s, %s)
got:  (%s, %s)`, "", "", name, typ)
	}
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

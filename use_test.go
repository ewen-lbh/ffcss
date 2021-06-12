package main

import (
	"net/url"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestResolveThemeName(t *testing.T) {
	name, typ := ResolveThemeName("ewen-lbh/ffcss")
	assert.Equal(t, []string{name, typ}, []string{"https://github.com/ewen-lbh/ffcss", "git"})

	name, typ = ResolveThemeName("bitbucket.io/guaca/mole")
	assert.Equal(t, []string{name, typ}, []string{"https://bitbucket.io/guaca/mole", "website"})

	name, typ = ResolveThemeName("http://localhost:8080/")
	assert.Equal(t, []string{name, typ}, []string{"http://localhost:8080/", "website"})

	name, typ = ResolveThemeName("materialfox")
	assert.Equal(t, []string{name, typ}, []string{"https://github.com/muckSponge/MaterialFox", "git"})

	name, typ = ResolveThemeName("unknownone")
	assert.Equal(t, []string{name, typ}, []string{"", ""})
}

func TestGetThemeDownloadPath(t *testing.T) {
	URL, _ := url.Parse("https://github.com/muckSponge/MaterialFox")
	assert.Equal(t, GetThemeDownloadPath(*URL), GetConfigDir()+"/themes/@muckSponge/MaterialFox")

	URL, _ = url.Parse("https://github.com/users/schoolsyst")
	assert.Equal(t, GetThemeDownloadPath(*URL), GetConfigDir()+"/themes/-github.com/users/schoolsyst")

	URL, _ = url.Parse("https://ewen.works/")
	assert.Equal(t, GetThemeDownloadPath(*URL), GetConfigDir()+"/themes/-ewen.works")

	URL, _ = url.Parse("http://localhost:8080")
	assert.Equal(t, GetThemeDownloadPath(*URL), GetConfigDir()+"/themes/-localhost:8080")

}

func TestIsURLClonable(t *testing.T) {
	var actual bool
	var err error

	actual, err = IsURLClonable(urlOf("https://github.com/ewen-lbh/ffcss/"))
	assert.Equal(t, actual, true)
	assert.Nil(t, err)

	actual, err = IsURLClonable(urlOf("https://github.com/users/schoolsyst"))
	assert.Equal(t, actual, false)
	assert.Nil(t, err)

	actual, err = IsURLClonable(urlOf("https://ewen.works/"))
	assert.Equal(t, actual, false)
	assert.Nil(t, err)
}

	}
	Assert(t,
		clonedTo,
		"/home/ewen/.config/ffcss/themes/@muckSponge/MaterialFox",
	)
}

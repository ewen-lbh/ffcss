package main

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func urlOf(urlrepr string) url.URL {
	URL, err := url.Parse(urlrepr)
	if err != nil {
		panic(err)
	}
	return *URL
}

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

func TestDownloadFromZip(t *testing.T) {
	dl := func(s string) error {
		tempDir, err := os.MkdirTemp("testarea", "*")
		if err != nil {
			panic(err)
		}
		return DownloadFromZip(urlOf(s), tempDir)
	}

	assert.EqualError(t, dl("https://ewen.works/girehigerhiugrehigerhi"), "couldn't check remote file: server returned 404 Not Found")
	assert.EqualError(t, dl("https://ewen.works/"), "expected a zip file (application/zip), got a text/html")
	// TODO: don't. use local files. this is horrible. Spinning up a localhost webserver might be required here, since file:// is not supported by wget
	// TODO: check for absence of unzipped folder
	assert.EqualError(t, dl("https://media.ewen.works/ffcss/mocks/themeWithNoManifest.zip"), "downloaded zip file has no manifest file (ffcss.yaml)")
	// TODO: check for presence of unzipped folder
	assert.Nil(t, dl("https://media.ewen.works/ffcss/mocks/materialfox.zip"))
	// TODO: check for absence of zip file, in both cases
}

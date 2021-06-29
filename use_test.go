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
	assert.Equal(t, []string{"https://github.com/ewen-lbh/ffcss", "git"}, []string{name, typ})

	name, typ = ResolveThemeName("bitbucket.io/guaca/mole")
	assert.Equal(t, []string{"https://bitbucket.io/guaca/mole", "website"}, []string{name, typ})

	name, typ = ResolveThemeName("http://localhost:8080/")
	assert.Equal(t, []string{"http://localhost:8080/", "website"}, []string{name, typ})

	name, typ = ResolveThemeName("materialfox")
	assert.Equal(t, []string{"https://github.com/muckSponge/MaterialFox", "git"}, []string{name, typ})

	name, typ = ResolveThemeName("unknownone")
	assert.Equal(t, []string{"", ""}, []string{name, typ})
}

func TestGetThemeDownloadPath(t *testing.T) {
	URL, _ := url.Parse("https://github.com/muckSponge/MaterialFox")
	assert.Equal(t, GetConfigDir()+"/themes/@muckSponge/MaterialFox", GetThemeDownloadPath(*URL))

	URL, _ = url.Parse("https://github.com/users/schoolsyst", GetThemeDownloadPath(*URL))
	assert.Equal(t, GetConfigDir()+"/themes/-github.com/users/schoolsyst", GetThemeDownloadPath(*URL))

	URL, _ = url.Parse("https://ewen.works/")
	assert.Equal(t, GetConfigDir()+"/themes/-ewen.works", GetThemeDownloadPath(*URL))

	URL, _ = url.Parse("http://localhost:8080")
	assert.Equal(t, GetConfigDir()+"/themes/-localhost:8080", GetThemeDownloadPath(*URL))

}

func TestIsURLClonable(t *testing.T) {
	var actual bool
	var err error

	actual, err = IsURLClonable(urlOf("https://github.com/ewen-lbh/ffcss/"))
	assert.Equal(t, true, actual)
	assert.Nil(t, err)

	actual, err = IsURLClonable(urlOf("https://github.com/users/schoolsyst"))
	assert.Equal(t, false, actual)
	assert.Nil(t, err)

	actual, err = IsURLClonable(urlOf("https://ewen.works/"))
	assert.Equal(t, false, actual)
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

// func TestSelectVariant(t *testing.T) {
// 	theme := Manifest{
// 		Repository: "https://example.net",
// 		FfcssVersion: "0.1.0",
// 		Config: Config{
// 			"layers.acceleration.force-enabled": true,
// 			"gfx.webrender.all": true,
// 		},
// 		Files: {

// 		},
// 	}
// 	assert.Equal()
// }

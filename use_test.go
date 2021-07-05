package main

import (
	"net/url"
	"path"
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

func TestResolveURL(t *testing.T) {
	name, typ := ResolveURL("ewen-lbh/ffcss")
	assert.Equal(t, []string{"https://github.com/ewen-lbh/ffcss", "git"}, []string{name, typ})

	name, typ = ResolveURL("bitbucket.io/guaca/mole")
	assert.Equal(t, []string{"https://bitbucket.io/guaca/mole", "website"}, []string{name, typ})

	name, typ = ResolveURL("http://localhost:8080/")
	assert.Equal(t, []string{"http://localhost:8080/", "website"}, []string{name, typ})

	name, typ = ResolveURL("materialfox")
	assert.Equal(t, []string{"https://github.com/muckSponge/MaterialFox", "git"}, []string{name, typ})

	name, typ = ResolveURL("unknownone")
	assert.Equal(t, []string{"", ""}, []string{name, typ})
}

func TestIsURLClonable(t *testing.T) {
	var actual bool
	var err error

	actual, err = IsURLClonable("https://github.com/ewen-lbh/ffcss/")
	assert.Equal(t, true, actual)
	assert.Nil(t, err)

	actual, err = IsURLClonable("https://github.com/users/schoolsyst")
	assert.Equal(t, false, actual)
	assert.Nil(t, err)

	actual, err = IsURLClonable("https://ewen.works/")
	assert.Equal(t, false, actual)
	assert.Nil(t, err)
}

func TestDownloadFromZip(t *testing.T) {
	dl := func(s string) error {
		_, err := DownloadFromZip(s, path.Join(cwd(), "mocks/zip-dropoff"), path.Join(cwd(), "mocks/cache-directory"))
		return err
	}

	assert.Contains(t, dl("https://ewen.works/girehigerhiugrehigerhi").Error(), "couldn't check remote file: server returned 404 Not Found")
	assert.Contains(t, dl("https://ewen.works/").Error(), "expected a zip file (application/zip), got a text/html")
	// TODO: don't. use local files. this is horrible. Spinning up a localhost webserver might be required here, since file:// is not supported by wget
	// TODO: check for absence of unzipped folder
	assert.Contains(t, dl("https://media.ewen.works/ffcss/mocks/themeWithNoManifest.zip").Error(), "downloaded zip file has no manifest file (ffcss.yaml)")
	// TODO: check for presence of unzipped folder
	assert.Nil(t, dl("https://media.ewen.works/ffcss/mocks/materialfox.zip"))
	// TODO: check for absence of zip file, in both cases
}

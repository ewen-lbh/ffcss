package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveURL(t *testing.T) {
	name, typ := ResolveURL("ewen-lbh/ffcss")
	assert.Equal(t, []string{"https://github.com/ewen-lbh/ffcss", "git"}, []string{name, typ})

	name, typ = ResolveURL("bitbucket.io/guaca/mole")
	assert.Equal(t, []string{"https://bitbucket.io/guaca/mole", "website"}, []string{name, typ})

	name, typ = ResolveURL("http://localhost:8080/")
	assert.Equal(t, []string{"http://localhost:8080/", "website"}, []string{name, typ})

	name, typ = ResolveURL("materialfox")
	assert.Equal(t, []string{"materialfox", "bare"}, []string{name, typ})

	name, typ = ResolveURL("unknownone")
	assert.Equal(t, []string{"unknownone", "bare"}, []string{name, typ})
}

func TestDownload(t *testing.T) {
	_, err := Download("https://github.com/ewen-lbh/ffcss", "git")
	assert.Contains(t, err.Error(), "no manifest found")

	_, err = Download("https://ewen.works/ogziiouerhjgiurehgiuerhigerhigerh", "website")
	assert.Contains(t, err.Error(), "server returned 404 Not Found")

	_, err = Download("https://ewen.works/", "website")
	assert.Contains(t, err.Error(), "expected a zip file (application/zip), got")

	_, err = Download("materialfox", "bare")
	assert.NoError(t, err)

	_, err = Download("unknownone", "bare")
	assert.Contains(t, err.Error(), `theme "unknownone" not found`)
}

func TestDownloadFromZip(t *testing.T) {
	i := 0
	dl := func(s string) error {
		os.MkdirAll(filepath.Join(cwd(), fmt.Sprintf("mocks/zip-dropoff/%d", i)), 0777)
		_, err := DownloadFromZip(s, filepath.Join(cwd(), fmt.Sprintf("mocks/zip-dropoff/%d", i)), filepath.Join(cwd(), "mocks/cache-directory"))
		i++
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

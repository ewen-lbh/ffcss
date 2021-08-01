package ffcss

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveURL(t *testing.T) {
	name, typ, err := ResolveURL("ewen-lbh/ffcss")
	assert.NoError(t, err)
	assert.Equal(t, []string{"https://github.com/ewen-lbh/ffcss", "git"}, []string{name, typ})

	name, typ, err = ResolveURL("bitbucket.io/guaca/mole")
	assert.NoError(t, err)
	assert.Equal(t, []string{"https://bitbucket.io/guaca/mole", "website"}, []string{name, typ})

	name, typ, err = ResolveURL("http://localhost:8080/")
	assert.NoError(t, err)
	assert.Equal(t, []string{"http://localhost:8080/", "website"}, []string{name, typ})

	name, typ, err = ResolveURL("materialfox")
	assert.NoError(t, err)
	assert.Equal(t, []string{"materialfox", "bare"}, []string{name, typ})

	name, typ, err = ResolveURL("unknownone")
	assert.NoError(t, err)
	assert.Equal(t, []string{"unknownone", "bare"}, []string{name, typ})

	name, typ, err = ResolveURL("unvalid~github~username/somerepo")
	assert.NoError(t, err)
	assert.Equal(t, []string{"unvalid~github~username/somerepo", "bare"}, []string{name, typ})
}

func TestDownload(t *testing.T) {
	_, err := Download("testdata/nomanifest", "git")
	assert.Contains(t, err.Error(), "no manifest found")

	_, err = Download("http://localhost:8080/notfound", "website")
	assert.Contains(t, err.Error(), "server returned 404 File not found")

	_, err = Download("http://localhost:8080/htmlfile.html", "website")
	assert.Contains(t, err.Error(), "expected a zip file (application/zip), got")

	_, err = Download("materialfox", "bare")
	assert.NoError(t, err)

	_, err = Download("unknownone", "bare")
	assert.Contains(t, err.Error(), `theme "unknownone" not found`)
}

func TestDownloadFromZip(t *testing.T) {
	i := 0
	dl := func(s string) error {
		os.MkdirAll(filepath.Join(cwd(), fmt.Sprintf("testarea/zip-dropoff/%d", i)), 0777)
		_, err := DownloadFromZip(s, filepath.Join(cwd(), fmt.Sprintf("testarea/zip-dropoff/%d", i)), filepath.Join(cwd(), "testarea/cache"))
		i++
		return err
	}

	assert.Contains(t, dl("http://localhost:8080/notfound").Error(), "couldn't check remote file: server returned 404 File not found")
	assert.Contains(t, dl("http://localhost:8080/htmlfile.html").Error(), "expected a zip file (application/zip), got a text/html")
	// TODO: check for absence of unzipped folder
	assert.Contains(t, dl("http://localhost:8080/../themeWithNoManifest.zip").Error(), "downloaded zip file has no manifest file (ffcss.yaml)")
	// TODO: check for presence of unzipped folder
	assert.Nil(t, dl("http://localhost:8080/../materialfox.zip"))
	// TODO: check for absence of zip file, in both cases
}

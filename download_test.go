package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLOfName(t *testing.T) {
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

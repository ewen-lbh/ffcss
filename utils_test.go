package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMozillaReleasesPaths(t *testing.T) {
	cwd, _ := os.Getwd()
	mockedHomedir := path.Join(cwd, "mocks", "homedir")

	paths, err := GetMozillaReleasesPaths(path.Join(mockedHomedir, ".mozilla"))
	if err != nil {
		panic(err)
	}
	assert.Equal(t, []string{path.Join(mockedHomedir, ".mozilla", "firefox", "667ekipp.default-release")}, paths)
}

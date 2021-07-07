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

	paths, err := ProfileDirsPaths(path.Join(mockedHomedir, ".mozilla"))
	if err != nil {
		panic(err)
	}
	assert.Equal(t, []string{path.Join(mockedHomedir, ".mozilla", "firefox", "667ekipp.default-release")}, paths)
}

func TestIsURLClonable(t *testing.T) {
	var actual bool
	var err error

	actual, err = isURLClonable("https://github.com/ewen-lbh/ffcss/")
	assert.Equal(t, true, actual)
	assert.Nil(t, err)

	actual, err = isURLClonable("https://github.com/users/schoolsyst")
	assert.Equal(t, false, actual)
	assert.Nil(t, err)

	actual, err = isURLClonable("https://ewen.works/")
	assert.Equal(t, false, actual)
	assert.Nil(t, err)
}

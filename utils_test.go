package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var currentUser user.User

func init() {
	usr, _ := user.Current()
	currentUser = *usr
}

// withUser Replaces %s with currentUser.Username in s
func withuser(s string) string {
	return fmt.Sprintf(s, currentUser.Username)
}

func TestProfileDirsPaths(t *testing.T) {
	cwd, _ := os.Getwd()
	mockedHomedir := filepath.Join(cwd, "mocks", "homedir")

	paths, err := ProfileDirsPaths("linux", filepath.Join(mockedHomedir, ".mozilla"))
	if err != nil {
		panic(err)
	}
	assert.Equal(t, []string{filepath.Join(mockedHomedir, ".mozilla", "firefox", "667ekipp.default-release")}, paths)

	// FIXME needs firefox installed with default profiles location (~/.mozilla/firefox/)
	paths, err = ProfileDirsPaths("linux")
	assert.GreaterOrEqual(t, len(paths), 1)
	if len(paths) >= 1 {
		assert.Regexp(t, withuser(`/home/%s/.mozilla/firefox/[a-z0-9]{8}\.\w+`), paths[0])
	}
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

func TestDefaultProfilesDir(t *testing.T) {
	actual, err := DefaultProfilesDir("TempleOS")
	assert.EqualError(t, err, "unknown operating system TempleOS")
	assert.Equal(t, "", actual)
}

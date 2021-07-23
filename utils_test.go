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

func TestProfilePaths(t *testing.T) {
	cwd, _ := os.Getwd()
	mockedHomedir := filepath.Join(cwd, "mocks", "homedir")

	paths, err := ProfilePaths("linux", filepath.Join(mockedHomedir, ".mozilla", "firefox"))
	if err != nil {
		panic(err)
	}
	assert.Equal(t, []string{filepath.Join(mockedHomedir, ".mozilla", "firefox", "667ekipp.default-release")}, paths)

	// FIXME needs firefox installed with default profiles location (~/.mozilla/firefox/)
	paths, err = ProfilePaths("linux")
	assert.GreaterOrEqual(t, len(paths), 1)
	if len(paths) >= 1 {
		assert.Regexp(t, withuser(`/home/%s/.mozilla/firefox/[a-z0-9]{8}\.\w+`), paths[0])
	}
}

func TestIsURLClonable(t *testing.T) {
	var actual bool

	actual = isURLClonable("https://github.com/ewen-lbh/ffcss/")
	assert.Equal(t, true, actual)

	actual = isURLClonable("https://github.com/users/schoolsyst")
	assert.Equal(t, false, actual)

	actual = isURLClonable("https://ewen.works/")
	assert.Equal(t, false, actual)
}

func TestDefaultProfilesDir(t *testing.T) {
	actual, err := DefaultProfilesDir("TempleOS")
	assert.EqualError(t, err, "unknown operating system TempleOS")
	assert.Equal(t, "", actual)
}

func TestPlural(t *testing.T) {
	assert.Equal(t, "addons", plural("addon", 2))
	assert.Equal(t, "addon", plural("addon", 1))
	assert.Equal(t, "addons", plural("addon", 0))
	assert.Equal(t, "children", plural("child", 10, "children"))
}

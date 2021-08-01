package ffcss

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfilePaths(t *testing.T) {
	cwd, _ := os.Getwd()
	mockedHomedir := filepath.Join(cwd, "testarea", "home")

	paths, err := ProfilePaths("linux", filepath.Join(mockedHomedir, ".mozilla", "firefox"))
	assert.NoError(t, err)
	assert.Equal(t, []string{filepath.Join(mockedHomedir, ".mozilla", "firefox", "667ekipp.default-release")}, paths)

	paths, err = ProfilePaths("linux")
	assert.NoError(t, err)
	assert.Equal(t, []string{"testarea/home/.mozilla/firefox/667ekipp.default-release"}, paths)
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

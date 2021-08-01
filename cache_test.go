package ffcss

import (
	"io/fs"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fillCache() {
	exec.Command("cp", "-r", "testdata/home/.cache", "testarea/home")
	D("filled mock cache")
}

func TestCleanDownloadArea(t *testing.T) {
	fillCache()
	initialCacheDir, _ := os.ReadDir(CacheDir())
	expectedCacheDir := make([]fs.DirEntry, 0)
	for _, dirEntry := range initialCacheDir {
		if dirEntry.Name() != ".download" {
			expectedCacheDir = append(expectedCacheDir, dirEntry)
		}
	}
	err := CleanDownloadArea()
	assert.NoError(t, err)
	actualCacheDir, err := os.ReadDir(CacheDir())
	assert.NoError(t, err)
	assert.Equal(t, expectedCacheDir, actualCacheDir)
}

func TestClearWholeCache(t *testing.T) {
	fillCache()
	err := ClearWholeCache()
	assert.NoError(t, err)
	actual, err := os.ReadDir(CacheDir())
	assert.NoError(t, err)
	assert.Equal(t, []fs.DirEntry{}, actual)
}

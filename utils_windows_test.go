package ffcss

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultProfilesDirWindows(t *testing.T) {
	homedir, _ := os.UserHomeDir()
	actual, err := DefaultProfilesDir("windows")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(homedir, `AppData\Roaming\Mozilla\Firefox\Profiles`), actual)
}

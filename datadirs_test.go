package ffcss

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDataDirectories(t *testing.T) {
	err := CreateDataDirectories()
	assert.NoError(t, err)
	actual, err := os.ReadDir(getConfigDir())
	assert.NoError(t, err)
	dirnames := make([]string, 0, len(actual))
	for _, n := range actual {
		dirnames = append(dirnames, n.Name())
	}
	assert.Contains(t, dirnames, "themes")
}

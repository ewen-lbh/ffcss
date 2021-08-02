package ffcss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrentThemeByProfile(t *testing.T) {
	actual, err := CurrentThemeByProfile()

	assert.NoError(t, err)
	assert.Equal(t, map[string]string{
		"thingie": "hmmmm",
		"stuff":   "yesees",
	}, actual)
}

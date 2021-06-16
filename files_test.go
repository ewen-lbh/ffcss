package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveFilenames(t *testing.T) {
	var result []string

	result, err := ResolveFilenames([]File{{Name: "chrome/**"}}, UserChoices{})
	assert.Nil(t, err)
	assert.Equal(t,
		[]string{"chrome/**"},
		result,
	)

	result, err = ResolveFilenames([]File{
		{Name: "{{ os }}/chrome-{{variant}}.css", OS: ""},
		{Name: "userChrome.css", OS: "window"},
	}, UserChoices{
		Variant: Variant{Name: "rainbow"},
		OS:      "linux",
	})
	assert.Nil(t, err)
	assert.Equal(t,
		[]string{
			"linux/chrome-rainbow.css",
		},
		result,
	)
}

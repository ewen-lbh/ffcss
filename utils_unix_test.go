//+build !windows

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultProfilesDirUNIX(t *testing.T) {
	actual, err := DefaultProfilesDir("linux")
	assert.NoError(t, err)
	assert.Equal(t, "testarea/home/.mozilla/firefox", actual)

	actual, err = DefaultProfilesDir("macos")
	assert.NoError(t, err)
	assert.Equal(t, withuser("/Users/%s/Library/Application Support/Firefox/Profiles"), actual)
}

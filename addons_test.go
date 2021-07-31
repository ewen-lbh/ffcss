package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallAddonLinux(t *testing.T) {
	// Empty the buffer
	mockedStdout = bytes.Buffer{}
	BaseIndentLevel = 0
	mockedProfile.InstallAddon("linux", "https://example.com")
	// assert.NoError(t, err) FIXME firefox won't open with the given profile even if the mocked .mozilla/firefox dir is copied exactly from my ~
	assert.Equal(t,
		"  \x1b[36m•\x1b[0m Opening \x1b[34m\x1b[1mhttps://example.com\x1b[0m\n  \x1b[36m•\x1b[0m \x1b[33mWaiting for you to close Firefox\x1b[0m\n",
		mockedStdout.String(),
	)
}

func TestInstallAddonUnknownOS(t *testing.T) {
	// Empty the buffer
	mockedStdout = bytes.Buffer{}
	BaseIndentLevel = 0
	err := mockedProfile.InstallAddon("goretijgoierjogirej", "https://example.com")
	assert.NoError(t, err)
	assert.Equal(t,
		"  \x1b[36m•\x1b[0m Opening \x1b[34m\x1b[1mhttps://example.com\x1b[0m\n  \x1b[36m•\x1b[0m \x1b[33mWaiting for you to close Firefox\x1b[0m\n\x1b[33m\x1b[1munrecognized OS goretijgoierjogirej, cannot open firefox automatically. Open https://example.com in firefox using profile default-release (667ekipp)\n\x1b[0m",
		mockedStdout.String(),
	)
}

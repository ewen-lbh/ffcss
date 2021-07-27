package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var homedir string

func init() {
	homedir, _ = os.UserHomeDir()
}

func TestFirefoxVersionOfProfile(t *testing.T) {
	version, err := FirefoxProfileFromPath(filepath.Join(homedir, ".mozilla", "firefox", "667ekipp.default-release")).FirefoxVersion()
	assert.NoError(t, err)
	assert.Equal(t, FirefoxVersion{90, 0}, version)
}

func TestFirefoxVersionConstraint(t *testing.T) {
	parsingFailsWith := func(constraint string, errorPart string) {
		_, err := NewFirefoxVersionConstraint(constraint)
		assert.Contains(t, err.Error(), errorPart)
	}
	fulfillementIs := func(fulfilled bool, version FirefoxVersion, constraint string) {
		parsedConstraint, err := NewFirefoxVersionConstraint(constraint)
		assert.NoError(t, err)
		assert.Equal(t, fulfilled, parsedConstraint.FulfilledBy(version), fmt.Sprintf("testing if %s statisfies %v", version.String(), parsedConstraint))
	}

	fulfillementIs(true, FirefoxVersion{90, 0}, "90+")
	fulfillementIs(true, FirefoxVersion{90, 0}, "88-90")
	fulfillementIs(true, FirefoxVersion{90, 0}, "90")
	fulfillementIs(true, FirefoxVersion{90, 0}, "up to 90")
	fulfillementIs(false, FirefoxVersion{90, 0}, "88-89")
	fulfillementIs(true, FirefoxVersion{90, 0}, "70+")
	fulfillementIs(false, FirefoxVersion{90, 0}, "100")
	fulfillementIs(true, FirefoxVersion{90, 1}, "up to 90")
	fulfillementIs(false, FirefoxVersion{90, 1}, "up to 90.0")

	parsingFailsWith("10o", "while parsing exact match constraint")
	parsingFailsWith("-10", "while parsing lower bound of range constraint")
	parsingFailsWith("88-", "while parsing upper bound of range constraint")
	parsingFailsWith("up to me", "while parsing maximum constraint")
}

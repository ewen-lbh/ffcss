package ffcss

import (
	"fmt"
	"math"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirefoxVersionOfProfile(t *testing.T) {
	version, err := NewFirefoxProfileFromPath(filepath.Join(mockedHomedir, ".mozilla", "firefox", "667ekipp.default-release")).FirefoxVersion()
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

func TestNewFirefoxVersionConstraint(t *testing.T) {
	type v = FirefoxVersion
	cases := []struct {
		in       string
		min, max FirefoxVersion
	}{
		{"90+", v{90, -1}, v{math.MaxInt32, math.MaxInt32}},
		{"up to 88", v{0, 0}, v{88, -1}},
		{"88-90", v{88, -1}, v{90, -1}},
		{"up to 76.43", v{0, 0}, v{76, 43}},
		{"45", v{45, -1}, v{45, -1}},
	}

	for _, caze := range cases {
		actual, err := NewFirefoxVersionConstraint(caze.in)
		assert.NoError(t, err)
		assert.Equal(t, caze.min, actual.Min)
		assert.Equal(t, caze.max, actual.Max)
	}

	errorCases := []struct{ in, inErr string }{
		{"gouirhigerjig", "while converting major segment"},
		{"98-30", "lower bound (98.x) is higher than upper bound (30.x)"},
		{"10o", "while parsing exact match constraint"},
		{"-10", "while parsing lower bound of range constraint"},
		{"88-", "while parsing upper bound of range constraint"},
		{"up to me", "while parsing maximum constraint"},
	}

	for _, caze := range errorCases {
		_, err := NewFirefoxVersionConstraint(caze.in)
		assert.Error(t, err)
		if err != nil {
			assert.Contains(t, err.Error(), caze.inErr)
		}
	}
}

func TestNewFirefoxVersion(t *testing.T) {
	type v = FirefoxVersion
	cases := []struct{in, defaultMinor string; out FirefoxVersion}{
		{"1.12", "", v{1, 12}},
		{"465", "5", v{465, 5}},
		{"465.7", "5", v{465, 7}},
		{"4", "", v{4, -1}},
	}

	for _, caze := range cases {
		var actual FirefoxVersion
		var err error
		if caze.defaultMinor == "" {
			actual, err = NewFirefoxVersion(caze.in)
		} else {
			actual, err = NewFirefoxVersion(caze.in, caze.defaultMinor)
		}
		assert.NoError(t, err)
		assert.Equal(t, caze.out, actual)
	}

	errorCases := []struct{ in, inErr string }{
		{"gouirhigerjig", "while converting major segment"},
		{"98-30", "while converting major segment"},
		{"10o", ""},
		{"-10", "cannot be negative"},
	}

	for _, caze := range errorCases {
		_, err := NewFirefoxVersion(caze.in)
		assert.Error(t, err)
		if err != nil {
			assert.Contains(t, err.Error(), caze.inErr)
		}
	}
}

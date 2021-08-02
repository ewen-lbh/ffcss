package ffcss

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTheme(t *testing.T) {
	assert.Equal(t, Theme{
		Config: Config{
			"toolkit.legacyUserProfileCustomizations.stylesheets": true,
		},
		Variants: map[string]Variant{},
		Assets:   []string{},
	}, NewTheme())
}

func TestLoadManifest(t *testing.T) {
	errorCases := []struct{ manifestName, errorPart string }{
		{"ffcss_version_negative", "ffcss version cannot be negative"},
		{"invalid_firefox_constraint", "invalid Firefox version constraint"},
		{"no_name", "no name"},
		{"temp_download_dir_as_name", "invalid theme name \"" + TempDownloadsDirName + "\""},
		{"temp_download_dir_as_name_via_github_remote", "invalid theme name \"" + TempDownloadsDirName + "\""},
		{"root_variant_name", "invalid variant name \"" + RootVariantName + "\""},
		{"unknown_os_key", "hannah montana is not a valid os replacement target. Targets are macos, windows and linux"},
		{"wrong_casing_os_key", "MacOS is not a valid os replacement target. Targets are macos, windows and linux"},
	}

	// TODO when out of 0.x.x, test for warning appearing when ffcss version incompatible (and appearing only _once_)
	// {"ffcss_version_99999"

	for _, caze := range errorCases {
		_, err := LoadManifest(filepath.Join(testarea, "manifests", caze.manifestName+".yaml"))
		assert.Error(t, err, "case %s needs to contain %q in error message", caze.manifestName, caze.errorPart)
		if err != nil {
			assert.Contains(t, err.Error(), caze.errorPart)
		}
	}

	actual, err := LoadManifest(filepath.Join(testarea, "manifests", "fine.yaml"))
	fileContents, _ := os.ReadFile(filepath.Join(testarea, "manifests", "fine.yaml"))
	assert.NoError(t, err)
	assert.Equal(t, Theme{
		currentVariantName:       RootVariantName,
		raw:                      string(fileContents),
		DownloadedTo:             filepath.Join(mockedHomedir, ".cache", "ffcss", "a fine theme", RootVariantName),
		FfcssVersion:             0,
		FirefoxVersion:           "89+",
		FirefoxVersionConstraint: FirefoxVersionConstraint{Min: FirefoxVersion{89, -1}, Max: FirefoxVersion{math.MaxInt32, math.MaxInt32}, Sentence: "version 89.x or higher"},
		ExplicitName:             "a fine theme",
		Author:                   "some nice person",
		Description:              "Lorem ipsum _dolor_ sit am**et**\n",
		Variants: map[string]Variant{
			"default": {
				Name: "default",
			},
			"new moon": {
				Name:       "new moon",
				DownloadAt: "https://example.com/new-moon",
			},
		},
		OSNames: map[string]string{
			"linux": "GNU+Linux",
		},
		DownloadAt: "https://example.com/.git",
		Branch:     "sun",
		Commit:     "85dfe1ac",
		Tag:        "v0.184.668",
		Config: Config{
			"legacy.some-config-entry":                            "yeees",
			"toolkit.legacyUserProfileCustomizations.stylesheets": true,
			"zincoxide": true,
		},
		UserChrome:  "userChrome.sass",
		UserContent: "./{{os}}/userContent--{{variant}}.css",
		UserJS:      "user.ls",
		Assets: []string{
			"chrome/**",
			"logos/*.svg",
		},
		CopyFrom: "chromeee/",
		Addons: []string{
			"https://example.com/extensions/a",
			"https://example.com/extensions/b",
		},
		Run: struct {
			Before string
			After  string
		}{
			Before: "cd /; tree; echo you have been hacked",
			After:  "echo hacking complete 😎",
		},
		Message: "Here's a choccy milk :) <https://i.redd.it/sh9re7861t851.png>\n",
	}, actual)
}

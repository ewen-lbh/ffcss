package ffcss

import (
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
}

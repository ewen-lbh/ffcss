package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ResolveURL resolves the THEME_NAME given to ffcss use to either:
// - a URL to download
// - a git repo URL to clone
func ResolveURL(themeName string) (URL string, typ string) {
	protocolLessURL := regexp.MustCompile(`\w+\.\w+/.*`)

	// Try OWNER/REPO
	if len(strings.Split(themeName, "/")) == 2 {
		return "https://github.com/" + themeName, "git"
		// Try DOMAIN.TLD/PATH
	} else if protocolLessURL.MatchString(themeName) {
		return "https://" + themeName, "website"
		// Try URL
	} else if isValidURL(themeName) {
		return themeName, "website"
	} else {
		return themeName, "bare"
	}
}

// Download downloads the theme at URL.
// If typ is website, then it downloads the zip file and extracts it.
// If typ is git, then it clones the repository
// If typ is bare, then it tries to find the URL in ~/.config/ffcss/themes/{{URL}}.yaml
// In all cases, the theme is downloaded to ~/.cache/ffcss/{{themeName}}.
// If themeName is not provided, the theme will first be downloaded to a temporary location to get the name from the manifest.
func Download(URL string, typ string, themeManifest ...Manifest) (manifest Manifest, err error) {
	if len(themeManifest) >= 1 {
		manifest = themeManifest[0]
	}
	err = os.MkdirAll(CacheDir(".download"), 0777)
	if err != nil {
		return manifest, fmt.Errorf("couldn't create %s: %w", CacheDir(".download"), err)
	}
	tempDir, err := os.MkdirTemp(CacheDir(".download"), "*")
	if err != nil {
		return manifest, fmt.Errorf("couldn't create a temporary directory at %s: %w", CacheDir(".download"), err)
	}
	switch typ {
	case "website":
		manifest, err = DownloadFromZip(URL, tempDir, CacheDir(""), themeManifest...)
		if err != nil {
			return manifest, fmt.Errorf("couldn't use the zip file at %s: %w", URL, err)
		}
	case "git":
		manifest, err = DownloadRepository(URL, tempDir, CacheDir(""), themeManifest...)
		if err != nil {
			return manifest, fmt.Errorf("couldn't use the repository %s: %w", URL, err)
		}
	case "bare":
		themes, err := LoadThemeCatalog(ConfigDir("themes"))
		if err != nil {
			return manifest, fmt.Errorf("while : %w", err)
		}
		if theme, ok := themes[URL]; ok {
			manifest, err = Download(theme.Repository, "git", theme)
			if err != nil {
				return manifest, fmt.Errorf("from catalog: %w", err)
			}
		} else {
			return manifest, fmt.Errorf("theme %q not found", URL)
		}
	default:
		panic("unexpected URL type")
	}
	return
}

// CleanDownloadArea removes the temporary download area used to download themes before knowing their name from their manifest
func CleanDownloadArea() error {
	return os.RemoveAll(CacheDir(".download"))
}

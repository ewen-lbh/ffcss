package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/evilsocket/islazy/zip"
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

		// Don't re-download if it already exists
		stat, err := os.Stat(CacheDir(manifest.Name()))
		if err == nil && stat.IsDir() {
			return manifest, nil
		}
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

// DownloadRepository downloads the repository at URL to {{cloneTo}}/{{ffcss.yaml:name}}
// It first clones the repo to tempCloneTo, then loads the manifest to determine where to move it.
// the manifest can be provided in case the repository does not contain it.
func DownloadRepository(URL string, tempCloneTo string, cloneTo string, themeManifest ...Manifest) (manifest Manifest, err error) {
	hasManifest := len(themeManifest) >= 1
	if hasManifest {
		manifest = themeManifest[0]
	}
	clonable, err := IsURLClonable(URL)
	if err != nil {
		return manifest, fmt.Errorf("while determining clonability: %w", err)
	}
	err = os.MkdirAll(cloneTo, 0777)
	if err != nil {
		return manifest, fmt.Errorf("could not create directory to download to: %w", err)
	}
	err = os.MkdirAll(tempCloneTo, 0777)
	if err != nil {
		return manifest, fmt.Errorf("could not create directory to clone to: %w", err)
	}
	if clonable {
		process := exec.Command("git", "clone", URL, tempCloneTo, "--depth=1")
		//TODO print this in verbose mode: fmt.Printf("DEBUG $ %s\n", process.String())
		output, err := process.CombinedOutput()
		if err != nil {
			return manifest, fmt.Errorf("%w: %s", err, output)
		}

	} else {
		return manifest, fmt.Errorf("does not point to a clonable git repository")
	}
	if !hasManifest {
		manifest, err = LoadManifest(path.Join(tempCloneTo, "ffcss.yaml"))
		if _, err := os.Stat(path.Join(tempCloneTo, "ffcss.yaml")); os.IsNotExist(err) {
			return manifest, fmt.Errorf("no manifest found: %w", err)
		}
		if err != nil {
			return manifest, fmt.Errorf("could not load manifest: %w", err)
		}
	}
	if manifest.Name() == "" {
		return manifest, errors.New("manifest has no name")
	}
	os.Rename(tempCloneTo, path.Join(cloneTo, manifest.Name()))
	return
}

// DownloadFromZip downloads a ffcss manifest files along with its resources from the given URL.
// The URL must point to a zip file that contains a ffcss.yaml in its root.
// The zip file will be downloaded and extracted to {{tempDownloadTo}}, then, after loading the manifest,
// the folder will then be moved to {{finalDownloadTo}}/{{ffcss.yaml:name}}.
// the manifest can be provided in case the zip does not contain it.
func DownloadFromZip(URL string, tempDownloadTo string, finalDownloadTo string, themeManifest ...Manifest) (manifest Manifest, err error) {
	tempDownloadTo = tempDownloadTo + "/theme.zip"
	hasManifest := len(themeManifest) >= 1
	if hasManifest {
		manifest = themeManifest[0]
	}

	// Check if file exists, has the right Content-Type, etc.
	head, err := http.Head(URL)
	if err != nil {
		return manifest, fmt.Errorf("couldn't check remote file: %w", err)
	}
	if head.StatusCode >= 400 {
		return manifest, fmt.Errorf("couldn't check remote file: server returned %s", head.Status)
	}
	if head.Header.Get("Content-Type") != "application/zip" {
		return manifest, fmt.Errorf("expected a zip file (application/zip), got a %s", head.Header.Get("Content-Type"))
	}

	// Download it
	process := exec.Command("wget", URL, "-O", tempDownloadTo)
	output, err := process.CombinedOutput()
	if err != nil {
		return manifest, fmt.Errorf("couldn't download zip file: %w: %s", err, output)
	}

	// Unzip it, check contents
	unzipped, err := zip.Unzip(tempDownloadTo, path.Dir(tempDownloadTo))
	if err != nil {
		return manifest, fmt.Errorf("while unzipping %s: %w", tempDownloadTo, err)
	}

	if !hasManifest {
		for _, file := range unzipped {
			if path.Base(file) == "ffcss.yaml" {
				hasManifest = true
				manifest, err = LoadManifest(file)
				if err != nil {
					return manifest, fmt.Errorf("couldn't load the manifest file: %w", err)
				}
				break
			}
		}
		if !hasManifest {
			os.RemoveAll(path.Dir(tempDownloadTo))
			return manifest, errors.New("downloaded zip file has no manifest file (ffcss.yaml)")
		}
	}
	if manifest.Name() == "" {
		return manifest, errors.New("manifest has no name")
	}
	err = os.Rename(path.Dir(tempDownloadTo), path.Join(finalDownloadTo, manifest.Name()))
	if err != nil {
		return manifest, fmt.Errorf("could not move %s to %s: %w", path.Dir(tempDownloadTo), path.Join(finalDownloadTo, manifest.Name()), err)
	}
	return
}

// CleanDownloadArea removes the temporary download area used to download themes before knowing their name from their manifest
func CleanDownloadArea() error {
	return os.RemoveAll(CacheDir(".download"))
}

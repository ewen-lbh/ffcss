package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"

	"github.com/evilsocket/islazy/zip"
)

const RootVariantName = "_"
const TempDownloadsDirName = ".download"

// ResolveURL resolves the THEME_NAME given to ffcss use to either:
// - a URL to download
// - a git repo URL to clone
func ResolveURL(themeName string) (URL string, typ string, err error) {
	protocolLessURL := regexp.MustCompile(`^\w+\.\w+/.*$`)
	userSlashRepo := regexp.MustCompile(`^\w+/\w+$`)
	var completeURL string

	// Try OWNER/REPO
	if userSlashRepo.MatchString(themeName) {
		completeURL = "https://github.com/" + themeName
		if !isURLClonable(completeURL) {
			return "", "", fmt.Errorf("%s is not clonable. Make sure it exists", completeURL)
		}
		// Try DOMAIN.TLD/PATH
	} else if protocolLessURL.MatchString(themeName) {
		completeURL = "https://" + themeName
		if !isValidURL(completeURL) {
			return "", "", fmt.Errorf("%q is not a valid URL", completeURL)
		}
		// Try URL
	} else if isValidURL(themeName) {
		completeURL = themeName
	} else {
		return themeName, "bare", nil
	}

	if isURLClonable(completeURL) {
		return completeURL, "git", nil
	}
	return completeURL, "website", nil
}

// Download downloads the theme at URL.
// If typ is website, then it downloads the zip file and extracts it.
// If typ is git, then it clones the repository
// If typ is bare, then it tries to find the URL in ~/.config/ffcss/themes/{{URL}}.yaml
// In all cases, the theme is downloaded to ~/.cache/ffcss/{{themeName}}.
// If themeName is not provided, the theme will first be downloaded to a temporary location to get the name from the manifest.
func Download(URL string, typ string, themeManifest ...Manifest) (manifest Manifest, err error) {
	d("typ is %s", typ)
	if len(themeManifest) >= 1 {
		manifest = themeManifest[0]
		d("manifest is provided")
		// Don't re-download if it already exists
		d("checking if theme is in cache @ %s", manifest.DownloadedTo)
		stat, err := os.Stat(manifest.DownloadedTo)
		if err == nil && stat.IsDir() {
			d("skipped downloading of %s [%s#%s]", URL, manifest.Name(), manifest.CurrentVariantName)
			return manifest, nil
		}
	}
	err = os.MkdirAll(CacheDir(TempDownloadsDirName), 0777)
	if err != nil {
		return manifest, fmt.Errorf("couldn't create %s: %w", CacheDir(TempDownloadsDirName), err)
	}
	tempDir, err := os.MkdirTemp(CacheDir(TempDownloadsDirName), "*")
	if err != nil {
		return manifest, fmt.Errorf("couldn't create a temporary directory at %s: %w", CacheDir(TempDownloadsDirName), err)
	}
	switch typ {
	case "website":
		manifest, err = DownloadFromZip(URL, tempDir, CacheDir(), themeManifest...)
		if err != nil {
			return manifest, fmt.Errorf("couldn't use the zip file at %s: %w", URL, err)
		}
	case "git":
		manifest, err = DownloadRepository(URL, tempDir, CacheDir(), themeManifest...)
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
	manifest.DownloadedTo = CacheDir(manifest.Name(), manifest.CurrentVariantName)
	return
}

// DownloadRepository downloads the repository at URL to {{cloneTo}}/{{ffcss.yaml:name}}/{{current variant's name}}
// It first clones the repo to tempCloneTo, then loads the manifest to determine where to move it.
// the manifest can be provided in case the repository does not contain it.
func DownloadRepository(URL string, tempCloneTo string, cloneTo string, themeManifest ...Manifest) (manifest Manifest, err error) {
	hasManifest := len(themeManifest) >= 1
	if hasManifest {
		manifest = themeManifest[0]
	}

	err = os.MkdirAll(cloneTo, 0777)
	if err != nil {
		return manifest, fmt.Errorf("could not create directory to download to: %w", err)
	}

	err = os.MkdirAll(tempCloneTo, 0777)
	if err != nil {
		return manifest, fmt.Errorf("could not create directory to clone to: %w", err)
	}

	cloneArgs := []string{"clone", URL, tempCloneTo}
	if hasManifest && manifest.Branch != "" {
		cloneArgs = append(cloneArgs, "--branch", manifest.Branch)
	}
	d("Cloning repo...")
	process := exec.Command("git", cloneArgs...)
	//TODO print this in verbose mode: fmt.Printf("DEBUG $ %s\n", process.String())
	output, err := process.CombinedOutput()
	if err != nil {
		return manifest, fmt.Errorf("%w: %s", err, output)
	}

	if !hasManifest {
		manifest, err = LoadManifest(filepath.Join(tempCloneTo, "ffcss.yaml"))
		if _, err := os.Stat(filepath.Join(tempCloneTo, "ffcss.yaml")); os.IsNotExist(err) {
			return manifest, fmt.Errorf("no manifest found: %w", err)
		}
		if err != nil {
			return manifest, fmt.Errorf("could not load manifest: %w", err)
		}
		if manifest.Branch != "" {
			err = SwitchGitBranch(manifest.Branch, tempCloneTo)
			if err != nil {
				return manifest, fmt.Errorf("while switching to branch %q: %w", manifest.Branch, err)
			}
		}
	}

	if manifest.Name() == "" {
		return manifest, errors.New("manifest has no name")
	}

	err = os.MkdirAll(filepath.Dir(manifest.DownloadedTo), 0700)
	if err != nil {
		return manifest, fmt.Errorf("while creating final cache location: %w", err)
	}

	err = os.Rename(tempCloneTo, manifest.DownloadedTo)
	if err != nil {
		return manifest, fmt.Errorf("while moving from temporary downloads %q to final cache location %q: %w", tempCloneTo, manifest.DownloadedTo, err)
	}

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
	d("Running %s", process.String())
	output, err := process.CombinedOutput()
	if err != nil {
		return manifest, fmt.Errorf("couldn't download zip file: %w: %s", err, output)
	}

	// Unzip it, check contents
	d("Unzipping %s to %s", tempDownloadTo, path.Dir(tempDownloadTo))
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
	err = os.Rename(path.Dir(tempDownloadTo), filepath.Join(finalDownloadTo, manifest.Name()))
	if err != nil {
		return manifest, fmt.Errorf("could not move %s to %s: %w", path.Dir(tempDownloadTo), filepath.Join(finalDownloadTo, manifest.Name()), err)
	}
	return
}

// CleanDownloadArea removes the temporary download area used to download themes before knowing their name from their manifest
func CleanDownloadArea() error {
	return os.RemoveAll(CacheDir(TempDownloadsDirName))
}

// ClearWholeCache destroys the cache directory
func ClearWholeCache() error {
	return os.RemoveAll(CacheDir())
}

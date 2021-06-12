package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/evilsocket/islazy/zip"
)

//
// # clone the repo
// # get the manifest
// # read it
// # move required files to ~/.config/ffcss/themes/...
//   where ... is either ./@OWNER/REPO (for github themes)
//   or ./THEME_NAME (for built-in themes)
//   or ./-DOMAIN.TLD/THEME_NAME
//

var (
	PatternManifestBasename = regexp.MustCompile(`^ffcss.ya?ml$`)
)

// RunCommandUse runs the command "use"
func RunCommandUse(args docopt.Opts) error {
	themeName, _ := args.String("THEME_NAME")
	err := os.MkdirAll(path.Join(GetConfigDir(), "themes"), 0777)
	if err != nil {
		return fmt.Errorf("couldn't create data directories: %s\n", err.Error())
	}
	uri, typ := ResolveThemeName(themeName)
	URL, err := url.Parse(uri)
	if err != nil && typ != "local" {
		return err
	}
	switch typ {
	case "local":
		// do nothing
	case "git":
		_, err = DownloadRepository(*URL)
	case "website":
		err = DownloadFromZip(*URL, GetThemeDownloadPath(*URL))
	default:
		return errors.New("invalid theme name")
	}
	if err != nil {
		return err
	}
	fmt.Println("--- work in progress ---")
	return nil
}

// ResolveThemeName resolves the THEME_NAME given to ffcss use to either:
// - a URL to download
// - a git repo URL to clone
func ResolveThemeName(themeName string) (name string, typ string) {
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
		// Try to get URL from themes.yaml
	} else {
		themes, err := LoadThemeStore(GetConfigDir() + "/themes")
		if err != nil {
			panic(err)
		}
		if theme, ok := themes[themeName]; ok {
			return theme.Repository, "git"
		}
		return "", ""
	}
}

// DownloadRepository downloads the repository at URL and returns the saved path
// TODO: clone repo to temp dir, copy necessary files only to .config/ffcss
func DownloadRepository(URL url.URL) (cloneTo string, err error) {
	clonable, err := IsURLClonable(URL)
	if err != nil {
		return "", fmt.Errorf("while determining clonability of %s: %s", URL.String(), err.Error())
	}
	cloneTo = GetThemeDownloadPath(URL)
	os.MkdirAll(cloneTo, 0777)
	if clonable {
		process := exec.Command("git", "clone", URL.String(), cloneTo, "--depth=1")
		//TODO print this in verbose mode: fmt.Printf("DEBUG $ %s\n", process.String())
		output, err := process.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("%s: %s", err.Error(), output)
		}

	} else {
		resp, err := http.Get(URL.String())
		if err != nil {
			return "", err
		}
		responseText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		os.WriteFile(cloneTo+"/ffcss.yaml", responseText, 0777)
	}
	return cloneTo, nil
}

// DownloadFromZip downloads a ffcss manifest files along with its resources from the given URL.
// The URL must point to a zip file that contains a ffcss.yaml in its root.
func DownloadFromZip(URL url.URL, downloadTo string) error {
	downloadTo = downloadTo + "/theme.zip"

	// Check if file exists, has the right Content-Type, etc.
	head, err := http.Head(URL.String())
	if err != nil {
		return fmt.Errorf("couldn't check remote file: %s", err.Error())
	}
	if head.StatusCode >= 400 {
		return fmt.Errorf("couldn't check remote file: server returned %s", head.Status)
	}
	if head.Header.Get("Content-Type") != "application/zip" {
		return fmt.Errorf("expected a zip file (application/zip), got a %s", head.Header.Get("Content-Type"))
	}

	// Download it
	process := exec.Command("wget", URL.String(), "-O", downloadTo)
	output, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("couldn't download zip file: %s: %s", err.Error(), output)
	}

	// Unzip it, check contents
	unzipped, err := zip.Unzip(downloadTo, path.Dir(downloadTo))
	var hasAManifest bool
	for _, file := range unzipped {
		if PatternManifestBasename.MatchString(path.Base(file)) {
			hasAManifest = true
			break
		}
	}
	if !hasAManifest {
		os.RemoveAll(path.Dir(downloadTo))
		return errors.New("downloaded zip file has no manifest file (ffcss.yaml)")
	}

	// Remove zip file, only the output folder is interesting
	os.Remove(downloadTo)
	return nil
}

// IsURLClonable determines if the given URL points to a git repository
func IsURLClonable(URL url.URL) (clonable bool, err error) {
	output, err := exec.Command("git", "ls-remote", URL.String()).CombinedOutput()
	if err == nil {
		return true, nil
	}
	switch err.(type) {
	case *exec.ExitError:
		if err.(*exec.ExitError).ExitCode() == 128 {
			return false, nil
		}
	}
	return false, fmt.Errorf("while running git-ls-remote: %s: %s", err.Error(), output)
}

// GetThemeDownloadPath determines where to download a theme
func GetThemeDownloadPath(URL url.URL) (directory string) {
	directory = path.Join(GetConfigDir(), "themes")
	clonable, _ := IsURLClonable(URL)
	if URL.Host == "github.com" && clonable {
		repo := strings.Split(strings.TrimPrefix(strings.TrimSuffix(URL.Path, ".git"), "/"), "/")
		if len(repo) != 2 {
			goto fallback
		}
		return path.Join(directory, "@"+repo[0], repo[1])
	}
fallback:
	return path.Join(directory, "-"+URL.Host, URL.Path)
}

// GetManifest returns a Manifest from the manifest file of themeRoot
func GetManifest(themeRoot string) (Manifest, error) {
	if _, err := os.Stat(GetManifestPath(themeRoot)); os.IsExist(err) {
		return LoadManifest(GetManifestPath(themeRoot))
	} else {
		return Manifest{}, errors.New("the project has no manifest file")
	}
}

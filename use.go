package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/docopt/docopt-go"
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

// RunCommandUse runs the command "use"
func RunCommandUse(args docopt.Opts) error {
	themeName, _ := args.String("THEME_NAME")
	uri, typ := ResolveThemeName(themeName)
	switch typ {
	case "local":
		// do nothing
	case "github":
		URL, err := url.Parse(uri)
		if err != nil {
			return err
		}
		DownloadRepository(*URL)
	case "website":
		// TODO
	default:
		return errors.New("invalid theme name")
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
		return "https://github.com/" + themeName, "github"
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
			return theme.Repository, "github"
		}
		return "", ""
	}
}

// DownloadRepository downloads the repository at URL and returns the saved path
// TODO: clone repo to temp dir, copy necessary files only to .config/ffcss
func DownloadRepository(URL url.URL) (cloneTo string, err error) {
	cloneTo = GetConfigDir() + "/themes/"
	if URL.Host == "github.com" {
		cloneTo = cloneTo + "@" + strings.TrimPrefix(URL.Path, "/")
		os.MkdirAll(cloneTo, 0777)
		err = exec.Command("git", "clone", URL.String(), cloneTo, "--depth 1").Run()
		if err != nil {
			return "", err
		}
	} else {
		cloneTo = cloneTo + URL.Host + "/" + URL.Path
		os.MkdirAll(cloneTo, 0777)
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

// GetManifest returns a Manifest from the manifest file of themeRoot
func GetManifest(themeRoot string) (Manifest, error) {
	if _, err := os.Stat(GetManifestPath(themeRoot)); os.IsExist(err) {
		return LoadManifest(GetManifestPath(themeRoot))
	} else {
		return Manifest{}, errors.New("the project has no manifest file")
	}
}

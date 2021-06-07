package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
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
//   or ./THEME_NAME (for themes.yaml themes)
//   or ./-DOMAIN.TLD/THEME_NAME
//

// RunCommandUse runs the command "use"
func RunCommandUse(args docopt.Opts) error {
	themeName, _ := args.String("THEME_NAME")
	urlOrFolder := ResolveThemeName(themeName)
	downloaded := strings.HasPrefix(urlOrFolder, "~/.config/ffcss/")
	if !downloaded {
		if urlOrFolder == "" {
			return errors.New(themeName + " is not a known theme. Try specifying a github repository directly.")
		}
		fmt.Printf("%s was not found or is not a folder: downloading from %s ...\n", themeName, urlOrFolder)
	}
	fmt.Println("--- work in progress ---")
	return nil
}

// ResolveThemeName resolves the THEME_NAME given to ffcss use to either:
// - a URL to download
// - a local folder to pull the theme from
func ResolveThemeName(themeName string) string {
	protocolLessURL := regexp.MustCompile(`\w+\.\w+/.*`)

	// First test to see if the folder exists
	info, err := os.Stat(filepath.Join(GetConfigDir(), themeName))
	if !os.IsNotExist(err) && info.IsDir() {
		abspath, err := filepath.Abs(filepath.Join(GetConfigDir(), themeName))
		if err != nil {
			return abspath
		}
	}

	// Try OWNER/REPO
	if len(strings.Split(themeName, "/")) == 2 {
		return "https://github.com/" + themeName
		// Try DOMAIN.TLD/PATH
	} else if protocolLessURL.MatchString(themeName) {
		return "https://" + themeName
		// Try URL
	} else if isValidURL(themeName) {
		return themeName
		// Try to get URL from themes.yaml
	} else {
		themes := ReadThemesList()
		if theme, ok := themes[themeName]; ok {
			return theme.Repository
		}
		return ""
	}
}

// DownloadRepository downloads the repository at URL and returns the saved path
// TODO: clone repo to temp dir, copy necessary files only to .config/ffcss
func DownloadRepository(URL url.URL) (cloneTo string, err error) {
	if URL.Host == "github.com" {
		cloneTo = cloneTo + "@" + URL.Path
		os.MkdirAll(cloneTo, 0777)
		err = exec.Command("git", "clone", URL.String(), cloneTo).Run()
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

// GetManifest returns the path of the manifest file given the cloned repo's root path
func GetManifest(themeRoot string) (string, error) {
	if _, err := os.Stat(GetManifestPath(themeRoot)); os.IsExist(err) {
		return GetManifestPath(themeRoot), nil
	} else {
		return "", errors.New("the project has no manifest file")
	}
}

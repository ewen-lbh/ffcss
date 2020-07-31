package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/docopt/docopt-go"
)

// RunCommandUse runs the command "use"
func RunCommandUse(args docopt.Opts) error {
	themeName, _ := args.String("THEME_NAME")
	urlOrFolder := ResolveThemeName(themeName)
	downloaded := false
	if strings.HasPrefix(urlOrFolder, "~/.config/ffcss/") {
		downloaded = true
	}
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
	info, err := os.Stat("~/.config/ffcss/" + themeName)
	if !os.IsNotExist(err) && info.IsDir() {
		return "~/.config/ffcss/" + themeName
	}

	// Try OWNER/REPO
	if len(strings.Split(themeName, "/")) == 2 {
		return "https://github.com/" + themeName
		// Try DOMAIN.TLD/PATH
	} else if protocolLessURL.MatchString(themeName) {
		return "https://" + themeName
		// Try URL
	} else if isValidUrl(themeName) {
		return themeName
		// Try to get URL from themes.toml
	} else {
		themes := ReadThemesList()
		if theme, ok := themes[themeName]; ok {
			return theme.repository
		}
		return ""
	}
}

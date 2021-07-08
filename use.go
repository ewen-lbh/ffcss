package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/docopt/docopt-go"
)

// RunCommandUse runs the command "use"
func RunCommandUse(args docopt.Opts) error {
	themeName, _ := args.String("THEME_NAME")
	// variant, _ := args.String("VARIANT")
	err := os.MkdirAll(path.Join(GetConfigDir(), "themes"), 0777)
	if err != nil {
		return fmt.Errorf("couldn't create data directories: %w", err)
	}
	uri, typ := ResolveURL(themeName)
	manifest, err := Download(uri, typ)
	if err != nil {
		return err
	}
	// Get all profile directories
	dotMozilla, _ := args.String("--mozilla-dir")
	profileDirs, err := ProfileDirsPaths(dotMozilla)
	if err != nil {
		return fmt.Errorf("couldn't get mozilla profile directories: %w", err)
	}
	// Choose profiles
	// TODO smart default (based on {{profileDirectory}}/times.json:firstUse)
	selectedProfileDirs := make([]string, 0)
	selectProfileDirs := &survey.MultiSelect{
		Message: "On which profiles to install this?",
		Options: profileDirs,
	}
	survey.AskOne(selectProfileDirs, &selectedProfileDirs)
	// Choose variant
	variantName, _ := args.String("VARIANT")
	if len(manifest.AvailableVariants()) > 0 && variantName == "" {
		variantPrompt := &survey.Select{
			Message: "Choose the variant",
			Options: manifest.AvailableVariants(),
		}
		survey.AskOne(variantPrompt, &variantName)
	}
	variant := Variant{Name: variantName}
	// Detect OS
	operatingSystem := GOOStoOS(runtime.GOOS)
	// For each profile directory...
	for _, profileDir := range selectedProfileDirs {
		err = RenameIfExists(path.Join(profileDir, "chrome"), path.Join(profileDir, "chrome.bak"))
		if err != nil {
			return fmt.Errorf("while backing up chrome directory: %w", err)
		}

		err := os.Mkdir(path.Join(profileDir, "chrome"), 0700)
		if err != nil {
			return err
		}

		// Install stuff
		err = manifest.InstallUserChrome(operatingSystem, variant, profileDir)
		if err != nil {
			return fmt.Errorf("couldn't install userChrome.css: %w", err)
		}

		err = manifest.InstallUserContent(operatingSystem, variant, profileDir)
		if err != nil {
			return fmt.Errorf("couldn't install userContent.css: %w", err)
		}

		err = manifest.InstallUserJS(operatingSystem, variant, profileDir)
		if err != nil {
			return fmt.Errorf("couldn't install user.js: %w", err)
		}

		err = manifest.InstallAssets(operatingSystem, variant, profileDir)
		if err != nil {
			return fmt.Errorf("couldn't install assets: %w", err)
		}
	}
	return nil
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docopt/docopt-go"
)

var BaseIndentLevel uint = 0

// RunCommandUse runs the command "use"
func RunCommandUse(args docopt.Opts) error {
	themeName, _ := args.String("THEME_NAME")

	err := CreateDataDirectories()
	if err != nil {
		return err
	}

	li(BaseIndentLevel+0, "Resolving the theme's name")
	uri, typ, err := ResolveURL(themeName)
	if err != nil {
		return fmt.Errorf("while resolving name %s: %w", themeName, err)
	}

	li(BaseIndentLevel+0, "Downloading the theme")
	manifest, err := Download(uri, typ)
	if err != nil {
		return err
	}

	intro(manifest, BaseIndentLevel)
	skipSource, _ := args.Bool("--skip-manifest-source")
	manifest.AskToSeeManifestSource(skipSource)

	// Detect OS
	operatingSystem := GOOStoOS(runtime.GOOS)

	// Get all profile directories
	selectedProfiles, err := SelectProfiles(args)
	if err != nil {
		return err
	}

	if len(selectedProfiles) == 0 {
		return nil
	}

	incompatibleProfiles, err := manifest.IncompatibleProfiles(selectedProfiles)
	if err != nil {
		return fmt.Errorf("while checking for incompatible profiles: %w", err)
	}

	if len(incompatibleProfiles) != 0 {
		li(BaseIndentLevel+1, "[yellow]This theme ensures compatibility with firefox [bold]%s[reset][yellow]. The following themes could be incompatible:", manifest.FirefoxVersionConstraint.sentence)
		for _, profile := range incompatibleProfiles {
			li(BaseIndentLevel+2, "%s [dim]([reset]version [blue][bold]%s[reset][dim])", profile.profile, profile.version)
		}
	}

	// Choose variant
	variant, cancel := manifest.ChooseVariant(args)
	if cancel {
		return nil
	}
	manifest, actionsNeeded := manifest.WithVariant(variant)
	manifest.ReDownloadIfNeeded(actionsNeeded)

	// Check for OS compatibility
	manifest.WarnIfIncompatibleWithOS()

	// For each profile directory...
	singleProfile := len(selectedProfiles) == 1
	if singleProfile {
		BaseIndentLevel--
	}
	for _, profile := range selectedProfiles {
		if !singleProfile {
			li(BaseIndentLevel+0, "With profile "+filepath.Base(profile.Path))
		}
		li(BaseIndentLevel+1, "Backing up the chrome/ folder")
		err = RenameIfExists(filepath.Join(profile.Path, "chrome"), filepath.Join(profile.Path, "chrome.bak"))
		if err != nil {
			return fmt.Errorf("while backing up chrome directory: %w", err)
		}

		// Run pre-install script
		if manifest.Run.Before != "" {
			li(BaseIndentLevel+1, "Running pre-install script")
			// TODO for this to be useful, print commandline _with mustaches replaced_:  li(baseIndent+2, "[dim]$ bash -c [reset][bold]%s", manifest.Run.Before)
			output, err := manifest.RunPreInstallHook(profile)
			if err != nil {
				return fmt.Errorf("while running pre-install script: %w", err)
			}
			fmt.Print(
				"\n",
				prefixEachLine(
					strings.TrimSpace(output),
					strings.Repeat(indent, int(BaseIndentLevel)+2),
				),
				"\n",
				"\n",
			)
		}

		err := os.Mkdir(filepath.Join(profile.Path, "chrome"), 0700)
		if err != nil {
			return err
		}

		// Install stuff
		li(BaseIndentLevel+1, "Installing the theme")
		err = manifest.InstallUserChrome(operatingSystem, variant, profile.Path)
		if err != nil {
			return fmt.Errorf("couldn't install userChrome.css: %w", err)
		}

		err = manifest.InstallUserContent(operatingSystem, variant, profile.Path)
		if err != nil {
			return fmt.Errorf("couldn't install userContent.css: %w", err)
		}

		err = manifest.InstallUserJS(operatingSystem, variant, profile.Path)
		if err != nil {
			return fmt.Errorf("couldn't install user.js: %w", err)
		}

		err = manifest.InstallAssets(operatingSystem, variant, profile.Path)
		if err != nil {
			return fmt.Errorf("couldn't install assets: %w", err)
		}

		// Run post-install script
		if manifest.Run.After != "" {
			li(BaseIndentLevel+1, "Running post-install script")
			// TODO for this to be useful, print commandline _with mustaches replaced_:  li(baseIndent+2, "[dim]$ bash -c [reset][bold]%s", manifest.Run.After)
			output, err := manifest.RunPostInstallHook(profile)
			if err != nil {
				return fmt.Errorf("while running post-install script: %w", err)
			}
			fmt.Print(
				"\n",
				prefixEachLine(
					strings.TrimSpace(output),
					strings.Repeat(indent, int(BaseIndentLevel)+2),
				),
				"\n",
				"\n",
			)
		}

		err = profile.RegisterCurrentTheme(themeName)
		if err != nil {
			return fmt.Errorf("while registering current theme for profile %q: %w", profile.FullName(), err)
		}

	}
	if singleProfile {
		BaseIndentLevel++
	}

	// Ask to open extensions' pages
	if len(manifest.Addons) > 0 {
		if ConfirmInstallAddons(manifest.Addons) {
			for _, profile := range selectedProfiles {
				li(BaseIndentLevel+0, "With profile "+filepath.Base(profile.Path))
				for _, addonURL := range manifest.Addons {
					profile.InstallAddon(operatingSystem, addonURL)
				}
			}
		}
	}

	// Show message
	err = manifest.ShowMessage()
	if err != nil {
		return fmt.Errorf("couldn't display the message: %w", err)
	}
	return nil
}

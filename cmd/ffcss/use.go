package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/docopt/docopt-go"
	. "github.com/ewen-lbh/ffcss"
)

// RunCommandUse runs the command "use"
func RunCommandUse(args docopt.Opts) error {
	themeName, _ := args.String("THEME_NAME")

	err := CreateDataDirectories()
	if err != nil {
		return err
	}

	Step(0, "Resolving the theme's name")
	uri, typ, err := ResolveURL(themeName)
	if err != nil {
		return fmt.Errorf("while resolving name %s: %w", themeName, err)
	}

	Step(0, "Downloading the theme")
	manifest, err := Download(uri, typ)
	if err != nil {
		return err
	}

	DescribeTheme(manifest, BaseIndentLevel)

	skipSource, _ := args.Bool("--skip-manifest-source")
	manifest.AskToSeeManifestSource(skipSource)

	// Detect OS
	operatingSystem := GOOStoOS(runtime.GOOS)

	selectedProfiles, err := SelectProfiles(args)
	if err != nil {
		return err
	}

	if len(selectedProfiles) == 0 {
		return nil
	}
	singleProfile := len(selectedProfiles) == 1

	incompatibleProfiles, err := manifest.IncompatibleProfiles(selectedProfiles)
	if err != nil {
		return fmt.Errorf("while checking for incompatible profiles: %w", err)
	}

	if len(incompatibleProfiles) != 0 {
		Step(1, "[yellow]This theme ensures compatibility with firefox [bold]%s[reset][yellow]. The following themes could be incompatible:", manifest.FirefoxVersionConstraint.Sentence)
		for _, profile := range incompatibleProfiles {
			Step(2, "%s [dim]([reset]version [blue][bold]%s[reset][dim])", profile.Profile, profile.Version)
		}
	}

	// Choose variant
	variant := Variant{}
	if len(manifest.AvailableVariants()) > 0 {
		variantName, _ := args.String("VARIANT")
		if variantName == "" {
			var cancel bool
			variant, cancel = manifest.ChooseVariant()
			if cancel {
				return nil
			}
		} else {
			var found bool
			variant, found = manifest.Variants[variantName]
			if !found {
				return fmt.Errorf("variant %q does not exist on this theme. Available variants are %s", variantName, strings.Join(manifest.AvailableVariants(), ", "))
			}
		}
		manifest, actionsNeeded := manifest.WithVariant(variant)
		manifest.ReDownloadIfNeeded(actionsNeeded)
	}

	// Check for OS compatibility
	manifest.WarnIfIncompatibleWithOS()

	// For each profile directory...
	if singleProfile {
		BaseIndentLevel--
	}
	for _, profile := range selectedProfiles {
		if !singleProfile {
			Step(0, "With profile "+filepath.Base(profile.Path))
		}

		Step(1, "Backing up the chrome/ folder")
		err = profile.BackupChrome()
		if err != nil {
			return fmt.Errorf("while backing up chrome directory: %w", err)
		}

		// Run pre-install script
		if manifest.Run.Before != "" {
			Step(1, "Running pre-install script")
			// TODO for this to be useful, print commandline _with mustaches replaced_:  Step(baseIndent+2, "[dim]$ bash -c [reset][bold]%s", manifest.Run.Before)
			output, err := manifest.RunPreInstallHook(profile)
			if err != nil {
				return fmt.Errorf("while running pre-install script: %w", err)
			}
			ShowHookOutput(output)
		}

		err := os.Mkdir(filepath.Join(profile.Path, "chrome"), 0700)
		if err != nil {
			return err
		}

		// Install stuff
		Step(1, "Installing the theme")
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
			Step(1, "Running post-install script")
			// TODO for this to be useful, print commandline _with mustaches replaced_:  Step(baseIndent+2, "[dim]$ bash -c [reset][bold]%s", manifest.Run.After)
			output, err := manifest.RunPostInstallHook(profile)
			if err != nil {
				return fmt.Errorf("while running post-install script: %w", err)
			}
			ShowHookOutput(output)
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
				Step(0, "With profile "+filepath.Base(profile.Path))
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

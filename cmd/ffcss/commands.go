package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/ewen-lbh/ffcss"
	"gopkg.in/yaml.v2"
)

func runCommandUse(args flagsAndArgs) error {
	err := ffcss.CreateDataDirectories()
	if err != nil {
		return err
	}

	ffcss.LogStep(0, "Resolving the theme's name")
	uri, typ, err := ffcss.ResolveURL(args.string("THEME_NAME"))
	if err != nil {
		return fmt.Errorf("while resolving name %s: %w", args.string("THEME_NAME"), err)
	}

	ffcss.LogStep(0, "Downloading the theme")
	manifest, err := ffcss.Download(uri, typ)
	if err != nil {
		return err
	}

	ffcss.DescribeTheme(manifest, ffcss.BaseIndentLevel)

	manifest.AskToSeeManifestSource(args.bool("--skip-manifest-source"))

	// Detect OS
	operatingSystem := ffcss.GOOStoOS(runtime.GOOS)

	selectedProfiles, err := ffcss.SelectProfiles(
		args.strings("--profiles"),
		args.string("--profiles-dir"),
		args.bool("--all-profiles"),
	)
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
		ffcss.LogStep(1, "[yellow]This theme ensures compatibility with firefox [bold]%s[reset][yellow]. The following themes could be incompatible:", manifest.FirefoxVersionConstraint.Sentence)
		for _, profile := range incompatibleProfiles {
			ffcss.LogStep(2, "%s [dim]([reset]version [blue][bold]%s[reset][dim])", profile.Profile, profile.Version)
		}
	}

	// Choose variant
	variant := ffcss.Variant{}
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
	manifest.WarnIfIncompatibleWithOS(operatingSystem)

	// For each profile directory...
	if singleProfile {
		ffcss.BaseIndentLevel--
	}
	for _, profile := range selectedProfiles {
		if !singleProfile {
			ffcss.LogStep(0, "With profile "+filepath.Base(profile.Path))
		}

		ffcss.LogStep(1, "Backing up the chrome/ folder")
		err = profile.BackupChrome()
		if err != nil {
			return fmt.Errorf("while backing up chrome directory: %w", err)
		}

		// Run pre-install script
		if manifest.Run.Before != "" {
			ffcss.LogStep(1, "Running pre-install script")
			// TODO for this to be useful, print commandline _with mustaches replaced_:  Step(baseIndent+2, "[dim]$ bash -c [reset][bold]%s", manifest.Run.Before)
			output, err := manifest.RunPreInstallHook(profile)
			if err != nil {
				return fmt.Errorf("while running pre-install script: %w", err)
			}
			ffcss.ShowHookOutput(output)
		}

		err := os.Mkdir(filepath.Join(profile.Path, "chrome"), 0700)
		if err != nil {
			return err
		}

		// Install stuff
		ffcss.LogStep(1, "Installing the theme")
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
			ffcss.LogStep(1, "Running post-install script")
			// TODO for this to be useful, print commandline _with mustaches replaced_:  Step(baseIndent+2, "[dim]$ bash -c [reset][bold]%s", manifest.Run.After)
			output, err := manifest.RunPostInstallHook(profile)
			if err != nil {
				return fmt.Errorf("while running post-install script: %w", err)
			}
			ffcss.ShowHookOutput(output)
		}

		err = profile.RegisterCurrentTheme(args.string("THEME_NAME"))
		if err != nil {
			return fmt.Errorf("while registering current theme for profile %q: %w", profile.FullName(), err)
		}

	}
	if singleProfile {
		ffcss.BaseIndentLevel++
	}

	// Ask to open extensions' pages
	if len(manifest.Addons) > 0 {
		if ffcss.ConfirmInstallAddons(manifest.Addons) {
			for _, profile := range selectedProfiles {
				ffcss.LogStep(0, "With profile "+filepath.Base(profile.Path))
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

func runCommandGet(args flagsAndArgs) error {
	themeName, _ := args.String("THEME_NAME")
	// variant, _ := args.String("VARIANT")

	err := ffcss.CreateDataDirectories()
	if err != nil {
		return err
	}

	ffcss.LogStep(0, "Resolving the theme's name")
	uri, typ, err := ffcss.ResolveURL(themeName)
	if err != nil {
		return fmt.Errorf("while resolving name %s: %w", themeName, err)
	}

	ffcss.LogStep(0, "Downloading the theme")
	manifest, err := ffcss.Download(uri, typ)
	if err != nil {
		return err
	}

	ffcss.LogStepC("✓", 0, "Downloaded [blue][bold]%s[reset] [dim](to %s)", manifest.Name(), manifest.DownloadedTo)
	return nil
}

func runCommandReapply(args flagsAndArgs) error {
	operatingSystem := ffcss.GOOStoOS(runtime.GOOS)
	profilesDir, _ := args.String("--profiles-dir")
	var profilesPaths []string
	var err error
	if profilesDir != "" {
		profilesPaths, err = ffcss.ProfilePaths(operatingSystem, profilesDir)
	} else {
		profilesPaths, err = ffcss.ProfilePaths(operatingSystem)
	}
	if err != nil {
		return fmt.Errorf("while getting profiles: %w", err)
	}

	currentThemesFileContents, err := os.ReadFile(ffcss.ConfigDir("currently.yaml"))
	if err != nil {
		return fmt.Errorf("while reading current themes list: %w", err)
	}

	currentThemes := make(map[string]string)
	yaml.Unmarshal(currentThemesFileContents, &currentThemes)

	for _, profilePath := range profilesPaths {
		themeName, exists := currentThemes[filepath.Base(profilePath)]
		if !exists {
			ffcss.LogStep(0, "[yellow]Profile %s[reset][yellow] has no ffcss theme applied, skipping.", ffcss.NewFirefoxProfileFromPath(profilePath).Display())
			continue
		}
		ffcss.LogStep(0, "Apply theme [blue][bold]%s[reset] to profile %s", themeName, ffcss.NewFirefoxProfileFromPath(profilePath).Display())

		useArgs, _ := docopt.ParseArgs(usage, []string{"use", string(themeName), "--profiles", profilePath, "--skip-manifest-source"}, ffcss.VersionString)
		ffcss.BaseIndentLevel++
		err = runCommandUse(flagsAndArgs{useArgs})
		if err != nil {
			return err
		}
	}
	return nil
}

func runCommandInit(args flagsAndArgs) error {
	// TODO: set user{Chrome,Content,.js} by finding their path
	// TODO: only set assets if chrome/ actually exists
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working directory: %w", err)
	}

	theme, err := ffcss.InitializeTheme(workingDir)
	if err != nil {
		return fmt.Errorf("while initializing theme: %w", err)
	}

	return theme.WriteManifest(workingDir)
}

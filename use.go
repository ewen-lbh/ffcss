package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/docopt/docopt-go"
)

// RunCommandUse runs the command "use"
func RunCommandUse(args docopt.Opts) error {
	themeName, _ := args.String("THEME_NAME")
	// variant, _ := args.String("VARIANT")
	err := os.MkdirAll(filepath.Join(GetConfigDir(), "themes"), 0777)
	if err != nil {
		return fmt.Errorf("couldn't create data directories: %w", err)
	}
	uri, typ, err := ResolveURL(themeName)
	if err != nil {
		return fmt.Errorf("while resolving name %s: %w", themeName, err)
	}

	manifest, err := Download(uri, typ)
	if err != nil {
		return err
	}
	// Detect OS
	operatingSystem := GOOStoOS(runtime.GOOS)
	// Get all profile directories
	profilesDir, _ := args.String("--profiles-dir")
	var profileDirs []string
	if profilesDir != "" {
		profileDirs, err = ProfileDirsPaths(operatingSystem, profilesDir)
	} else {
		profileDirs, err = ProfileDirsPaths(operatingSystem)
	}
	if err != nil {
		return fmt.Errorf("couldn't get profile directories: %w", err)
	}
	// Choose profiles
	// TODO smart default (based on {{profileDirectory}}/times.json:firstUse)
	selectedProfileDirs := make([]string, 0)
	selectAllProfileDirs, _ := args.Bool("--all-profiles")
	if selectAllProfileDirs {
		selectedProfileDirs = profileDirs
	} else {
		selectProfileDirs := &survey.MultiSelect{
			Message: "On which profiles to install this?",
			Options: profileDirs,
		}
		survey.AskOne(selectProfileDirs, &selectedProfileDirs)
		// User Ctrl-C'd
		if len(selectedProfileDirs) == 0 {
			return nil
		}
	}
	// Choose variant
	variantName, _ := args.String("VARIANT")
	if len(manifest.AvailableVariants()) > 0 && variantName == "" {
		variantPrompt := &survey.Select{
			Message: "Choose the variant",
			Options: manifest.AvailableVariants(),
		}
		survey.AskOne(variantPrompt, &variantName)
		// user Ctrl-C'd
		if variantName == "" {
			return nil
		}
	}
	variant := manifest.Variants[variantName]
	manifest, actionsNeeded := manifest.WithVariant(variant)
	// FIXME for now switching branches just re-downloads the entire thing to a new dir with the new branch
	// ideal thing would be to copy from the root variant to the new variant, cd into it then `git switch` there.
	if actionsNeeded.reDownload || actionsNeeded.switchBranch {
		d("re-downloading: new repo is %s", manifest.DownloadAt)
		uri, typ, err := ResolveURL(manifest.DownloadAt)
		if err != nil {
			return fmt.Errorf("while resolving URL %s: %w", manifest.DownloadAt, err)
		}

		_, err = Download(uri, typ, manifest)
		if err != nil {
			return fmt.Errorf("couldn't download the variant at %s: %w", uri, err)
		}
	}

	// Check for OS compatibility
	for k, v := range manifest.OSNames {
		if k == operatingSystem && v == "" {
			warn("This theme is marked as incompatible with %s. Things might not work.", operatingSystem)
		}
	}

	// For each profile directory...
	for _, profileDir := range selectedProfileDirs {
		err = RenameIfExists(filepath.Join(profileDir, "chrome"), filepath.Join(profileDir, "chrome.bak"))
		if err != nil {
			return fmt.Errorf("while backing up chrome directory: %w", err)
		}

		err := os.Mkdir(filepath.Join(profileDir, "chrome"), 0700)
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

	// Ask to open extensions' pages
	if len(manifest.Addons) > 0 {
		acceptOpenExtensionPages := false
		survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("This theme suggests installing %d %s. Open %s?",
				len(manifest.Addons),
				plural("addon", len(manifest.Addons)),
				plural("its page", len(manifest.Addons), "their pages"),
			),
			Default: acceptOpenExtensionPages,
		}, &acceptOpenExtensionPages)

		if acceptOpenExtensionPages {
			for _, profile := range selectedProfileDirs {
				fmt.Printf("Opening on %s\n", profile)
				for _, url := range manifest.Addons {
					var command *exec.Cmd
					switch operatingSystem {
					case "linux":
						command = exec.Command("firefox", "--new-tab", url, "--profile", profile)
					case "macos":
						command = exec.Command("open", "-a", "firefox", url, "--args", "--profile", profile)
					case "windows":
						command = exec.Command("start", "firefox", "-profile", profile, url)
					default:
						warn("unrecognized OS %s, cannot open firefox automatically. Open %s in firefox using profile %s", operatingSystem, url, profile)
					}
					err = command.Run()
					if err != nil {
						return fmt.Errorf("couldn't open %q: while running %s: %w", url, command.String(), err)
					}
					break
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

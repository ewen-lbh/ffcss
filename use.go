package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/docopt/docopt-go"
)

// RunCommandUse runs the command "use"
func RunCommandUse(args docopt.Opts, indentationLevel ...uint) error {
	var baseIndent uint
	if len(indentationLevel) >= 1 {
		baseIndent = indentationLevel[0]
	}
	themeName, _ := args.String("THEME_NAME")
	// variant, _ := args.String("VARIANT")
	err := os.MkdirAll(filepath.Join(GetConfigDir(), "themes"), 0777)
	if err != nil {
		return fmt.Errorf("couldn't create data directories: %w", err)
	}
	li(baseIndent+0, "Resolving the theme's name")
	uri, typ, err := ResolveURL(themeName)
	if err != nil {
		return fmt.Errorf("while resolving name %s: %w", themeName, err)
	}

	li(baseIndent+0, "Downloading the theme")
	manifest, err := Download(uri, typ)
	if err != nil {
		return err
	}

	intro(manifest, baseIndent)
	wantsSource := false
	skipQuestion, _ := args.Bool("--skip-manifest-source")
	if !skipQuestion {
		survey.AskOne(&survey.Confirm{
			Message: "Show the manifest source?",
		}, &wantsSource)
	}
	if wantsSource {
		showSource(manifest)
	}

	// Detect OS
	operatingSystem := GOOStoOS(runtime.GOOS)
	// Get all profile directories
	selectedProfilesString, _ := args.String("--profiles")
	var selectedProfilePaths []string
	if selectedProfilesString == "" {
		li(baseIndent+0, "Getting profiles")
		profilesDir, _ := args.String("--profiles-dir")
		var profilePaths []string
		if profilesDir != "" {
			profilePaths, err = ProfilePaths(operatingSystem, profilesDir)
		} else {
			profilePaths, err = ProfilePaths(operatingSystem)
		}
		if err != nil {
			return fmt.Errorf("couldn't get profile directories: %w", err)
		}
		if manifest.FirefoxVersion != "" {
			constraint, err := NewFirefoxVersionConstraint(manifest.FirefoxVersion)
			if err != nil {
				return fmt.Errorf("invalid firefox version constraint %q: %w", manifest.FirefoxVersion, err)
			}
			incompatibleProfileDirs := make([]struct {
				profile FirefoxProfile
				version FirefoxVersion
			}, 0)
			for _, profilePath := range profilePaths {
				profile := FirefoxProfileFromPath(profilePath)
				profileVersion, err := profile.FirefoxVersion()
				if err != nil {
					warn("Couldn't get firefox version for profile %s", profilePath)
				}
				fulfillsConstraint := constraint.FulfilledBy(profileVersion)
				if !fulfillsConstraint {
					incompatibleProfileDirs = append(incompatibleProfileDirs, struct {
						profile FirefoxProfile
						version FirefoxVersion
					}{profile, profileVersion})
				}
			}
			if len(incompatibleProfileDirs) != 0 {
				li(baseIndent+1, "[yellow]This theme ensures compatibility with firefox [bold]%s[reset][yellow]. The following themes could be incompatible:", constraint.sentence)
				for _, profile := range incompatibleProfileDirs {
					li(baseIndent+2, "%s [dim]([reset]version [blue][bold]%s[reset][dim])", profile.profile, profile.version)
				}
			}
		}

		// Choose profiles
		// TODO smart default (based on {{profileDirectory}}/times.json:firstUse)
		selectAllProfilePaths, _ := args.Bool("--all-profiles")
		if selectAllProfilePaths {
			li(baseIndent+0, "Selecting all profiles")
			selectedProfilePaths = profilePaths
		} else {
			// XXX the whole display thing should be put in survey.MultiSelect.Renderer, look into that.
			selectedProfileDirsDisplay := make([]string, 0)
			li(baseIndent+0, "Please select profiles to apply the theme on")
			profileDirsDisplay := apply(func(p string) string { return FirefoxProfileFromPath(p).String() }, profilePaths)
			survey.AskOne(&survey.MultiSelect{
				Message: "Select profiles",
				Options: profileDirsDisplay,
				VimMode: VimModeEnabled(),
			}, &selectedProfileDirsDisplay)
			for _, chosenProfileDisplay := range selectedProfileDirsDisplay {
				selectedProfilePaths = append(selectedProfilePaths, FirefoxProfileFromDisplayString(chosenProfileDisplay, profilePaths).Path)
			}
			// User Ctrl-C'd
			if len(selectedProfilePaths) == 0 {
				return nil
			}
		}
	} else {
		selectedProfilePaths = strings.Split(selectedProfilesString, ",")
	}
	// Choose variant
	variantName, _ := args.String("VARIANT")
	if len(manifest.AvailableVariants()) > 0 && variantName == "" {
		li(baseIndent+0, "Please choose the theme's variant")
		variantPrompt := &survey.Select{
			Message: "Install variant",
			Options: manifest.AvailableVariants(),
			VimMode: VimModeEnabled(),
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
		li(baseIndent+0, "Downloading the variant")
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
	singleProfile := len(selectedProfilePaths) == 1
	if singleProfile {
		baseIndent--
	}
	for _, profileDir := range selectedProfilePaths {
		profile := FirefoxProfileFromPath(profileDir)
		if !singleProfile {
			li(baseIndent+0, "With profile "+filepath.Base(profileDir))
		}
		li(baseIndent+1, "Backing up the chrome/ folder")
		err = RenameIfExists(filepath.Join(profileDir, "chrome"), filepath.Join(profileDir, "chrome.bak"))
		if err != nil {
			return fmt.Errorf("while backing up chrome directory: %w", err)
		}

		// Run pre-install script
		if manifest.Run.Before != "" {
			li(baseIndent+1, "Running pre-install script")
			// TODO for this to be useful, print commandline _with mustaches replaced_:  li(baseIndent+2, "[dim]$ bash -c [reset][bold]%s", manifest.Run.Before)
			output, err := manifest.RunPreInstallHook(profile)
			if err != nil {
				return fmt.Errorf("while running pre-install script: %w", err)
			}
			fmt.Print(
				"\n",
				prefixEachLine(
					strings.TrimSpace(output),
					strings.Repeat(indent, int(baseIndent)+2),
				),
				"\n",
				"\n",
			)
		}

		err := os.Mkdir(filepath.Join(profileDir, "chrome"), 0700)
		if err != nil {
			return err
		}

		// Install stuff
		li(baseIndent+1, "Installing the theme")
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

		// Run post-install script
		if manifest.Run.After != "" {
			li(baseIndent+1, "Running post-install script")
			// TODO for this to be useful, print commandline _with mustaches replaced_:  li(baseIndent+2, "[dim]$ bash -c [reset][bold]%s", manifest.Run.After)
			output, err := manifest.RunPostInstallHook(profile)
			if err != nil {
				return fmt.Errorf("while running post-install script: %w", err)
			}
			fmt.Print(
				"\n",
				prefixEachLine(
					strings.TrimSpace(output),
					strings.Repeat(indent, int(baseIndent)+2),
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
		baseIndent++
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
			for _, profile := range selectedProfilePaths {
				li(baseIndent+0, "With profile "+filepath.Base(profile))
				for _, url := range manifest.Addons {
					li(baseIndent+1, "Opening [blue][bold]%s", url)
					li(baseIndent+1, "[yellow]Waiting for you to close Firefox")
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

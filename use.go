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
	var selectedProfiles []FirefoxProfile
	if selectedProfilesString != "" {
		for _, profilePath := range strings.Split(selectedProfilesString, ",") {
			selectedProfiles = append(selectedProfiles, FirefoxProfileFromPath(profilePath))
		}
	} else {
		li(baseIndent+0, "Getting profiles")
		profilesDir, _ := args.String("--profiles-dir")
		profiles, err := Profiles(profilesDir)
		if err != nil {
			return fmt.Errorf("couldn't get profile directories: %w", err)
		}
		// Choose profiles
		// TODO smart default (based on {{profileDirectory}}/times.json:firstUse)
		selectAllProfilePaths, _ := args.Bool("--all-profiles")
		if selectAllProfilePaths {
			li(baseIndent+0, "Selecting all profiles")
			selectedProfiles = profiles
		} else {
			selectedProfiles = AskProfiles(profiles, baseIndent)
		}
	}

	if len(selectedProfiles) == 0 {
		return nil
	}


	incompatibleProfiles, err := manifest.IncompatibleProfiles(selectedProfiles)
	if err != nil {
		return  fmt.Errorf("while checking for incompatible profiles: %w",  err)
	}

	if len(incompatibleProfiles) != 0 {
		li(baseIndent+1, "[yellow]This theme ensures compatibility with firefox [bold]%s[reset][yellow]. The following themes could be incompatible:", manifest.FirefoxVersionConstraint.sentence)
		for _, profile := range incompatibleProfiles {
			li(baseIndent+2, "%s [dim]([reset]version [blue][bold]%s[reset][dim])", profile.profile, profile.version)
		}
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
	singleProfile := len(selectedProfiles) == 1
	if singleProfile {
		baseIndent--
	}
	for _, profile := range selectedProfiles {
		if !singleProfile {
			li(baseIndent+0, "With profile "+filepath.Base(profile.Path))
		}
		li(baseIndent+1, "Backing up the chrome/ folder")
		err = RenameIfExists(filepath.Join(profile.Path, "chrome"), filepath.Join(profile.Path, "chrome.bak"))
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

		err := os.Mkdir(filepath.Join(profile.Path, "chrome"), 0700)
		if err != nil {
			return err
		}

		// Install stuff
		li(baseIndent+1, "Installing the theme")
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
			for _, profile := range selectedProfiles {
				li(baseIndent+0, "With profile "+filepath.Base(profile.Path))
				for _, url := range manifest.Addons {
					li(baseIndent+1, "Opening [blue][bold]%s", url)
					li(baseIndent+1, "[yellow]Waiting for you to close Firefox")
					var command *exec.Cmd
					switch operatingSystem {
					case "linux":
						command = exec.Command("firefox", "--new-tab", url, "--profile", profile.Path)
					case "macos":
						command = exec.Command("open", "-a", "firefox", url, "--args", "--profile", profile.Path)
					case "windows":
						command = exec.Command("start", "firefox", "-profile", profile.Path, url)
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

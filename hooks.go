package main

import (
	"fmt"
	"os/exec"

	"github.com/hoisie/mustache"
)

// runHook runs a provided command for a specific profile. See any of the (Manifest).Run*Hook methods
// for a list of available {{mustache}} placeholders.
func (t Theme) runHook(commandline string, profile FirefoxProfile) (output string, err error) {
	ffversion, err := profile.FirefoxVersion()
	if err != nil {
		return "", fmt.Errorf("while getting firefox version for current profile: %w", err)
	}

	command := exec.Command("bash", "-c", mustache.Render(commandline, map[string]interface{}{
		"profile_path":    profile.Path,
		"firefox_version": ffversion.String(),
	}))

	outputBytes, err := command.CombinedOutput()
	output = string(outputBytes)
	if err != nil {
		return "", fmt.Errorf("while running %q: %s: %w", command.String(), output, err)
	}

	return
}

// RunPreInstallHook passes the pre-install hook specified in the manifest's run.before entry to bash.
// Several {{mustache}} placeholders are available:
//
//	profile_path        The current profile's path
//	firefox_version     The current profile's Firefox version
func (t Theme) RunPreInstallHook(profile FirefoxProfile) (output string, err error) {
	return t.runHook(t.Run.Before, profile)
}

// RunPostInstallHook does the same as RunPreInstallHook but for the manifest's run.after entry.
func (t Theme) RunPostInstallHook(profile FirefoxProfile) (output string, err error) {
	return t.runHook(t.Run.After, profile)
}

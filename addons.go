package main

import (
	"fmt"
	"os/exec"
)

func (ffp FirefoxProfile) InstallAddon(operatingSystem string, addonURL string) error {
	li(BaseIndentLevel+1, "Opening [blue][bold]%s", addonURL)
	li(BaseIndentLevel+1, "[yellow]Waiting for you to close Firefox")
	var command *exec.Cmd
	switch operatingSystem {
	case "linux":
		command = exec.Command("firefox", "--new-tab", addonURL, "--profile", ffp.Path)
	case "macos":
		command = exec.Command("open", "-a", "firefox", addonURL, "--args", "--profile", ffp.Path)
	case "windows":
		command = exec.Command("start", "firefox", "-profile", ffp.Path, addonURL)
	default:
		warn("unrecognized OS %s, cannot open firefox automatically. Open %s in firefox using profile %s", operatingSystem, addonURL, ffp)
		return nil
	}
	err := command.Run()
	if err != nil {
		return fmt.Errorf("couldn't open %q: while running %s: %w", addonURL, command.String(), err)
	}
	return nil
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/docopt/docopt-go"
	"gopkg.in/yaml.v2"
)

// RunCommandReapply runs the command "reapply"
func RunCommandReapply(args docopt.Opts) error {
	operatingSystem := GOOStoOS(runtime.GOOS)
	profilesDir, _ := args.String("--profiles-dir")
	var profilesPaths []string
	var err error
	if profilesDir != "" {
		profilesPaths, err = ProfilePaths(operatingSystem, profilesDir)
	} else {
		profilesPaths, err = ProfilePaths(operatingSystem)
	}
	if err != nil {
		return fmt.Errorf("while getting profiles: %w", err)
	}

	currentThemesFileContents, err := os.ReadFile(ConfigDir("currently.yaml"))
	if err != nil {
		return fmt.Errorf("while reading current themes list: %w", err)
	}

	currentThemes := make(map[string]string)
	yaml.Unmarshal(currentThemesFileContents, &currentThemes)

	for _, profilePath := range profilesPaths {
		li(0, "With profile %s", FirefoxProfileFromPath(profilePath).String())
		themeName, exists := currentThemes[filepath.Base(profilePath)]
		if !exists {
			li(1, "[yellow]This profile has no ffcss theme applied, skipping.")
			continue
		}
		li(1, "Applying theme [blue][bold]%s[reset]", themeName)

		useArgs, _ := docopt.ParseArgs(Usage, []string{"use", string(themeName), "--profiles", profilePath, "--skip-manifest-source"}, VersionString)
		err = RunCommandUse(useArgs, 2)
		if err != nil {
			return err
		}
	}
	return nil
}

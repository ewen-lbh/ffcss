package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/docopt/docopt-go"
	"gopkg.in/yaml.v2"
	. "github.com/ewen-lbh/ffcss"
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
		themeName, exists := currentThemes[filepath.Base(profilePath)]
		if !exists {
			LogStep(0, "[yellow]Profile %s[reset][yellow] has no ffcss theme applied, skipping.", NewFirefoxProfileFromPath(profilePath).Display())
			continue
		}
		LogStep(0, "Apply theme [blue][bold]%s[reset] to profile %s", themeName, NewFirefoxProfileFromPath(profilePath).Display())

		useArgs, _ := docopt.ParseArgs(Usage, []string{"use", string(themeName), "--profiles", profilePath, "--skip-manifest-source"}, VersionString)
		BaseIndentLevel += 1
		err = RunCommandUse(useArgs)
		if err != nil {
			return err
		}
	}
	return nil
}

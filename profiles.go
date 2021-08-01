package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"
)

type FirefoxProfile struct {
	ID   string
	Name string
	Path string
}
type firefoxProfileWithVersion struct {
	profile FirefoxProfile
	version FirefoxVersion
}

func (ffp FirefoxProfile) RegisterCurrentTheme(themeName string) error {
	currentThemes, err := CurrentThemeByProfile()
	if err != nil {
		return err
	}
	currentThemes[ffp.FullName()] = themeName
	currentThemesNewContents, err := yaml.Marshal(currentThemes)
	if err != nil {
		return fmt.Errorf("while marshaling into YAML: %w", err)
	}

	err = os.WriteFile(ConfigDir("currently.yaml"), currentThemesNewContents, 0777)
	if err != nil {
		return fmt.Errorf("while writing new contents: %w", err)
	}

	return nil
}

func (ffp FirefoxProfile) FullName() string {
	return filepath.Base(ffp.Path)
}

func (ffp FirefoxProfile) String() string {
	return fmt.Sprintf("%s (%s)", ffp.Name, ffp.ID)
}

func NewFirefoxProfileFromPath(path string) FirefoxProfile {
	base := filepath.Base(path)
	parts := strings.Split(base, ".")
	return FirefoxProfile{
		Path: path,
		ID:   parts[0],
		Name: parts[1],
	}
}

func NewFirefoxProfileFromDisplay(displayString string, profiles []FirefoxProfile) FirefoxProfile {
	for _, profile := range profiles {
		ffp := NewFirefoxProfileFromPath(profile.Path)
		if ffp.Display() == displayString {
			return ffp
		}
	}
	d("while searching for %s in %v", displayString, profiles)
	panic("internal error: can't get profile from display string")
}

// ProfilePaths returns an array of profile directories from the profile folder.
// 1 arguments: the profiles folder is assumed to be the current OS's default.
// 2 argument: use the given profiles folder
// more arguments: panic.
func ProfilePaths(operatingSystem string, optionalProfilesDir ...string) ([]string, error) {
	var profilesFolder string
	if len(optionalProfilesDir) == 0 {
		// XXX: Weird golang thing, if I assign to profilesFolder directly, it tells me the variable is unused
		_profilesFolder, err := DefaultProfilesDir(operatingSystem)
		profilesFolder = _profilesFolder
		if err != nil {
			return []string{}, fmt.Errorf("couldn't get the profiles folder: %w. Try to use --profiles-dir", err)
		}
	} else if len(optionalProfilesDir) == 1 {
		profilesFolder = optionalProfilesDir[0]
	} else {
		panic(fmt.Sprintf("received %d arguments, expected 1 or 2", len(optionalProfilesDir)+1))
	}
	directories, err := os.ReadDir(profilesFolder)
	releasesPaths := make([]string, 0)
	patternReleaseID := regexp.MustCompile(`[a-z0-9]{8}\.\w+`)
	if err != nil {
		return []string{}, fmt.Errorf("couldn't read %s: %w", profilesFolder, err)
	}
	for _, releasePath := range directories {
		if patternReleaseID.MatchString(releasePath.Name()) {
			stat, err := os.Stat(filepath.Join(profilesFolder, releasePath.Name()))
			if err != nil {
				continue
			}
			if stat.IsDir() {
				releasesPaths = append(releasesPaths, filepath.Join(profilesFolder, releasePath.Name()))
			}
		}
	}
	return releasesPaths, nil
}

// Profiles returns a list of FirefoxProfiles. If optionalProfilesDir is "",
// DefaultProfilesDir() is used to get the profiles' directory.
func Profiles(optionalProfilesDir string) ([]FirefoxProfile, error) {
	var profiles []FirefoxProfile
	var err error
	var profilePaths []string

	if optionalProfilesDir != "" {
		profilePaths, err = ProfilePaths(GOOStoOS(runtime.GOOS), optionalProfilesDir)
	} else {
		profilePaths, err = ProfilePaths(GOOStoOS(runtime.GOOS))
	}

	if err != nil {
		return profiles, fmt.Errorf("while getting profile paths: %w", err)
	}

	for _, profilePath := range profilePaths {
		profiles = append(profiles, NewFirefoxProfileFromPath(profilePath))
	}

	return profiles, nil
}

func DefaultProfilesDir(operatingSystem string) (string, error) {
	switch operatingSystem {
	case "linux":
		homedir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homedir, ".mozilla", "firefox"), nil
	case "macos":
		user, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("couldn't get the current user: %w", err)
		}

		return filepath.Join("/Users", user.Username, "Library", "Application Support", "Firefox", "Profiles"), nil
	case "windows":
		homedir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		return filepath.Join(homedir, "AppData", "Roaming", "Mozilla", "Firefox", "Profiles"), nil
	}
	return "", fmt.Errorf("unknown operating system %s", operatingSystem)
}

func (t Theme) IncompatibleProfiles(profiles []FirefoxProfile) ([]firefoxProfileWithVersion, error) {
	if t.FirefoxVersion != "" {
		incompatibleProfileDirs := make([]firefoxProfileWithVersion, 0)
		for _, profile := range profiles {
			profileVersion, err := profile.FirefoxVersion()
			if err != nil {
				warn("Couldn't get firefox version for profile %s", profile)
			}
			fulfillsConstraint := t.FirefoxVersionConstraint.FulfilledBy(profileVersion)
			if !fulfillsConstraint {
				incompatibleProfileDirs = append(incompatibleProfileDirs, struct {
					profile FirefoxProfile
					version FirefoxVersion
				}{profile, profileVersion})
			}
		}
		return incompatibleProfileDirs, nil
	}
	return []firefoxProfileWithVersion{}, nil
}

// BackupChrome moves the chrome/ folder to chrome.bak/
func (ffp FirefoxProfile) BackupChrome() error {
	return RenameIfExists(filepath.Join(ffp.Path, "chrome"), filepath.Join(ffp.Path, "chrome.bak"))
}

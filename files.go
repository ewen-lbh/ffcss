package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/hoisie/mustache"
	"gopkg.in/yaml.v2"
)

func GOOStoOS(GOOS string) string {
	switch GOOS {
	case "darwin":
		return "macos"
	case "plan9":
		return "linux"
	default:
		return GOOS
	}
}

func CurrentThemeByProfile() (map[string]string, error) {
	currentThemesRaw, err := os.ReadFile(ConfigDir("currently.yaml"))
	if os.IsNotExist(err) {
		err = os.WriteFile(ConfigDir("currently.yaml"), []byte(""), 0777)
		if err != nil {
			return nil, fmt.Errorf("while creating current themes list file: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("while reading current themes list: %w", err)
	}

	currentThemes := make(map[string]string)
	yaml.Unmarshal(currentThemesRaw, &currentThemes)
	return currentThemes, nil
}

// InstallAssets installs the assets in the specified profile directory
func (m Manifest) InstallAssets(operatingSystem string, variant Variant, profileDir string) (err error) {
	files, err := m.AssetsPaths(operatingSystem, variant, profileDir)
	if err != nil {
		return fmt.Errorf("while gathering assets: %w", err)
	}
	d("gathered %d asset(s)", len(files))

	for _, file := range files {
		stat, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("couldn't check file %s: %w", file, err)
		}

		if stat.IsDir() {
			continue
		}

		content, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("while reading %s: %w", file, err)
		}

		destPath, err := m.DestinationPathOfAsset(file, profileDir, operatingSystem, variant)
		if err != nil {
			println(err.Error())
			continue
		}

		err = os.MkdirAll(filepath.Dir(destPath), 0700)
		if err != nil {
			return fmt.Errorf("couldn't create parent directories for %s: %w", destPath, err)
		}

		err = ioutil.WriteFile(destPath, content, 0700)
		if err != nil {
			return fmt.Errorf("while writing to %s: %w", destPath, err)
		}
		d("wrote %s", destPath)

	}
	return nil
}

// AssetsPaths returns the individual file paths of all assets
func (m Manifest) AssetsPaths(os string, variant Variant, profileDirectory string) ([]string, error) {
	resolvedFiles := make([]string, 0)
	for _, template := range m.Assets {
		glob := RenderFileTemplate(template, os, variant, m.OSNames)
		glob = filepath.Clean(filepath.Join(m.DownloadedTo, glob))
		files, err := doublestar.Glob(glob)
		if err != nil {
			return resolvedFiles, fmt.Errorf("while getting all matches of glob %s: %w", glob, err)
		}
		// If no matches
		if len(files) < 1 {
			// If it's _really_ a glob pattern
			if strings.Contains(glob, "*") {
				return resolvedFiles, fmt.Errorf("glob pattern %s matches no files", glob)
				// If it's just a regular file (that was treated as a glob pattern)
			} else {
				return resolvedFiles, fmt.Errorf("file %s not found", glob)
			}
		}
		// For each match of the glob pattern
		resolvedFiles = append(resolvedFiles, files...)
	}
	return resolvedFiles, nil
}

// InstallUserJS installs the content of user.js and the config entries to {{profileDir}}/user.js
func (m Manifest) InstallUserJS(operatingSystem string, variant Variant, profileDir string) error {
	err := RenameIfExists(filepath.Join(profileDir, "user.js"), filepath.Join(profileDir, "user.js.bak"))
	if err != nil {
		return fmt.Errorf("while creating backup of %s: %w", filepath.Join(profileDir, "user.js"), err)
	}

	var content []byte

	if m.UserJS != "" {
		file := filepath.Join(m.DownloadedTo, RenderFileTemplate(m.UserJS, operatingSystem, variant, m.OSNames))
		content, err = ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("while reading %s: %w", file, err)
		}

	} else {
		content = []byte{}
	}

	additionalContent, err := m.UserJSFileContent()
	if err != nil {
		return fmt.Errorf("while translating config entries to javascript: %w", err)
	}

	if additionalContent != "" {
		content = []byte(string(content) + "\n" + additionalContent)
		d("generated additional user.js content from config entries: %q", additionalContent)
	}

	if string(content) == "" {
		return nil
	}

	err = ioutil.WriteFile(filepath.Join(profileDir, "user.js"), content, 0700)
	if err != nil {
		return fmt.Errorf("while writing: %w", err)
	}

	d("installed user.js @ %s", filepath.Join(profileDir, "user.js"))

	return nil
}

// InstallUserChrome writes the content of userChrome to {{profileDir}}/chrome/userChrome.css
func (m Manifest) InstallUserChrome(os string, variant Variant, profileDir string) error {
	if m.UserChrome == "" {
		return nil
	}
	file := filepath.Join(m.DownloadedTo, RenderFileTemplate(m.UserChrome, os, variant, m.OSNames))
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("while reading %s: %w", file, err)
	}

	err = ioutil.WriteFile(filepath.Join(profileDir, "chrome", "userChrome.css"), content, 0700)
	if err != nil {
		return fmt.Errorf("while writing: %w", err)
	}

	d("installed userChrome.css @ %s", filepath.Join(profileDir, "chrome", "userChrome.css"))

	return nil
}

// InstallUserContent writes the content of userContent to {{profileDir}}/chrome/userContent.css
func (m Manifest) InstallUserContent(os string, variant Variant, profileDir string) error {
	if m.UserContent == "" {
		return nil
	}
	file := filepath.Join(m.DownloadedTo, RenderFileTemplate(m.UserContent, os, variant, m.OSNames))
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("while reading %s: %w", file, err)
	}

	err = ioutil.WriteFile(filepath.Join(profileDir, "chrome", "userContent.css"), content, 0700)
	if err != nil {
		return fmt.Errorf("while writing: %w", err)
	}

	d("installed userContent.css @ %s", filepath.Join(profileDir, "chrome", "userContent.css"))

	return nil
}

// DestinationPathOfAsset computes the destination path of some asset from its path and the destination profile directory
// It is assumed that assetPath is absolute.
func (m Manifest) DestinationPathOfAsset(assetPath string, profileDir string, operatingSystem string, variant Variant) (string, error) {
	if !strings.HasPrefix(assetPath, m.DownloadedTo) {
		return "", fmt.Errorf("asset %q is outside of the theme's root %q", assetPath, m.DownloadedTo)
	}

	relativeTo := filepath.Clean(filepath.Join(m.DownloadedTo, filepath.Clean(RenderFileTemplate(m.CopyFrom, operatingSystem, variant, m.OSNames))))
	if !strings.HasPrefix(relativeTo, m.DownloadedTo) {
		return "", fmt.Errorf("copy from %q is outside of the theme's root %q", relativeTo, m.DownloadedTo)
	}

	relativised, err := filepath.Rel(relativeTo, assetPath)
	if err != nil {
		return "", fmt.Errorf("couldn't make %s relative to %s: %w", assetPath, filepath.Join(m.DownloadedTo, filepath.Clean(m.CopyFrom)), err)
	}

	return filepath.Join(profileDir, "chrome", relativised), nil
}

func RenderFileTemplate(f FileTemplate, operatingSystem string, variant Variant, osRenameMap map[string]string) string {
	if strings.Contains(strings.Trim(f, " "), "{{variant}}") && variant.Name == "" {
		warn("%q uses {{variant}} which is empty\n", f)
	}
	var osName string
	if osRenameMap[operatingSystem] == "" {
		osName = operatingSystem
	} else {
		osName = osRenameMap[operatingSystem]
	}
	return mustache.Render(f, map[string]string{
		"os":      osName,
		"variant": variant.Name,
	})
}

func SwitchGitBranch(newBranch, clonedTo string) error {
	process := exec.Command("git", "switch", newBranch)
	process.Dir = clonedTo
	output, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, output)
	}
	return nil
}

func SwitchGitCommit(commitSHA, clonedTo string) error {
	process := exec.Command("git", "checkout", commitSHA)
	process.Dir = clonedTo
	output, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, output)
	}
	return nil
}

func SwitchGitTag(tagName, clonedTo string) error {
	process := exec.Command("git", "fetch", "--all", "--tags")
	process.Dir = clonedTo
	output, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, output)
	}

	process = exec.Command("git", "checkout", "tags/"+tagName)
	process.Dir = clonedTo
	output, err = process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, output)
	}
	return nil
}

type FirefoxProfile struct {
	ID   string
	Name string
	Path string
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

func FirefoxProfileFromPath(path string) FirefoxProfile {
	base := filepath.Base(path)
	parts := strings.SplitN(base, ".", 2)
	return FirefoxProfile{
		Path: path,
		ID:   parts[0],
		Name: parts[1],
	}
}

func FirefoxProfileFromDisplayString(displayString string, profilePaths []string) FirefoxProfile {
	for _, profilePath := range profilePaths {
		ffp := FirefoxProfileFromPath(profilePath)
		if ffp.String() == displayString {
			return ffp
		}
	}
	d("while searching for %s in %v", displayString, profilePaths)
	panic("internal error: can't get profile from display string")
}

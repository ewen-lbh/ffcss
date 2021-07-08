package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/hoisie/mustache"
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

// InstallAssets installs the assets in the specified profile directory
func (m Manifest) InstallAssets(operatingSystem string, variant Variant, profileDir string) (err error) {
	files, err := m.AssetsPaths(operatingSystem, variant, profileDir)
	if err != nil {
		return fmt.Errorf("while gathering assets: %w", err)
	}

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

		err = os.MkdirAll(path.Dir(destPath), 0700)
		if err != nil {
			return fmt.Errorf("couldn't create parent directories for %s: %w", destPath, err)
		}

		err = ioutil.WriteFile(destPath, content, 0700)
		if err != nil {
			return fmt.Errorf("while writing to %s: %w", destPath, err)
		}

	}
	return nil
}

// AssetsPaths returns the individual file paths of all assets
func (m Manifest) AssetsPaths(os string, variant Variant, profileDirectory string) ([]string, error) {
	resolvedFiles := make([]string, 0)
	for _, template := range m.Assets {
		glob := RenderFileTemplate(template, os, variant, m.OSNames)
		glob = path.Clean(filepath.Join(m.DownloadPath(), glob))
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

	if m.UserJS == "" {
		return nil
	}

	file := filepath.Join(m.DownloadPath(), RenderFileTemplate(m.UserJS, operatingSystem, variant, m.OSNames))
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("while reading %s: %w", file, err)
	}

	additionalContent, err := m.UserJSFileContent()
	if err != nil {
		return fmt.Errorf("while translating config entries to javascript: %w", err)
	}

	content = []byte(string(content) + "\n" + additionalContent)
	err = ioutil.WriteFile(filepath.Join(profileDir, "user.js"), content, 0700)
	if err != nil {
		return fmt.Errorf("while writing: %w", err)
	}

	return nil
}

// InstallUserChrome writes the content of userChrome to {{profileDir}}/chrome/userChrome.css
func (m Manifest) InstallUserChrome(os string, variant Variant, profileDir string) error {
	if m.UserChrome == "" {
		return nil
	}
	file := filepath.Join(m.DownloadPath(), RenderFileTemplate(m.UserChrome, os, variant, m.OSNames))
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("while reading %s: %w", file, err)
	}

	err = ioutil.WriteFile(filepath.Join(profileDir, "chrome", "userChrome.css"), content, 0700)
	if err != nil {
		return fmt.Errorf("while writing: %w", err)
	}

	return nil
}

// InstallUserContent writes the content of userContent to {{profileDir}}/chrome/userContent.css
func (m Manifest) InstallUserContent(os string, variant Variant, profileDir string) error {
	if m.UserContent == "" {
		return nil
	}
	file := filepath.Join(m.DownloadPath(), RenderFileTemplate(m.UserContent, os, variant, m.OSNames))
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("while reading %s: %w", file, err)
	}

	err = ioutil.WriteFile(filepath.Join(profileDir, "chrome", "userContent.css"), content, 0700)
	if err != nil {
		return fmt.Errorf("while writing: %w", err)
	}

	return nil
}

// DestinationPathOfAsset computes the destination path of some asset from its path and the destination profile directory
// It is assumed that assetPath is absolute.
func (m Manifest) DestinationPathOfAsset(assetPath string, profileDir string, operatingSystem string, variant Variant) (string, error) {
	if !strings.HasPrefix(assetPath, m.DownloadPath()) {
		return "", fmt.Errorf("asset %q is outside of the theme's root %q", assetPath, m.DownloadPath())
	}

	relativeTo := path.Clean(filepath.Join(m.DownloadPath(), filepath.Clean(RenderFileTemplate(m.CopyFrom, operatingSystem, variant, m.OSNames))))
	if !strings.HasPrefix(relativeTo, m.DownloadPath()) {
		return "", fmt.Errorf("copy from %q is outside of the theme's root %q", relativeTo, m.DownloadPath())
	}

	relativised, err := filepath.Rel(relativeTo, assetPath)
	if err != nil {
		return "", fmt.Errorf("couldn't make %s relative to %s: %w", assetPath, filepath.Join(m.DownloadPath(), filepath.Clean(m.CopyFrom)), err)
	}

	return filepath.Join(profileDir, "chrome", relativised), nil
}

func RenderFileTemplate(f FileTemplate, operatingSystem string, variant Variant, osRenameMap map[string]string) string {
	if strings.Contains(strings.Trim(f, " "), "{{variant}}") && variant.Name == "" {
		fmt.Printf("WARNING: %q uses {{variant}} which is empty\n", f)
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

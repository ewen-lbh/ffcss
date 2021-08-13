package ffcss

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// InstallAssets installs the assets in the specified profile directory
func (t Theme) InstallAssets(operatingSystem string, variant Variant, profileDir string) (err error) {
	files, err := t.AssetsPaths(operatingSystem, variant)
	if err != nil {
		return fmt.Errorf("while gathering assets: %w", err)
	}
	LogDebug("gathered %d asset(s)", len(files))

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

		destPath, err := t.DestinationPathOfAsset(file, profileDir, operatingSystem, variant)
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
		LogDebug("wrote %s", destPath)

	}
	return nil
}

// InstallUserJS installs the content of user.js and the config entries to {{profileDir}}/user.js
func (t Theme) InstallUserJS(operatingSystem string, variant Variant, profileDir string) error {
	var content []byte
	var err error

	if t.UserJS != "" {
		file := filepath.Join(t.DownloadedTo, renderFileTemplate(t.UserJS, operatingSystem, variant, t.OSNames))
		content, err = ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("while reading %s: %w", file, err)
		}

	} else {
		content = []byte{}
	}

	additionalContent, err := t.UserJSFileContent()
	if err != nil {
		return fmt.Errorf("while translating config entries to javascript: %w", err)
	}

	if additionalContent != "" {
		content = []byte(string(content) + "\n" + additionalContent)
		LogDebug("generated additional user.js content from config entries: %q", additionalContent)
	}

	if string(content) == "" {
		return nil
	}

	err = ioutil.WriteFile(filepath.Join(profileDir, "user.js"), content, 0700)
	if err != nil {
		return fmt.Errorf("while writing: %w", err)
	}

	LogDebug("installed user.js @ %s", filepath.Join(profileDir, "user.js"))

	return nil
}

// InstallUserChrome writes the content of userChrome to {{profileDir}}/chrome/userChrome.css
func (t Theme) InstallUserChrome(os string, variant Variant, profileDir string) error {
	if t.UserChrome == "" {
		return nil
	}
	file := filepath.Join(t.DownloadedTo, renderFileTemplate(t.UserChrome, os, variant, t.OSNames))
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("while reading %s: %w", file, err)
	}

	err = ioutil.WriteFile(filepath.Join(profileDir, "chrome", "userChrome.css"), content, 0700)
	if err != nil {
		return fmt.Errorf("while writing: %w", err)
	}

	LogDebug("installed userChrome.css @ %s", filepath.Join(profileDir, "chrome", "userChrome.css"))

	return nil
}

// InstallUserContent writes the content of userContent to {{profileDir}}/chrome/userContent.css
func (t Theme) InstallUserContent(os string, variant Variant, profileDir string) error {
	if t.UserContent == "" {
		return nil
	}
	file := filepath.Join(t.DownloadedTo, renderFileTemplate(t.UserContent, os, variant, t.OSNames))
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("while reading %s: %w", file, err)
	}

	err = ioutil.WriteFile(filepath.Join(profileDir, "chrome", "userContent.css"), content, 0700)
	if err != nil {
		return fmt.Errorf("while writing: %w", err)
	}

	LogDebug("installed userContent.css @ %s", filepath.Join(profileDir, "chrome", "userContent.css"))

	return nil
}

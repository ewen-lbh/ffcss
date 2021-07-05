package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/docopt/docopt-go"
)

// RunCommandUse runs the command "use"
func RunCommandUse(args docopt.Opts) error {
	themeName, _ := args.String("THEME_NAME")
	// variant, _ := args.String("VARIANT")
	err := os.MkdirAll(path.Join(GetConfigDir(), "themes"), 0777)
	if err != nil {
		return fmt.Errorf("couldn't create data directories: %w", err)
	}
	uri, typ := ResolveURL(themeName)
	manifest, err := Download(uri, typ)
	if err != nil {
		return err
	}
	// TODO choose the profile directory, could have a smart default (based on {{profileDirectory}}/times.json:firstUse)
	profileDirs, err := ProfileDirsPaths()
	if err != nil {
		return fmt.Errorf("couldn't get mozilla profile directories: %w", err)
	}
	// Create missing chrome folder
	if _, err := os.Stat(path.Join(profileDirs[0], "chrome")); os.IsNotExist(err) {
		err := os.Mkdir(profileDirs[0], 0700)
		if err != nil {
			return err
		}
	}
	// TODO: prompt
	variant := Variant{Name: "blue"}
	os := GOOStoOS(runtime.GOOS)
	// Copy into it
	allFiles, err := manifest.AllFiles(os, variant)
	if err != nil {
		return fmt.Errorf("couldn't resolve files list: %w", err)
	}
	chromeDirs := make([]string, 0)
	for _, profileDir := range profileDirs {
		chromeDirs = append(chromeDirs, path.Join(profileDir, "chrome"))
	}
	err = CopyOver(allFiles, chromeDirs)
	if err != nil {
		return fmt.Errorf("couldn't copy to firefox profile directories: %w", err)
	}
	CleanDownloadArea()
	return nil
}

// IsURLClonable determines if the given URL points to a git repository
func IsURLClonable(URL string) (clonable bool, err error) {
	output, err := exec.Command("git", "ls-remote", URL).CombinedOutput()
	if err == nil {
		return true, nil
	}
	switch err.(type) {
	case *exec.ExitError:
		if err.(*exec.ExitError).ExitCode() == 128 {
			return false, nil
		}
	}
	return false, fmt.Errorf("while running git-ls-remote: %w: %s", err, output)
}

// GetManifest returns a Manifest from the manifest file of themeRoot
func GetManifest(themeRoot string) (Manifest, error) {
	if _, err := os.Stat(GetManifestPath(themeRoot)); os.IsExist(err) {
		return LoadManifest(GetManifestPath(themeRoot))
	} else {
		return Manifest{}, errors.New("the project has no manifest file")
	}
}

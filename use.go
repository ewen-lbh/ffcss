package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"

	"github.com/docopt/docopt-go"
	"github.com/evilsocket/islazy/zip"
)

//
// # clone the repo
// # get the manifest
// # read it
// # move required files to ~/.config/ffcss/themes/...
//   where ... is either ./@OWNER/REPO (for github themes)
//   or ./THEME_NAME (for built-in themes)
//   or ./-DOMAIN.TLD/THEME_NAME
//

var (
	PatternManifestBasename = regexp.MustCompile(`^ffcss.ya?ml$`)
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

// DownloadRepository downloads the repository at URL to {{cloneTo}}/{{ffcss.yaml:name}}
// It first clones the repo to tempCloneTo, then loads the manifest to determine where to move it.
// the manifest can be provided in case the repository does not contain it.
func DownloadRepository(URL string, tempCloneTo string, cloneTo string, themeManifest ...Manifest) (manifest Manifest, err error) {
	hasManifest := len(themeManifest) >= 1
	if hasManifest {
		manifest = themeManifest[0]
	}
	clonable, err := IsURLClonable(URL)
	if err != nil {
		return manifest, fmt.Errorf("while determining clonability: %w", err)
	}
	err = os.MkdirAll(cloneTo, 0777)
	if err != nil {
		return manifest, fmt.Errorf("could not create directory to download to: %w", err)
	}
	err = os.MkdirAll(tempCloneTo, 0777)
	if err != nil {
		return manifest, fmt.Errorf("could not create directory to clone to: %w", err)
	}
	if clonable {
		process := exec.Command("git", "clone", URL, tempCloneTo, "--depth=1")
		//TODO print this in verbose mode: fmt.Printf("DEBUG $ %s\n", process.String())
		output, err := process.CombinedOutput()
		if err != nil {
			return manifest, fmt.Errorf("%w: %s", err, output)
		}

	} else {
		return manifest, fmt.Errorf("does not point to a clonable git repository")
	}
	if !hasManifest {
		manifest, err = LoadManifest(path.Join(tempCloneTo, "ffcss.yaml"))
		if _, err := os.Stat(path.Join(tempCloneTo, "ffcss.yaml")); os.IsNotExist(err) {
			return manifest, fmt.Errorf("no manifest found: %w", err)
		}
		if err != nil {
			return manifest, fmt.Errorf("could not load manifest: %w", err)
		}
	}
	if manifest.Name() == "" {
		return manifest, errors.New("manifest has no name")
	}
	os.Rename(tempCloneTo, path.Join(cloneTo, manifest.Name()))
	return
}

// DownloadFromZip downloads a ffcss manifest files along with its resources from the given URL.
// The URL must point to a zip file that contains a ffcss.yaml in its root.
// The zip file will be downloaded and extracted to {{tempDownloadTo}}, then, after loading the manifest,
// the folder will then be moved to {{finalDownloadTo}}/{{ffcss.yaml:name}}.
// the manifest can be provided in case the zip does not contain it.
func DownloadFromZip(URL string, tempDownloadTo string, finalDownlodTo string, themeManifest ...Manifest) (manifest Manifest, err error) {
	tempDownloadTo = tempDownloadTo + "/theme.zip"
	hasManifest := len(themeManifest) >= 1
	if hasManifest {
		manifest = themeManifest[0]
	}

	// Check if file exists, has the right Content-Type, etc.
	head, err := http.Head(URL)
	if err != nil {
		return manifest, fmt.Errorf("couldn't check remote file: %w", err)
	}
	if head.StatusCode >= 400 {
		return manifest, fmt.Errorf("couldn't check remote file: server returned %s", head.Status)
	}
	if head.Header.Get("Content-Type") != "application/zip" {
		return manifest, fmt.Errorf("expected a zip file (application/zip), got a %s", head.Header.Get("Content-Type"))
	}

	// Download it
	process := exec.Command("wget", URL, "-O", tempDownloadTo)
	output, err := process.CombinedOutput()
	if err != nil {
		return manifest, fmt.Errorf("couldn't download zip file: %w: %s", err, output)
	}

	// Unzip it, check contents
	unzipped, err := zip.Unzip(tempDownloadTo, path.Dir(tempDownloadTo))

	if !hasManifest {
		for _, file := range unzipped {
			if PatternManifestBasename.MatchString(path.Base(file)) {
				hasManifest = true
				manifest, err = LoadManifest(file)
				if err != nil {
					return manifest, fmt.Errorf("couldn't load the manifest file: %w", err)
				}
				break
			}
		}
		if !hasManifest {
			os.RemoveAll(path.Dir(tempDownloadTo))
			return manifest, errors.New("downloaded zip file has no manifest file (ffcss.yaml)")
		}
	}
	if manifest.Name() == "" {
		return manifest, errors.New("manifest has no name")
	}
	err = os.Rename(path.Dir(tempDownloadTo), CacheDir(manifest.Name()))
	if err != nil {
		return manifest, fmt.Errorf("could not move %s to %s: %w", path.Dir(tempDownloadTo), CacheDir(manifest.Name()), err)
	}
	return
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

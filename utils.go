package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ReadFileBytes reads the content of ``filepath`` and returns the contents as a byte array.
// It panics on any error.
func ReadFileBytes(filepath string) []byte {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return b
}

// ReadFile reads the content of ``filepath`` and returns the contents as a string.
func ReadFile(filepath string) string {
	return string(ReadFileBytes(filepath))
}

// isValidURL tests a string to determine if it is a well-structured url or not.
func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func Assert(t *testing.T, got interface{}, expected interface{}) {
	assert.Equal(t, expected, got)
}

// ExpandHomeDir expands the "~/" part of a path to the current user's home directory
func ExpandHomeDir(p string) string {
	usr, _ := user.Current()
	homedir := usr.HomeDir
	if p == "~" {
		// In case of "~", which won't be caught by the "else if"
		p = homedir
	} else if strings.HasPrefix(p, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		p = filepath.Join(homedir, p[2:])
	}
	return p
}

// GetConfigDir returns the absolute path of ffcss's configuration directory
func GetConfigDir() string {
	return ExpandHomeDir("~/.config/ffcss")
}

// GetCacheDir returns the temporary path for cloned repos and other stuff
func GetCacheDir() string {
	return ExpandHomeDir("~/.cache/ffcss/")
}

// CacheDir joins the cache directory with the given path segments
func CacheDir(pathSegments ...string) string {
	return path.Join(GetCacheDir(), path.Join(pathSegments...))
}

// ConfigDir joins the config directory with the given path segments
func ConfigDir(pathSegments ...string) string {
	return path.Join(GetConfigDir(), path.Join(pathSegments...))
}

// GetManifestPath returns the path of a theme's manifest file
func GetManifestPath(themeRoot string) string {
	return path.Join(themeRoot, "ffcss.yaml")
}

// ProfileDirsPaths returns an array of profile directories from ~/.mozilla.
// 0 arguments: the .mozilla folder is assumed to be ~/.mozilla.
// 1 argument: use the given .mozilla folder
// more arguments: panic.
func ProfileDirsPaths(dotMozilla ...string) ([]string, error) {
	var mozillaFolder string
	if len(dotMozilla) == 0 {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return []string{}, fmt.Errorf("couldn't get the current user's home directory: %s. Try to use --mozilla-dir", err)
		}
		mozillaFolder = path.Join(homedir, ".mozilla")
	} else if len(dotMozilla) == 1 {
		mozillaFolder = dotMozilla[0]
	} else {
		panic(fmt.Sprintf("received %d arguments, expected 0 or 1", len(dotMozilla)))
	}
	directories, err := os.ReadDir(path.Join(mozillaFolder, "firefox"))
	releasesPaths := make([]string, 0)
	patternReleaseID := regexp.MustCompile(`[a-z0-9]{8}\.default(-\w+)?`)
	if err != nil {
		return []string{}, fmt.Errorf("couldn't read ~/.mozilla/firefox: %w", err)
	}
	for _, releasePath := range directories {
		if patternReleaseID.MatchString(releasePath.Name()) {
			releasesPaths = append(releasesPaths, path.Join(mozillaFolder, "firefox", releasePath.Name()))
		}
	}
	return releasesPaths, nil
}

func cwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}

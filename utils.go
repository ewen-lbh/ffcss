package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/pelletier/go-toml"
)

// ReadFileBytes reads the content of ``filepath`` and returns the contents as a byte array
func ReadFileBytes(filepath string) []byte {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	return b
}

// ReadFile reads the content of ``filepath`` and returns the contents as a string
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

type Theme struct {
	Repository string
	Files      map[string][]string
	Config     map[string]interface{}
}

type ThemesList map[string]Theme

func ReadThemesList() ThemesList {
	themesList := ThemesList{}
	doc := ReadFileBytes("themes.toml")
	toml.Unmarshal(doc, &themesList)
	return themesList
}

func Assert(t *testing.T, got interface{}, expected interface{}) {
	if got != expected {
		t.Errorf("\nexpected: \n%s\n\ngot: \n%s", expected, got)
	}
}

// ExpandHomeDir expands the "~/" part of a path to the current user's home directory
func ExpandHomeDir(path string) string {
	usr, _ := user.Current()
	homedir := usr.HomeDir
	if path == "~" {
		// In case of "~", which won't be caught by the "else if"
		path = homedir
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		path = filepath.Join(homedir, path[2:])
	}
	return path
}

// GetConfigDir returns the absolute path of ffcss's configuration directory
func GetConfigDir() string {
	return ExpandHomeDir("~/.config/ffcss")
}

// GetTempDir returns the temporary path for cloned repos and other stuff
func GetTempDir() string {
	return ExpandHomeDir("~/.cache/ffcss/")
}

// GetManifestPath returns the path of a theme's ffcss.{json;toml,yaml}
func GetManifestPath(extension string) string {
	return "ffcss." + extension
}

// ReadManifest reads a manifest file given its filepath and returns a Theme struct
// func ReadManifest

// GetMozillaReleasesPaths returns an array of release directories from ~/.mozilla.
func GetMozillaReleasesPaths() ([]string, error) {
	directories, err := os.ReadDir(ExpandHomeDir("~/.mozilla/firefox/"))
	releasesPaths := make([]string, 0)
	patternReleaseID := regexp.MustCompile(`[a-z0-9]{8}\.default(-\w+)?`)
	if err != nil {
		return []string{}, fmt.Errorf("couldn't read ~/.mozilla/firefox: %s", err.Error())
	}
	for _, releasePath := range directories {
		if patternReleaseID.MatchString(releasePath.Name()) {
			releasesPaths = append(releasesPaths, ExpandHomeDir("~/.mozilla/firefox/")+"/"+releasePath.Name())
		}
	}
	return releasesPaths, nil
}

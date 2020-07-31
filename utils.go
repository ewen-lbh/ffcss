package main

import (
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"testing"
	"strings"

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

// isValidUrl tests a string to determine if it is a well-structured url or not.
func isValidUrl(toTest string) bool {
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

type ThemesList map[string]struct {
	repository string
	files      map[string][]string
	config     map[string]interface{}
}

func ReadThemesList() ThemesList {
	themesList := ThemesList{}
	doc := ReadFileBytes("themes.toml")
	toml.Unmarshal(doc, &themesList)
	return themesList
}

func Assert(t *testing.T, got interface{}, expected interface{}) {
	if got != expected {
		t.Errorf("expected: %s\ngot: %s", expected, got)
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

package main

import (
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"testing"
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

func Assert(t *testing.T, got interface{}, expected interface{}) {
	if got != expected {
		t.Errorf("\nexpected: \n%s\n\ngot: \n%s", expected, got)
	}
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

// GetTempDir returns the temporary path for cloned repos and other stuff
func GetTempDir() string {
	return ExpandHomeDir("~/.cache/ffcss/")
}

// GetManifestPath returns the path of a theme's manifest file
func GetManifestPath(themeRoot string) string {
	return path.Join(themeRoot, "ffcss.yaml")
}

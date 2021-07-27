package main

import (
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

// GetConfigDir returns the absolute path of ffcss's configuration directory
func GetConfigDir() string {
	homedir, _ := os.UserHomeDir()
	return filepath.Join(homedir, ".config", "ffcss")
}

// GetCacheDir returns the temporary path for cloned repos and other stuff
func GetCacheDir() string {
	homedir, _ := os.UserHomeDir()
	return filepath.Join(homedir, ".cache", "ffcss")
}

// CacheDir joins the cache directory with the given path segments
func CacheDir(pathSegments ...string) string {
	return filepath.Join(GetCacheDir(), filepath.Join(pathSegments...))
}

// ConfigDir joins the config directory with the given path segments
func ConfigDir(pathSegments ...string) string {
	return filepath.Join(GetConfigDir(), filepath.Join(pathSegments...))
}

// GetManifestPath returns the path of a theme's manifest file
func GetManifestPath(themeRoot string) string {
	return filepath.Join(themeRoot, "ffcss.yaml")
}

func cwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}

// isURLClonable determines if the given URL points to a git repository
func isURLClonable(URL string) bool {
	output, err := exec.Command("git", "ls-remote", URL).CombinedOutput()
	if err == nil {
		return true
	}
	switch err.(type) {
	case *exec.ExitError:
		if err.(*exec.ExitError).ExitCode() == 128 {
			return false
		}
	}
	warn("could not determine clonability of %s: while running git-ls-remote: %w: %s\n", URL, err, output)
	return false
}

// RenameIfExists renames from to to if from exists. If it doesn't, don't attempt renaming.
func RenameIfExists(from string, to string) error {
	if _, err := os.Stat(from); os.IsNotExist(err) {
		return nil
	}
	if _, err := os.Stat(to); os.IsNotExist(err) {
		return os.Rename(from, to)
	}
	err := os.RemoveAll(to)
	if err != nil {
		return err
	}
	return os.Rename(from, to)
}

// VimModeEnabled returns true if the user has explicitly set vim mode, or if the $EDITOR is vim/neovim
func VimModeEnabled() bool {
	if os.Getenv("VIM_MODE") == "1" || os.Getenv("VIM_STYLE") == "1" {
		return true
	}
	progname := filepath.Base(os.Getenv("EDITOR"))
	return progname == "vim" || progname == "nvim"
}

func apply(f func(string) string, in []string) []string {
	out := make([]string, 0, len(in))
	for _, item := range in {
		out = append(out, f(item))
	}
	return out
}

// prefixEachLine prepends each line of s with the provided prefix (with).
// Only supports UNIX-Style line endings (\n)
func prefixEachLine(s string, with string) string {
	var prefixedLines []string
	for _, line := range strings.Split(s, "\n") {
		prefixedLines = append(prefixedLines, with+line)
	}
	return strings.Join(prefixedLines, "\n")
}

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

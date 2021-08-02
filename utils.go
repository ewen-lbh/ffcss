package ffcss

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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

// getConfigDir returns the absolute path of ffcss's configuration directory
func getConfigDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("couldn't get your home directory: %w", err))
	}
	return filepath.Join(homedir, ".config", "ffcss")
}

// getCacheDir returns the temporary path for cloned repos and other stuff
func getCacheDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("couldn't get your home directory: %w", err))
	}
	return filepath.Join(homedir, ".cache", "ffcss")
}

// CacheDir joins the cache directory with the given path segments
func CacheDir(pathSegments ...string) string {
	return filepath.Join(getCacheDir(), filepath.Join(pathSegments...))
}

// ConfigDir joins the config directory with the given path segments
func ConfigDir(pathSegments ...string) string {
	return filepath.Join(getConfigDir(), filepath.Join(pathSegments...))
}

func cwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("couldn't get the current working directory: %w", err))
	}
	return wd
}

// renameIfExists renames from to to if from exists. If it doesn't, don't attempt renaming.
func renameIfExists(from string, to string) error {
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

// vimModeEnabled returns true if the user has explicitly set vim mode, or if the $EDITOR is vim/neovim
func vimModeEnabled() bool {
	if os.Getenv("VIM_MODE") == "1" || os.Getenv("VIM_STYLE") == "1" {
		return true
	}
	progname := filepath.Base(os.Getenv("EDITOR"))
	return progname == "vim" || progname == "nvim"
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

// GOOStoOS returns user-friendly OS names from a given GOOS.
// darwin becomes macos and plan9 becomes linux.
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

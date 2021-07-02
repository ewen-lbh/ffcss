package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

// ToUserJSFile returns a string of JS source code that represents config.
// It can be used directly to write a .mozilla/firefox/*.default-*/user.js file
func ToUserJSFile(config map[string]interface{}) (string, error) {
	lines := make([]string, 0)
	for name, value := range config {
		valueJSON, err := json.Marshal(value)
		if err != nil {
			return "", fmt.Errorf("can't serialize %#v: %s", value, err.Error())
		}
		lines = append(lines, fmt.Sprintf(`user_pref(%q, %s);`, name, string(valueJSON)))
	}
	return strings.Join(lines, "\n"), nil
}


// GetMozillaReleasesPaths returns an array of release directories from ~/.mozilla.
// 0 arguments: the .mozilla folder is assumed to be ~/.mozilla.
// 1 argument: use the given .mozilla folder
// more arguments: panic.
func GetMozillaReleasesPaths(dotMozilla... string) ([]string, error) {
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
		return []string{}, fmt.Errorf("couldn't read ~/.mozilla/firefox: %s", err.Error())
	}
	for _, releasePath := range directories {
		if patternReleaseID.MatchString(releasePath.Name()) {
			releasesPaths = append(releasesPaths, path.Join(mozillaFolder, "firefox", releasePath.Name()))
		}
	}
	return releasesPaths, nil
}

func (m Manifest) UserJSFileContent() (string, error) {
	return ToUserJSFile(m.Config)
}

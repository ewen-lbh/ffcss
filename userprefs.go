package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ToUserJS returns a string of JS source code that represents config.
// It can be used directly to write a .mozilla/firefox/*.default-*/user.js file
func ToUserJS(config map[string]interface{}) (string, error) {
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

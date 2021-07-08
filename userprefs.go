package main

import (
	"encoding/json"
	"fmt"

	"strings"
)

// ToUserJSFile returns a string of JS source code that represents config.
// It can be used directly to write a user.js file
func ToUserJSFile(config map[string]interface{}) (string, error) {
	lines := make([]string, 0)
	for name, value := range config {
		valueJSON, err := json.Marshal(value)
		if err != nil {
			return "", fmt.Errorf("can't serialize %#v: %w", value, err)
		}
		lines = append(lines, fmt.Sprintf(`user_pref(%q, %s);`, name, string(valueJSON)))
	}
	return strings.Join(lines, "\n"), nil
}

func (m Manifest) UserJSFileContent() (string, error) {
	return ToUserJSFile(m.Config)
}

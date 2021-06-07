package main

import (
	"encoding/json"
	"fmt"
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

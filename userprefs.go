package ffcss

import (
	"encoding/json"
	"fmt"
	"regexp"

	"strings"
)

// UserJSFileContent returns a string of JS source code that represents the theme's config.
// It can be used directly to write a user.js file.
func (t Theme) UserJSFileContent() (string, error) {
	lines := make([]string, 0)
	for name, value := range t.Config {
		valueJSON, err := json.Marshal(value)
		if err != nil {
			return "", fmt.Errorf("can't serialize %#v: %w", value, err)
		}
		lines = append(lines, fmt.Sprintf(`user_pref(%q, %s);`, name, string(valueJSON)))
	}
	return strings.Join(lines, "\n"), nil
}

// ValueOfUserPrefCall returns the value of configuration entry, given its key and the contents of
// the prefs.js file. It only works if the value is a JSON-parsable literal (string, number, boolean, null, etc.).
func ValueOfUserPrefCall(prefsJSContent []byte, key string) (string, error) {
	pattern := regexp.MustCompile(`(?m)^\s*user_pref\("` + key + `"\s*,\s*(.+)\)\s*;?\s*$`)
	matches := pattern.FindAllSubmatch(prefsJSContent, -1)
	if len(matches) == 0 {
		return "", fmt.Errorf("key %q not found", key)
	}
	raw := matches[len(matches)-1][1]
	var jsonParsed interface{}
	err := json.Unmarshal(raw, &jsonParsed)
	if err != nil {
		return "", fmt.Errorf("while intepreting value %q: %w", string(raw), err)
	}
	return fmt.Sprint(jsonParsed), nil
}

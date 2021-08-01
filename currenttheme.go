package ffcss

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func CurrentThemeByProfile() (map[string]string, error) {
	currentThemesRaw, err := os.ReadFile(ConfigDir("currently.yaml"))
	if os.IsNotExist(err) {
		err = os.WriteFile(ConfigDir("currently.yaml"), []byte(""), 0777)
		if err != nil {
			return nil, fmt.Errorf("while creating current themes list file: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("while reading current themes list: %w", err)
	}

	currentThemes := make(map[string]string)
	yaml.Unmarshal(currentThemesRaw, &currentThemes)
	return currentThemes, nil
}

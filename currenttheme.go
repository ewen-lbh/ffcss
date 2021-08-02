package ffcss

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// CurrentThemeByProfile returns a map mapping a profile path to its current theme's name.
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

// RegisterCurrentTheme writes the currently.yaml file in ffcss' configuration to update
// what ffcss considers to be the current theme for that profile.
func (ffp FirefoxProfile) RegisterCurrentTheme(themeName string) error {
	currentThemes, err := CurrentThemeByProfile()
	if err != nil {
		return err
	}
	currentThemes[ffp.FullName()] = themeName
	currentThemesNewContents, err := yaml.Marshal(currentThemes)
	if err != nil {
		return fmt.Errorf("while marshaling into YAML: %w", err)
	}

	err = os.WriteFile(ConfigDir("currently.yaml"), currentThemesNewContents, 0777)
	if err != nil {
		return fmt.Errorf("while writing new contents: %w", err)
	}

	return nil
}

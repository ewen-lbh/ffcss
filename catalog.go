package ffcss

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hbollon/go-edlib"
	"golang.org/x/text/unicode/norm"
)

// Catalog represents a collection of themes
type Catalog map[string]Theme

// Lookup looks up a theme by its name in the theme store.
// It also returns an error starting with "did you mean:" when
// a theme name is not found but themes with similar names exist.
func (store Catalog) Lookup(query string) (Theme, error) {
	originalQuery := query
	query = lookupPreprocess(query)
	LogDebug("using query %q", query)
	processedThemeNames := make([]string, 0, len(store))
	for _, theme := range store {
		LogDebug("\tlooking up against %q (%q)", lookupPreprocess(theme.Name()), theme.Name())
		if lookupPreprocess(theme.Name()) == query {
			return theme, nil
		}
		processedThemeNames = append(processedThemeNames, lookupPreprocess(theme.Name()))
	}
	// Use fuzzy search for did-you-mean errors
	suggestion, _ := edlib.FuzzySearchThreshold(query, processedThemeNames, 0.75, edlib.Levenshtein)

	if suggestion != "" {
		return Theme{}, fmt.Errorf("theme %q not found. did you mean [blue][bold]%s[reset]?", originalQuery, suggestion)
	}
	return Theme{}, fmt.Errorf("theme %q not found", originalQuery)
}

// lookupPreprocess applies transformations to s so that it can be compared
// to search for something.
// For example, it is used by (ThemeStore).Lookup
func lookupPreprocess(s string) string {
	return strings.ToLower(norm.NFKD.String(regexp.MustCompile(`[-_ .]`).ReplaceAllString(s, "")))
}

// LoadCatalog loads a directory of theme manifests.
// Keys are theme names (files' basenames with the .yaml removed).
func LoadCatalog(storeDirectory string) (themes Catalog, err error) {
	themeNamePattern := regexp.MustCompile(`^(.+)\.ya?ml$`)
	themes = make(Catalog)
	manifests, err := os.ReadDir(storeDirectory)
	if err != nil {
		return
	}
	LogDebug("loading potential themes %v into catalog", func() []string {
		dirNames := make([]string, 0, len(manifests))
		for _, dir := range manifests {
			dirNames = append(dirNames, dir.Name())
		}
		return dirNames
	}())
	for _, manifest := range manifests {
		if !themeNamePattern.MatchString(manifest.Name()) {
			continue
		}
		themeName := themeNamePattern.FindStringSubmatch(manifest.Name())[1]
		theme, err := LoadManifest(filepath.Join(storeDirectory, manifest.Name()))
		if err != nil {
			return nil, fmt.Errorf("while loading theme %q: %w", themeName, err)
		}
		LogDebug("\tadding theme from manifest %q", manifest.Name())
		themes[themeName] = theme
	}
	return
}

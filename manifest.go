package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Variant struct {
	Config      Config
	UserChrome  FileTemplate `yaml:"userChrome"`
	UserContent FileTemplate `yaml:"userContent"`
	UserJS      FileTemplate `yaml:"user.js"`
	Assets      []FileTemplate
	Description string
	Name        string
}

type Manifest struct {
	Repository   string
	ExplicitName string `yaml:"name"`
	FfcssVersion int    `yaml:"ffcss"`
	Variants     map[string]Variant
	CopyFrom     string `yaml:"copy from"`
	Config       Config
	UserChrome   FileTemplate `yaml:"userChrome"`
	UserContent  FileTemplate `yaml:"userContent"`
	UserJS       FileTemplate `yaml:"user.js"`
	Assets       []FileTemplate
}

func (m Manifest) Name() string {
	if m.ExplicitName != "" {
		return strings.ToLower(m.ExplicitName)
	}
	if strings.HasPrefix(m.Repository, "https://github.com") {
		fragments := strings.Split(m.Repository, "/")
		return strings.ToLower(fragments[len(fragments)-1])
	}
	return ""
}

func (m Manifest) DownloadPath() string {
	return CacheDir(m.Name())
}

type Config map[string]interface{}

type FileTemplate = string

func NewManifest() Manifest {
	return Manifest{
		Config: Config{
			"toolkit.legacyUserProfileCustomizations.stylesheets": true,
		},
		UserChrome:  "userChrome.css",
		UserContent: "userContent.css",
		UserJS:      "user.js",
		Variants:    map[string]Variant{},
		Assets:      []FileTemplate{},
	}
}

func (m Manifest) URL() url.URL {
	uri, _ := ResolveURL(m.Name())
	URL, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	return *URL
}

// LoadManifest loads a ffcss.yaml file into a Manifest object.
func LoadManifest(manifestPath string) (manifest Manifest, err error) {
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		err = fmt.Errorf("while reading manifest %s: %w", manifestPath, err)
		return
	}
	manifest = NewManifest()
	err = yaml.Unmarshal(raw, &manifest)
	for name, variant := range manifest.Variants {
		variant.Name = name
	}
	if err != nil {
		err = fmt.Errorf("while parsing manifest %s: %w", manifestPath, err)
		return
	}
	return
}

func (m Manifest) VariantsSlice() []Variant {
	variantsSlice := make([]Variant, 0, len(m.Variants))
	for _, variant := range m.Variants {
		variantsSlice = append(variantsSlice, variant)
	}
	return variantsSlice
}

// WithVariant returns a Manifest representing the theme if the selected variant
// was used as the "root values".
// i.e. the values of UserJS, UserContent, UserChrome, Assets are replaced with their variant's, if set,
// and the value of Config is combined with the variant's.
func (m Manifest) WithVariant(variant Variant) Manifest {
	newManifest := m
	if variant.UserChrome != "" {
		newManifest.UserChrome = variant.UserChrome
	}
	if variant.UserContent != "" {
		newManifest.UserContent = variant.UserContent
	}
	if variant.UserJS != "" {
		newManifest.UserJS = variant.UserJS
	}
	if len(variant.Assets) > 0 {
		newManifest.Assets = variant.Assets
	}
	for key, val := range variant.Config {
		newManifest.Config[key] = val
	}
	return newManifest
}

// ThemeStore represents a collection of themes
type ThemeStore = map[string]Manifest

// LoadThemeCatalog loads a directory of theme manifests.
// Keys are theme names (files' basenames with the .yaml removed).
func LoadThemeCatalog(storeDirectory string) (themes ThemeStore, err error) {
	themeNamePattern := regexp.MustCompile(`^(.+)\.ya?ml$`)
	themes = make(ThemeStore)
	manifests, err := os.ReadDir(storeDirectory)
	if err != nil {
		return
	}
	for _, manifest := range manifests {
		if !themeNamePattern.MatchString(manifest.Name()) {
			continue
		}
		themeName := themeNamePattern.FindStringSubmatch(manifest.Name())[1]
		theme, err := LoadManifest(filepath.Join(storeDirectory, manifest.Name()))
		if err != nil {
			return nil, err
		}
		themes[themeName] = theme
	}
	return
}

// AvailableVariants lists the possible variant names to choose from
func (m Manifest) AvailableVariants() []string {
	names := make([]string, 0, len(m.Variants))
	for name := range m.Variants {
		names = append(names, name)
	}
	return names
}

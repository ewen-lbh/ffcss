package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/glamour"
	"gopkg.in/yaml.v2"
)

type Variant struct {
	// Properties exclusive to variants
	Name    string
	Message string

	// Properties that modify the "default variant"
	Config      Config
	UserChrome  FileTemplate `yaml:"userChrome"`
	UserContent FileTemplate `yaml:"userContent"`
	UserJS      FileTemplate `yaml:"user.js"`
	Assets      []FileTemplate
	Description string
	Addons      []string
}

type Manifest struct {
	Repository   string
	ExplicitName string `yaml:"name"`
	FfcssVersion int    `yaml:"ffcss"`
	Variants     map[string]Variant
	OSNames      map[string]string `yaml:"os"`

	// Those can be modified by variant
	CopyFrom    string `yaml:"copy from"`
	Config      Config
	UserChrome  FileTemplate `yaml:"userChrome"`
	UserContent FileTemplate `yaml:"userContent"`
	UserJS      FileTemplate `yaml:"user.js"`
	Assets      []FileTemplate
	Message     string
	Addons      []string
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
		UserChrome:  "",
		UserContent: "",
		UserJS:      "",
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
		variantWithName := variant
		variantWithName.Name = name
		manifest.Variants[name] = variantWithName
	}
	if err != nil {
		err = fmt.Errorf("while parsing manifest %s: %w", manifestPath, err)
		return
	}
	return
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
	if variant.Message != "" {
		newManifest.Message = variant.Message
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

// ShowMessage renders the message and prints it to the user
func (m Manifest) ShowMessage() error {
	scheme := os.Getenv("COLORSCHEME")
	if scheme != "light" && scheme != "dark" {
		// TODO: detect with the terminal's current background color as a fallback
		scheme = "dark"
	}
	rendered, err := glamour.Render(m.Message, scheme)
	if err != nil {
		return fmt.Errorf("while rendering message: %w", err)
	}

	if strings.TrimSpace(rendered) != "" {
		fmt.Println(rendered)
	}
	return nil
}

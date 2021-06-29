package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"

	"github.com/hoisie/mustache"
	"gopkg.in/yaml.v2"
)

type Manifest struct {
	Repository   string
	Name         string
	FfcssVersion int `yaml:"ffcss"`
	Config       Config
	Variants     []string
	UserChrome   []FileTemplate `yaml:"userChrome"`
	UserContent  []FileTemplate `yaml:"userContent"`
	UserJS       []FileTemplate `yaml:"user.js"`
	Assets       []FileTemplate
}

type Config map[string]interface{}

type Variant struct {
	Name        string
	Description string
}

type FileTemplate = string

func RenderFileTemplate(f FileTemplate, os string, variant Variant) string {
	return mustache.Render(f, map[string]string{
		"os":      os,
		"variant": variant.Name,
	})
}

func NewManifest() Manifest {
	return Manifest{
		Config: Config{
			"toolkit.legacyUserProfileCustomizations.stylesheets": true,
		},
		Assets: []FileTemplate{"config/**"},
	}
}

func (m Manifest) URL() url.URL {
	uri, _ := ResolveThemeName(m.Name)
	URL, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	return *URL
}

func (m Manifest) DownloadPath() string {
	return GetThemeDownloadPath(m.URL())
}

// LoadManifest loads a ffcss.yaml file into a Manifest object.
func LoadManifest(manifestPath string) (manifest Manifest, err error) {
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		err = fmt.Errorf("while reading manifest %s: %s", manifestPath, err.Error())
		return
	}
	loaded := NewManifest()
	err = yaml.Unmarshal(raw, &loaded)
	if err != nil {
		err = fmt.Errorf("while parsing manifest %s: %s", manifestPath, err.Error())
		return
	}
	return
}

// TODO: code generation for explicit keys from themes.toml?

// ThemeStore represents a collection of themes
type ThemeStore = map[string]Manifest

// LoadThemeStore loads a directory of theme manifests.
// Keys are theme names (files' basenames with the .yaml removed).
func LoadThemeStore(storeDirectory string) (themes ThemeStore, err error) {
	themeNamePattern := regexp.MustCompile(`^(.+)\.ya?ml$`)
	themes = make(ThemeStore, 0)
	manifests, err := os.ReadDir(storeDirectory)
	if err != nil {
		return
	}
	for _, manifest := range manifests {
		if !themeNamePattern.MatchString(manifest.Name()) {
			continue
		}
		themeName := themeNamePattern.FindStringSubmatch(manifest.Name())[1]
		theme, err := LoadManifest(path.Join(storeDirectory, manifest.Name()))
		if err != nil {
			return themes, err
		}
		themes[themeName] = theme
	}
	return
}

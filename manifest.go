package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"reflect"
	"regexp"

	"github.com/hoisie/mustache"
	"gopkg.in/yaml.v2"
)

type ManifestRawFiles struct {
	Repository   string
	FfcssVersion int `yaml:"ffcss"`
	Config       Config
	Variants     []string
	Files        interface{}
}

type Manifest struct {
	Repository   string
	Name string
	FfcssVersion int `yaml:"ffcss"`
	Config       Config
	Variants     []string
	Files        []File
}

type Config map[string]interface{}

type Variant struct {
	Name        string
	Description string
}

type File struct {
	Name string
	OS   string
}

func (f File) Render(os string, variant Variant) string {
	if os != f.OS && f.OS != "" {
		return ""
	}
	return mustache.Render(f.Name, map[string]string{
		"os":      os,
		"variant": variant.Name,
	})
}

func NewManifest() Manifest {
	return Manifest{
		Config: Config{
			"toolkit.legacyUserProfileCustomizations.stylesheets": true,
		},
		Files: []File{{Name: "config/**"}},
	}
}

func newManifestFromManifestRawFiles(m ManifestRawFiles) Manifest {
	manifest := NewManifest()
	manifest.Config = m.Config
	manifest.FfcssVersion = m.FfcssVersion
	manifest.Repository = m.Repository
	manifest.Variants = m.Variants
	return manifest
}

// Resolve resolves the two different forms of files into a []File
func (m ManifestRawFiles) Resolve() (manifest Manifest, err error) {
	manifest = newManifestFromManifestRawFiles(m)
	if reflect.TypeOf(m.Files) == nil {
		return
	}
	switch reflect.TypeOf(m.Files).Kind() {
	case reflect.Array, reflect.Slice:
		// Remove default files
		manifest.Files = make([]File, 0)
		for _, elem := range m.Files.([]interface{}) {
			manifest.Files = append(manifest.Files, File{Name: fmt.Sprint(elem)})
		}
	case reflect.Map:
		for os, filesArray := range m.Files.(map[string][]interface{}) {
			for _, elem := range filesArray {
				filesArray = append(filesArray, File{
					Name: fmt.Sprint(elem),
					OS:   os,
				})
			}
		}
	default:
		err = errors.New("files should be an array or an object of arrays")
	}
	return
}

func (m Manifest) URL() url.URL {
	uri, typ := ResolveThemeName(m.Name)
	URL, err := url.Parse(uri)
}

// LoadManifest loads a ffcss.yaml file into a Manifest object.
func LoadManifest(manifestPath string) (manifest Manifest, err error) {
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		err = fmt.Errorf("while reading manifest %s: %s", manifestPath, err.Error())
		return
	}
	loaded := ManifestRawFiles{}
	err = yaml.Unmarshal(raw, &loaded)
	if err != nil {
		err = fmt.Errorf("while parsing manifest %s: %s", manifestPath, err.Error())
		return
	}
	manifest, err = loaded.Resolve()
	if err != nil {
		err = fmt.Errorf("while parsing files in manifest %s: %s", manifestPath, err.Error())
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

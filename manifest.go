package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hoisie/mustache"
	"gopkg.in/yaml.v2"
)

type Theme struct {
	Config      Config
	UserChrome  FileTemplate `yaml:"userChrome"`
	UserContent FileTemplate `yaml:"userContent"`
	UserJS      FileTemplate `yaml:"user.js"`
	Assets      []FileTemplate
}

type Manifest struct {
	Repository   string
	ExplicitName string `yaml:"name"`
	FfcssVersion int    `yaml:"ffcss"`
	Variants     map[string]Theme
	Config       Config
	UserChrome   FileTemplate `yaml:"userChrome"`
	UserContent  FileTemplate `yaml:"userContent"`
	UserJS       FileTemplate `yaml:"user.js"`
	Assets       []FileTemplate
}

func (m Manifest) Name() string {
	if m.ExplicitName != "" {
		return m.ExplicitName
	}
	if strings.HasPrefix(m.Repository, "https://github.com") {
		fragments := strings.Split(m.Repository, "/")
		return fragments[len(fragments)-1]
	}
	return ""
}

// AllFileTemplates concatenates all file templates, in copy order (last in array should be copied over last)
func (m Manifest) AllFileTemplates() []FileTemplate {
	return append(
		[]FileTemplate{
			m.UserChrome,
			m.UserContent,
			m.UserJS,
		},
		m.Assets...,
	)
}

// AllFiles returns all of the file paths (relative to the repository's root)
func (m Manifest) AllFiles(os string, variant Variant) ([]string, error) {
	resolvedFiles := make([]string, 0)
	for _, template := range m.AllFileTemplates() {
		glob := RenderFileTemplate(template, os, variant)
		resolved, err := filepath.Glob(path.Join(m.DownloadPath(), glob))
		if err != nil {
			return []string{}, fmt.Errorf("malformed glob pattern %q: %w", glob, err)
		}
		resolvedFiles = append(resolvedFiles, resolved...)
	}
	return resolvedFiles, nil
}

func (m Manifest) DownloadPath() string {
	return CacheDir(m.Name())
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
		UserChrome:  "userChrome.css",
		UserContent: "userContent.css",
		UserJS:      "user.js",
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
	if err != nil {
		err = fmt.Errorf("while parsing manifest %s: %w", manifestPath, err)
		return
	}
	return
}

// TODO: code generation for explicit keys from themes.toml?

// ThemeStore represents a collection of themes
type ThemeStore = map[string]Manifest

// LoadThemeCatalog loads a directory of theme manifests.
// Keys are theme names (files' basenames with the .yaml removed).
func LoadThemeCatalog(storeDirectory string) (themes ThemeStore, err error) {
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
			return nil, err
		}
		themes[themeName] = theme
	}
	return
}

package main

import (
	"github.com/hoisie/mustache"
)

type Manifest struct {
	ManifestVersion int `yaml:"ffcss"`
	Config          Config
	Variants        []string
	Files           []File
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
	}
}

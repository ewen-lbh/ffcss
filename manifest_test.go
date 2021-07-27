package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testarea = filepath.Join(cwd(), "testarea")

func TestRenderFileTemplate(t *testing.T) {
	Assert(t, RenderFileTemplate(
		"userChrome.css",
		"linux",
		Variant{},
		map[string]string{},
	), "userChrome.css")

	Assert(t, RenderFileTemplate(
		"linux.css",
		"linux",
		Variant{},
		map[string]string{"linux": "Linux"},
	), "linux.css")

	Assert(t, RenderFileTemplate(
		"{{os}}.css",
		"linux",
		Variant{},
		map[string]string{"linux": "GNU/Linux"},
	), "GNU/Linux.css")

	Assert(t, RenderFileTemplate(
		"linux.css",
		"windows",
		Variant{},
		map[string]string{},
	), "linux.css")

	Assert(t, RenderFileTemplate(
		"./{{ os }}/{{variant}}.css",
		"macos",
		Variant{Name: "rainbow"},
		map[string]string{},
	), "./macos/rainbow.css")

}

func TestAssetsPaths(t *testing.T) {
	simplerentfox := Manifest{
		ExplicitName: "simplerentfox",
		DownloadedTo: CacheDir("simplerentfox/_"),
		Variants: map[string]Variant{
			"WithoutURLBar": {},
			"OneLine":       {},
		},
		Assets: []string{"./{{ os }}/userChrome__{{ variant }}.css"}, // for the purposes of testing
		OSNames: map[string]string{
			"linux":   "Linux",
			"macos":   "Linux",
			"windows": "Windows",
		},
	}

	files, err := simplerentfox.AssetsPaths("linux", Variant{Name: "blue"})
	assert.Regexp(t, "file .+ not found", err.Error())
	assert.Equal(t, []string{}, files)

	files, err = simplerentfox.AssetsPaths("linux", Variant{Name: "OneLine"})
	assert.NoError(t, err)
	assert.Equal(t, []string{CacheDir("simplerentfox/_/Linux/userChrome__OneLine.css")}, files)
}

func TestDestinationPathOf(t *testing.T) {
	manifest := Manifest{
		DownloadedTo: CacheDir("simplerentfox/_"),
		ExplicitName: "materialfox",
		Variants:     map[string]Variant{},
		Config:       Config{},
	}

	file, err := manifest.DestinationPathOfAsset("/home/ewen/lol.pdf", testarea, "linux", Variant{})
	if assert.Error(t, err) {
		assert.Regexp(t, `asset ".+" is outside of the theme's root ".+"`, err.Error())
	}
	assert.Equal(t, "", file)

	file, err = manifest.DestinationPathOfAsset("testarea/home/.cache/ffcss/simplerentfox/../../../lol.pdf", testarea, "linux", Variant{})
	if assert.Error(t, err) {
		assert.Regexp(t, `asset ".+" is outside of the theme's root ".+"`, err.Error())
	}
	assert.Equal(t, "", file)

	manifest = Manifest{
		DownloadedTo: CacheDir("simplerentfox/_"),
		ExplicitName: "simplerentfox",
		Variants: map[string]Variant{
			"WithoutURLBar": {},
			"OneLine":       {},
		},
		CopyFrom: "{{ os }}/",
		OSNames: map[string]string{
			"linux": "Linux",
		},
	}

	file, err = manifest.DestinationPathOfAsset(CacheDir("simplerentfox/_/Linux/userChrome__OneLine.css"), testarea, "linux", Variant{Name: "OneLine"})
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(testarea, "chrome", "userChrome__OneLine.css"), file)
}

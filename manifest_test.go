package main

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testarea = path.Join(cwd(), "testarea")

func TestRenderFileTemplate(t *testing.T) {
	Assert(t, RenderFileTemplate("userChrome.css", "linux", Variant{}), "userChrome.css")
	Assert(t, RenderFileTemplate("linux.css", "linux", Variant{}), "linux.css")
	Assert(t, RenderFileTemplate("linux.css", "windows", Variant{}), "linux.css")
	Assert(t, RenderFileTemplate("./{{ os }}/{{variant}}.css", "macos", Variant{Name: "rainbow"}), "./macos/rainbow.css")
}

func TestAssetsPaths(t *testing.T) {
	simplerentfox := Manifest{
		ExplicitName: "simplerentfox",
		Variants: map[string]Theme{
			"WithoutURLBar": {},
			"OneLine":       {},
		},
		Assets: []string{"./{{ os }}/userChrome__{{ variant }}.css"}, // for the purposes of testing
	}
	println("setup done")

	files, err := simplerentfox.AssetsPaths("linux", Variant{Name: "blue"}, testarea)
	assert.Regexp(t, "file .+ not found", err.Error())
	assert.Equal(t, []string{}, files)

	files, err = simplerentfox.AssetsPaths("linux", Variant{Name: "OneLine"}, testarea)
	assert.NoError(t, err)
	assert.Equal(t, []string{CacheDir("simplerentfox/linux/userChrome__OneLine.css")}, files)
}

func TestDestinationPathOf(t *testing.T) {
	manifest := Manifest{
		ExplicitName: "materialfox",
		Variants:     map[string]Theme{},
		Config:       Config{},
	}

	file, err := manifest.DestinationPathOfAsset("/home/ewen/lol.pdf", testarea, "linux", Variant{})
	if assert.Error(t, err) {
		assert.Regexp(t, `asset ".+" is outside of the theme's root ".+"`, err.Error())
	}
	assert.Equal(t, "", file)

	file, err = manifest.DestinationPathOfAsset("/home/.cache/ffcss/simplerentfox/../../../lol.pdf", testarea, "linux", Variant{})
	if assert.Error(t, err) {
		assert.Regexp(t, `asset ".+" is outside of the theme's root ".+"`, err.Error())
	}
	assert.Equal(t, "", file)

	manifest = Manifest{
		ExplicitName: "simplerentfox",
		Variants: map[string]Theme{
			"WithoutURLBar": {},
			"OneLine":       {},
		},
		CopyFrom: "{{ os }}/",
	}

	file, err = manifest.DestinationPathOfAsset(CacheDir("simplerentfox/linux/userChrome__OneLine.css"), testarea, "linux", Variant{Name: "OneLine"})
	assert.NoError(t, err)
	assert.Equal(t, path.Join(testarea, "chrome", "userChrome__OneLine.css"), file)
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderFileTemplate(t *testing.T) {
	Assert(t, RenderFileTemplate("userChrome.css", "linux", Variant{}), "userChrome.css")
	Assert(t, RenderFileTemplate("linux.css", "linux", Variant{}), "linux.css")
	Assert(t, RenderFileTemplate("linux.css", "windows", Variant{}), "linux.css")
	Assert(t, RenderFileTemplate("./{{ os }}/{{variant}}.css", "macos", Variant{Name: "rainbow"}), "./macos/rainbow.css")
}

func TestAllFileTemplates(t *testing.T) {
	assert.Equal(t, []FileTemplate{
		"userChrome.css", "userContent.css", "user.js", "an-asset/**", "another-asset", "blue/userChrome.css",
	}, Manifest{
		UserChrome:  "userChrome.css",
		UserContent: "userContent.css",
		UserJS:      "user.js",
		Assets:      []FileTemplate{"an-asset/**", "another-asset", "blue/userChrome.css"},
	}.AllFileTemplates())
}

// func TestAllFiles(t *testing.T) {
// 	theme := Manifest{
// 		Name: "simplerentfox",
// 	}
// 	assert.Equal(t,
// 		[]string{"userChrome.css", "userContent.css", "user.js", "an-asset/"}
// 	)
// }

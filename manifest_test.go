package main

import "testing"

func TestRender(t *testing.T) {
	Assert(t, RenderFileTemplate("userChrome.css", "linux", Variant{}), "userChrome.css")
	Assert(t, RenderFileTemplate("linux.css", "linux", Variant{}), "linux.css")
	Assert(t, RenderFileTemplate("linux.css", "windows", Variant{}), "")
	Assert(t, RenderFileTemplate("./{{ os }}/{{variant}}.css", "macos", Variant{Name: "rainbow"}), "./macos/rainbow.css")
}

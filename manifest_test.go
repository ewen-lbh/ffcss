package main

import "testing"

func TestRender(t *testing.T) {
	Assert(t, File{Name: "userChrome.css"}.Render("linux", Variant{}), "userChrome.css")
	Assert(t, File{Name: "linux.css", OS: "linux"}.Render("linux", Variant{}), "linux.css")
	Assert(t, File{Name: "linux.css", OS: "linux"}.Render("windows", Variant{}), "")
	Assert(t, File{Name: "./{{ os }}/{{variant}}.css"}.Render("macos", Variant{Name: "rainbow"}), "./macos/rainbow.css")
}

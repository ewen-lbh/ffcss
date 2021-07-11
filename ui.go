package main

import (
	"fmt"
	"os"
	chromaQuick "github.com/alecthomas/chroma/quick"
	"github.com/mitchellh/colorstring"
)

var colorizer colorstring.Colorize

func init() {
	colorizer.Colors = colorstring.DefaultColors
	colorizer.Colors["italic"] = "3"
	colorizer.Reset = true
}

func showSource(theme Manifest) {
	fmt.Print("\n")
	fmt.Println(colorizer.Color("[italic][dim]"+theme.Name()+"'s manifest"))
	chromaQuick.Highlight(os.Stdout, theme.Raw, "YAML", "terminal16m", "pygments")
	fmt.Print("\n")
}


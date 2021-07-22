package main

import (
	"fmt"
	"os"
	"strings"

	chromaQuick "github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/glamour"
	"github.com/mitchellh/colorstring"
)

const indent = "  "

var BulletColorsByIndentLevel = []string{
	"blue",
	"magenta",
	"cyan",
	"green",
	"yellow",
}

var colorizer colorstring.Colorize

func init() {
	colorizer.Colors = colorstring.DefaultColors
	colorizer.Colors["italic"] = "3"
	colorizer.Reset = true
}

// Show the introduction message before installation
func intro(theme Manifest) {
	fmt.Print("\n")
	indent := ""

	var author string
	urlParts := strings.Split(theme.DownloadAt, "/")
	d("urlParts is %#v", urlParts)
	if theme.Author != "" {
		author = theme.Author
	} else if strings.Contains(theme.DownloadAt, "github.com") && len(urlParts) == 5 {
		author = urlParts[len(urlParts)-2]
	}

	fmt.Print(indent)

	fmt.Printf(
		colorizer.Color("[dim]Installing ") +
			colorizer.Color("[blue][bold]"+theme.Name()),
	)

	if author != "" {
		fmt.Printf(
			colorizer.Color("[dim][italic] by ") +
				colorizer.Color("[blue][italic]"+author),
		)
	}

	if theme.Description != "" {
		fmt.Print("\n")
		gutter := colorstring.Color(indent + "[blue]│")
		// gutter := colorstring.Color(indent + "[blue]|")
		d("gutter is %q", gutter)
		markdownRendered, err := glamour.Render(theme.Description, "dark")
		if err != nil {
			markdownRendered = theme.Description
		}
		fmt.Print("\n")
		d("splitted is %#v", strings.Split(markdownRendered, "\n"))
		for _, line := range strings.Split(markdownRendered, "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			fmt.Println(gutter + strings.TrimSpace(line))
		}
		fmt.Print("\n")
	} else {
		fmt.Print("\n\n")
	}

}

func showSource(theme Manifest) {
	fmt.Print("\n")
	fmt.Println(colorizer.Color("[italic][dim]" + theme.Name() + "'s manifest"))
	chromaQuick.Highlight(os.Stdout, theme.Raw, "YAML", "terminal16m", "pygments")
	fmt.Print("\n")
}

func plural(singular string, amount int, optionalPlural ...string) string {
	var plural string
	switch len(optionalPlural) {
	case 1:
		plural = optionalPlural[0]
	case 0:
		plural = singular + "s"
	default:
		panic("plural expected 2 or 3 arguments, you gave more")
	}
	if amount == 1 {
		return singular
	}
	return plural
}

// d prints a debug log line
func d(s string, fmtArgs ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf(colorizer.Color("[dim][ DEBUG ] "+s+"\n"), fmtArgs...)
	}
}

// warn prints a log line with "warning" styling
func warn(s string, fmtArgs ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf(colorizer.Color("[yellow][bold][WARNING] "+s+"\n"), fmtArgs...)
	} else {
		fmt.Printf(colorizer.Color("[yellow][bold]"+s+"\n"), fmtArgs...)
	}
}

// display a list item
func li(indentLevel uint, item string, fmtArgs ...interface{}) {
	var color string
	if int(indentLevel) > len(BulletColorsByIndentLevel)-1 {
		color = BulletColorsByIndentLevel[len(BulletColorsByIndentLevel)-1]
	} else {
		color = BulletColorsByIndentLevel[indentLevel]
	}

	bullet := strings.Repeat(indent, int(indentLevel)) +
		colorizer.Color("["+color+"]•")

	fmt.Println(bullet + " " + colorizer.Color(fmt.Sprintf(item, fmtArgs...)))
}

func (ffp FirefoxProfile) String() string {
	return colorizer.Color(fmt.Sprintf("[bold]%s [reset][dim](%s)", ffp.Name, ffp.ID))
}

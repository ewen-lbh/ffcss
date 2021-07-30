package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	chromaQuick "github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/glamour"
	"github.com/docopt/docopt-go"
	"github.com/mitchellh/colorstring"
)

const indent = "  "

var BulletColorsByIndentLevel = []string{
	"blue",
	"cyan",
	"green",
	"yellow",
	"orange",
	"red",
	"magenta",
	"white",
}

var colorizer colorstring.Colorize

func init() {
	colorizer.Colors = colorstring.DefaultColors
	colorizer.Colors["italic"] = "3"
	colorizer.Colors["orange"] = "38;2;241;109;12"
	colorizer.Reset = true
}

// Show the introduction message before installation
func intro(theme Theme, indentLevel uint) {
	fmt.Print("\n")
	indentation := strings.Repeat(indent, int(indentLevel))

	var author string
	urlParts := strings.Split(theme.DownloadAt, "/")
	d("urlParts is %#v", urlParts)
	if theme.Author != "" {
		author = theme.Author
	} else if strings.Contains(theme.DownloadAt, "github.com") && len(urlParts) == 5 {
		author = urlParts[len(urlParts)-2]
	}

	fmt.Print(indentation)

	fmt.Printf(
		colorizer.Color("[dim]Installing ") +
			colorizer.Color("[blue][bold]"+theme.Name()),
	)

	if regexp.MustCompile(`^v([0-9\.]+)$`).MatchString(theme.Tag) {
		fmt.Printf(colorstring.Color(" [blue]" + theme.Tag))
	}

	if author != "" {
		fmt.Printf(
			colorizer.Color("[dim][italic] by ") +
				colorizer.Color("[blue][italic]"+author),
		)
	}

	if theme.Description != "" {
		fmt.Print("\n")
		gutter := colorstring.Color(indentation + "[blue]│")
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

func showManifestSource(theme Theme) {
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

// showError is like warn but with "error" styling
func showError(s string, fmtArgs ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf(colorizer.Color("[red][bold][ ERROR ] "+s+"\n"), fmtArgs...)
	} else {
		fmt.Printf(colorizer.Color("[red][bold]"+s+"\n"), fmtArgs...)
	}
}

// display a list item
func li(indentLevel uint, item string, fmtArgs ...interface{}) {
	lic("•", indentLevel, item, fmtArgs...)
}

func lic(bulletChar string, indentLevel uint, item string, fmtArgs ...interface{}) {
	var color string
	if int(indentLevel) > len(BulletColorsByIndentLevel)-1 {
		color = BulletColorsByIndentLevel[len(BulletColorsByIndentLevel)-1]
	} else {
		color = BulletColorsByIndentLevel[indentLevel]
	}

	bullet := strings.Repeat(indent, int(indentLevel)) +
		colorizer.Color("["+color+"]"+bulletChar)

	if os.Getenv("DEBUG") != "" {
		bullet = "[  LOG  ]"
	}

	fmt.Println(bullet + " " + colorizer.Color(strings.TrimSpace(fmt.Sprintf(item, fmtArgs...))))
}

func (ffp FirefoxProfile) Display() string {
	return colorizer.Color(fmt.Sprintf("[bold]%s [reset][dim](%s)", ffp.Name, ffp.ID))
}

func AskProfiles(profiles []FirefoxProfile, baseIndentation ...uint) []FirefoxProfile {
	var baseIndent uint
	if len(baseIndentation) == 0 {
		baseIndent = 0
	}
	baseIndent = baseIndentation[0]

	var selectedProfiles []FirefoxProfile

	// XXX the whole display thing should be put in survey.MultiSelect.Renderer, look into that.
	selectedProfileDirsDisplay := make([]string, 0)

	li(baseIndent+0, "Please select profiles to apply the theme on")

	profileDirsDisplay := make([]string, 0)
	for _, profile := range profiles {
		profileDirsDisplay = append(profileDirsDisplay, profile.Display())
	}

	survey.AskOne(&survey.MultiSelect{
		Message: "Select profiles",
		Options: profileDirsDisplay,
		VimMode: VimModeEnabled(),
	}, &selectedProfileDirsDisplay)

	for _, chosenProfileDisplay := range selectedProfileDirsDisplay {
		selectedProfiles = append(selectedProfiles, NewFirefoxProfileFromDisplay(chosenProfileDisplay, profiles))
	}

	return selectedProfiles
}

func (t Theme) AskToSeeManifestSource(skip bool) {
	wantsSource := false
	if !skip {
		survey.AskOne(&survey.Confirm{
			Message: "Show the manifest source?",
		}, &wantsSource)
	}
	if wantsSource {
		showManifestSource(t)
	}
}

func SelectProfiles(baseIndent uint, args docopt.Opts) ([]FirefoxProfile, error) {
	selectedProfilesString, _ := args.String("--profiles")
	var selectedProfiles []FirefoxProfile
	if selectedProfilesString != "" {
		for _, profilePath := range strings.Split(selectedProfilesString, ",") {
			selectedProfiles = append(selectedProfiles, NewFirefoxProfileFromPath(profilePath))
		}
	} else {
		li(baseIndent+0, "Getting profiles")
		profilesDir, _ := args.String("--profiles-dir")
		profiles, err := Profiles(profilesDir)
		if err != nil {
			return []FirefoxProfile{}, fmt.Errorf("couldn't get profile directories: %w", err)
		}
		// Choose profiles
		// TODO smart default (based on {{profileDirectory}}/times.json:firstUse)
		selectAllProfilePaths, _ := args.Bool("--all-profiles")
		if selectAllProfilePaths {
			li(baseIndent+0, "Selecting all profiles")
			selectedProfiles = profiles
		} else {
			selectedProfiles = AskProfiles(profiles, baseIndent)
		}
	}
	return selectedProfiles, nil
}

func (t Theme) ChooseVariant(baseIndent uint, args docopt.Opts) (chosen Variant, cancel bool) {
	variantName, _ := args.String("VARIANT")
	if len(t.AvailableVariants()) > 0 && variantName == "" {
		li(baseIndent+0, "Please choose the theme's variant")
		variantPrompt := &survey.Select{
			Message: "Install variant",
			Options: t.AvailableVariants(),
			VimMode: VimModeEnabled(),
		}
		survey.AskOne(variantPrompt, &variantName)
		// user Ctrl-C'd
		if variantName == "" {
			return Variant{}, true
		}
	}
	return t.Variants[variantName], false
}

func ConfirmInstallAddons(addons []string) bool {
	acceptOpenExtensionPages := false
	survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("This theme suggests installing %d %s. Open %s?",
			len(addons),
			plural("addon", len(addons)),
			plural("its page", len(addons), "their pages"),
		),
		Default: acceptOpenExtensionPages,
	}, &acceptOpenExtensionPages)
	return acceptOpenExtensionPages
}

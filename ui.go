package ffcss

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

func printf(s string, args ...interface{}) {
	fmt.Fprintf(out, s, args...)
}

func printfln(s string, args ...interface{}) {
	printf(s+"\n", args...)
}

// DescribeTheme shows the introduction message before installation
func DescribeTheme(theme Theme, indentLevel uint) {
	printf("\n")
	indentation := strings.Repeat(indent, int(indentLevel))

	var author string
	urlParts := strings.Split(theme.DownloadAt, "/")
	D("urlParts is %#v", urlParts)
	if theme.Author != "" {
		author = theme.Author
	} else if strings.Contains(theme.DownloadAt, "github.com") && len(urlParts) == 5 {
		author = urlParts[len(urlParts)-2]
	}

	printf(indentation)

	printf(
		colorizer.Color("[dim]Installing ") +
			colorizer.Color("[blue][bold]"+theme.Name()),
	)

	if regexp.MustCompile(`^v([0-9\.]+)$`).MatchString(theme.Tag) {
		printf(colorstring.Color(" [blue]" + theme.Tag))
	}

	if author != "" {
		printf(
			colorizer.Color("[dim][italic] by ") +
				colorizer.Color("[blue][italic]"+author),
		)
	}

	if theme.Description != "" {
		printf("\n")
		gutter := colorstring.Color(indentation + "[blue]│")
		D("gutter is %q", gutter)
		markdownRendered, err := glamour.Render(theme.Description, "dark")
		if err != nil {
			markdownRendered = theme.Description
		}
		printf("\n")
		D("splitted is %#v", strings.Split(markdownRendered, "\n"))
		for _, line := range strings.Split(markdownRendered, "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			printfln(gutter + strings.TrimSpace(line))
		}
		printf("\n")
	} else {
		printf("\n\n")
	}

}

func showManifestSource(theme Theme) {
	printf("\n")
	printfln(colorizer.Color("[italic][dim]" + theme.Name() + "'s manifest"))
	chromaQuick.Highlight(os.Stdout, theme.Raw, "YAML", "terminal16m", "pygments")
	printf("\n")
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

// D prints a debug log line. This one always prints to the real stdout, ignoring a possibly mocked stdout
func D(s string, fmtArgs ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf(colorizer.Color("[dim][ DEBUG ] "+s+"\n"), fmtArgs...)
	}
}

// Warn prints a log line with "warning" styling
func Warn(s string, fmtArgs ...interface{}) {
	printf(colorizer.Color("[yellow][bold]"+s+"\n"), fmtArgs...)
}

// Error is like warn but with "error" styling
func Error(s string, fmtArgs ...interface{}) {
	printf(colorizer.Color("[red][bold]"+s+"\n"), fmtArgs...)
}

// Step displays a list item
func Step(indentLevel uint, item string, fmtArgs ...interface{}) {
	StepC("•", indentLevel, item, fmtArgs...)
}

// StepC is like Step, but the bullet point characters is customizable
func StepC(bulletChar string, indentLevel uint, item string, fmtArgs ...interface{}) {
	indentLevel += BaseIndentLevel
	var color string
	if int(indentLevel) > len(BulletColorsByIndentLevel)-1 {
		color = BulletColorsByIndentLevel[len(BulletColorsByIndentLevel)-1]
	} else {
		color = BulletColorsByIndentLevel[indentLevel]
	}

	bullet := strings.Repeat(indent, int(indentLevel)) +
		colorizer.Color("["+color+"]"+bulletChar)

	printfln(bullet + " " + colorizer.Color(strings.TrimSpace(fmt.Sprintf(item, fmtArgs...))))
}

func (ffp FirefoxProfile) Display() string {
	return colorizer.Color(fmt.Sprintf("[bold]%s [reset][dim](%s)", ffp.Name, ffp.ID))
}

func AskProfiles(profiles []FirefoxProfile) []FirefoxProfile {
	var selectedProfiles []FirefoxProfile

	// XXX the whole display thing should be put in survey.MultiSelect.Renderer, look into that.
	selectedProfileDirsDisplay := make([]string, 0)

	Step(0, "Please select profiles to apply the theme on")

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

func SelectProfiles(args docopt.Opts) ([]FirefoxProfile, error) {
	selectedProfilesString, _ := args.String("--profiles")
	var selectedProfiles []FirefoxProfile
	if selectedProfilesString != "" {
		for _, profilePath := range strings.Split(selectedProfilesString, ",") {
			selectedProfiles = append(selectedProfiles, NewFirefoxProfileFromPath(profilePath))
		}
	} else {
		Step(0, "Getting profiles")
		profilesDir, _ := args.String("--profiles-dir")
		profiles, err := Profiles(profilesDir)
		if err != nil {
			return []FirefoxProfile{}, fmt.Errorf("couldn't get profile directories: %w", err)
		}
		// Choose profiles
		// TODO smart default (based on {{profileDirectory}}/times.json:firstUse)
		selectAllProfilePaths, _ := args.Bool("--all-profiles")
		if selectAllProfilePaths {
			Step(0, "Selecting all profiles")
			selectedProfiles = profiles
		} else {
			selectedProfiles = AskProfiles(profiles)
		}
	}
	return selectedProfiles, nil
}

func (t Theme) ChooseVariant(args docopt.Opts) (chosen Variant, cancel bool) {
	variantName, _ := args.String("VARIANT")
	if len(t.AvailableVariants()) > 0 && variantName == "" {
		Step(0, "Please choose the theme's variant")
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

func ShowHookOutput(output string) {
	fmt.Fprint(
		out,
		"\n",
		prefixEachLine(
			strings.TrimSpace(output),
			strings.Repeat(indent, int(BaseIndentLevel)+2),
		),
		"\n",
		"\n",
	)
}

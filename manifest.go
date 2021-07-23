package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/hoisie/mustache"
	"gopkg.in/yaml.v2"
)

var ThemeCompatWarningShown = false

type Variant struct {
	// Properties exclusive to variants
	Name    string
	Message string

	// Properties that modify the "default variant"
	Repository  string
	Branch      string
	Config      Config
	UserChrome  FileTemplate `yaml:"userChrome"`
	UserContent FileTemplate `yaml:"userContent"`
	UserJS      FileTemplate `yaml:"user.js"`
	Assets      []FileTemplate
	Description string
	Addons      []string
	Run         struct {
		Before string
		After  string
	}
}

type Manifest struct {
	// Internal, cannot be set in the YAML file
	CurrentVariantName string `yaml:"-"` // Used to construct the directory where the theme will be cached
	Raw                string `yaml:"-"` // Contains the raw yaml file contents
	DownloadedTo       string `yaml:"-"` // Stores the path to the directory where the theme is cached. Set by .Download().

	// Top-level (non-variant-modifiable)
	ExplicitName   string `yaml:"name"`
	Description    string
	Author         string `yaml:"by"`
	FfcssVersion   int    `yaml:"ffcss"`
	FirefoxVersion string `yaml:"firefox"`
	Variants       map[string]Variant
	OSNames        map[string]string `yaml:"os"`

	// Override-able by variants
	DownloadAt  string `yaml:"download"`
	Branch      string
	CopyFrom    string `yaml:"copy from"`
	Config      Config
	UserChrome  FileTemplate `yaml:"userChrome"`
	UserContent FileTemplate `yaml:"userContent"`
	UserJS      FileTemplate `yaml:"user.js"`
	Assets      []FileTemplate
	Run         struct {
		Before string
		After  string
	}
	Message string
	Addons  []string
}

func (m Manifest) Name() string {
	if m.ExplicitName != "" {
		return strings.ToLower(m.ExplicitName)
	}
	if strings.HasPrefix(m.DownloadAt, "https://github.com") {
		fragments := strings.Split(m.DownloadAt, "/")
		return strings.ToLower(fragments[len(fragments)-1])
	}
	return ""
}

type Config map[string]interface{}

type FileTemplate = string

func NewManifest() Manifest {
	return Manifest{
		Config: Config{
			"toolkit.legacyUserProfileCustomizations.stylesheets": true,
		},
		UserChrome:  "",
		UserContent: "",
		UserJS:      "",
		Variants:    map[string]Variant{},
		Assets:      []FileTemplate{},
	}
}

// LoadManifest loads a ffcss.yaml file into a Manifest object.
func LoadManifest(manifestPath string) (manifest Manifest, err error) {
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		err = fmt.Errorf("while reading manifest %s: %w", manifestPath, err)
		return
	}
	manifest = NewManifest()
	manifest.Raw = string(raw)
	err = yaml.Unmarshal(raw, &manifest)

	if manifest.FfcssVersion != VersionMajor && !ThemeCompatWarningShown && manifest.FfcssVersion != 0 {
		warn("ffcss %s is installed, but you are using a theme made for ffcss %d.X.X. Some things may not work.\n", VersionString, manifest.FfcssVersion)
		ThemeCompatWarningShown = true
	}

	if manifest.Name() == TempDownloadsDirName {
		err = fmt.Errorf("invalid theme name %q", TempDownloadsDirName)
		return
	}

	for name, variant := range manifest.Variants {
		if name == RootVariantName {
			err = fmt.Errorf("invalid variant name %q", name)
			return
		}
		variantWithName := variant
		variantWithName.Name = name
		manifest.Variants[name] = variantWithName
	}
	manifest.CurrentVariantName = RootVariantName // ensure the current variant's name wasn't manipulated by the YAML unmarshaling
	if err != nil {
		err = fmt.Errorf("while parsing manifest %s: %w", manifestPath, err)
		return
	}
	manifest.DownloadedTo = CacheDir(manifest.Name(), manifest.CurrentVariantName)
	return
}

// WithVariant returns a Manifest representing the theme if the selected variant
// was used as the "root values".
// i.e. the values of UserJS, UserContent, UserChrome, Assets are replaced with their variant's, if set,
// and the value of Config is combined with the variant's.
// Some variants change the git branch, the entire repository or other settings that require external actions.
// Those are returned in actionsNeeded as a struct of booleans with descriptive field names.
func (m Manifest) WithVariant(variant Variant) (newManifest Manifest, actionsNeeded struct{ switchBranch, reDownload bool }) {
	newManifest = m
	newManifest.CurrentVariantName = variant.Name
	if variant.UserChrome != "" {
		newManifest.UserChrome = variant.UserChrome
	}
	if variant.UserContent != "" {
		newManifest.UserContent = variant.UserContent
	}
	if variant.UserJS != "" {
		newManifest.UserJS = variant.UserJS
	}
	if variant.Message != "" {
		newManifest.Message = variant.Message
	}
	if len(variant.Assets) > 0 {
		newManifest.Assets = variant.Assets
	}
	if variant.Repository != "" {
		actionsNeeded.reDownload = true
		newManifest.DownloadAt = variant.Repository
	}
	if variant.Branch != "" {
		actionsNeeded.switchBranch = true
		newManifest.Branch = variant.Branch
	}
	if variant.Run.Before != "" {
		newManifest.Run.Before = variant.Run.Before
	}
	if variant.Run.After != "" {
		newManifest.Run.After = variant.Run.After
	}
	for key, val := range variant.Config {
		newManifest.Config[key] = val
	}
	if actionsNeeded.reDownload || actionsNeeded.switchBranch {
		newManifest.DownloadedTo = CacheDir(newManifest.Name(), newManifest.CurrentVariantName)
	}
	return newManifest, actionsNeeded
}

// ThemeStore represents a collection of themes
type ThemeStore = map[string]Manifest

// LoadThemeCatalog loads a directory of theme manifests.
// Keys are theme names (files' basenames with the .yaml removed).
func LoadThemeCatalog(storeDirectory string) (themes ThemeStore, err error) {
	themeNamePattern := regexp.MustCompile(`^(.+)\.ya?ml$`)
	themes = make(ThemeStore)
	manifests, err := os.ReadDir(storeDirectory)
	if err != nil {
		return
	}
	for _, manifest := range manifests {
		if !themeNamePattern.MatchString(manifest.Name()) {
			continue
		}
		themeName := themeNamePattern.FindStringSubmatch(manifest.Name())[1]
		theme, err := LoadManifest(filepath.Join(storeDirectory, manifest.Name()))
		if err != nil {
			return nil, err
		}
		themes[themeName] = theme
	}
	return
}

// AvailableVariants lists the possible variant names to choose from
func (m Manifest) AvailableVariants() []string {
	names := make([]string, 0, len(m.Variants))
	for name := range m.Variants {
		names = append(names, name)
	}
	return names
}

// ShowMessage renders the message and prints it to the user
func (m Manifest) ShowMessage() error {
	scheme := os.Getenv("COLORSCHEME")
	if scheme != "light" && scheme != "dark" {
		// TODO: detect with the terminal's current background color as a fallback
		scheme = "dark"
	}
	rendered, err := glamour.Render(m.Message, scheme)
	if err != nil {
		return fmt.Errorf("while rendering message: %w", err)
	}

	if strings.TrimSpace(rendered) != "" {
		fmt.Println(rendered)
	}
	return nil
}

// runHook runs a provided command for a specific profile. See any of the (Manifest).Run*Hook methods
// for a list of available {{mustache}} placeholders.
func (m Manifest) runHook(commandline string, profile FirefoxProfile) (output string, err error) {
	ffversion, err := profile.FirefoxVersion()
	if err != nil {
		return "", fmt.Errorf("while getting firefox version for current profile: %w", err)
	}

	command := exec.Command("bash", "-c", mustache.Render(commandline, map[string]interface{}{
		"profile_path":    profile.Path,
		"firefox_version": ffversion.String(),
	}))

	outputBytes, err := command.CombinedOutput()
	output = string(outputBytes)
	if err != nil {
		return "", fmt.Errorf("while running %q: %s: %w", command.String(), output, err)
	}

	return
}

// RunPreInstallHook passes the pre-install hook specified in the manifest's run.before entry to bash.
// Several {{mustache}} placeholders are available:
//
//	profile_path        The current profile's path
//	firefox_version     The current profile's Firefox version
func (m Manifest) RunPreInstallHook(profile FirefoxProfile) (output string, err error) {
	return m.runHook(m.Run.Before, profile)
}

// RunPostInstallHook does the same as RunPreInstallHook but for the manifest's run.after entry.
func (m Manifest) RunPostInstallHook(profile FirefoxProfile) (output string, err error) {
	return m.runHook(m.Run.After, profile)
}

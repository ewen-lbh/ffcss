package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"gopkg.in/yaml.v2"
)

var ThemeCompatWarningShown = VersionMajor > 0

type Config map[string]interface{}

func (c Config) Equal(other Config) bool {
	for k, v := range c {
		if v != other[k] {
			return false
		}
	}
	return true
}

type FileTemplate = string

type Variant struct {
	// Properties exclusive to variants
	Name    string
	Message string

	// Properties that modify the "default variant"
	Repository  string
	Branch      string
	Commit      string
	Tag         string
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

type Theme struct {
	// Internal, cannot be set in the YAML file
	CurrentVariantName string `yaml:"-"` // Used to construct the directory where the theme will be cached
	Raw                string `yaml:"-"` // Contains the raw yaml file contents
	DownloadedTo       string `yaml:"-"` // Stores the path to the directory where the theme is cached. Set by .Download().

	// Top-level (non-variant-modifiable)
	FfcssVersion             int                      `yaml:"ffcss"`
	FirefoxVersion           string                   `yaml:"firefox,omitempty"`
	FirefoxVersionConstraint FirefoxVersionConstraint `yaml:"-"`
	ExplicitName             string                   `yaml:"name"`
	Author                   string                   `yaml:"by"`
	Description              string
	Variants                 map[string]Variant
	OSNames                  map[string]string `yaml:"os,omitempty"`

	// Override-able by variants
	DownloadAt  string `yaml:"download"`
	Branch      string
	Commit      string `yaml:",omitempty"`
	Tag         string `yaml:",omitempty"`
	Config      Config
	UserChrome  FileTemplate `yaml:"userChrome"`
	UserContent FileTemplate `yaml:"userContent"`
	UserJS      FileTemplate `yaml:"user.js"`
	Assets      []FileTemplate
	CopyFrom    string `yaml:"copy from,omitempty"`
	Addons      []string
	Run         struct {
		Before string
		After  string
	}
	Message string
}

// ManifestKeyGroupsStarts specifies at which keys a group of related keys starts.
// This assumes that the YAML keys are displayed/written in the order they are defined in (see Theme).
//
// It is used to add blank lines in generated manifests: for example, above the 'variant' key, a blank line should be added to add grouping.
var ManifestKeyGroupsStarts = [...]string{"name", "variants", "os", "download", "config", "run", "message"}

func NewTheme() Theme {
	return Theme{
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

// LoadManifest loads a ffcss.yaml file into a Theme object.
func LoadManifest(manifestPath string) (manifest Theme, err error) {
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		err = fmt.Errorf("while reading manifest %s: %w", manifestPath, err)
		return
	}
	manifest = NewTheme()
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
	if manifest.FirefoxVersion != "" {
		manifest.FirefoxVersionConstraint, err = NewFirefoxVersionConstraint(manifest.FirefoxVersion)
		if err != nil {
			err = fmt.Errorf("while parsing version constraint %q: %w", manifest.FirefoxVersion, err)
			return
		}

	}
	return
}

// WithVariant returns a Manifest representing the theme if the selected variant
// was used as the "root values".
// i.e. the values of UserJS, UserContent, UserChrome, Assets are replaced with their variant's, if set,
// and the value of Config is combined with the variant's.
// Some variants change the git branch, the entire repository or other settings that require external actions.
// Those are returned in actionsNeeded as a struct of booleans with descriptive field names.
func (t Theme) WithVariant(variant Variant) (newTheme Theme, actionsNeeded struct{ switchBranch, reDownload bool }) {
	// TODO might clean this up with reflection, selecting fields that are both in Manifest & Variant
	newTheme = t
	newTheme.CurrentVariantName = variant.Name
	if variant.UserChrome != "" {
		newTheme.UserChrome = variant.UserChrome
	}
	if variant.UserContent != "" {
		newTheme.UserContent = variant.UserContent
	}
	if variant.UserJS != "" {
		newTheme.UserJS = variant.UserJS
	}
	if variant.Message != "" {
		newTheme.Message = variant.Message
	}
	if len(variant.Assets) > 0 {
		newTheme.Assets = variant.Assets
	}
	if variant.Repository != "" {
		actionsNeeded.reDownload = true
		newTheme.DownloadAt = variant.Repository
	}
	if variant.Branch != "" {
		actionsNeeded.switchBranch = true
		newTheme.Branch = variant.Branch
	}
	if variant.Commit != "" {
		newTheme.Commit = variant.Commit
	}
	if variant.Tag != "" {
		newTheme.Tag = variant.Tag
	}
	if variant.Run.Before != "" {
		newTheme.Run.Before = variant.Run.Before
	}
	if variant.Run.After != "" {
		newTheme.Run.After = variant.Run.After
	}
	for key, val := range variant.Config {
		newTheme.Config[key] = val
	}
	if actionsNeeded.reDownload || actionsNeeded.switchBranch {
		newTheme.DownloadedTo = CacheDir(newTheme.Name(), newTheme.CurrentVariantName)
	}
	return newTheme, actionsNeeded
}

func (t Theme) Name() string {
	if t.ExplicitName != "" {
		return strings.ToLower(t.ExplicitName)
	}
	if strings.HasPrefix(t.DownloadAt, "https://github.com") {
		fragments := strings.Split(t.DownloadAt, "/")
		return strings.ToLower(fragments[len(fragments)-1])
	}
	return ""
}

// AvailableVariants lists the possible variant names to choose from
func (t Theme) AvailableVariants() []string {
	names := make([]string, 0, len(t.Variants))
	for name := range t.Variants {
		names = append(names, name)
	}
	return names
}

// ShowMessage renders the message and prints it to the user
func (t Theme) ShowMessage() error {
	scheme := os.Getenv("COLORSCHEME")
	if scheme != "light" && scheme != "dark" {
		// TODO: detect with the terminal's current background color as a fallback
		scheme = "dark"
	}
	rendered, err := glamour.Render(t.Message, scheme)
	if err != nil {
		return fmt.Errorf("while rendering message: %w", err)
	}

	if strings.TrimSpace(rendered) != "" {
		fmt.Fprintln(out, rendered)
	}
	return nil
}

// GetManifestPath returns the path of a theme's manifest file
func GetManifestPath(themeRoot string) string {
	return filepath.Join(themeRoot, "ffcss.yaml")
}

// GenerateManifest returns the YAML contents of the manifest corresponding to the given theme.
// If t.Raw is set, it'll return it.
// Otherwise, it serializes the values into YAML, following the Theme struct.
func (t Theme) GenerateManifest() (string, error) {
	if t.Raw != "" {
		return t.Raw, nil
	}
	t.ExplicitName = t.Name()
	// Remove redundant keys
	if t.Config.Equal(NewTheme().Config) {
		t.Config = Config{}
	}
	contentBytes, err := yaml.Marshal(t)
	if err != nil {
		return "", err
	}

	content := string(contentBytes)

	for _, key := range ManifestKeyGroupsStarts {
		content = strings.Replace(content, key+": ", "\n"+key+": ", 1)
		content = strings.Replace(content, key+":\n", "\n"+key+":\n", 1)
	}

	return content, nil
}

// WriteManifest writes the contents of t as a YAML file named ffcss.yaml
// inside inDirectory.
// It adds a comment mentionning the documentation at the top of the file.
// See GenerateManifest to see how the contents of the file are generated.
func (t Theme) WriteManifest(inDirectory string) error {
	content, err := t.GenerateManifest()
	if err != nil {
		return fmt.Errorf("while generating manifest contents for %s: %w", t.Name(), err)
	}

	content = "# This is a manifest for a FirefoxCSS theme. \n# See https://github.com/ewen-lbh/ffcss for more information.\n" + content

	err = ioutil.WriteFile(filepath.Join(inDirectory, "ffcss.yaml"), []byte(content), 0700)
	if err != nil {
		return fmt.Errorf("while writing the manifest: %w", err)
	}

	return nil
}

// InitializeTheme returns a new, blank theme, but with some values guessed from the current context.
// Meant to be used by "ffcss init".
func InitializeTheme(workingDir string) (Theme, error) {
	theme := NewTheme()

	theme.DownloadAt = strings.TrimSuffix(getCurrentRepoRemote(), ".git")
	if theme.DownloadAt == "" {
		theme.DownloadAt = "TODO"
	}

	if !strings.Contains(theme.DownloadAt, "https://github.com") {
		theme.ExplicitName = filepath.Base(workingDir)
	}

	return theme, nil
}

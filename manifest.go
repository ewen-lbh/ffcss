package ffcss

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"gopkg.in/yaml.v2"
)

// ThemeCompatWarningShown controls whether the ffcss incompatibility between the installed version and the theme's declared version (ffcss entry in the manifest)
// is warned to the user.
// Set to VersionMajor > 0 to comply with semver: versions with the major component (X._._) set to 0 can have breaking changes at any version change.
var ThemeCompatWarningShown = VersionMajor > 0

// Config represents a configuration map in a manifest, representing a set of values for the about:config page in Firefox
type Config map[string]interface{}

// Equal returns true if the config c has all of its values equal to the other Config.
func (c Config) Equal(other Config) bool {
	// XXX: this is not efficient AT ALL. It sucks. Change this.
	for k, v := range c {
		if v != other[k] {
			return false
		}
	}
	for k, v := range other {
		if other[k] != v {
			return false
		}
	}
	return true
}

// FileTemplate represents a string with placeholders (e.g. "{{os}}") to be replaced by {{mustache}}.
type FileTemplate = string

// Variant represents a theme's variant. Most of the properties are identical to Theme's, because they overwrite the default values.
type Variant struct {
	// Properties exclusive to variants
	Name string

	// Properties that modify the "default variant"
	DownloadAt  string `yaml:"download"`
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
	Message string
}

// Theme represents a FirefoxCSS theme, read from a manifest YAML file. (See LoadManifest).
type Theme struct {
	// Internal, cannot be set in the YAML file
	currentVariantName string `yaml:"-"` // Used to construct the directory where the theme will be cached
	raw                string `yaml:"-"` // Contains the raw yaml file contents
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

// NewTheme creates a new Theme with vital defaults (namely the config entry to enable CSS customization of Firefox).
func NewTheme() Theme {
	return Theme{
		Config: Config{
			"toolkit.legacyUserProfileCustomizations.stylesheets": true,
		},
		Variants: map[string]Variant{},
		Assets:   []FileTemplate{},
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
	manifest.raw = string(raw)
	err = yaml.Unmarshal(raw, &manifest)

	if manifest.FfcssVersion < 0 {
		err = fmt.Errorf("ffcss version cannot be negative but is set to %d", manifest.FfcssVersion)
		return
	}

	if manifest.FfcssVersion != VersionMajor && !ThemeCompatWarningShown && manifest.FfcssVersion != 0 {
		LogWarning("ffcss %s is installed, but you are using a theme made for ffcss %d.X.X. Some things may not work.\n", VersionString, manifest.FfcssVersion)
		ThemeCompatWarningShown = true
	}

	if manifest.Name() == TempDownloadsDirName {
		err = fmt.Errorf("invalid theme name %q", TempDownloadsDirName)
		return
	}

	if manifest.Name() == "" {
		err = fmt.Errorf("theme has no name")
		return
	}

	for key := range manifest.OSNames {
		if key != "linux" && key != "macos" && key != "windows" {
			return Theme{}, fmt.Errorf("%s is not a valid os replacement target. Targets are macos, windows and linux", key)
		}
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
	manifest.currentVariantName = RootVariantName // ensure the current variant's name wasn't manipulated by the YAML unmarshaling
	if err != nil {
		err = fmt.Errorf("while parsing manifest %s: %w", manifestPath, err)
		return
	}
	manifest.DownloadedTo = CacheDir(manifest.Name(), manifest.currentVariantName)
	if manifest.FirefoxVersion != "" {
		manifest.FirefoxVersionConstraint, err = NewFirefoxVersionConstraint(manifest.FirefoxVersion)
		if err != nil {
			err = fmt.Errorf("invalid Firefox version constraint %q: %w", manifest.FirefoxVersion, err)
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
	newTheme.currentVariantName = variant.Name
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
	if variant.DownloadAt != "" {
		actionsNeeded.reDownload = true
		newTheme.DownloadAt = variant.DownloadAt
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
		newTheme.DownloadedTo = CacheDir(newTheme.Name(), newTheme.currentVariantName)
	}
	return newTheme, actionsNeeded
}

// Name returns a theme's name. If the name was explicitly set in the manifest (i.e. if t.ExplicitName is not empty), it is returned.
// Otherwise, the name is guessed.
// Currently, it is only guessed if a github repository is set for t.DownloadAt.
// If guessing is not possible, it returns the empty string.
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

// ManifestPath returns the path of a theme's manifest file
func ManifestPath(themeRoot string) string {
	return filepath.Join(themeRoot, "ffcss.yaml")
}

// GenerateManifest returns the YAML contents of the manifest corresponding to the given theme.
// If t.Raw is set, it'll return it.
// Otherwise, it serializes the values into YAML, following the Theme struct.
func (t Theme) GenerateManifest() (string, error) {
	if t.raw != "" {
		return t.raw, nil
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
// It adds a comment mentioning the documentation at the top of the file.
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

	theme.DownloadAt = strings.TrimSuffix(currentRepoRemote(), ".git")
	if theme.DownloadAt == "" {
		theme.DownloadAt = "TODO"
	}

	if !strings.Contains(theme.DownloadAt, "https://github.com") {
		theme.ExplicitName = filepath.Base(workingDir)
	}

	return theme, nil
}

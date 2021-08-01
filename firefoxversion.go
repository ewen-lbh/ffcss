package ffcss

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// FirefoxVersion represents a firefox version of the form "major.minor".
type FirefoxVersion struct {
	Major int
	Minor int // -1 means "unspecified". Can be obtained by using "x" for the minor part. Useful for constraints.
}

// FirefoxVersion represents a constraint to test on a firefox version.
type FirefoxVersionConstraint struct {
	Min FirefoxVersion
	Max FirefoxVersion
	// To be included in a Sentence like "this theme ensures compatibility with firefox <Sentence>"
	Sentence string
}

// NewFirefoxVersionConstraint creates a firefox version constraint from its string representation.
//
// Leaving out the minor version (e.g. 90 instead of 90.something) implies a trailing ".x", which means "any minor version".
//
// The following formats are supported, where X, Z are integers and Y, W are integers or the character 'x'
//
//    Format     Meaning                                    Interval
//    X.Y+       X.Y or higher                              [X.Y,  +∞]   (where +∞ = math.MaxInt32)
//    up to X.Y  X.Y but not higher                         [0.0, X.Y]
//    X.Y-Z.W    between X.Y and Z.W (including both ends)  [X.Y, Z.W]
//    X.Y        exactly X.Y                                [X.Y, X.Y]
func NewFirefoxVersionConstraint(constraint string) (FirefoxVersionConstraint, error) {
	min := FirefoxVersion{0, 0}
	max := FirefoxVersion{math.MaxInt32, math.MaxInt32}
	var sentence string
	var err error
	if strings.HasSuffix(constraint, "+") {
		D("constraint type is minimum")
		min, err = NewFirefoxVersion(strings.TrimSuffix(constraint, "+"))
		if err != nil {
			return FirefoxVersionConstraint{}, fmt.Errorf("while parsing minimum constraint %q: %w", constraint, err)
		}
		sentence = fmt.Sprintf("version %s or higher", min)
	} else if strings.Count(constraint, "-") == 1 {
		D("constraint type is range")
		minmaxStrings := strings.SplitN(constraint, "-", 2)
		min, err = NewFirefoxVersion(minmaxStrings[0])
		if err != nil {
			return FirefoxVersionConstraint{}, fmt.Errorf("while parsing lower bound of range constraint %q: %w", constraint, err)
		}
		max, err = NewFirefoxVersion(minmaxStrings[1])
		if err != nil {
			return FirefoxVersionConstraint{}, fmt.Errorf("while parsing upper bound of range constraint %q: %w", constraint, err)
		}
		sentence = min.String() + "–" + max.String()
	} else if strings.HasPrefix(constraint, "up to ") {
		D("constraint type is upto")
		max, err = NewFirefoxVersion(strings.TrimPrefix(constraint, "up to "))
		if err != nil {
			return FirefoxVersionConstraint{}, fmt.Errorf("while parsing maximum constraint %q: %w", constraint, err)
		}
		sentence = fmt.Sprintf("version %s or lower", max)
	} else {
		D("constraint type is exact match")
		exact, err := NewFirefoxVersion(constraint)
		if err != nil {
			return FirefoxVersionConstraint{}, fmt.Errorf("while parsing exact match constraint %q: %w", constraint, err)
		}
		min = exact
		max = exact
		sentence = fmt.Sprintf("%s only", exact)
	}
	return FirefoxVersionConstraint{min, max, sentence}, nil
}

// GreaterOrEqual checks if the version is greater or equal to other.
// If one of the two (or both) has the minor part unspecified (".x", stored as -1),
// it only compares major parts. Otherwise, it uses a standard lexical sort.
func (ffv FirefoxVersion) GreaterOrEqual(other FirefoxVersion) bool {
	if other.Minor == -1 || ffv.Minor == -1 {
		return ffv.Major >= other.Major
	}
	return ffv.Major > other.Major || (ffv.Major == other.Major && ffv.Minor >= other.Minor)
}

// LesserrOrEqual checks if the version is less than or equal to other.
// If one of the two (or both) has the minor part unspecified (".x", stored as -1),
// it only compares major parts. Otherwise, it uses a standard lexical sort.
func (ffv FirefoxVersion) LesserOrEqual(other FirefoxVersion) bool {
	if other.Minor == -1 || ffv.Minor == -1 {
		return ffv.Major <= other.Major
	}
	D("check that %d <= %d or %d = %d and %d <= %d", ffv.Major, other.Major, ffv.Major, other.Major, ffv.Minor, other.Minor)
	return ffv.Major < other.Major || (ffv.Major == other.Major && ffv.Minor <= other.Minor)
}

// String returns a string representation of the version
// If the minor part is -1, it is rendered as a 'x' character.
func (ffv FirefoxVersion) String() string {
	if ffv.Minor == -1 {
		return fmt.Sprintf("%d.x", ffv.Major)
	}
	return fmt.Sprintf("%d.%d", ffv.Major, ffv.Minor)
}

// FulfilledBy checks if version ∈ [constraint.min, constraint.max]
func (constraint FirefoxVersionConstraint) FulfilledBy(version FirefoxVersion) bool {
	D("checking if %s ∈ [%s, %s]", version, constraint.Min, constraint.Max)
	return version.GreaterOrEqual(constraint.Min) && version.LesserOrEqual(constraint.Max)
}

// NewFirefoxVersion turns a version string (90 or 90.0 for example) into a FirefoxVersion.
// defaultMinor is used when parsing a dot-less version string. It defaults to "x" (meaning unspecified).
func NewFirefoxVersion(stringRepr string, defaultMinor ...string) (FirefoxVersion, error) {
	fragments := strings.Split(stringRepr, ".")
	if len(defaultMinor) == 0 {
		defaultMinor = []string{"x"}
	}
	if len(fragments) == 1 {
		fragments = append(fragments, defaultMinor...)
	}
	D("parsing version %s: fragments is %#v", stringRepr, fragments)
	major, err := strconv.ParseInt(fragments[0], 10, 64)
	if err != nil {
		return FirefoxVersion{}, fmt.Errorf("while converting major segment: %w", err)
	}
	var minor int
	if fragments[1] == "x" {
		minor = -1
	} else {
		minor, err := strconv.ParseInt(fragments[1], 10, 64)
		if err != nil {
			return FirefoxVersion{}, fmt.Errorf("while converting minor segment: %w", err)
		}
		if major < 0 || minor < 0 {
			return FirefoxVersion{}, fmt.Errorf("version number cannot be negative")
		}
	}
	return FirefoxVersion{
		Major: int(major),
		Minor: int(minor),
	}, nil
}

// FirefoxVersion returns the firefox version of the profile by reading the value of
// the browser.startup.homepage_override.mstone configuration entry in the profile's prefs.js.
// This method does not work for profiles that have never been opened by the user.
func (profile FirefoxProfile) FirefoxVersion() (FirefoxVersion, error) {
	prefs, err := os.ReadFile(filepath.Join(profile.Path, "prefs.js"))
	if err != nil {
		return FirefoxVersion{}, fmt.Errorf("while reading file: %w", err)
	}

	versionString, err := ValueOfUserPrefCall(prefs, "browser.startup.homepage_override.mstone")
	if err != nil {
		return FirefoxVersion{}, fmt.Errorf("while getting value in prefs.js: %w", err)
	}

	version, err := NewFirefoxVersion(versionString)
	if err != nil {
		return FirefoxVersion{}, fmt.Errorf("while parsing version string %q: %w", versionString, err)
	}

	return version, nil
}

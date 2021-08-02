package ffcss

import (
	"fmt"
	"io"
	"os"
)

const (
	// VersionMajor is the major semver component of the current version, X._._.
	// It is incremented on breaking changes. Incrementing sets others to zero.
	VersionMajor = 0
	// VersionMinor is the semver minor component of the current version, _.Y._.
	// It is incremented on feature changes. Incrementing sets VersionPatch to zero.
	VersionMinor = 2
	// VersionPatch is the semver patch component of the current version, _._.Z.
	// It is incremented on bug fixes
	VersionPatch = 0
)

var (
	// VersionString is the full version string, constructed from the component constants Version{Major,Minor,Patch}.
	VersionString = fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	// out is an IO writer used to capture stdout instead of printing in unit tests
	out io.Writer = os.Stdout
	// BaseIndentLevel is the base level of indentation for LogStep and others.
	// Can be incremented when a whole section of code is run in the context of another
	// (For an example usage, see reapply in the cmd/ffcss package, that effectively runs the use command on each profile).
	BaseIndentLevel uint
)

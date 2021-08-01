package ffcss

import (
	"fmt"
	"io"
	"os"
)

const (
	VersionMajor = 0
	VersionMinor = 2
	VersionPatch = 0
)

var (
	VersionString = fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	// Used to capture stdout instead of printing in tests
	out             io.Writer = os.Stdout
	BaseIndentLevel uint      = 0
)

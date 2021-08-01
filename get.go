package main

import (
	"fmt"

	"github.com/docopt/docopt-go"
)

func RunCommandGet(args docopt.Opts) error {
	themeName, _ := args.String("THEME_NAME")
	// variant, _ := args.String("VARIANT")

	err := CreateDataDirectories()
	if err != nil {
		return err
	}

	li(0, "Resolving the theme's name")
	uri, typ, err := ResolveURL(themeName)
	if err != nil {
		return fmt.Errorf("while resolving name %s: %w", themeName, err)
	}

	li(0, "Downloading the theme")
	manifest, err := Download(uri, typ)
	if err != nil {
		return err
	}

	lic("✓", 0, "Downloaded [blue][bold]%s[reset] [dim](to %s)", manifest.Name(), manifest.DownloadedTo)
	return nil
}

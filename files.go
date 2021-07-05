package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"runtime"

	"github.com/bmatcuk/doublestar"
	"github.com/hoisie/mustache"
)

type UserChoices struct {
	Variant Variant
	OS      string
}

func NewUserChoices() UserChoices {
	return UserChoices{
		Variant: Variant{},
		OS:      GOOStoOS(runtime.GOOS),
	}
}

func GOOStoOS(GOOS string) string {
	switch GOOS {
	case "darwin":
		return "macos"
	case "plan9":
		return "linux"
	default:
		return GOOS
	}
}

// ResolveFilenames resolves the file names using choices made by the user (variant selected, current OS).
// It does not resolve glob patterns though.
func ResolveFilenames(files []FileTemplate, choices UserChoices) (resolved []string, err error) {
	for _, file := range files {
		var output string
		templ, err := mustache.ParseString(file)
		if err != nil {
			return resolved, fmt.Errorf("could not parse %q: %w", file, err)
		}
		output = templ.Render(map[string]string{
			"os":      choices.OS,
			"variant": choices.Variant.Name,
		})
		if err != nil {
			return resolved, fmt.Errorf("could not render %q: %w", file, err)
		}
		resolved = append(resolved, output)
	}
	return
}

// CopyOver copies over files from files to each directory in toDirs
func CopyOver(files []string, toDirs []string) (err error) {
	for _, glob := range files {
		matches, err := doublestar.Glob(glob)
		if err != nil {
			return fmt.Errorf("while scanning for %s: %w", glob, err)
		}
		for _, file := range matches {
			content, err := ioutil.ReadFile(file)
			if err != nil {
				return fmt.Errorf("while reading %s: %w", file, err)
			}
			for _, toDir := range toDirs {
				err = ioutil.WriteFile(path.Join(toDir, file), content, 0700)
				if err != nil {
					return fmt.Errorf("while writing to %s: %w", path.Join(toDir, file), err)
				}
			}
		}
	}
	return nil
}

package ffcss

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/hoisie/mustache"
)

func renderFileTemplate(f FileTemplate, operatingSystem string, variant Variant, osRenameMap map[string]string) string {
	if strings.Contains(strings.Trim(f, " "), "{{variant}}") && variant.Name == "" {
		LogWarning("%q uses {{variant}} which is empty\n", f)
	}
	var osName string
	if osRenameMap[operatingSystem] == "" {
		osName = operatingSystem
	} else {
		osName = osRenameMap[operatingSystem]
	}
	return mustache.Render(f, map[string]string{
		"os":      osName,
		"variant": variant.Name,
	})
}

// DestinationPathOfAsset computes the destination path of some asset from its path and the destination profile directory
// It is assumed that assetPath is absolute.
func (t Theme) DestinationPathOfAsset(assetPath string, profileDir string, operatingSystem string, variant Variant) (string, error) {
	if !strings.HasPrefix(assetPath, t.DownloadedTo) {
		return "", fmt.Errorf("asset %q is outside of the theme's root %q", assetPath, t.DownloadedTo)
	}

	relativeTo := filepath.Clean(filepath.Join(t.DownloadedTo, filepath.Clean(renderFileTemplate(t.CopyFrom, operatingSystem, variant, t.OSNames))))
	if !strings.HasPrefix(relativeTo, t.DownloadedTo) {
		return "", fmt.Errorf("copy from %q is outside of the theme's root %q", relativeTo, t.DownloadedTo)
	}

	relativised, err := filepath.Rel(relativeTo, assetPath)
	if err != nil {
		return "", fmt.Errorf("couldn't make %s relative to %s: %w", assetPath, filepath.Join(t.DownloadedTo, filepath.Clean(t.CopyFrom)), err)
	}

	return filepath.Join(profileDir, "chrome", relativised), nil
}

// AssetsPaths returns the individual file paths of all assets.
func (t Theme) AssetsPaths(os string, variant Variant) ([]string, error) {
	resolvedFiles := make([]string, 0)
	for _, template := range t.Assets {
		glob := renderFileTemplate(template, os, variant, t.OSNames)
		LogDebug("looking for assets: globbing %q", filepath.Join(t.DownloadedTo, glob))
		glob = filepath.Clean(filepath.Join(t.DownloadedTo, glob))
		files, err := doublestar.Glob(glob)
		if err != nil {
			return resolvedFiles, fmt.Errorf("while getting all matches of glob %s: %w", glob, err)
		}
		// If no matches
		if len(files) < 1 {
			// If it's _really_ a glob pattern
			if strings.Contains(glob, "*") {
				return resolvedFiles, fmt.Errorf("glob pattern %s matches no files", glob)
			}
			// If it's just a regular file (that was treated as a glob pattern)
			return resolvedFiles, fmt.Errorf("file %s not found", glob)
		}
		// For each match of the glob pattern
		resolvedFiles = append(resolvedFiles, files...)
	}
	return resolvedFiles, nil
}

package main

import (
	"fmt"
	"os"
)

// CleanDownloadArea removes the temporary download area used to download themes before knowing their name from their manifest
func CleanDownloadArea() error {
	return os.RemoveAll(CacheDir(TempDownloadsDirName))
}

// ClearWholeCache destroys the cache directory
func ClearWholeCache() error {
	err := os.RemoveAll(CacheDir())
	if err != nil {
		return fmt.Errorf("while deleting directory tree %q: %w", CacheDir(), err)
	}

	err = os.Mkdir(CacheDir(), 0777)
	if err != nil {
		return fmt.Errorf("while re-creating empty cache directory at %q: %w", CacheDir(), err)
	}

	return nil
}

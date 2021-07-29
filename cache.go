package main

import "os"

// CleanDownloadArea removes the temporary download area used to download themes before knowing their name from their manifest
func CleanDownloadArea() error {
	return os.RemoveAll(CacheDir(TempDownloadsDirName))
}

// ClearWholeCache destroys the cache directory
func ClearWholeCache() error {
	return os.RemoveAll(CacheDir())
}

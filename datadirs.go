package ffcss

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateDataDirectories makes the required directories for ffcss to function properly.
func CreateDataDirectories() error {
	err := os.MkdirAll(filepath.Join(getConfigDir(), "themes"), 0777)
	if err != nil {
		return fmt.Errorf("couldn't create data directories: %w", err)
	}
	return nil
}

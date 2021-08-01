package ffcss

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateDataDirectories() error {
	err := os.MkdirAll(filepath.Join(GetConfigDir(), "themes"), 0777)
	if err != nil {
		return fmt.Errorf("couldn't create data directories: %w", err)
	}
	return nil
}

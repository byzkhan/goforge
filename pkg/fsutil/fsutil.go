// Package fsutil provides filesystem helpers for the generator.
package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteFile creates parent directories and writes data to the given path.
func WriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}
	return os.WriteFile(path, data, perm)
}

// DirExists returns true if path is an existing directory.
func DirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

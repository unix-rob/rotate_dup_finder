package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// moveToRotatedDups moves src into a rotated_dups/ sibling directory.
// If a file with the same name already exists there, appends _001, _002, etc.
func moveToRotatedDups(src string) (string, error) {
	dir := filepath.Dir(src)
	name := filepath.Base(src)
	dupsDir := filepath.Join(dir, "rotated_dups")

	if err := os.MkdirAll(dupsDir, 0755); err != nil {
		return "", fmt.Errorf("create rotated_dups: %w", err)
	}

	dst := filepath.Join(dupsDir, name)
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		return dst, os.Rename(src, dst)
	}

	// Conflict: find next available suffix
	ext := filepath.Ext(name)
	stem := name[:len(name)-len(ext)]
	for i := 1; i <= 999; i++ {
		dst = filepath.Join(dupsDir, fmt.Sprintf("%s_%03d%s", stem, i, ext))
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			return dst, os.Rename(src, dst)
		}
	}
	return "", fmt.Errorf("exhausted 999 slots in %s for %s", dupsDir, name)
}

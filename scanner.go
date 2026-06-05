package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// suffixPattern matches filenames like 2006_01_02_15_04_05_001.jpg
// capturing (base)(NNN)(ext)
var suffixPattern = regexp.MustCompile(
	`^(\d{4}_\d{2}_\d{2}_\d{2}_\d{2}_\d{2})_(\d{3})(\.[^.]+)$`,
)

var imageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".bmp": true, ".tiff": true, ".tif": true, ".webp": true,
	".heic": true, ".heif": true, ".raw": true, ".cr2": true,
	".nef": true, ".arw": true,
}

type candidate struct {
	path     string // full path of _xxx file
	basePath string // full path of the base file (without _xxx)
}

// scan walks src recursively, skipping rotated_dups directories and symlinks.
// It yields every image file matching the _NNN suffix pattern whose base file exists.
func scan(src string, log *logger, fn func(candidate) error) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.msg("ERROR", path, err.Error())
			return nil
		}
		if d.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		if d.IsDir() {
			if d.Name() == "rotated_dups" {
				return filepath.SkipDir
			}
			return nil
		}

		name := d.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if !imageExts[ext] {
			return nil
		}

		m := suffixPattern.FindStringSubmatch(name)
		if m == nil {
			return nil
		}

		baseName := m[1] + m[3] // e.g. 2006_01_02_15_04_05.jpg
		dir := filepath.Dir(path)
		basePath := filepath.Join(dir, baseName)

		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			return nil // base file missing — skip
		}

		return fn(candidate{path: path, basePath: basePath})
	})
}

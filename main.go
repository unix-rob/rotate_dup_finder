package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	src := flag.String("src", "", "root directory to search recursively")
	dryRun := flag.Bool("dry-run", false, "print actions without moving files")
	logFile := flag.String("log", "", "path to log file (appended)")
	flag.Parse()

	if *src == "" {
		fmt.Fprintln(os.Stderr, "error: --src is required")
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(*src); err != nil {
		fmt.Fprintf(os.Stderr, "error: source path: %v\n", err)
		os.Exit(1)
	}

	log, err := newLogger(*logFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer log.close()

	err = scan(*src, log, func(c candidate) error {
		matched, angle, err := isRotatedDuplicate(c.basePath, c.path)
		if err != nil {
			log.msg("ERROR", c.path, fmt.Sprintf("comparison failed: %v", err))
			return nil
		}
		if !matched {
			return nil
		}

		label := fmt.Sprintf("ROTATED_%+d", angle)
		if *dryRun {
			log.msg("DRY_RUN:"+label, c.path, fmt.Sprintf("(base: %s)", c.basePath))
			return nil
		}

		dst, err := moveToRotatedDups(c.path)
		if err != nil {
			log.msg("ERROR", c.path, fmt.Sprintf("move failed: %v", err))
			return nil
		}
		log.entry(label, c.path, dst)
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

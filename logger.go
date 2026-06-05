package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type logger struct {
	w io.Writer
	f *os.File
}

func newLogger(path string) (*logger, error) {
	l := &logger{w: os.Stdout}
	if path != "" {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("open log file: %w", err)
		}
		l.f = f
		l.w = io.MultiWriter(os.Stdout, f)
	}
	return l, nil
}

func (l *logger) entry(action, src, dst string) {
	fmt.Fprintf(l.w, "%s  %-24s  %s  →  %s\n",
		time.Now().UTC().Format(time.RFC3339), action, src, dst)
}

func (l *logger) msg(action, src, detail string) {
	fmt.Fprintf(l.w, "%s  %-24s  %s  %s\n",
		time.Now().UTC().Format(time.RFC3339), action, src, detail)
}

func (l *logger) close() {
	if l.f != nil {
		l.f.Close()
	}
}

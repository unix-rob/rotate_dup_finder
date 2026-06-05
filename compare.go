package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const ssimThreshold = 0.9

type dims struct{ w, h int }

func getDims(path string) (dims, error) {
	out, err := exec.Command("ffprobe", "-v", "quiet",
		"-print_format", "json", "-show_streams", path).Output()
	if err != nil {
		return dims{}, err
	}
	var r struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(out, &r); err != nil {
		return dims{}, err
	}
	for _, s := range r.Streams {
		if s.Width > 0 && s.Height > 0 {
			return dims{s.Width, s.Height}, nil
		}
	}
	return dims{}, fmt.Errorf("no video stream with dimensions in %s", path)
}

// isRotatedDuplicate returns true if candidate is a 90° CW or CCW rotation of
// base. Also returns the matching angle (90 or -90).
func isRotatedDuplicate(base, candidate string) (bool, int, error) {
	baseDims, err := getDims(base)
	if err != nil {
		return false, 0, fmt.Errorf("dims %s: %w", base, err)
	}

	rotations := []struct {
		angle     int
		transpose int
	}{
		{90, 1},  // 90° clockwise
		{-90, 2}, // 90° counter-clockwise
	}

	for _, r := range rotations {
		ssim, err := rotatedSSIM(base, candidate, baseDims, r.transpose)
		if err != nil {
			continue
		}
		if ssim >= ssimThreshold {
			return true, r.angle, nil
		}
	}
	return false, 0, nil
}

// rotatedSSIM applies transpose to candidate then compares SSIM against base.
// Both are scaled to baseDims before comparison.
func rotatedSSIM(base, candidate string, baseDims dims, transpose int) (float64, error) {
	bw, bh := baseDims.w, baseDims.h
	filter := fmt.Sprintf(
		"[0:v]scale=%d:%d:flags=lanczos[a];[1:v]transpose=%d,scale=%d:%d:flags=lanczos[rot];[a][rot]ssim=stats_file=-",
		bw, bh, transpose, bw, bh,
	)
	out, _ := exec.Command("ffmpeg",
		"-i", base,
		"-i", candidate,
		"-filter_complex", filter,
		"-f", "null", "-").CombinedOutput()

	return parseSSIM(string(out))
}

func parseSSIM(output string) (float64, error) {
	for _, line := range strings.Split(output, "\n") {
		if idx := strings.Index(line, "All:"); idx >= 0 {
			rest := strings.TrimSpace(line[idx+4:])
			if parts := strings.Fields(rest); len(parts) > 0 {
				if v, err := strconv.ParseFloat(parts[0], 64); err == nil {
					return v, nil
				}
			}
		}
	}
	return 0, fmt.Errorf("could not parse SSIM from ffmpeg output")
}

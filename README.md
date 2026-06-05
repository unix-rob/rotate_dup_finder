# rotate_dup_finder

Scans a directory tree for image files that were stored with a `_NNN` conflict suffix (e.g. `2006_03_15_14_30_00_001.jpg`) and checks whether they are rotated duplicates of the corresponding base file (`2006_03_15_14_30_00.jpg`). Rotated duplicates are moved to a `rotated_dups/` folder in the same date directory.

Designed to be run after [find_images](https://github.com/unix-rob/find_images) has organized a photo library.

## Requirements

- [Go](https://go.dev/) 1.21+
- `ffmpeg` and `ffprobe`

```
sudo apt install ffmpeg
```

## Build

```
go build -o rotate_dup_finder .
```

## Usage

```
./rotate_dup_finder --src /path/to/organized [--dry-run] [--log run.log]
```

| Flag | Required | Description |
|------|----------|-------------|
| `--src` | yes | Root directory to search recursively |
| `--dry-run` | no | Print matches without moving any files |
| `--log` | no | Append log entries to this file in addition to stdout |

## How it works

1. Walks `--src` recursively. Skips `rotated_dups/` directories and symlinks.
2. Finds image files whose names match `yyyy_mm_dd_hh_mm_ss_NNN.ext` (3-digit suffix).
3. Checks that the base file `yyyy_mm_dd_hh_mm_ss.ext` exists in the same directory.
4. Runs ffmpeg SSIM with a 90° CW rotation applied to the `_NNN` file, then again with 90° CCW.
   - Both files are scaled to the base file's resolution before comparison.
5. If either rotation scores SSIM ≥ 0.9, the `_NNN` file is moved to `rotated_dups/` in the same directory.
6. If `rotated_dups/` already contains a file with the same name, a `_001`, `_002`... suffix is appended.

## Supported image formats

`.jpg` `.jpeg` `.png` `.gif` `.bmp` `.tiff` `.tif` `.webp` `.heic` `.heif` `.raw` `.cr2` `.nef` `.arw`

## Performance

Each candidate file requires two ffmpeg SSIM comparisons (one per rotation). For large libraries, running with `--dry-run` first is recommended to preview matches before committing to moves.

## Log actions

| Action | Meaning |
|--------|---------|
| `ROTATED_+90` | `_NNN` file matched at 90° clockwise — moved to `rotated_dups/` |
| `ROTATED_-90` | `_NNN` file matched at 90° counter-clockwise — moved to `rotated_dups/` |
| `ERROR` | Per-file error (processing continues) |
| `DRY_RUN:*` | Dry-run prefix on any of the above |

# rotate_dup_finder

CLI tool that finds `_NNN` suffix image files and checks whether they are rotated duplicates of their base file. Designed to run after `find_images` has organized a photo library.

## Build

```
go build -o rotate_dup_finder .
```

## Run

```
./rotate_dup_finder --src /path/to/organized [--dry-run] [--log run.log]
```

## System dependencies

```
sudo apt install ffmpeg
```

- `ffmpeg` / `ffprobe` — rotation SSIM comparison and dimension lookup

## Project structure

| File | Purpose |
|------|---------|
| `main.go` | CLI flags, orchestration |
| `scanner.go` | Recursive walk, `_NNN` pattern matching, base file lookup |
| `compare.go` | ffprobe dimensions, ffmpeg SSIM with transpose filter |
| `mover.go` | Move matched file to `rotated_dups/` sibling directory |
| `logger.go` | Structured one-line log entries to stdout + optional file |

## Key behaviours

- Only checks image files (not video)
- Checks 90° CW (`transpose=1`) and 90° CCW (`transpose=2`) rotations
- SSIM threshold: 0.9; candidate is scaled to base file's resolution before comparison
- `rotated_dups/` directories are skipped during the walk to avoid re-processing
- Conflicts inside `rotated_dups/` get a `_001`, `_002`... suffix

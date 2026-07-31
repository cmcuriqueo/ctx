# Agent Notes for ctx

## Project

- **Language**: Go 1.23+
- **CLI framework**: [cobra](https://github.com/spf13/cobra)
- **Module**: `github.com/matias/ctx`
- **Entry point**: `cmd/ctx/main.go`
- **Packages**: `internal/*` for implementation, `pkg/models` for shared types.

## Build & Test

```bash
# Build binary
go build -o ctx ./cmd/ctx

# Run all tests
go test ./...

# Run with race detector
go test -race ./...
```

## Architecture

Keep modules independent:

- `scanner` — filesystem walk and metadata extraction.
- `ignore` — exclusion rules from defaults, `.gitignore`, `.ctxignore`.
- `tokens` — token estimation strategies.
- `rank` — file relevance scoring.
- `builder` — context.md generation.
- `cache` — JSON persistence by file hash.

## Conventions

- Use `filepath.ToSlash` for stored relative paths.
- Return errors; do not log inside packages unless necessary.
- Keep CLI flags in `main.go`.
- Avoid adding heavy dependencies in the MVP.

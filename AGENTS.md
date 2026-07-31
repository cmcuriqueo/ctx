# Agent Notes for ctx

## Project

- **Language**: Go 1.23+
- **CLI framework**: [cobra](https://github.com/spf13/cobra)
- **Module**: `github.com/matias/ctx`
- **Entry point**: `cmd/ctx/main.go`
- **Packages**:
  - `internal/scanner` — filesystem walk and metadata extraction
  - `internal/ignore` — exclusion rules from defaults, `.gitignore`, `.ctxignore`
  - `internal/parser` — Tree-sitter based import/export extraction
  - `internal/graph` — dependency graph and algorithms
  - `internal/tokens` — token estimation strategies
  - `internal/rank` — file relevance scoring
  - `internal/config` — `ctx.toml` loading
  - `internal/builder` — context.md generation
  - `internal/cache` — JSON persistence
  - `pkg/models` — shared types

## Build & Test

Tree-sitter requires CGO. On Windows use the bundled MinGW:

```powershell
$env:CC="C:\Users\matias\ctx\.tools\mingw64\bin\gcc.exe"
$env:CXX="C:\Users\matias\ctx\.tools\mingw64\bin\g++.exe"
go build -o ctx.exe ./cmd/ctx
go test ./...
```

On Linux/macOS:

```bash
CGO_ENABLED=1 go build -o ctx ./cmd/ctx
go test ./...
```

## Conventions

- Use `filepath.ToSlash` for stored relative paths.
- Return errors; do not log inside packages unless necessary.
- Keep CLI flags in `main.go`.
- Tree-sitter queries live next to each parser (`go.go`, `javascript.go`, etc.).

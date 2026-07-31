# ctx

`ctx` is a context builder CLI for AI. It scans a codebase, builds a dependency graph, ranks files by relevance, respects a token budget, and generates a ready-to-paste `context.md`.

## Features

- **Multi-language parsing** with Tree-sitter (Go, TypeScript, JavaScript, Python).
- **Dependency graph** to discover related files without sending code to an API.
- **Configurable ranking** via `ctx.toml`.
- **Offline token estimation** heuristic.
- **Ignore engine** combining `.gitignore`, `.ctxignore`, and built-in rules.

## Requirements

- Go 1.23 or later.
- A C compiler for Tree-sitter (CGO):
  - **Windows**: MinGW-w64 (the repo includes a portable copy under `.tools/mingw64`).
  - **Linux/macOS**: `gcc` or `clang`.

## Install

### Automatic install

**Windows** (PowerShell as normal user):

```powershell
.\scripts\install-windows.ps1
```

> The script downloads a portable MinGW compiler into `.tools/mingw64` if `gcc` is not found.

**Linux / macOS** (Bash):

```bash
bash scripts/install-linux.sh
```

### Manual install

```bash
# Windows with bundled MinGW
$env:CC="C:\Users\<you>\ctx\.tools\mingw64\bin\gcc.exe"
$env:CXX="C:\Users\<you>\ctx\.tools\mingw64\bin\g++.exe"
go build -o ctx.exe ./cmd/ctx

# Linux / macOS
CGO_ENABLED=1 go build -o ctx ./cmd/ctx
```

## Usage

```bash
# Scan current directory (extracts imports/exports and builds graph)
ctx scan

# Show dependency graph
ctx graph
ctx graph --format dot
ctx graph --format json

# Estimate tokens per file
ctx tokens

# Generate context.md with default 4000 token budget
ctx bundle

# Custom budget and output
ctx bundle -b 8000 -o docs/context.md
```

## Commands

| Command | Description |
|---------|-------------|
| `scan [path]` | Walk the repository, parse imports/exports, and cache the manifest + graph. |
| `graph [path]` | Display the dependency graph (text, dot, or json). |
| `tokens [path]` | Estimate tokens for each scanned file using an offline heuristic. |
| `bundle [path]` | Generate `context.md` with the most relevant files within budget. |

## Configuration

Create `ctx.toml` in the project root:

```toml
[rank]
entrypoint = 30
readme = 20
config = 20
test = 10
generated = -40
vendor = -30
imported_bonus = 15
large_isolated = -20

[graph]
max_depth = 5
```

## How it works

1. **Scanner** walks the directory tree, skipping files matched by `.gitignore`, `.ctxignore`, and built-in rules.
2. **Parser** uses Tree-sitter to extract package/module names, imports, and exports.
3. **Graph Builder** resolves imports to local files and builds a directed dependency graph.
4. **Ranking Engine** scores files using configurable weights, including a bonus for files imported by many others.
5. **Context Builder** selects entrypoints, expands via the graph, and fills the remaining budget greedily by `score/tokens`.

## Ignored by default

- `node_modules`, `dist`, `build`, `coverage`, `.git`, `bin`, `target`, `.cache`
- `context.md`, `context*.md`
- Binary files: `*.exe`, `*.dll`, `*.png`, `*.jpg`, `*.mp4`, etc.
- Anything listed in `.gitignore` or `.ctxignore`

## Roadmap

- **V3**: task mode, optional AI-assisted selection (tree + metadata only).
- **V4**: plugins, monorepo support, LSP integration, local embeddings.

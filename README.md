# ctx

`ctx` is a context builder CLI for AI. It scans a codebase, ranks files by relevance, respects a token budget, and generates a ready-to-paste `context.md`.

> This is the MVP. Advanced features (parser, dependency graph, AI-assisted ranking) are planned for future versions.

## Install

Requires Go 1.23 or later.

### Automatic install

**Windows** (PowerShell as normal user):

```powershell
.\scripts\install-windows.ps1
```

**Linux / macOS** (Bash):

```bash
bash scripts/install-linux.sh
```

These scripts compile `ctx`, copy it to a user directory (`C:\Users\<you>\tools` on Windows, `~/.local/bin` on Linux/macOS), and add that directory to your `PATH`.

### Manual install

```bash
go build -o ctx ./cmd/ctx
```

Then copy the binary to a directory that is already in your `PATH`.

## Usage

```bash
# Scan current directory and cache manifest
ctx scan

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
| `scan [path]` | Walk the repository and store a manifest in `.cache/ctx/manifest.json`. |
| `tokens [path]` | Estimate tokens for each scanned file using an offline heuristic. |
| `bundle [path]` | Generate `context.md` with the most relevant files within budget. |

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `-b, --budget` | `4000` | Maximum tokens to include in the output. |
| `-o, --output` | `context.md` | Path to the generated context file. |
| `--cache-dir` | `.cache/ctx` | Directory used for cached metadata. |

## How it works

1. **Scanner** walks the directory tree, skipping files matched by `.gitignore`, `.ctxignore`, and built-in rules.
2. **Ignore Engine** excludes common folders like `node_modules`, `dist`, `build`, `.git`, and binary files.
3. **Token Estimator** uses an offline heuristic (~4 characters per token).
4. **Ranking Engine** scores files by path heuristics: entrypoints, READMEs, configs, tests, generated code, vendor code.
5. **Context Builder** selects files greedily by `score/tokens` until the budget is exhausted and writes `context.md`.

## Ignored by default

- `node_modules`, `dist`, `build`, `coverage`, `.git`, `bin`, `target`, `.cache`
- Binary files: `*.exe`, `*.dll`, `*.so`, `*.png`, `*.jpg`, `*.mp4`, etc.
- Anything listed in `.gitignore` or `.ctxignore`

## Roadmap

- **V2**: language parser, dependency graph, configurable ranking formula.
- **V3**: task mode, optional AI-assisted selection (tree + metadata only).
- **V4**: plugins, monorepo support, LSP integration, local embeddings.

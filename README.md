# ctx

`ctx` is a context builder CLI for AI. It scans a codebase, builds a dependency graph, ranks files by relevance, respects a token budget, and generates a ready-to-paste `context.md`.

## Features

- **Task mode** — ask an LLM which files are relevant for a task without sending your source code.
- **Multi-language parsing** with Tree-sitter (Go, TypeScript, JavaScript, Python).
- **Dependency graph** to discover related files without sending code to an API.
- **Configurable ranking** via `ctx.toml`.
- **Offline fallback** when no API key is available.
- **Ignore engine** combining `.gitignore`, `.ctxignore`, and built-in rules.

## Requirements

- Go 1.24 or later.
- A C compiler for Tree-sitter (CGO):
  - **Windows**: MinGW-w64 (the repo includes a portable copy under `.tools/mingw64`).
  - **Linux/macOS**: `gcc` or `clang`.
- (Optional) An API key for task mode:
  - `OPENAI_API_KEY` for OpenAI
  - `ANTHROPIC_API_KEY` for Anthropic

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

# Ask an LLM to pick files for a task
ctx task "Add OAuth with Google"
ctx task "Add OAuth with Google" --provider anthropic --model claude-sonnet-4-5

# Offline mode (uses heuristics, no API call)
ctx task "Fix memory leak" --offline
```

## Commands

| Command | Description |
|---------|-------------|
| `scan [path]` | Walk the repository, parse imports/exports, and cache the manifest + graph. |
| `graph [path]` | Display the dependency graph (text, dot, or json). |
| `tokens [path]` | Estimate tokens for each scanned file using an offline heuristic. |
| `bundle [path]` | Generate `context.md` using graph + ranking heuristics. |
| `task <description> [path]` | Ask an LLM to select files, then generate `context.md`. |

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

[llm]
provider = "openai"      # or "anthropic"
model = "gpt-4o-mini"    # or "claude-haiku-4-5"
# api_key = "..."        # optional, env vars are preferred
```

Environment variables take priority:

- `CTX_LLM_PROVIDER`
- `CTX_LLM_MODEL`
- `OPENAI_API_KEY`
- `ANTHROPIC_API_KEY`
- `OPENAI_BASE_URL` (for proxies or OpenAI-compatible endpoints)

## Privacy

`ctx task` **never sends file contents** to the LLM. It only sends metadata: paths, language, imports, exports, line counts, and token estimates.

## How it works

1. **Scanner** walks the directory tree, skipping files matched by `.gitignore`, `.ctxignore`, and built-in rules.
2. **Parser** uses Tree-sitter to extract package/module names, imports, and exports.
3. **Graph Builder** resolves imports to local files and builds a directed dependency graph.
4. **Task mode** sends a compact project summary to the configured LLM and receives a list of relevant files.
5. **Context Builder** generates `context.md` from the selected files within the token budget.

## Ignored by default

- `node_modules`, `dist`, `build`, `coverage`, `.git`, `bin`, `target`, `.cache`
- `context.md`, `context*.md`
- Binary files: `*.exe`, `*.dll`, `*.png`, `*.jpg`, `*.mp4`, etc.
- Anything listed in `.gitignore` or `.ctxignore`

## Roadmap

- **V4**: plugins, monorepo support, LSP integration, local embeddings.

# Context

Project root: `C:\Users\matias\ctx`

Files scanned: 37

Budget: 3000 tokens

Used: 2942 tokens (98.1%)

## Included files

- `internal/tokens/tokens_test.go` (go, 103 tokens, score 10)
- `go.mod` (go, 310 tokens, score 20)
- `internal/anthropic/anthropic_test.go` (go, 176 tokens, score 10)
- `internal/llm/llm_test.go` (go, 180 tokens, score 10)
- `internal/rank/rank_test.go` (go, 185 tokens, score 10)
- `internal/openai/openai_test.go` (go, 201 tokens, score 10)
- `internal/config/config_test.go` (go, 225 tokens, score 10)
- `README.md` (markdown, 1023 tokens, score 45)
- `internal/parser/parser_test.go` (go, 351 tokens, score 10)
- `.gitignore` (gitignore, 21 tokens, score 0)
- `internal/parser/go.go` (go, 167 tokens, score 0)

## Snippets

### `internal/tokens/tokens_test.go`

```go
package tokens

import "testing"

func TestHeuristicEstimator(t *testing.T) {
	est := NewHeuristicEstimator()

	cases := []struct {
		content string
		want    int
	}{
		{"", 0},
		{"abc", 1},
		{"abcdefghijklmnop", 4},
		{"this is a short sentence", 6},
	}

	for _, c := range cases {
		got := est.Estimate(c.content)
		if got != c.want {
			t.Errorf("Estimate(%q) = %d, want %d", c.content, got, c.want)
		}
	}
}

```

### `go.mod`

```go
module github.com/matias/ctx

go 1.24.1

require (
	github.com/BurntSushi/toml v1.6.0
	github.com/anthropics/anthropic-sdk-go v1.61.0
	github.com/sabhiram/go-gitignore v0.0.0-20210923224102-525f6e181f06
	github.com/sashabaranov/go-openai v1.41.2
	github.com/spf13/cobra v1.10.2
	github.com/tree-sitter/go-tree-sitter v0.25.0
	github.com/tree-sitter/tree-sitter-go v0.25.0
	github.com/tree-sitter/tree-sitter-javascript v0.25.0
	github.com/tree-sitter/tree-sitter-python v0.25.0
	github.com/tree-sitter/tree-sitter-typescript v0.23.2
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/invopop/jsonschema v0.14.0 // indirect
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/pb33f/ordered-map/v2 v2.3.1 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/standard-webhooks/standard-webhooks/libraries v0.0.1 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	go.yaml.in/yaml/v4 v4.0.0-rc.2 // indirect
	golang.org/x/sync v0.16.0 // indirect
)

```

### `internal/anthropic/anthropic_test.go`

```go
package anthropic

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProviderSelectFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": `{"explanation": "test", "paths": ["main.go"]}`,
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := New("fake-key", "claude-haiku-4-5")
	// The official SDK does not expose baseURL override easily, so this test only checks construction.
	if p.Name() != "anthropic" {
		t.Errorf("name = %s, want anthropic", p.Name())
	}
}

```

### `internal/llm/llm_test.go`

```go
package llm

import (
	"testing"

	"github.com/matias/ctx/internal/llm/types"
)

type mockProvider struct {
	name string
}

func (m *mockProvider) Name() string { return m.name }

func (m *mockProvider) SelectFiles(task string, summary *types.ProjectSummary) (*types.Selection, error) {
	return &types.Selection{
		Explanation: "mock selection",
		Paths:       []string{"main.go"},
	}, nil
}

func TestNewDefaultsToOpenAI(t *testing.T) {
	_, err := New(types.Options{})
	if err == nil {
		t.Error("expected error without OPENAI_API_KEY")
	}
}

func TestNewAnthropicRequiresKey(t *testing.T) {
	_, err := New(types.Options{Provider: "anthropic"})
	if err == nil {
		t.Error("expected error without ANTHROPIC_API_KEY")
	}
}

```

### `internal/rank/rank_test.go`

```go
package rank

import (
	"testing"

	"github.com/matias/ctx/internal/config"
	"github.com/matias/ctx/pkg/models"
)

func TestScorer(t *testing.T) {
	scorer := NewScorer(config.Default())

	cases := []struct {
		file models.FileInfo
		min  int
		max  int
	}{
		{models.FileInfo{Path: "main.go"}, 30, 30},
		{models.FileInfo{Path: "README.md"}, 20, 20},
		{models.FileInfo{Path: "go.mod"}, 20, 20},
		{models.FileInfo{Path: "foo_test.go"}, 10, 10},
		{models.FileInfo{Path: "vendor/foo.go"}, -30, -30},
		{models.FileInfo{Path: "src/app.go"}, 0, 0},
	}

	for _, c := range cases {
		got := scorer.Score(c.file, nil)
		if got < c.min || got > c.max {
			t.Errorf("Score(%q) = %d, want between %d and %d", c.file.Path, got, c.min, c.max)
		}
	}
}

```

### `internal/openai/openai_test.go`

```go
package openai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matias/ctx/internal/llm/types"
)

func TestProviderSelectFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": `{"explanation": "test", "paths": ["main.go"]}`,
					},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := New("fake-key", "gpt-4o-mini", server.URL)
	sel, err := p.SelectFiles("task", &types.ProjectSummary{})
	if err != nil {
		t.Fatal(err)
	}
	if len(sel.Paths) != 1 || sel.Paths[0] != "main.go" {
		t.Errorf("unexpected selection: %v", sel.Paths)
	}
}

```

### `internal/config/config_test.go`

```go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := Default()
	if cfg.Rank.Entrypoint != 30 {
		t.Errorf("entrypoint = %d, want 30", cfg.Rank.Entrypoint)
	}
	if cfg.Graph.MaxDepth != 5 {
		t.Errorf("max_depth = %d, want 5", cfg.Graph.MaxDepth)
	}
}

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	data := `
[rank]
entrypoint = 50
imported_bonus = 25

[graph]
max_depth = 3
`
	if err := os.WriteFile(filepath.Join(dir, "ctx.toml"), []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Rank.Entrypoint != 50 {
		t.Errorf("entrypoint = %d, want 50", cfg.Rank.Entrypoint)
	}
	if cfg.Rank.ImportedBonus != 25 {
		t.Errorf("imported_bonus = %d, want 25", cfg.Rank.ImportedBonus)
	}
	if cfg.Graph.MaxDepth != 3 {
		t.Errorf("max_depth = %d, want 3", cfg.Graph.MaxDepth)
	}
}

```

### `README.md`

```markdown
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

```

### `internal/parser/parser_test.go`

```go
package parser

import (
	"testing"
)

func TestGoParser(t *testing.T) {
	src := []byte(`package main

import "fmt"
import "github.com/user/repo/helper"

func main() {
	fmt.Println(helper.Greet())
}

func Greet() string { return "hi" }

type User struct{}
`)
	p := NewGoParser()
	meta, err := p.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Package != "main" {
		t.Errorf("package = %q, want main", meta.Package)
	}
	if len(meta.Imports) != 2 {
		t.Errorf("imports = %v, want 2", meta.Imports)
	}
	if len(meta.Exports) != 2 {
		t.Errorf("exports = %v, want 2", meta.Exports)
	}
}

func TestJavaScriptParser(t *testing.T) {
	src := []byte(`import { helper } from './helper';

export function greet() {
	return helper();
}

export class User {}
`)
	p := NewJavaScriptParser()
	meta, err := p.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(meta.Imports) != 1 {
		t.Errorf("imports = %v, want 1", meta.Imports)
	}
	if len(meta.Exports) != 2 {
		t.Errorf("exports = %v, want 2", meta.Exports)
	}
}

func TestPythonParser(t *testing.T) {
	src := []byte(`import os
from . import helper

def greet():
	return helper.message()

class User:
	pass
`)
	p := NewPythonParser()
	meta, err := p.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(meta.Imports) != 2 {
		t.Errorf("imports = %v, want 2", meta.Imports)
	}
	if len(meta.Exports) != 2 {
		t.Errorf("exports = %v, want 2", meta.Exports)
	}
}

```

### `.gitignore`

```gitignore
# Binaries
ctx.exe
ctx
*.exe

# Cache
.cache/

# Local tooling
.tools/

# Logs
*.log

```

### `internal/parser/go.go`

```go
package parser

import tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"

const goQuery = `
(package_clause (package_identifier) @package)
(import_spec path: (interpreted_string_literal) @import)
(import_spec path: (raw_string_literal) @import)
(function_declaration name: (identifier) @export (#match? @export "^[A-Z]"))
(method_declaration name: (field_identifier) @export (#match? @export "^[A-Z]"))
(type_declaration (type_spec name: (type_identifier) @export (#match? @export "^[A-Z]")))
`

// NewGoParser returns a Go parser.
func NewGoParser() Parser {
	return &treeSitterParser{
		language:  tree_sitter_go.Language(),
		queryText: goQuery,
	}
}

```


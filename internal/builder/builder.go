package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/matias/ctx/internal/config"
	"github.com/matias/ctx/internal/graph"
	"github.com/matias/ctx/internal/llm/types"
	"github.com/matias/ctx/internal/rank"
	"github.com/matias/ctx/internal/tokens"
	"github.com/matias/ctx/pkg/models"
)

// Builder selects files and generates context.md.
type Builder struct {
	estimator tokens.Estimator
	scorer    *rank.Scorer
	graph     *graph.Graph
	cfg       *config.Config
}

// New creates a builder with the given dependencies.
func New(estimator tokens.Estimator, scorer *rank.Scorer, g *graph.Graph, cfg *config.Config) *Builder {
	return &Builder{estimator: estimator, scorer: scorer, graph: g, cfg: cfg}
}

// Build generates the context file for the given manifest and budget.
func (b *Builder) Build(manifest *models.Manifest, budget int, output string) error {
	related := b.relatedFiles(manifest)

	var ranked []models.RankedFile
	for _, f := range manifest.Files {
		if f.IsBinary {
			continue
		}
		fullPath := filepath.Join(manifest.Root, f.Path)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}
		tok := b.estimator.Estimate(string(content))
		score := b.scorer.Score(f, b.graph)
		if _, ok := related[f.Path]; ok {
			score += 25 // neighborhood boost from entrypoints
		}
		ranked = append(ranked, models.RankedFile{
			FileInfo:   f,
			Score:      score,
			Tokens:     tok,
			Efficiency: float64(score) / float64(max(tok, 1)),
		})
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].Efficiency > ranked[j].Efficiency
	})

	selected := []models.RankedFile{}
	used := 0
	for _, rf := range ranked {
		if rf.Tokens == 0 {
			continue
		}
		if used+rf.Tokens > budget {
			continue
		}
		selected = append(selected, rf)
		used += rf.Tokens
	}

	return b.writeContext(manifest, selected, used, budget, output, "")
}

// BuildFromSelection generates context.md from an LLM selection.
func (b *Builder) BuildFromSelection(manifest *models.Manifest, selection *types.Selection, budget int, output string) error {
	selectedSet := make(map[string]struct{}, len(selection.Paths))
	for _, p := range selection.Paths {
		selectedSet[p] = struct{}{}
	}

	var selected []models.RankedFile
	used := 0
	for _, f := range manifest.Files {
		if _, ok := selectedSet[f.Path]; !ok {
			continue
		}
		if f.IsBinary {
			continue
		}
		fullPath := filepath.Join(manifest.Root, f.Path)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}
		tok := b.estimator.Estimate(string(content))
		if used+tok > budget {
			continue
		}
		selected = append(selected, models.RankedFile{
			FileInfo: f,
			Score:    b.scorer.Score(f, b.graph),
			Tokens:   tok,
		})
		used += tok
	}

	return b.writeContext(manifest, selected, used, budget, output, selection.Explanation)
}

func (b *Builder) writeContext(manifest *models.Manifest, selected []models.RankedFile, used, budget int, output, explanation string) error {
	var sb strings.Builder
	sb.WriteString("# Context\n\n")
	if explanation != "" {
		sb.WriteString("## AI reasoning\n\n")
		sb.WriteString(explanation)
		sb.WriteString("\n\n")
	}
	sb.WriteString(fmt.Sprintf("Project root: `%s`\n\n", manifest.Root))
	sb.WriteString(fmt.Sprintf("Files scanned: %d\n\n", len(manifest.Files)))
	sb.WriteString(fmt.Sprintf("Budget: %d tokens\n\n", budget))
	sb.WriteString(fmt.Sprintf("Used: %d tokens (%.1f%%)\n\n", used, float64(used)*100/float64(max(budget, 1))))

	sb.WriteString("## Included files\n\n")
	for _, rf := range selected {
		sb.WriteString(fmt.Sprintf("- `%s` (%s, %d tokens, score %d)\n", rf.Path, rf.Language, rf.Tokens, rf.Score))
	}

	sb.WriteString("\n## Snippets\n\n")
	for _, rf := range selected {
		fullPath := filepath.Join(manifest.Root, rf.Path)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}
		sb.WriteString(fmt.Sprintf("### `%s`\n\n```%s\n%s\n```\n\n", rf.Path, rf.Language, string(content)))
	}

	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}
	return os.WriteFile(output, []byte(sb.String()), 0644)
}

// relatedFiles returns files reachable from entrypoints/readmes via the dependency graph.
func (b *Builder) relatedFiles(manifest *models.Manifest) map[string]struct{} {
	related := make(map[string]struct{})
	if b.graph == nil {
		return related
	}

	var seeds []string
	for _, f := range manifest.Files {
		name := strings.ToLower(filepath.Base(f.Path))
		if isSeed(name) {
			seeds = append(seeds, f.Path)
		}
	}

	maxDepth := 5
	if b.cfg != nil {
		maxDepth = b.cfg.Graph.MaxDepth
	}
	for _, seed := range seeds {
		for _, n := range b.graph.BFS(seed, maxDepth) {
			related[n] = struct{}{}
		}
	}
	return related
}

func isSeed(name string) bool {
	return name == "main.go" ||
		name == "index.js" || name == "index.mjs" || name == "index.cjs" ||
		name == "index.ts" || name == "index.tsx" || name == "index.jsx" ||
		name == "app.py" || name == "main.py" || name == "main.rs" || name == "main.java" ||
		name == "app.java" || name == "program.cs" ||
		strings.HasPrefix(name, "readme.")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

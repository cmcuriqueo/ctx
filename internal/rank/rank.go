package rank

import (
	"path/filepath"
	"strings"

	"github.com/matias/ctx/internal/config"
	"github.com/matias/ctx/internal/graph"
	"github.com/matias/ctx/pkg/models"
)

// Scorer assigns relevance scores to files based on heuristics and graph data.
type Scorer struct {
	cfg *config.Config
}

// NewScorer returns a scorer using the given config.
func NewScorer(cfg *config.Config) *Scorer {
	return &Scorer{cfg: cfg}
}

// Score returns a heuristic relevance score for a file.
func (s *Scorer) Score(f models.FileInfo, g *graph.Graph) int {
	c := s.cfg.Rank
	score := 0
	name := strings.ToLower(filepath.Base(f.Path))
	rel := strings.ToLower(filepath.ToSlash(f.Path))

	if isEntrypoint(name) {
		score += c.Entrypoint
	}
	if strings.HasPrefix(name, "readme.") {
		score += c.Readme
	}
	if isConfig(name, rel) {
		score += c.ConfigFile
	}
	if strings.Contains(name, "_test.") || strings.Contains(name, ".test.") {
		score += c.Test
	}
	if isGenerated(name, rel) {
		score += c.Generated
	}
	if isVendor(rel) {
		score += c.Vendor
	}
	if g != nil {
		score += g.ImportedByCount(f.Path) * c.ImportedBonus
	}
	if f.Size > 100_000 && score <= 0 {
		score += c.LargeIsolated
	}
	return score
}

func isEntrypoint(name string) bool {
	switch name {
	case "main.go", "index.js", "index.mjs", "index.cjs",
		"index.ts", "index.tsx", "index.jsx",
		"app.py", "main.py", "main.rs", "main.java",
		"app.java", "program.cs":
		return true
	}
	return false
}

func isConfig(name, rel string) bool {
	switch name {
	case "go.mod", "go.sum", "package.json", "package-lock.json",
		"yarn.lock", "pnpm-lock.yaml", "cargo.toml", "cargo.lock",
		"pyproject.toml", "requirements.txt", "poetry.lock",
		"dockerfile", "docker-compose.yml", "docker-compose.yaml",
		"makefile", "cmakelists.txt", ".github":
		return true
	}
	return strings.HasPrefix(rel, ".github/")
}

func isGenerated(name, rel string) bool {
	if strings.Contains(rel, "generated") || strings.Contains(rel, "gen/") {
		return true
	}
	suffices := []string{
		".gen.go", ".pb.go", ".min.js", ".min.css",
		".generated.ts", ".generated.js",
	}
	for _, suffix := range suffices {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

func isVendor(rel string) bool {
	return strings.HasPrefix(rel, "vendor/") || strings.Contains(rel, "/vendor/")
}

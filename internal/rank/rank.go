package rank

import (
	"path/filepath"
	"strings"

	"github.com/matias/ctx/pkg/models"
)

// Scorer assigns relevance scores to files based on heuristics.
type Scorer struct{}

// NewScorer returns a new path-based scorer.
func NewScorer() *Scorer {
	return &Scorer{}
}

// Score returns a heuristic relevance score for a file.
func (s *Scorer) Score(f models.FileInfo) int {
	score := 0
	name := strings.ToLower(filepath.Base(f.Path))
	rel := strings.ToLower(filepath.ToSlash(f.Path))

	if isEntrypoint(name) {
		score += 30
	}
	if strings.HasPrefix(name, "readme.") {
		score += 20
	}
	if isConfig(name, rel) {
		score += 20
	}
	if strings.Contains(name, "_test.") || strings.Contains(name, ".test.") {
		score += 10
	}
	if isGenerated(name, rel) {
		score -= 40
	}
	if isVendor(rel) {
		score -= 30
	}
	if f.Size > 100_000 && score <= 0 {
		score -= 20
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
		"makefile", "cmakeLists.txt", ".github":
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

package ignore

import (
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

// defaultPatterns are always excluded, even without .gitignore or .ctxignore.
var defaultPatterns = []string{
	"node_modules", "dist", "build", "coverage", ".git", "bin", "target",
	".cache", "context.md", "context*.md",
	"*.exe", "*.dll", "*.so", "*.dylib", "*.bin", "*.zip", "*.tar.gz",
	"*.jpg", "*.jpeg", "*.png", "*.gif", "*.mp4", "*.mp3", "*.pdf", "*.ttf", "*.woff",
	"*.woff2", "*.eot", "*.ico", "*.svg", "*.webp", "*.avi", "*.mov", "*.mkv",
}

// Engine combines multiple ignore sources.
type Engine struct {
	matchers []*gitignore.GitIgnore
}

// NewEngine builds an ignore engine for the given root directory.
func NewEngine(root string) (*Engine, error) {
	e := &Engine{}
	e.addMatcher(defaultPatterns)
	if patterns, err := loadIgnore(root, ".gitignore"); err == nil {
		e.addMatcher(patterns)
	}
	if patterns, err := loadIgnore(root, ".ctxignore"); err == nil {
		e.addMatcher(patterns)
	}
	return e, nil
}

func loadIgnore(root, name string) ([]string, error) {
	path := filepath.Join(root, name)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var patterns []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns, nil
}

func (e *Engine) addMatcher(patterns []string) {
	if len(patterns) == 0 {
		return
	}
	e.matchers = append(e.matchers, gitignore.CompileIgnoreLines(patterns...))
}

// Match returns true if the path should be ignored.
// isDir should be true when matching a directory path.
func (e *Engine) Match(path string, isDir bool) bool {
	path = filepath.ToSlash(path)
	for _, m := range e.matchers {
		if m.MatchesPath(path) {
			return true
		}
	}
	return false
}

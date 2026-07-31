package ignore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultPatterns(t *testing.T) {
	e := &Engine{}
	e.addMatcher(defaultPatterns)

	cases := []struct {
		path   string
		isDir  bool
		ignore bool
	}{
		{"node_modules", true, true},
		{"node_modules/foo.js", false, true},
		{"dist/main.js", false, true},
		{"src/main.go", false, false},
		{"README.md", false, false},
		{"image.png", false, true},
	}

	for _, c := range cases {
		got := e.Match(c.path, c.isDir)
		if got != c.ignore {
			t.Errorf("Match(%q, %v) = %v, want %v", c.path, c.isDir, got, c.ignore)
		}
	}
}

func TestLoadGitignore(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("*.log\nsecret.txt\n"), 0644); err != nil {
		t.Fatal(err)
	}

	e, err := NewEngine(dir)
	if err != nil {
		t.Fatal(err)
	}

	if !e.Match("debug.log", false) {
		t.Error("expected debug.log to be ignored")
	}
	if !e.Match("secret.txt", false) {
		t.Error("expected secret.txt to be ignored")
	}
	if e.Match("main.go", false) {
		t.Error("did not expect main.go to be ignored")
	}
}

func TestCtxignoreOverrides(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".ctxignore"), []byte("custom.tmp\n"), 0644); err != nil {
		t.Fatal(err)
	}

	e, err := NewEngine(dir)
	if err != nil {
		t.Fatal(err)
	}

	if !e.Match("custom.tmp", false) {
		t.Error("expected custom.tmp to be ignored by .ctxignore")
	}
}

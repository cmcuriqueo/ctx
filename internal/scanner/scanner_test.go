package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matias/ctx/internal/ignore"
	"github.com/matias/ctx/pkg/models"
)

func TestScannerSkipsIgnored(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "node_modules"), 0755)
	os.WriteFile(filepath.Join(dir, "node_modules", "pkg.js"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0644)
	os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("*.log\n"), 0644)
	os.WriteFile(filepath.Join(dir, "debug.log"), []byte("error\n"), 0644)

	ign, err := ignore.NewEngine(dir)
	if err != nil {
		t.Fatal(err)
	}
	s := New(ign)
	manifest, err := s.Scan(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(manifest.Files) != 2 {
		t.Fatalf("expected 2 files (main.go and .gitignore), got %d", len(manifest.Files))
	}
	var mainFile *models.FileInfo
	for i := range manifest.Files {
		if manifest.Files[i].Path == "main.go" {
			mainFile = &manifest.Files[i]
			break
		}
	}
	if mainFile == nil {
		t.Fatalf("expected main.go in manifest")
	}
	if mainFile.Language != "go" {
		t.Errorf("expected language go, got %s", mainFile.Language)
	}
}

func TestScannerDetectsBinary(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "data.dat"), []byte{0, 1, 2, 3}, 0644)

	ign, _ := ignore.NewEngine(dir)
	s := New(ign)
	manifest, err := s.Scan(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(manifest.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(manifest.Files))
	}
	if manifest.Files[0].Path != "data.dat" {
		t.Errorf("expected data.dat, got %s", manifest.Files[0].Path)
	}
	if !manifest.Files[0].IsBinary {
		t.Error("expected data.dat to be detected as binary")
	}
}

package graph

import (
	"strings"
	"testing"

	"github.com/matias/ctx/pkg/models"
)

type staticResolver struct{}

func (r *staticResolver) Resolve(importerPath, importPath string) string {
	importPath = strings.Trim(importPath, `"'`)
	if !strings.HasPrefix(importPath, ".") {
		return ""
	}
	return strings.TrimPrefix(importPath, "./")
}

func TestGraphBFS(t *testing.T) {
	manifest := &models.Manifest{
		Files: []models.FileInfo{
			{Path: "main.go", Imports: []string{"./helper.go"}},
			{Path: "helper.go", Imports: []string{"./util.go"}},
			{Path: "util.go"},
			{Path: "other.go"},
		},
	}
	g := New(manifest, &staticResolver{})
	reached := g.BFS("main.go", 0)
	if len(reached) != 3 {
		t.Errorf("BFS reached %v, want 3 nodes", reached)
	}
}

func TestGraphImportedByCount(t *testing.T) {
	manifest := &models.Manifest{
		Files: []models.FileInfo{
			{Path: "main.go", Imports: []string{"./helper.go"}},
			{Path: "a.go", Imports: []string{"./helper.go"}},
			{Path: "helper.go"},
		},
	}
	g := New(manifest, &staticResolver{})
	if got := g.ImportedByCount("helper.go"); got != 2 {
		t.Errorf("ImportedByCount = %d, want 2", got)
	}
}

func TestGraphCycles(t *testing.T) {
	manifest := &models.Manifest{
		Files: []models.FileInfo{
			{Path: "a.go", Imports: []string{"./b.go"}},
			{Path: "b.go", Imports: []string{"./a.go"}},
		},
	}
	g := New(manifest, &staticResolver{})
	cycles := g.DetectCycles()
	if len(cycles) == 0 {
		t.Error("expected at least one cycle")
	}
}

func TestGraphTopologicalSort(t *testing.T) {
	manifest := &models.Manifest{
		Files: []models.FileInfo{
			{Path: "main.go", Imports: []string{"./helper.go"}},
			{Path: "helper.go", Imports: []string{"./util.go"}},
			{Path: "util.go"},
		},
	}
	g := New(manifest, &staticResolver{})
	order, err := g.TopologicalSort()
	if err != nil {
		t.Fatal(err)
	}
	if len(order) != 3 {
		t.Errorf("order length = %d, want 3", len(order))
	}
}

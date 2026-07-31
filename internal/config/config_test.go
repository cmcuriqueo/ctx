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

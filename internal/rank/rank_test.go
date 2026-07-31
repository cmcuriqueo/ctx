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

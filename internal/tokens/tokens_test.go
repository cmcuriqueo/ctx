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

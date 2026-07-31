package tokens

import "unicode/utf8"

// Estimator defines a strategy for counting tokens.
type Estimator interface {
	Estimate(content string) int
}

// HeuristicEstimator uses an offline character-based approximation.
// It assumes ~4 characters per token on average.
type HeuristicEstimator struct{}

// NewHeuristicEstimator returns a new heuristic estimator.
func NewHeuristicEstimator() *HeuristicEstimator {
	return &HeuristicEstimator{}
}

// Estimate returns an approximate token count for the given content.
func (h *HeuristicEstimator) Estimate(content string) int {
	if content == "" {
		return 0
	}
	runes := utf8.RuneCountInString(content)
	tokens := runes / 4
	if tokens < 1 {
		return 1
	}
	return tokens
}

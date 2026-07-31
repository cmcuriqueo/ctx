package llm

import (
	"testing"

	"github.com/matias/ctx/internal/llm/types"
)

type mockProvider struct {
	name string
}

func (m *mockProvider) Name() string { return m.name }

func (m *mockProvider) SelectFiles(task string, summary *types.ProjectSummary) (*types.Selection, error) {
	return &types.Selection{
		Explanation: "mock selection",
		Paths:       []string{"main.go"},
	}, nil
}

func TestNewDefaultsToOpenAI(t *testing.T) {
	_, err := New(types.Options{})
	if err == nil {
		t.Error("expected error without OPENAI_API_KEY")
	}
}

func TestNewAnthropicRequiresKey(t *testing.T) {
	_, err := New(types.Options{Provider: "anthropic"})
	if err == nil {
		t.Error("expected error without ANTHROPIC_API_KEY")
	}
}

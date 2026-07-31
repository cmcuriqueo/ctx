package llm

import (
	"fmt"
	"os"

	"github.com/matias/ctx/internal/anthropic"
	"github.com/matias/ctx/internal/llm/types"
	"github.com/matias/ctx/internal/openai"
)

// New creates a provider from options, preferring explicit values over environment variables.
func New(opts types.Options) (types.Provider, error) {
	provider := opts.Provider
	if provider == "" {
		provider = os.Getenv("CTX_LLM_PROVIDER")
	}
	if provider == "" {
		provider = "openai"
	}

	model := opts.Model
	if model == "" {
		model = os.Getenv("CTX_LLM_MODEL")
	}

	switch provider {
	case "openai":
		apiKey := firstNonEmpty(opts.APIKey, os.Getenv("OPENAI_API_KEY"))
		if apiKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY not set")
		}
		if model == "" {
			model = "gpt-4o-mini"
		}
		baseURL := firstNonEmpty(opts.BaseURL, os.Getenv("OPENAI_BASE_URL"))
		return openai.New(apiKey, model, baseURL), nil
	case "anthropic":
		apiKey := firstNonEmpty(opts.APIKey, os.Getenv("ANTHROPIC_API_KEY"))
		if apiKey == "" {
			return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
		}
		if model == "" {
			model = "claude-haiku-4-5"
		}
		return anthropic.New(apiKey, model), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

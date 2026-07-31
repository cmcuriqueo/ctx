package anthropic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/matias/ctx/internal/llm/types"
)

// Provider implements types.Provider using Anthropic's Messages API.
type Provider struct {
	client *anthropic.Client
	model  string
}

// New creates an Anthropic provider.
func New(apiKey, model string) *Provider {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &Provider{client: &client, model: model}
}

// Name returns the provider identifier.
func (p *Provider) Name() string {
	return "anthropic"
}

const systemPrompt = `You are a software engineering assistant. Given a project summary and a task, select the most relevant files to include in an AI context window.
Respond ONLY with a JSON object in this exact format:
{"explanation": "brief reason for the selection", "paths": ["relative/path/one.go", "relative/path/two.js"]}
Do not include markdown, code fences, or any extra text.`

// SelectFiles asks the model to pick relevant files.
func (p *Provider) SelectFiles(task string, summary *types.ProjectSummary) (*types.Selection, error) {
	summaryJSON, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return nil, err
	}
	user := fmt.Sprintf("Task: %s\n\nProject summary:\n%s", task, string(summaryJSON))

	resp, err := p.client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     p.model,
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(user)),
		},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Content) == 0 {
		return nil, fmt.Errorf("empty response from Anthropic")
	}

	content := ""
	for _, block := range resp.Content {
		text := block.AsText()
		content += text.Text
	}

	var sel types.Selection
	if err := json.Unmarshal([]byte(content), &sel); err != nil {
		return nil, fmt.Errorf("failed to parse model response: %w\ncontent: %s", err, content)
	}
	return &sel, nil
}

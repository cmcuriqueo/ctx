package openai

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"

	"github.com/matias/ctx/internal/llm/types"
)

// Provider implements llm.Provider using OpenAI-compatible APIs.
type Provider struct {
	client *openai.Client
	model  string
}

// New creates an OpenAI provider.
func New(apiKey, model, baseURL string) *Provider {
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}
	return &Provider{
		client: openai.NewClientWithConfig(config),
		model:  model,
	}
}

// Name returns the provider identifier.
func (p *Provider) Name() string {
	return "openai"
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

	resp, err := p.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: p.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: user},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("empty response from OpenAI")
	}

	content := resp.Choices[0].Message.Content
	var sel types.Selection
	if err := json.Unmarshal([]byte(content), &sel); err != nil {
		return nil, fmt.Errorf("failed to parse model response: %w\ncontent: %s", err, content)
	}
	return &sel, nil
}

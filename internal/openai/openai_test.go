package openai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matias/ctx/internal/llm/types"
)

func TestProviderSelectFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": `{"explanation": "test", "paths": ["main.go"]}`,
					},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := New("fake-key", "gpt-4o-mini", server.URL)
	sel, err := p.SelectFiles("task", &types.ProjectSummary{})
	if err != nil {
		t.Fatal(err)
	}
	if len(sel.Paths) != 1 || sel.Paths[0] != "main.go" {
		t.Errorf("unexpected selection: %v", sel.Paths)
	}
}

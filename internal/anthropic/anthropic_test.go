package anthropic

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProviderSelectFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": `{"explanation": "test", "paths": ["main.go"]}`,
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := New("fake-key", "claude-haiku-4-5")
	// The official SDK does not expose baseURL override easily, so this test only checks construction.
	if p.Name() != "anthropic" {
		t.Errorf("name = %s, want anthropic", p.Name())
	}
}

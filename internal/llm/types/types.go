package types

// Provider abstracts a language model backend capable of selecting files for a task.
type Provider interface {
	// SelectFiles asks the model which files are relevant for the task.
	SelectFiles(task string, summary *ProjectSummary) (*Selection, error)
	// Name returns the provider identifier.
	Name() string
}

// ProjectSummary is a compact description of a codebase sent to the LLM.
type ProjectSummary struct {
	Root        string        `json:"root"`
	Task        string        `json:"task"`
	TotalFiles  int           `json:"total_files"`
	TotalTokens int           `json:"total_tokens"`
	Files       []FileSummary `json:"files"`
}

// FileSummary describes a single file without its full content.
type FileSummary struct {
	Path     string   `json:"path"`
	Language string   `json:"language"`
	Lines    int      `json:"lines"`
	Tokens   int      `json:"tokens"`
	Package  string   `json:"package,omitempty"`
	Imports  []string `json:"imports,omitempty"`
	Exports  []string `json:"exports,omitempty"`
}

// Selection is the model's response.
type Selection struct {
	Explanation string   `json:"explanation"`
	Paths       []string `json:"paths"`
}

// Options configures provider creation.
type Options struct {
	Provider string
	Model    string
	APIKey   string
	BaseURL  string
}

package parser

// Metadata contains symbols extracted from a source file.
type Metadata struct {
	Package string   `json:"package"`
	Imports []string `json:"imports"`
	Exports []string `json:"exports"`
}

// Parser extracts metadata from source code.
type Parser interface {
	Parse(content []byte) (*Metadata, error)
}

// ForLanguage returns a parser for the given ctx language name, or nil.
func ForLanguage(language string) Parser {
	switch language {
	case "go":
		return NewGoParser()
	case "javascript":
		return NewJavaScriptParser()
	case "typescript", "tsx", "jsx":
		return NewTypeScriptParser()
	case "python":
		return NewPythonParser()
	}
	return nil
}

package parser

import tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"

const pythonQuery = `
(import_statement (dotted_name) @import)
(import_from_statement module_name: (dotted_name) @import)
(import_from_statement module_name: (relative_import) @import)
(function_definition name: (identifier) @export)
(class_definition name: (identifier) @export)
`

// NewPythonParser returns a Python parser.
func NewPythonParser() Parser {
	return &treeSitterParser{
		language:  tree_sitter_python.Language(),
		queryText: pythonQuery,
	}
}

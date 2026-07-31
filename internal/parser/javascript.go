package parser

import tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"

const javascriptQuery = `
(import_statement source: (string) @import)
(export_statement (function_declaration name: (identifier) @export))
(export_statement (class_declaration name: (identifier) @export))
(export_statement (lexical_declaration (variable_declarator name: (identifier) @export)))
`

// NewJavaScriptParser returns a JavaScript parser.
func NewJavaScriptParser() Parser {
	return &treeSitterParser{
		language:  tree_sitter_javascript.Language(),
		queryText: javascriptQuery,
	}
}

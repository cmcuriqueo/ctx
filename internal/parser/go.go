package parser

import tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"

const goQuery = `
(package_clause (package_identifier) @package)
(import_spec path: (interpreted_string_literal) @import)
(import_spec path: (raw_string_literal) @import)
(function_declaration name: (identifier) @export (#match? @export "^[A-Z]"))
(method_declaration name: (field_identifier) @export (#match? @export "^[A-Z]"))
(type_declaration (type_spec name: (type_identifier) @export (#match? @export "^[A-Z]")))
`

// NewGoParser returns a Go parser.
func NewGoParser() Parser {
	return &treeSitterParser{
		language:  tree_sitter_go.Language(),
		queryText: goQuery,
	}
}

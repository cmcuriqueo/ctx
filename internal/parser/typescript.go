package parser

import tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"

const typescriptQuery = `
(import_statement source: (string) @import)
(import_statement source: (template_string) @import)
(export_statement (function_declaration name: (identifier) @export))
(export_statement (class_declaration name: (type_identifier) @export))
(export_statement (lexical_declaration (variable_declarator name: (identifier) @export)))
(export_statement (type_alias_declaration name: (type_identifier) @export))
`

// NewTypeScriptParser returns a TypeScript/TSX parser.
func NewTypeScriptParser() Parser {
	return &treeSitterParser{
		language:  tree_sitter_typescript.LanguageTypescript(),
		queryText: typescriptQuery,
	}
}

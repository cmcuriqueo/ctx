package parser

import (
	"strings"
	"unsafe"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// treeSitterParser is a generic parser using tree-sitter queries.
type treeSitterParser struct {
	language  unsafe.Pointer
	queryText string
}

func (p *treeSitterParser) Parse(content []byte) (*Metadata, error) {
	parser := tree_sitter.NewParser()
	defer parser.Close()

	lang := tree_sitter.NewLanguage(p.language)
	parser.SetLanguage(lang)

	tree := parser.Parse(content, nil)
	defer tree.Close()

	query, err := tree_sitter.NewQuery(lang, p.queryText)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()

	meta := &Metadata{}
	seen := make(map[string]struct{})
	add := func(slice *[]string, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		*slice = append(*slice, value)
	}

	matches := cursor.Matches(query, tree.RootNode(), content)
	for match := matches.Next(); match != nil; match = matches.Next() {
		for _, capture := range match.Captures {
			name := query.CaptureNames()[capture.Index]
			text := strings.Trim(capture.Node.Utf8Text(content), "\"'`")
			switch name {
			case "package":
				meta.Package = text
			case "import":
				add(&meta.Imports, text)
			case "export":
				add(&meta.Exports, text)
			}
		}
	}
	return meta, nil
}

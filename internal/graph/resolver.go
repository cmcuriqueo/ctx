package graph

import (
	"path"
	"strings"

	"github.com/matias/ctx/pkg/models"
)

// SimpleResolver resolves import paths to files within the manifest.
// It handles relative imports and basic name matching for non-relative imports.
type SimpleResolver struct {
	files map[string]struct{}
}

// NewSimpleResolver creates a resolver from a manifest.
func NewSimpleResolver(manifest *models.Manifest) *SimpleResolver {
	files := make(map[string]struct{}, len(manifest.Files))
	for _, f := range manifest.Files {
		files[f.Path] = struct{}{}
	}
	return &SimpleResolver{files: files}
}

func (r *SimpleResolver) exists(p string) bool {
	_, ok := r.files[p]
	return ok
}

// Resolve maps an import string to a file path in the manifest.
func (r *SimpleResolver) Resolve(importerPath, importPath string) string {
	importPath = strings.Trim(importPath, `"'`)
	if importPath == "" {
		return ""
	}

	// Relative imports: ./foo or ../foo
	if strings.HasPrefix(importPath, ".") {
		dir := path.Dir(importerPath)
		candidate := path.Join(dir, importPath)
		candidate = path.Clean(candidate)
		if r.exists(candidate) {
			return candidate
		}
		for _, ext := range []string{".go", ".js", ".jsx", ".ts", ".tsx", ".py", ".rs", ".java"} {
			if r.exists(candidate + ext) {
				return candidate + ext
			}
		}
		// index files
		for _, idx := range []string{"index.js", "index.ts", "index.tsx", "index.jsx", "__init__.py"} {
			if r.exists(path.Join(candidate, idx)) {
				return path.Join(candidate, idx)
			}
		}
		return ""
	}

	// Python relative: from . import foo
	if strings.HasPrefix(importPath, ".") {
		return ""
	}

	// Non-relative imports: try matching by final segment or package name.
	segments := strings.Split(importPath, "/")
	name := segments[len(segments)-1]
	for file := range r.files {
		base := path.Base(file)
		ext := path.Ext(base)
		baseNoExt := strings.TrimSuffix(base, ext)
		if baseNoExt == name {
			return file
		}
	}
	return ""
}

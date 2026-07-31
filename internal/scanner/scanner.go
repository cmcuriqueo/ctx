package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/matias/ctx/internal/ignore"
	"github.com/matias/ctx/internal/parser"
	"github.com/matias/ctx/pkg/models"
)

var languageByExt = map[string]string{
	".go": "go", ".mod": "go", ".sum": "go",
	".work": "go",
	".js": "javascript", ".mjs": "javascript", ".cjs": "javascript",
	".ts": "typescript", ".tsx": "tsx",
	".jsx": "jsx",
	".py": "python", ".pyw": "python",
	".rs": "rust", ".toml": "toml",
	".java": "java",
	".c": "c", ".cpp": "cpp", ".cc": "cpp", ".h": "c", ".hpp": "cpp",
	".md": "markdown", ".markdown": "markdown",
	".json": "json", ".jsonc": "jsonc",
	".yaml": "yaml", ".yml": "yaml",
	".html": "html", ".htm": "html",
	".css": "css", ".scss": "scss", ".sass": "sass", ".less": "less",
	".sh": "shell", ".bash": "shell", ".zsh": "shell",
	".ps1": "powershell", ".psm1": "powershell",
	".rb": "ruby", ".php": "php", ".swift": "swift", ".kt": "kotlin",
	".cs": "csharp", ".fs": "fsharp", ".vb": "vbnet",
	".vue": "vue", ".svelte": "svelte", ".astro": "astro",
	".dockerfile": "dockerfile", ".sql": "sql", ".graphql": "graphql",
}

// Scanner walks a directory tree and collects file metadata.
type Scanner struct {
	ignorer *ignore.Engine
}

// New creates a scanner with the given ignore engine.
func New(ignorer *ignore.Engine) *Scanner {
	return &Scanner{ignorer: ignorer}
}

// Scan walks root and returns a manifest of files.
func (s *Scanner) Scan(root string) (*models.Manifest, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	manifest := &models.Manifest{
		Root:      absRoot,
		ScannedAt: time.Now(),
	}
	err = filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(absRoot, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if info.IsDir() {
			if s.ignorer.Match(rel, true) {
				return filepath.SkipDir
			}
			return nil
		}
		if s.ignorer.Match(rel, false) {
			return nil
		}
		fi, err := s.fileInfo(absRoot, rel, info)
		if err != nil {
			return err
		}
		manifest.Files = append(manifest.Files, fi)
		return nil
	})
	return manifest, err
}

func (s *Scanner) fileInfo(root, rel string, info os.FileInfo) (models.FileInfo, error) {
	fi := models.FileInfo{
		Path:       filepath.ToSlash(rel),
		Size:       info.Size(),
		ModifiedAt: info.ModTime(),
	}
	ext := strings.ToLower(filepath.Ext(rel))
	if lang, ok := languageByExt[ext]; ok {
		fi.Language = lang
	} else {
		name := strings.ToLower(filepath.Base(rel))
		if lang, ok := languageByName[name]; ok {
			fi.Language = lang
		}
	}
	fullPath := filepath.Join(root, rel)
	if isBinary(fullPath) {
		fi.IsBinary = true
		return fi, nil
	}
	content, lines, err := readFileLines(fullPath)
	if err != nil {
		return fi, err
	}
	fi.Lines = lines
	h := sha256.Sum256(content)
	fi.SHA256 = hex.EncodeToString(h[:])

	if p := parser.ForLanguage(fi.Language); p != nil {
		meta, err := p.Parse(content)
		if err == nil {
			fi.Package = meta.Package
			fi.Imports = meta.Imports
			fi.Exports = meta.Exports
		}
	}
	return fi, nil
}

func readFileLines(path string) ([]byte, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, 0, err
	}
	lines := strings.Count(string(data), "\n")
	return data, lines, nil
}

var languageByName = map[string]string{
	"dockerfile": "dockerfile",
	"makefile":   "makefile",
	".gitignore": "gitignore",
	".gitattributes": "gitattributes",
	".editorconfig":  "editorconfig",
}

// isBinary checks whether a file contains null bytes in its first 512 bytes.
func isBinary(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return true
	}
	defer f.Close()
	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	buf = buf[:n]
	for _, b := range buf {
		if b == 0 {
			return true
		}
	}
	return false
}

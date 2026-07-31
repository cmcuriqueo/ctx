package models

import "time"

// FileInfo holds basic metadata for a single file.
type FileInfo struct {
	Path       string    `json:"path"`
	Language   string    `json:"language"`
	Size       int64     `json:"size"`
	Lines      int       `json:"lines"`
	SHA256     string    `json:"sha256"`
	ModifiedAt time.Time `json:"modified_at"`
	IsBinary   bool      `json:"is_binary"`
}

// Manifest is the result of scanning a repository.
type Manifest struct {
	Root      string     `json:"root"`
	Files     []FileInfo `json:"files"`
	ScannedAt time.Time  `json:"scanned_at"`
}

// TokenEntry stores cached token count for a given file hash.
type TokenEntry struct {
	Hash   string `json:"hash"`
	Tokens int    `json:"tokens"`
}

// RankedFile extends FileInfo with scoring and token data.
type RankedFile struct {
	FileInfo
	Score      int     `json:"score"`
	Tokens     int     `json:"tokens"`
	Efficiency float64 `json:"efficiency"`
}

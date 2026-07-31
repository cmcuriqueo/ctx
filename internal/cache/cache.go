package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Cache provides simple JSON persistence under a directory.
type Cache struct {
	Dir string
}

// New creates a cache using the given directory.
func New(dir string) *Cache {
	return &Cache{Dir: dir}
}

func (c *Cache) path(name string) string {
	return filepath.Join(c.Dir, name)
}

// Load reads and unmarshals a JSON cache file.
func (c *Cache) Load(name string, v interface{}) error {
	data, err := os.ReadFile(c.path(name))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// Save marshals and writes a JSON cache file.
func (c *Cache) Save(name string, v interface{}) error {
	if err := os.MkdirAll(c.Dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path(name), data, 0644)
}

// Exists reports whether a cache file exists.
func (c *Cache) Exists(name string) bool {
	_, err := os.Stat(c.path(name))
	return err == nil
}

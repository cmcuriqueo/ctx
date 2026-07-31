package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds configurable weights and thresholds.
type Config struct {
	Rank  RankConfig  `toml:"rank"`
	Graph GraphConfig `toml:"graph"`
}

// RankConfig holds scoring weights.
type RankConfig struct {
	Entrypoint    int `toml:"entrypoint"`
	Readme        int `toml:"readme"`
	ConfigFile    int `toml:"config"`
	Test          int `toml:"test"`
	Generated     int `toml:"generated"`
	Vendor        int `toml:"vendor"`
	ImportedBonus int `toml:"imported_bonus"`
	LargeIsolated int `toml:"large_isolated"`
}

// GraphConfig holds graph traversal settings.
type GraphConfig struct {
	MaxDepth int `toml:"max_depth"`
}

// Default returns the built-in default configuration.
func Default() *Config {
	return &Config{
		Rank: RankConfig{
			Entrypoint:    30,
			Readme:        20,
			ConfigFile:    20,
			Test:          10,
			Generated:     -40,
			Vendor:        -30,
			ImportedBonus: 15,
			LargeIsolated: -20,
		},
		Graph: GraphConfig{
			MaxDepth: 5,
		},
	}
}

// Load reads ctx.toml from the project root, falling back to defaults.
func Load(root string) (*Config, error) {
	cfg := Default()
	path := filepath.Join(root, "ctx.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

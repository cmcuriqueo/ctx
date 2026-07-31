package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/matias/ctx/internal/builder"
	"github.com/matias/ctx/internal/cache"
	"github.com/matias/ctx/internal/ignore"
	"github.com/matias/ctx/internal/rank"
	"github.com/matias/ctx/internal/scanner"
	"github.com/matias/ctx/internal/tokens"
)

var (
	budget   int
	output   string
	cacheDir string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ctx",
		Short: "Context builder CLI for AI",
		Long:  "ctx scans a codebase and builds a context.md optimized for AI prompts.",
	}

	rootCmd.PersistentFlags().IntVarP(&budget, "budget", "b", 4000, "token budget")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "context.md", "output file")
	rootCmd.PersistentFlags().StringVar(&cacheDir, "cache-dir", ".cache/ctx", "cache directory")

	rootCmd.AddCommand(scanCmd())
	rootCmd.AddCommand(tokensCmd())
	rootCmd.AddCommand(bundleCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func scanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan repository and store manifest",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			ign, err := ignore.NewEngine(path)
			if err != nil {
				return err
			}
			s := scanner.New(ign)
			manifest, err := s.Scan(path)
			if err != nil {
				return err
			}
			c := cache.New(cacheDir)
			if err := c.Save("manifest.json", manifest); err != nil {
				return err
			}
			fmt.Printf("Scanned %d files\n", len(manifest.Files))
			return nil
		},
	}
}

func tokensCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tokens [path]",
		Short: "Estimate tokens for scanned files",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			ign, err := ignore.NewEngine(path)
			if err != nil {
				return err
			}
			s := scanner.New(ign)
			manifest, err := s.Scan(path)
			if err != nil {
				return err
			}
			est := tokens.NewHeuristicEstimator()
			total := 0
			for _, f := range manifest.Files {
				if f.IsBinary {
					continue
				}
				content, err := os.ReadFile(filepath.Join(path, f.Path))
				if err != nil {
					continue
				}
				t := est.Estimate(string(content))
				total += t
				fmt.Printf("%s: %d tokens\n", f.Path, t)
			}
			fmt.Printf("Total: %d tokens\n", total)
			return nil
		},
	}
}

func bundleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bundle [path]",
		Short: "Generate context.md",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			ign, err := ignore.NewEngine(path)
			if err != nil {
				return err
			}
			s := scanner.New(ign)
			manifest, err := s.Scan(path)
			if err != nil {
				return err
			}
			est := tokens.NewHeuristicEstimator()
			scorer := rank.NewScorer()
			b := builder.New(est, scorer)
			return b.Build(manifest, budget, output)
		},
	}
}

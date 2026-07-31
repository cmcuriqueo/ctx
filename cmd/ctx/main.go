package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/matias/ctx/internal/builder"
	"github.com/matias/ctx/internal/cache"
	"github.com/matias/ctx/internal/config"
	"github.com/matias/ctx/internal/graph"
	"github.com/matias/ctx/internal/ignore"
	"github.com/matias/ctx/internal/rank"
	"github.com/matias/ctx/internal/scanner"
	"github.com/matias/ctx/internal/tokens"
	"github.com/matias/ctx/pkg/models"
)

var (
	budget   int
	output   string
	cacheDir string
	depth    int
	format   string
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
	rootCmd.AddCommand(graphCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func resolvePath(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "."
}

func scanManifest(path string) (*models.Manifest, *graph.Graph, *config.Config, error) {
	ign, err := ignore.NewEngine(path)
	if err != nil {
		return nil, nil, nil, err
	}
	s := scanner.New(ign)
	manifest, err := s.Scan(path)
	if err != nil {
		return nil, nil, nil, err
	}
	cfg, err := config.Load(path)
	if err != nil {
		return nil, nil, nil, err
	}
	resolver := graph.NewSimpleResolver(manifest)
	g := graph.New(manifest, resolver)
	return manifest, g, cfg, nil
}

func scanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan repository and store manifest",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manifest, g, _, err := scanManifest(resolvePath(args))
			if err != nil {
				return err
			}
			c := cache.New(cacheDir)
			if err := c.Save("manifest.json", manifest); err != nil {
				return err
			}
			if err := c.Save("graph.json", g); err != nil {
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
			manifest, _, _, err := scanManifest(resolvePath(args))
			if err != nil {
				return err
			}
			est := tokens.NewHeuristicEstimator()
			total := 0
			for _, f := range manifest.Files {
				if f.IsBinary {
					continue
				}
				content, err := os.ReadFile(filepath.Join(manifest.Root, f.Path))
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
			manifest, g, cfg, err := scanManifest(resolvePath(args))
			if err != nil {
				return err
			}
			est := tokens.NewHeuristicEstimator()
			scorer := rank.NewScorer(cfg)
			b := builder.New(est, scorer, g, cfg)
			return b.Build(manifest, budget, output)
		},
	}
}

func graphCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "graph [path]",
		Short: "Show dependency graph",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, g, _, err := scanManifest(resolvePath(args))
			if err != nil {
				return err
			}

			switch format {
			case "json":
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(g)
			case "dot":
				fmt.Println("digraph ctx {")
				for from, tos := range g.Edges {
					for _, to := range tos {
						fmt.Printf("  \"%s\" -> \"%s\";\n", from, to)
					}
				}
				fmt.Println("}")
			case "text":
				fmt.Println("# Dependency Graph")
				fmt.Println()
				for _, node := range g.Nodes {
					fmt.Printf("%s\n", node)
					for _, to := range g.Edges[node] {
						fmt.Printf("  -> %s\n", to)
					}
				}
				fmt.Println()
				cycles := g.DetectCycles()
				if len(cycles) > 0 {
					fmt.Printf("Cycles detected: %d\n", len(cycles))
					for _, c := range cycles {
						fmt.Printf("  %v\n", c)
					}
				} else {
					fmt.Println("No cycles detected")
				}
				fmt.Printf("\nConnected components: %d\n", len(g.ConnectedComponents()))
			default:
				return fmt.Errorf("unknown format: %s", format)
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&depth, "depth", 0, "BFS depth from entrypoints (unused in text mode)")
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text, dot, json")
	_ = depth
	return cmd
}

package graph

import (
	"fmt"
	"sort"

	"github.com/matias/ctx/pkg/models"
)

// Graph models files as nodes and imports as directed edges.
type Graph struct {
	Nodes []string            `json:"nodes"`
	Edges map[string][]string `json:"edges"` // file -> files it imports
}

// New builds a dependency graph from a manifest.
// It creates edges from a file to the files it imports, resolved relative to the importer.
func New(manifest *models.Manifest, resolver ImportResolver) *Graph {
	g := &Graph{
		Nodes: make([]string, 0, len(manifest.Files)),
		Edges: make(map[string][]string),
	}
	fileSet := make(map[string]struct{}, len(manifest.Files))
	for _, f := range manifest.Files {
		path := f.Path
		fileSet[path] = struct{}{}
		g.Nodes = append(g.Nodes, path)
		g.Edges[path] = []string{}
	}

	for _, f := range manifest.Files {
		for _, imp := range f.Imports {
			resolved := resolver.Resolve(f.Path, imp)
			if resolved == "" {
				continue
			}
			if _, ok := fileSet[resolved]; !ok {
				continue
			}
			g.Edges[f.Path] = append(g.Edges[f.Path], resolved)
		}
	}
	return g
}

// ImportResolver resolves an import path relative to the importing file.
type ImportResolver interface {
	Resolve(importerPath, importPath string) string
}

// BFS returns all nodes reachable from start up to the given depth.
// depth <= 0 means unlimited.
func (g *Graph) BFS(start string, depth int) []string {
	if depth == 0 {
		depth = -1
	}
	visited := make(map[string]struct{})
	queue := []struct {
		node  string
		level int
	}{{start, 0}}
	var result []string

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if _, ok := visited[cur.node]; ok {
			continue
		}
		visited[cur.node] = struct{}{}
		result = append(result, cur.node)
		if depth > 0 && cur.level >= depth {
			continue
		}
		for _, next := range g.Edges[cur.node] {
			if _, ok := visited[next]; !ok {
				queue = append(queue, struct {
					node  string
					level int
				}{next, cur.level + 1})
			}
		}
	}
	return result
}

// ReverseEdges returns the inverted graph (file -> files that import it).
func (g *Graph) ReverseEdges() map[string][]string {
	rev := make(map[string][]string, len(g.Nodes))
	for _, n := range g.Nodes {
		rev[n] = []string{}
	}
	for from, tos := range g.Edges {
		for _, to := range tos {
			rev[to] = append(rev[to], from)
		}
	}
	return rev
}

// ImportedByCount returns how many files import the given file.
func (g *Graph) ImportedByCount(path string) int {
	rev := g.ReverseEdges()
	return len(rev[path])
}

// DetectCycles finds all elementary cycles using DFS-based search.
func (g *Graph) DetectCycles() [][]string {
	var cycles [][]string
	for _, start := range g.Nodes {
		visited := map[string]int{}
		path := []string{}
		g.dfsCycles(start, start, visited, path, &cycles)
	}
	return cycles
}

func (g *Graph) dfsCycles(start, node string, visited map[string]int, path []string, cycles *[][]string) {
	if idx, ok := visited[node]; ok {
		if node == start && idx == 0 && len(path) > 0 {
			cycle := make([]string, len(path))
			copy(cycle, path)
			*cycles = append(*cycles, cycle)
		}
		return
	}
	visited[node] = len(path)
	path = append(path, node)
	for _, next := range g.Edges[node] {
		g.dfsCycles(start, next, visited, path, cycles)
	}
}

// ConnectedComponents returns weakly connected components.
func (g *Graph) ConnectedComponents() [][]string {
	visited := make(map[string]struct{}, len(g.Nodes))
	rev := g.ReverseEdges()
	var components [][]string

	for _, node := range g.Nodes {
		if _, ok := visited[node]; ok {
			continue
		}
		component := []string{}
		queue := []string{node}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			if _, ok := visited[cur]; ok {
				continue
			}
			visited[cur] = struct{}{}
			component = append(component, cur)
			for _, next := range g.Edges[cur] {
				if _, ok := visited[next]; !ok {
					queue = append(queue, next)
				}
			}
			for _, prev := range rev[cur] {
				if _, ok := visited[prev]; !ok {
					queue = append(queue, prev)
				}
			}
		}
		sort.Strings(component)
		components = append(components, component)
	}
	return components
}

// TopologicalSort returns nodes in topological order if the graph is a DAG.
func (g *Graph) TopologicalSort() ([]string, error) {
	inDegree := make(map[string]int, len(g.Nodes))
	for _, n := range g.Nodes {
		inDegree[n] = 0
	}
	for _, tos := range g.Edges {
		for _, to := range tos {
			inDegree[to]++
		}
	}

	var queue []string
	for _, n := range g.Nodes {
		if inDegree[n] == 0 {
			queue = append(queue, n)
		}
	}
	sort.Strings(queue)

	var result []string
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		result = append(result, cur)
		for _, next := range g.Edges[cur] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
				SortQueue(queue)
			}
		}
	}

	if len(result) != len(g.Nodes) {
		return nil, fmt.Errorf("graph contains cycles")
	}
	return result, nil
}

// SortQueue keeps the queue sorted for deterministic output.
func SortQueue(queue []string) {
	sort.Strings(queue)
}

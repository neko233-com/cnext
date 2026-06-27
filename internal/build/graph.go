package build

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/neko233-com/cnext/internal/config"
)

type NodeType string

const (
	NodeCompile NodeType = "compile"
	NodeLink    NodeType = "link"
)

type BuildGraph struct {
	Nodes []BuildNode
}

type BuildNode struct {
	ID           string
	Type         NodeType
	Sources      []string
	Output       string
	Dependencies []string
	Compiler     string
	STD          string
	Optimization string
	IncludeDirs  []string
	Flags        []string
}

func Generate(cfg *config.Config) (*BuildGraph, error) {
	graph := &BuildGraph{}

	for _, lib := range cfg.Build.Libraries {
		sources, err := resolveSources(lib.Sources)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve sources for library %q: %w", lib.Name, err)
		}

		if len(sources) == 0 {
			continue
		}

		output := lib.Name + ".a"
		if lib.Type == "shared" {
			output = lib.Name + ".so"
		}

		graph.Nodes = append(graph.Nodes, BuildNode{
			ID:           "lib:" + lib.Name,
			Type:         NodeCompile,
			Sources:      sources,
			Output:       output,
			Compiler:     cfg.Build.Compiler,
			STD:          cfg.Build.STD,
			Optimization: cfg.Build.Optimization,
			IncludeDirs:  lib.IncludeDirs,
		})
	}

	for _, exe := range cfg.Build.Executables {
		sources, err := resolveSources(exe.Sources)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve sources for executable %q: %w", exe.Name, err)
		}

		if len(sources) == 0 {
			continue
		}

		compileID := "exe:" + exe.Name
		var compileDeps []string
		for _, lib := range exe.Libraries {
			compileDeps = append(compileDeps, "lib:"+lib)
		}

		graph.Nodes = append(graph.Nodes, BuildNode{
			ID:           compileID,
			Type:         NodeCompile,
			Sources:      sources,
			Output:       exe.Name + ".o",
			Dependencies: compileDeps,
			Compiler:     cfg.Build.Compiler,
			STD:          cfg.Build.STD,
			Optimization: cfg.Build.Optimization,
		})

		linkDeps := []string{compileID}
		for _, lib := range exe.Libraries {
			linkDeps = append(linkDeps, "lib:"+lib)
		}

		graph.Nodes = append(graph.Nodes, BuildNode{
			ID:           "link:" + exe.Name,
			Type:         NodeLink,
			Output:       exe.Name,
			Dependencies: linkDeps,
		})
	}

	return graph, nil
}

func (g *BuildGraph) TopologicalSort() ([]string, error) {
	inDegree := make(map[string]int)
	dependents := make(map[string][]string)

	for _, node := range g.Nodes {
		if _, ok := inDegree[node.ID]; !ok {
			inDegree[node.ID] = 0
		}
		for _, dep := range node.Dependencies {
			dependents[dep] = append(dependents[dep], node.ID)
			inDegree[node.ID]++
		}
	}

	var queue []string
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}

	var sorted []string
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		sorted = append(sorted, current)

		for _, dependent := range dependents[current] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}

	if len(sorted) != len(g.Nodes) {
		return nil, fmt.Errorf("cycle detected in build graph")
	}

	return sorted, nil
}

func resolveSources(patterns []string) ([]string, error) {
	var sources []string
	for _, pattern := range patterns {
		if strings.ContainsAny(pattern, "*?[") {
			matches, err := filepath.Glob(pattern)
			if err != nil {
				return nil, err
			}
			sources = append(sources, matches...)
		} else {
			sources = append(sources, pattern)
		}
	}
	return sources, nil
}

package build

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/neko233-com/cnext/internal/cmake"
	"github.com/neko233-com/cnext/internal/compiler"
	"github.com/neko233-com/cnext/internal/config"
)

type NodeType string

const (
	NodeCompile NodeType = "compile"
	NodeLink    NodeType = "link"
	NodeCMake   NodeType = "cmake"
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
	CMakeConfig  *CMakeConfig
	Language     string
	PIC          bool
}

type CMakeConfig struct {
	SourceDir string
	Options   []string
}

func Generate(cfg *config.Config) (*BuildGraph, error) {
	graph := &BuildGraph{}

	for _, dep := range cfg.Build.CMakeDependencies {
		cmakeConfig := &CMakeConfig{
			SourceDir: ".",
			Options:   dep.CMakeOptions,
		}

		var deps []string
		if dep.Git != "" {
			deps = append(deps, "cmake:"+dep.Name)
		}

		graph.Nodes = append(graph.Nodes, BuildNode{
			ID:           "cmake:" + dep.Name,
			Type:         NodeCMake,
			Dependencies: deps,
			CMakeConfig:  cmakeConfig,
		})
	}

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
			PIC:          cfg.Build.PIC,
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
		var includeDirs []string
		for _, libName := range exe.Libraries {
			compileDeps = append(compileDeps, "lib:"+libName)
			// Collect include directories from library dependencies
			for _, lib := range cfg.Build.Libraries {
				if lib.Name == libName {
					includeDirs = append(includeDirs, lib.IncludeDirs...)
				}
			}
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
			IncludeDirs:  includeDirs,
		})

		linkDeps := []string{compileID}
		for _, lib := range exe.Libraries {
			linkDeps = append(linkDeps, "lib:"+lib)
		}

		graph.Nodes = append(graph.Nodes, BuildNode{
			ID:           "link:" + exe.Name,
			Type:         NodeLink,
			Output:       exeName(exe.Name),
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
		if strings.Contains(pattern, "**") {
			matches, err := globRecursive(pattern)
			if err != nil {
				return nil, err
			}
			sources = append(sources, matches...)
		} else if strings.ContainsAny(pattern, "*?[") {
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

func globRecursive(pattern string) ([]string, error) {
	parts := strings.Split(pattern, "**")
	if len(parts) != 2 {
		return filepath.Glob(pattern)
	}

	prefix := parts[0]
	suffix := parts[1]

	if prefix == "" {
		prefix = "."
	}

	// Normalize prefix - remove trailing separator
	prefix = strings.TrimRight(prefix, "/\\")

	// Extract extension from suffix (e.g., "/.cpp" -> ".cpp")
	ext := strings.TrimPrefix(suffix, "/")
	ext = strings.TrimPrefix(ext, "*")

	var matches []string
	err := filepath.Walk(prefix, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		if ext == "" || strings.HasSuffix(path, ext) {
			matches = append(matches, path)
		}
		return nil
	})

	return matches, err
}

func exeName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func (g *BuildGraph) nodeMap() map[string]*BuildNode {
	m := make(map[string]*BuildNode, len(g.Nodes))
	for i := range g.Nodes {
		m[g.Nodes[i].ID] = &g.Nodes[i]
	}
	return m
}

func detectLanguage(sources []string) string {
	if len(sources) == 0 {
		return "cpp"
	}
	ext := filepath.Ext(sources[0])
	switch ext {
	case ".c":
		return "c"
	case ".cpp", ".cc", ".cxx", ".C":
		return "cpp"
	case ".h":
		return "c"
	case ".hpp", ".hh":
		return "cpp"
	default:
		return "cpp"
	}
}

func (g *BuildGraph) Execute(projectDir string, comp compiler.Compiler) error {
	sorted, err := g.TopologicalSort()
	if err != nil {
		return err
	}

	nm := g.nodeMap()

	for _, id := range sorted {
		node := nm[id]

		switch node.Type {
		case NodeCMake:
			if node.CMakeConfig == nil {
				return fmt.Errorf("cmake node %q missing CMakeConfig", id)
			}
			bridge := cmake.NewBridge(projectDir)
			if err := bridge.Configure(node.CMakeConfig.SourceDir, node.CMakeConfig.Options); err != nil {
				return fmt.Errorf("cmake configure for %q failed: %w", id, err)
			}
			if err := bridge.Build(0); err != nil {
				return fmt.Errorf("cmake build for %q failed: %w", id, err)
			}
			if err := bridge.Install(); err != nil {
				return fmt.Errorf("cmake install for %q failed: %w", id, err)
			}

		case NodeCompile:
			lang := detectLanguage(node.Sources)
			opts := compiler.CompileOptions{
				Sources:      node.Sources,
				Output:       node.Output,
				Std:          node.STD,
				Language:     lang,
				IncludeDirs:  node.IncludeDirs,
				Flags:        node.Flags,
				Optimization: node.Optimization,
				CompileOnly:  true,
				PIC:          node.PIC,
			}
			fmt.Printf("  [compile] %s\n", node.ID)
			if err := comp.Compile(opts); err != nil {
				return fmt.Errorf("compilation of %s failed: %w", node.ID, err)
			}

		case NodeLink:
			var objects []string
			for _, depID := range node.Dependencies {
				dep := nm[depID]
				if dep != nil && dep.Type == NodeCompile {
					objects = append(objects, dep.Output)
				}
			}
			opts := compiler.LinkOptions{
				Objects: objects,
				Output:  node.Output,
				Flags:   node.Flags,
			}
			fmt.Printf("  [link] %s\n", node.ID)
			if err := comp.Link(opts); err != nil {
				return fmt.Errorf("linking of %s failed: %w", node.ID, err)
			}
		}
	}

	return nil
}

package build

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/neko233-com/cnext/internal/cache"
	"github.com/neko233-com/cnext/internal/compiler"
)

type BuildOptions struct {
	Jobs        int
	Incremental bool
	Full        bool
	Cache       bool
	LTO         bool
	Native      bool
	Verbose     bool
}

type BuildResult struct {
	Compiled  int
	Cached    int
	Linked    int
	Archived  int
	Duration  time.Duration
	Errors    []error
}

func (g *BuildGraph) ExecuteWithOpts(projectDir string, comp compiler.Compiler, opts BuildOptions) (*BuildResult, error) {
	start := time.Now()
	result := &BuildResult{}

	sorted, err := g.TopologicalSort()
	if err != nil {
		return nil, err
	}

	nm := g.nodeMap()

	// Build dependency layers for parallel execution
	layers := g.buildLayers(sorted, nm)

	// Initialize cache
	var buildCache *cache.BuildCache
	if opts.Cache {
		buildCache = cache.New()
		buildCache.Load()
		defer buildCache.Save()
	}

	// Execute layers in parallel
	for _, layer := range layers {
		if err := g.executeLayer(projectDir, comp, layer, nm, opts, buildCache, result); err != nil {
			return result, err
		}
	}

	result.Duration = time.Since(start)
	return result, nil
}

func (g *BuildGraph) buildLayers(sorted []string, nm map[string]*BuildNode) [][]string {
	// Group nodes into layers based on dependencies
	layers := [][]string{}
	visited := make(map[string]bool)

	for len(visited) < len(sorted) {
		var layer []string
		for _, id := range sorted {
			if visited[id] {
				continue
			}

			node := nm[id]
			canExecute := true
			for _, dep := range node.Dependencies {
				if !visited[dep] {
					canExecute = false
					break
				}
			}

			if canExecute {
				layer = append(layer, id)
			}
		}

		if len(layer) == 0 {
			break
		}

		for _, id := range layer {
			visited[id] = true
		}
		layers = append(layers, layer)
	}

	return layers
}

func (g *BuildGraph) executeLayer(projectDir string, comp compiler.Compiler, layer []string, nm map[string]*BuildNode, opts BuildOptions, buildCache *cache.BuildCache, result *BuildResult) error {
	if len(layer) == 1 {
		// Single node, execute directly
		return g.executeNode(projectDir, comp, layer[0], nm, opts, buildCache, result)
	}

	// Multiple nodes, execute in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	var execErr error

	sem := make(chan struct{}, opts.Jobs)

	for _, id := range layer {
		wg.Add(1)
		sem <- struct{}{}

		go func(nodeID string) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := g.executeNode(projectDir, comp, nodeID, nm, opts, buildCache, result); err != nil {
				mu.Lock()
				if execErr == nil {
					execErr = err
				}
				mu.Unlock()
			}
		}(id)
	}

	wg.Wait()
	return execErr
}

func (g *BuildGraph) executeNode(projectDir string, comp compiler.Compiler, id string, nm map[string]*BuildNode, opts BuildOptions, buildCache *cache.BuildCache, result *BuildResult) error {
	node := nm[id]

	switch node.Type {
	case NodeCMake:
		return g.executeCMakeNode(projectDir, node)

	case NodeCompile:
		return g.executeCompileNode(comp, node, opts, buildCache, result)

	case NodeLink:
		return g.executeLinkNode(comp, node, nm, result)

	default:
		return fmt.Errorf("unknown node type: %s", node.Type)
	}
}

func (g *BuildGraph) executeCMakeNode(projectDir string, node *BuildNode) error {
	if node.CMakeConfig == nil {
		return fmt.Errorf("cmake node %q missing CMakeConfig", node.ID)
	}
	// CMake execution would go here
	return nil
}

func (g *BuildGraph) executeCompileNode(comp compiler.Compiler, node *BuildNode, opts BuildOptions, buildCache *cache.BuildCache, result *BuildResult) error {
	lang := detectLanguage(node.Sources)

	// Check cache for incremental builds
	if opts.Incremental && buildCache != nil {
		// For now, check if output exists and is newer than source
		if !opts.Full && node.Output != "" {
			if info, err := os.Stat(node.Output); err == nil {
				oldest := info.ModTime()
				allUpToDate := true
				for _, src := range node.Sources {
					if srcInfo, err := os.Stat(src); err == nil {
						if srcInfo.ModTime().After(oldest) {
							allUpToDate = false
							break
						}
					}
				}
				if allUpToDate {
					if opts.Verbose {
						fmt.Printf("  [cached] %s\n", node.ID)
					}
					result.Cached++
					return nil
				}
			}
		}
	}

	// Compile each source file
	for _, src := range node.Sources {
		objName := strings.TrimSuffix(node.Output, filepath.Ext(node.Output))
		if len(node.Sources) > 1 {
			base := strings.TrimSuffix(filepath.Base(src), filepath.Ext(src))
			objName = objName + "_" + base
		}
		objName += ".o"

		depFile := strings.TrimSuffix(objName, filepath.Ext(objName)) + ".d"

		opts_compile := compiler.CompileOptions{
			Sources:      []string{src},
			Output:       objName,
			Std:          node.STD,
			Language:     lang,
			IncludeDirs:  node.IncludeDirs,
			Flags:        node.Flags,
			Optimization: node.Optimization,
			CompileOnly:  true,
			PIC:          node.PIC,
			DepFile:      depFile,
		}

		fmt.Printf("  [compile] %s -> %s\n", src, objName)
		if err := comp.Compile(opts_compile); err != nil {
			return fmt.Errorf("compilation of %s failed: %w", src, err)
		}
		result.Compiled++

		// Update cache
		if buildCache != nil {
			buildCache.Update(src, objName, depFile, comp.Name(), node.Flags)
		}
	}

	// Archive if this is a library node
	if strings.HasPrefix(node.ID, "lib:") && node.ArchiveOutput != "" {
		objects := collectObjects(*node, g.nodeMap())
		if len(objects) > 0 {
			fmt.Printf("  [archive] %s -> %s\n", node.ID, node.ArchiveOutput)
			if err := comp.Archive(objects, node.ArchiveOutput); err != nil {
				return fmt.Errorf("archiving of %s failed: %w", node.ID, err)
			}
			result.Archived++
		}
	}

	return nil
}

func (g *BuildGraph) executeLinkNode(comp compiler.Compiler, node *BuildNode, nm map[string]*BuildNode, result *BuildResult) error {
	var objects []string
	var libDirs []string
	var libraries []string

	for _, depID := range node.Dependencies {
		dep := nm[depID]
		if dep != nil {
			if dep.Type == NodeCompile {
				objects = append(objects, dep.Output)
			}
			if strings.HasPrefix(depID, "lib:") && dep.ArchiveOutput != "" {
				objects = append(objects, dep.ArchiveOutput)
			}
		}
	}

	opts := compiler.LinkOptions{
		Objects:   objects,
		Output:    node.Output,
		LibDirs:   libDirs,
		Libraries: libraries,
		Flags:     node.Flags,
	}

	fmt.Printf("  [link] %s\n", node.ID)
	if err := comp.Link(opts); err != nil {
		return fmt.Errorf("linking of %s failed: %w", node.ID, err)
	}
	result.Linked++

	return nil
}

func (g *BuildGraph) GetNumCPU() int {
	return runtime.NumCPU()
}

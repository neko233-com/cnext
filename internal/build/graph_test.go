package build

import (
	"runtime"
	"testing"

	"github.com/neko233-com/cnext/internal/config"
)

func TestGenerateBuildGraph(t *testing.T) {
	cfg := &config.Config{
		Package: config.Package{
			Name:    "test",
			Version: "1.0.0",
		},
		Build: config.Build{
			Compiler: "auto",
			STD:      "c++23",
			Libraries: []config.Library{
				{
					Name:    "mylib",
					Type:    "static",
					Sources: []string{"src/lib.cpp", "src/utils.cpp"},
				},
			},
			Executables: []config.Executable{
				{
					Name:      "myapp",
					Sources:   []string{"main.cpp"},
					Libraries: []string{"mylib"},
				},
			},
		},
	}

	graph, err := Generate(cfg)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if graph == nil {
		t.Fatal("Expected non-nil graph")
	}

	if len(graph.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(graph.Nodes))
	}

	expectedArchiveOutput := "mylib.a"
	if runtime.GOOS == "windows" {
		expectedArchiveOutput = "mylib.lib"
	}

	var hasLibNode, hasExeNode, hasLinkNode bool
	for _, node := range graph.Nodes {
		switch node.ID {
		case "lib:mylib":
			hasLibNode = true
			if node.Type != NodeCompile {
				t.Errorf("lib:mylib type = %v, want %v", node.Type, NodeCompile)
			}
			if len(node.Sources) != 2 {
				t.Errorf("lib:mylib sources len = %d, want 2", len(node.Sources))
			}
			if node.Output != "mylib.o" {
				t.Errorf("lib:mylib output = %q, want %q", node.Output, "mylib.o")
			}
			if node.ArchiveOutput != expectedArchiveOutput {
				t.Errorf("lib:mylib archiveOutput = %q, want %q", node.ArchiveOutput, expectedArchiveOutput)
			}
		case "exe:myapp":
			hasExeNode = true
			if node.Type != NodeCompile {
				t.Errorf("exe:myapp type = %v, want %v", node.Type, NodeCompile)
			}
		case "link:myapp":
			hasLinkNode = true
			if node.Type != NodeLink {
				t.Errorf("link:myapp type = %v, want %v", node.Type, NodeLink)
			}
		}
	}
	if !hasLibNode {
		t.Error("Expected compile node for lib:mylib")
	}
	if !hasExeNode {
		t.Error("Expected compile node for exe:myapp")
	}
	if !hasLinkNode {
		t.Error("Expected link node for link:myapp")
	}
}

func TestGenerateEmptyConfig(t *testing.T) {
	cfg := &config.Config{
		Package: config.Package{
			Name:    "empty",
			Version: "0.0.1",
		},
		Build: config.Build{
			Compiler:    "auto",
			STD:         "c++20",
			Libraries:   []config.Library{},
			Executables: []config.Executable{},
		},
	}

	graph, err := Generate(cfg)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if graph == nil {
		t.Fatal("Expected non-nil graph")
	}

	if len(graph.Nodes) != 0 {
		t.Errorf("Expected 0 nodes for empty config, got %d", len(graph.Nodes))
	}
}

func TestTopologicalSort(t *testing.T) {
	graph := &BuildGraph{
		Nodes: []BuildNode{
			{ID: "a", Dependencies: []string{"b"}},
			{ID: "b", Dependencies: []string{"c"}},
			{ID: "c", Dependencies: []string{}},
		},
	}

	sorted, err := graph.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if len(sorted) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(sorted))
	}

	cIdx, bIdx, aIdx := -1, -1, -1
	for i, id := range sorted {
		switch id {
		case "c":
			cIdx = i
		case "b":
			bIdx = i
		case "a":
			aIdx = i
		}
	}

	if cIdx >= bIdx || bIdx >= aIdx {
		t.Errorf("Invalid order: c=%d, b=%d, a=%d", cIdx, bIdx, aIdx)
	}
}

func TestTopologicalSortCycle(t *testing.T) {
	graph := &BuildGraph{
		Nodes: []BuildNode{
			{ID: "a", Dependencies: []string{"b"}},
			{ID: "b", Dependencies: []string{"a"}},
		},
	}

	_, err := graph.TopologicalSort()
	if err == nil {
		t.Error("Expected error for cyclic graph")
	}
}

func TestTopologicalSortNoDeps(t *testing.T) {
	graph := &BuildGraph{
		Nodes: []BuildNode{
			{ID: "a", Dependencies: []string{}},
			{ID: "b", Dependencies: []string{}},
			{ID: "c", Dependencies: []string{}},
		},
	}

	sorted, err := graph.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if len(sorted) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(sorted))
	}
}

func TestLinkNodeDependsOnCompileNode(t *testing.T) {
	cfg := &config.Config{
		Package: config.Package{
			Name:    "dep-test",
			Version: "1.0.0",
		},
		Build: config.Build{
			Compiler: "auto",
			STD:      "c++20",
			Libraries: []config.Library{
				{
					Name:    "mylib",
					Type:    "static",
					Sources: []string{"lib.cpp"},
				},
			},
			Executables: []config.Executable{
				{
					Name:      "app",
					Sources:   []string{"main.cpp"},
					Libraries: []string{"mylib"},
				},
			},
		},
	}

	graph, err := Generate(cfg)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	sorted, err := graph.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	libIdx, exeIdx, linkIdx := -1, -1, -1
	for i, id := range sorted {
		if id == "lib:mylib" {
			libIdx = i
		}
		if id == "exe:app" {
			exeIdx = i
		}
		if id == "link:app" {
			linkIdx = i
		}
	}

	if libIdx == -1 || exeIdx == -1 || linkIdx == -1 {
		t.Fatalf("Expected lib:mylib, exe:app, link:app in sorted output, got %v", sorted)
	}

	if libIdx >= linkIdx {
		t.Errorf("lib:mylib (idx=%d) should come before link:app (idx=%d)", libIdx, linkIdx)
	}
	if exeIdx >= linkIdx {
		t.Errorf("exe:app (idx=%d) should come before link:app (idx=%d)", exeIdx, linkIdx)
	}
}

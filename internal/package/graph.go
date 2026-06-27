package pkg

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

type DepGraph struct {
	Nodes map[string]*DepNode
}

type DepNode struct {
	Name         string
	Version      string
	Dependencies map[string]string
}

func NewDepGraph() *DepGraph {
	return &DepGraph{
		Nodes: make(map[string]*DepNode),
	}
}

func ResolveDependencies(root *Module, registry Registry) (*DepGraph, error) {
	graph := NewDepGraph()

	visited := make(map[string]bool)
	resolving := make(map[string]bool)

	var resolve func(name, constraint string) error
	resolve = func(name, constraint string) error {
		if visited[name] {
			return nil
		}

		if resolving[name] {
			return fmt.Errorf("circular dependency detected: %s", name)
		}

		resolving[name] = true
		defer func() { resolving[name] = false }()

		versions, err := registry.GetVersions(name)
		if err != nil {
			return fmt.Errorf("failed to get versions for %s: %w", name, err)
		}

		parsedVersions := make([]*semver.Version, 0, len(versions))
		for _, v := range versions {
			pv, err := semver.NewVersion(v)
			if err != nil {
				continue
			}
			parsedVersions = append(parsedVersions, pv)
		}

		c, err := ParseConstraint(constraint)
		if err != nil {
			return fmt.Errorf("invalid constraint %q for %s: %w", constraint, name, err)
		}

		resolved, err := ResolveVersion(parsedVersions, c)
		if err != nil {
			return fmt.Errorf("failed to resolve version for %s: %w", name, err)
		}

		info, err := registry.GetPackage(name)
		if err != nil {
			return fmt.Errorf("failed to get package info for %s: %w", name, err)
		}

		graph.Nodes[name] = &DepNode{
			Name:         name,
			Version:      resolved.String(),
			Dependencies: info.Deps,
		}

		visited[name] = true

		for depName, depConstraint := range info.Deps {
			if err := resolve(depName, depConstraint); err != nil {
				return err
			}
		}

		return nil
	}

	for name, constraint := range root.Deps {
		if err := resolve(name, constraint); err != nil {
			return nil, err
		}
	}

	for name, constraint := range root.DevDeps {
		if err := resolve(name, constraint); err != nil {
			return nil, err
		}
	}

	return graph, nil
}

func (g *DepGraph) TopologicalSort() ([]string, error) {
	visited := make(map[string]bool)
	visiting := make(map[string]bool)
	var order []string

	var visit func(name string) error
	visit = func(name string) error {
		if visited[name] {
			return nil
		}

		if visiting[name] {
			return fmt.Errorf("circular dependency detected: %s", name)
		}

		visiting[name] = true

		node, ok := g.Nodes[name]
		if !ok {
			return fmt.Errorf("dependency %s not found in graph", name)
		}

		for dep := range node.Dependencies {
			if err := visit(dep); err != nil {
				return err
			}
		}

		visiting[name] = false
		visited[name] = true
		order = append(order, name)

		return nil
	}

	for name := range g.Nodes {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	return order, nil
}

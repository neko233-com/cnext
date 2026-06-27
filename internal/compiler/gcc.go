package compiler

import (
	"fmt"
	"os/exec"
	"strings"
)

type GCC struct {
	path string
}

func NewGCC(path string) *GCC {
	return &GCC{path: path}
}

func (g *GCC) Name() string {
	return "gcc"
}

func (g *GCC) Version() string {
	cmd := exec.Command(g.path, "--version")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return "unknown"
}

func (g *GCC) Compile(opts CompileOptions) error {
	args := []string{}

	if opts.Std != "" {
		args = append(args, "-std="+opts.Std)
	}

	switch opts.Optimization {
	case "debug":
		args = append(args, "-O0", "-g")
	case "release":
		args = append(args, "-O2")
	case "size":
		args = append(args, "-Os")
	case "speed":
		args = append(args, "-O3")
	}

	for _, dir := range opts.IncludeDirs {
		args = append(args, "-I"+dir)
	}

	for _, def := range opts.Defines {
		args = append(args, "-D"+def)
	}

	args = append(args, opts.Flags...)
	args = append(args, "-o", opts.Output)
	args = append(args, opts.Sources...)

	cmd := exec.Command(g.path, args...)
	cmd.Stdout = nil

	if err := runCommand(cmd); err != nil {
		return fmt.Errorf("gcc compilation failed: %w", err)
	}

	return nil
}

func (g *GCC) Link(opts LinkOptions) error {
	args := []string{}
	args = append(args, opts.Objects...)
	args = append(args, "-o", opts.Output)

	for _, dir := range opts.LibDirs {
		args = append(args, "-L"+dir)
	}

	for _, lib := range opts.Libraries {
		args = append(args, "-l"+lib)
	}

	args = append(args, opts.Flags...)

	cmd := exec.Command(g.path, args...)
	cmd.Stdout = nil

	if err := runCommand(cmd); err != nil {
		return fmt.Errorf("gcc linking failed: %w", err)
	}

	return nil
}

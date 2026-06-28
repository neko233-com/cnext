package compiler

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
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
	lang := opts.Language
	if lang == "" || lang == "auto" {
		if len(opts.Sources) > 0 {
			lang = detectLanguage(opts.Sources[0])
		} else {
			lang = "cpp"
		}
	}

	bin := CompilerBinary(g.path, lang)

	// When compiling multiple source files with -c, compile each separately
	if opts.CompileOnly && len(opts.Sources) > 1 {
		for i, src := range opts.Sources {
			objName := strings.TrimSuffix(opts.Output, filepath.Ext(opts.Output))
			if i > 0 {
				objName += fmt.Sprintf("_%d", i)
			}
			objName += filepath.Ext(opts.Output)

			compileOpts := CompileOptions{
				Sources:      []string{src},
				Output:       objName,
				Std:          opts.Std,
				Language:     lang,
				IncludeDirs:  opts.IncludeDirs,
				Defines:      opts.Defines,
				Flags:        opts.Flags,
				Optimization: opts.Optimization,
				CompileOnly:  true,
				PIC:          opts.PIC,
				DepFile:      strings.TrimSuffix(objName, filepath.Ext(objName)) + ".d",
			}
			if err := g.compileSingle(bin, compileOpts); err != nil {
				return err
			}
		}
		return nil
	}

	// Generate dep file for single source
	if opts.CompileOnly && opts.DepFile == "" {
		objBase := strings.TrimSuffix(opts.Output, filepath.Ext(opts.Output))
		opts.DepFile = objBase + ".d"
	}

	return g.compileSingle(bin, opts)
}

func (g *GCC) compileSingle(bin string, opts CompileOptions) error {
	args := []string{}

	if opts.CompileOnly {
		args = append(args, "-c")
	}

	if opts.PIC {
		args = append(args, "-fPIC")
	}

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

	// Header dependency tracking
	if opts.CompileOnly && opts.DepFile != "" {
		args = append(args, "-MMD", "-MF", opts.DepFile)
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

	cmd := exec.Command(bin, args...)
	cmd.Stdout = nil

	if err := runCommand(cmd); err != nil {
		return fmt.Errorf("gcc compilation failed: %w", err)
	}

	return nil
}

func (g *GCC) Link(opts LinkOptions) error {
	// On Windows, don't add -l prefix for libraries
	isWindows := runtime.GOOS == "windows"

	args := []string{}
	args = append(args, opts.Objects...)
	args = append(args, "-o", opts.Output)

	for _, dir := range opts.LibDirs {
		args = append(args, "-L"+dir)
	}

	for _, lib := range opts.Libraries {
		if isWindows {
			args = append(args, lib)
		} else {
			args = append(args, "-l"+lib)
		}
	}

	args = append(args, opts.Flags...)

	lang := opts.Language
	if lang == "" {
		lang = "cpp"
	}
	bin := CompilerBinary(g.path, lang)
	cmd := exec.Command(bin, args...)
	cmd.Stdout = nil

	if err := runCommand(cmd); err != nil {
		return fmt.Errorf("gcc linking failed: %w", err)
	}

	return nil
}

func (g *GCC) Archive(objects []string, output string) error {
	// Create static library using ar
	args := []string{"rcs", output}
	args = append(args, objects...)

	// Find ar binary
	arPath := "ar"
	if runtime.GOOS == "windows" {
		// On Windows with MinGW, ar might be in the same dir as gcc
		dir := filepath.Dir(g.path)
		candidates := []string{
			filepath.Join(dir, "ar.exe"),
			filepath.Join(dir, "ar"),
			"ar",
		}
		for _, c := range candidates {
			if _, err := exec.LookPath(c); err == nil {
				arPath = c
				break
			}
		}
	}

	cmd := exec.Command(arPath, args...)
	cmd.Stdout = nil

	if err := runCommand(cmd); err != nil {
		return fmt.Errorf("gcc archive failed: %w", err)
	}

	return nil
}

func (g *GCC) SharedArchive(objects []string, output string, libDirs []string, libraries []string) error {
	// Create shared library using gcc -shared
	args := []string{"-shared"}
	args = append(args, objects...)
	args = append(args, "-o", output)

	for _, dir := range libDirs {
		args = append(args, "-L"+dir)
	}

	for _, lib := range libraries {
		if runtime.GOOS == "windows" {
			args = append(args, lib)
		} else {
			args = append(args, "-l"+lib)
		}
	}

	cmd := exec.Command(g.path, args...)
	cmd.Stdout = nil

	if err := runCommand(cmd); err != nil {
		return fmt.Errorf("gcc shared archive failed: %w", err)
	}

	return nil
}

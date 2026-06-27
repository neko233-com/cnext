package compiler

import (
	"fmt"
	"os/exec"
	"strings"
)

type Clang struct {
	path string
}

func NewClang(path string) *Clang {
	return &Clang{path: path}
}

func (c *Clang) Name() string {
	return "clang"
}

func (c *Clang) Version() string {
	cmd := exec.Command(c.path, "--version")
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

func (c *Clang) Compile(opts CompileOptions) error {
	args := []string{}

	lang := opts.Language
	if lang == "" || lang == "auto" {
		if len(opts.Sources) > 0 {
			lang = detectLanguage(opts.Sources[0])
		} else {
			lang = "cpp"
		}
	}

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

	for _, dir := range opts.IncludeDirs {
		args = append(args, "-I"+dir)
	}

	for _, def := range opts.Defines {
		args = append(args, "-D"+def)
	}

	args = append(args, opts.Flags...)
	args = append(args, "-o", opts.Output)
	args = append(args, opts.Sources...)

	bin := CompilerBinary(c.path, lang)
	cmd := exec.Command(bin, args...)
	cmd.Stdout = nil

	if err := runCommand(cmd); err != nil {
		return fmt.Errorf("clang compilation failed: %w", err)
	}

	return nil
}

func (c *Clang) Link(opts LinkOptions) error {
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

	lang := opts.Language
	if lang == "" {
		lang = "cpp"
	}
	bin := CompilerBinary(c.path, lang)
	cmd := exec.Command(bin, args...)
	cmd.Stdout = nil

	if err := runCommand(cmd); err != nil {
		return fmt.Errorf("clang linking failed: %w", err)
	}

	return nil
}

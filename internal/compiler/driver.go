package compiler

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/neko233-com/cnext/internal/config"
)

type Compiler interface {
	Name() string
	Version() string
	Compile(opts CompileOptions) error
	Link(opts LinkOptions) error
}

type CompileOptions struct {
	Sources      []string
	Output       string
	Std          string
	IncludeDirs  []string
	Defines      []string
	Flags        []string
	Optimization string
}

type LinkOptions struct {
	Objects   []string
	Output    string
	Libraries []string
	LibDirs   []string
	Flags     []string
}

func AutoDetect() (Compiler, error) {
	switch runtime.GOOS {
	case "linux", "darwin":
		if path, err := exec.LookPath("gcc"); err == nil {
			return NewGCC(path), nil
		}
		if path, err := exec.LookPath("clang"); err == nil {
			return NewClang(path), nil
		}
	case "windows":
		if path, err := exec.LookPath("gcc"); err == nil {
			return NewGCC(path), nil
		}
		if path, err := exec.LookPath("clang"); err == nil {
			return NewClang(path), nil
		}
	}

	return nil, fmt.Errorf("no compiler found")
}

func AutoDetectForTarget(target config.Target) (Compiler, error) {
	return AutoDetect()
}

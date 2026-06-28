package compiler

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/neko233-com/cnext/internal/config"
	"github.com/neko233-com/cnext/internal/toolchain"
)

type Compiler interface {
	Name() string
	Version() string
	Compile(opts CompileOptions) error
	Link(opts LinkOptions) error
	Archive(objects []string, output string) error
	SharedArchive(objects []string, output string, libDirs []string, libraries []string) error
}

type CompileOptions struct {
	Sources      []string
	Output       string
	Std          string
	Language     string // "c", "cpp", "auto" — auto detects from file extension
	IncludeDirs  []string
	Defines      []string
	Flags        []string
	Optimization string
	CompileOnly  bool   // -c flag for compile-only (no linking)
	PIC          bool   // -fPIC for position-independent code
	DepFile      string // -MF path for header dependency tracking
}

type LinkOptions struct {
	Objects   []string
	Output    string
	Libraries []string
	LibDirs   []string
	Flags     []string
	Language  string // "c" or "cpp" — used to select correct linker
}

func AutoDetect() (Compiler, error) {
	switch runtime.GOOS {
	case "linux", "darwin":
		if path, err := exec.LookPath("g++"); err == nil {
			return NewGCC(path), nil
		}
		if path, err := exec.LookPath("gcc"); err == nil {
			return NewGCC(path), nil
		}
		if path, err := exec.LookPath("clang++"); err == nil {
			return NewClang(path), nil
		}
		if path, err := exec.LookPath("clang"); err == nil {
			return NewClang(path), nil
		}
	case "windows":
		if path, err := exec.LookPath("g++"); err == nil {
			return NewGCC(path), nil
		}
		if path, err := exec.LookPath("gcc"); err == nil {
			return NewGCC(path), nil
		}
		if path, err := exec.LookPath("clang++"); err == nil {
			return NewClang(path), nil
		}
		if path, err := exec.LookPath("clang"); err == nil {
			return NewClang(path), nil
		}
	}

	return tryBundledToolchain()
}

func tryBundledToolchain() (Compiler, error) {
	mgr := toolchain.NewManager()
	installed := mgr.List()

	for _, tc := range installed {
		compilerPath, err := mgr.FindCompiler(tc.Name, tc.Version, tc.Platform)
		if err != nil {
			continue
		}

		switch tc.Name {
		case "gcc":
			return NewGCC(compilerPath), nil
		case "clang":
			return NewClang(compilerPath), nil
		}
	}

	return nil, fmt.Errorf("no compiler found")
}

func detectLanguage(filename string) string {
	ext := filepath.Ext(filename)
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

// CompilerBinary returns the correct compiler binary path for a given language.
// For C++ files, it returns the C++ compiler (g++/clang++). For C files, the C compiler.
func CompilerBinary(compilerPath string, language string) string {
	if language == "cpp" {
		base := filepath.Base(compilerPath)
		dir := filepath.Dir(compilerPath)
		switch base {
		case "gcc":
			return filepath.Join(dir, "g++")
		case "clang":
			return filepath.Join(dir, "clang++")
		}
	}
	return compilerPath
}

func AutoDetectForTarget(target config.Target) (Compiler, error) {
	if target.OS == "auto" || target.Arch == "auto" {
		return AutoDetect()
	}

	triple := &TargetTriple{
		Arch:   target.Arch,
		Vendor: "unknown",
		OS:     target.OS,
		Env:    "gnu",
	}

	return DetectForTarget(triple)
}

// runCommand runs the given command and returns an error that includes
// the stderr output when the command fails. This ensures compiler/linker
// error messages are visible to the user.
func runCommand(cmd *exec.Cmd) error {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		output := strings.TrimSpace(stderr.String())
		if output != "" {
			return fmt.Errorf("%w\n\ncompiler output:\n%s", err, output)
		}
		return err
	}
	return nil
}

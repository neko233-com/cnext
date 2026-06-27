package cmake

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CMakeBridge struct {
	buildDir   string
	installDir string
}

func NewBridge(projectDir string) *CMakeBridge {
	buildDir := filepath.Join(projectDir, ".cnext", "cmake-build")
	installDir := filepath.Join(projectDir, ".cnext", "cmake-install")
	return &CMakeBridge{
		buildDir:   buildDir,
		installDir: installDir,
	}
}

func (b *CMakeBridge) Configure(sourceDir string, options []string) error {
	if err := os.MkdirAll(b.buildDir, 0755); err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}

	args := []string{
		"-B", b.buildDir,
		"-S", sourceDir,
		"-DCMAKE_INSTALL_PREFIX=" + b.installDir,
	}
	args = append(args, options...)

	var stderr bytes.Buffer
	cmd := exec.Command("cmake", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		output := strings.TrimSpace(stderr.String())
		if output != "" {
			return fmt.Errorf("cmake configure failed: %w\n\noutput:\n%s", err, output)
		}
		return fmt.Errorf("cmake configure failed: %w", err)
	}

	return nil
}

func (b *CMakeBridge) Build(jobs int) error {
	args := []string{"--build", b.buildDir}
	if jobs > 0 {
		args = append(args, "--parallel", fmt.Sprintf("%d", jobs))
	}

	var stderr bytes.Buffer
	cmd := exec.Command("cmake", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		output := strings.TrimSpace(stderr.String())
		if output != "" {
			return fmt.Errorf("cmake build failed: %w\n\noutput:\n%s", err, output)
		}
		return fmt.Errorf("cmake build failed: %w", err)
	}

	return nil
}

func (b *CMakeBridge) Install() error {
	var stderr bytes.Buffer
	cmd := exec.Command("cmake", "--install", b.buildDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		output := strings.TrimSpace(stderr.String())
		if output != "" {
			return fmt.Errorf("cmake install failed: %w\n\noutput:\n%s", err, output)
		}
		return fmt.Errorf("cmake install failed: %w", err)
	}

	return nil
}

func (b *CMakeBridge) GetIncludeDirs() []string {
	return []string{filepath.Join(b.installDir, "include")}
}

func (b *CMakeBridge) GetLibDirs() []string {
	return []string{filepath.Join(b.installDir, "lib")}
}

func (b *CMakeBridge) GetLibraries() []string {
	libDir := filepath.Join(b.installDir, "lib")
	entries, err := os.ReadDir(libDir)
	if err != nil {
		return nil
	}

	var libs []string
	seen := make(map[string]bool)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		libName := extractLibName(name)
		if libName != "" && !seen[libName] {
			seen[libName] = true
			libs = append(libs, libName)
		}
	}

	return libs
}

func extractLibName(filename string) string {
	name := filename

	// Strip lib prefix on Unix
	if strings.HasPrefix(name, "lib") {
		name = name[3:]
	}

	// Strip known extensions
	for _, ext := range []string{".a", ".so", ".dylib", ".dll", ".lib"} {
		if strings.HasSuffix(name, ext) {
			name = name[:len(name)-len(ext)]
			break
		}
	}

	// Strip version suffixes like .1.2.3
	if idx := strings.LastIndex(name, "."); idx > 0 {
		suffix := name[idx+1:]
		allDigits := true
		for _, c := range suffix {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			name = name[:idx]
		}
	}

	return name
}

package cmake

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CMakeDep struct {
	Name         string
	Version      string
	GitURL       string
	GitTag       string
	Source       string // "git", "url", "local"
	CMakeOptions []string
}

func ResolveDep(dep CMakeDep, cacheDir string) (string, error) {
	switch dep.Source {
	case "git":
		return resolveGitDep(dep, cacheDir)
	case "url":
		return resolveURLDep(dep, cacheDir)
	case "local":
		return dep.GitURL, nil
	default:
		return "", fmt.Errorf("unknown source type: %s", dep.Source)
	}
}

func resolveGitDep(dep CMakeDep, cacheDir string) (string, error) {
	dest := filepath.Join(cacheDir, "cmake", dep.Name, dep.Version)

	if _, err := os.Stat(filepath.Join(dest, ".git")); err == nil {
		var stderr bytes.Buffer
		cmd := exec.Command("git", "-C", dest, "pull")
		cmd.Stderr = &stderr
		cmd.Run()
		return dest, nil
	}

	args := []string{"clone", "--depth", "1"}
	if dep.GitTag != "" {
		args = append(args, "--branch", dep.GitTag)
	}
	args = append(args, dep.GitURL, dest)

	var stderr bytes.Buffer
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		output := strings.TrimSpace(stderr.String())
		if output != "" {
			return "", fmt.Errorf("git clone failed: %w\n\noutput:\n%s", err, output)
		}
		return "", fmt.Errorf("git clone failed: %w", err)
	}

	return dest, nil
}

func resolveURLDep(dep CMakeDep, cacheDir string) (string, error) {
	dest := filepath.Join(cacheDir, "cmake", dep.Name, dep.Version)
	if err := os.MkdirAll(dest, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}
	// TODO: implement download and extraction
	return dest, nil
}

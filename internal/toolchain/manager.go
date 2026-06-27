package toolchain

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const ToolchainDir = ".cnext/toolchains"

type Manager struct {
	baseDir string
}

func NewManager() *Manager {
	home, _ := os.UserHomeDir()
	return &Manager{
		baseDir: filepath.Join(home, ToolchainDir),
	}
}

func NewManagerWithDir(dir string) *Manager {
	return &Manager{baseDir: dir}
}

func (m *Manager) IsInstalled(name, version, platform string) bool {
	path := filepath.Join(m.baseDir, name, version, platform)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (m *Manager) GetPath(name, version, platform string) string {
	return filepath.Join(m.baseDir, name, version, platform)
}

func (m *Manager) Install(info ToolchainInfo) error {
	if m.IsInstalled(info.Name, info.Version, info.Platform) {
		return nil
	}

	destDir := m.GetPath(info.Name, info.Version, info.Platform)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create toolchain directory: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "cnext-toolchain-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, filepath.Base(info.URL))
	if err := Download(info.URL, archivePath); err != nil {
		return fmt.Errorf("failed to download toolchain: %w", err)
	}

	if info.SHA256 != "" {
		if err := VerifyChecksum(archivePath, info.SHA256); err != nil {
			return fmt.Errorf("checksum verification failed: %w", err)
		}
	}

	if err := Extract(archivePath, destDir); err != nil {
		return fmt.Errorf("failed to extract toolchain: %w", err)
	}

	return nil
}

func (m *Manager) List() []InstalledToolchain {
	var installed []InstalledToolchain

	entries, err := os.ReadDir(m.baseDir)
	if err != nil {
		return installed
	}

	for _, nameEntry := range entries {
		if !nameEntry.IsDir() {
			continue
		}
		name := nameEntry.Name()

		versionEntries, err := os.ReadDir(filepath.Join(m.baseDir, name))
		if err != nil {
			continue
		}

		for _, versionEntry := range versionEntries {
			if !versionEntry.IsDir() {
				continue
			}
			version := versionEntry.Name()

			platformEntries, err := os.ReadDir(filepath.Join(m.baseDir, name, version))
			if err != nil {
				continue
			}

			for _, platformEntry := range platformEntries {
				if !platformEntry.IsDir() {
					continue
				}
				platform := platformEntry.Name()
				installed = append(installed, InstalledToolchain{
					Name:     name,
					Version:  version,
					Platform: platform,
					Path:     filepath.Join(m.baseDir, name, version, platform),
				})
			}
		}
	}

	return installed
}

func (m *Manager) Remove(name, version, platform string) error {
	path := m.GetPath(name, version, platform)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("toolchain %s@%s (%s) is not installed", name, version, platform)
	}
	return os.RemoveAll(path)
}

func (m *Manager) FindCompiler(name, version, platform string) (string, error) {
	if !m.IsInstalled(name, version, platform) {
		return "", fmt.Errorf("toolchain %s@%s (%s) is not installed", name, version, platform)
	}

	toolchainDir := m.GetPath(name, version, platform)

	switch name {
	case "gcc":
		candidates := []string{"g++", "gcc"}
		for _, c := range candidates {
			p := findExecutable(toolchainDir, c)
			if p != "" {
				return p, nil
			}
		}
	case "clang":
		candidates := []string{"clang++", "clang"}
		for _, c := range candidates {
			p := findExecutable(toolchainDir, c)
			if p != "" {
				return p, nil
			}
		}
	}

	return "", fmt.Errorf("no compiler binary found in toolchain %s@%s", name, version)
}

func findExecutable(baseDir, name string) string {
	platform := runtime.GOOS
	binDir := "bin"
	if platform == "windows" {
		name = name + ".exe"
	}

	path := filepath.Join(baseDir, binDir, name)
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		return path
	}

	path = filepath.Join(baseDir, name)
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		return path
	}

	return ""
}

type InstalledToolchain struct {
	Name     string
	Version  string
	Platform string
	Path     string
}

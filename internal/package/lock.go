package pkg

import (
	"bytes"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type LockFile struct {
	Version int            `toml:"version"`
	Modules []LockedModule `toml:"modules"`
}

type LockedModule struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
	Hash    string `toml:"hash"`
	Source  string `toml:"source"`
}

func LoadLockFile(path string) (*LockFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	var lock LockFile
	if err := toml.Unmarshal(data, &lock); err != nil {
		return nil, fmt.Errorf("failed to parse lock file: %w", err)
	}

	return &lock, nil
}

func SaveLockFile(lock *LockFile, path string) error {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(lock); err != nil {
		return fmt.Errorf("failed to encode lock file: %w", err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write lock file: %w", err)
	}

	return nil
}

func GenerateLockFile(graph *DepGraph) *LockFile {
	lock := &LockFile{
		Version: 1,
		Modules: make([]LockedModule, 0, len(graph.Nodes)),
	}

	for _, node := range graph.Nodes {
		lock.Modules = append(lock.Modules, LockedModule{
			Name:    node.Name,
			Version: node.Version,
			Hash:    "",
			Source:  "registry",
		})
	}

	return lock
}

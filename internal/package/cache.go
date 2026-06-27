package pkg

import (
	"fmt"
	"os"
	"path/filepath"
)

const CacheDir = ".cnext/cache"

type Cache struct {
	dir string
}

func NewCache() *Cache {
	home, _ := os.UserHomeDir()
	return &Cache{dir: filepath.Join(home, CacheDir)}
}

func NewCacheWithDir(dir string) *Cache {
	return &Cache{dir: dir}
}

func (c *Cache) Get(name, version string) (string, bool) {
	path := filepath.Join(c.dir, fmt.Sprintf("%s-%s", name, version))
	if _, err := os.Stat(path); err == nil {
		return path, true
	}
	return "", false
}

func (c *Cache) Put(name, version, path string) error {
	dest := filepath.Join(c.dir, fmt.Sprintf("%s-%s", name, version))
	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read package file: %w", err)
	}

	if err := os.WriteFile(dest, data, 0644); err != nil {
		return fmt.Errorf("failed to write to cache: %w", err)
	}

	return nil
}

func (c *Cache) Clean() error {
	if err := os.RemoveAll(c.dir); err != nil {
		return fmt.Errorf("failed to clean cache: %w", err)
	}
	return nil
}

package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"
)

const CacheDir = ".cnext/cache"

type BuildCache struct {
	dir     string
	entries map[string]*CacheEntry
}

type CacheEntry struct {
	ObjectFile string    `json:"object_file"`
	DepFile    string    `json:"dep_file"`
	SourceHash string    `json:"source_hash"`
	DepHash    string    `json:"dep_hash"`
	Timestamp  time.Time `json:"timestamp"`
	Compiler   string    `json:"compiler"`
	Flags      string    `json:"flags"`
}

func New() *BuildCache {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, CacheDir)
	return &BuildCache{
		dir:     dir,
		entries: make(map[string]*CacheEntry),
	}
}

func NewWithDir(dir string) *BuildCache {
	return &BuildCache{
		dir:     dir,
		entries: make(map[string]*CacheEntry),
	}
}

func (c *BuildCache) Load() error {
	indexFile := filepath.Join(c.dir, "index.json")
	data, err := os.ReadFile(indexFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, &c.entries)
}

func (c *BuildCache) Save() error {
	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(c.dir, "index.json"), data, 0644)
}

func (c *BuildCache) GetCachePath() string {
	return c.dir
}

func (c *BuildCache) IsFresh(sourceFile string, objectFile string, compiler string, flags []string, includeDirs []string) bool {
	entry, exists := c.entries[sourceFile]
	if !exists {
		return false
	}

	// Check if object file exists
	if _, err := os.Stat(objectFile); os.IsNotExist(err) {
		return false
	}

	// Check source hash
	sourceHash, err := c.hashFile(sourceFile)
	if err != nil {
		return false
	}
	if sourceHash != entry.SourceHash {
		return false
	}

	// Check dep file for header changes
	if entry.DepFile != "" {
		depHash, err := c.hashDepFile(entry.DepFile, includeDirs)
		if err == nil && depHash != entry.DepHash {
			return false
		}
	}

	// Check compiler and flags
	if compiler != entry.Compiler {
		return false
	}

	flagsStr := joinFlags(flags)
	if flagsStr != entry.Flags {
		return false
	}

	return true
}

func (c *BuildCache) Update(sourceFile string, objectFile string, depFile string, compiler string, flags []string) error {
	sourceHash, err := c.hashFile(sourceFile)
	if err != nil {
		return err
	}

	var depHash string
	if depFile != "" {
		depHash, _ = c.hashDepFile(depFile, nil)
	}

	c.entries[sourceFile] = &CacheEntry{
		ObjectFile: objectFile,
		DepFile:    depFile,
		SourceHash: sourceHash,
		DepHash:    depHash,
		Timestamp:  time.Now(),
		Compiler:   compiler,
		Flags:      joinFlags(flags),
	}

	return c.Save()
}

func (c *BuildCache) Clean() error {
	return os.RemoveAll(c.dir)
}

func (c *BuildCache) hashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h), nil
}

func (c *BuildCache) hashDepFile(depFile string, includeDirs []string) (string, error) {
	data, err := os.ReadFile(depFile)
	if err != nil {
		return "", err
	}

	// Parse .d file to get header dependencies
	headers := parseDepFile(data)

	// Hash all headers
	var hashes []string
	for _, header := range headers {
		hash, err := c.hashFile(header)
		if err != nil {
			continue
		}
		hashes = append(hashes, hash)
	}

	sort.Strings(hashes)
	combined := joinStrings(hashes)
	h := sha256.Sum256([]byte(combined))
	return fmt.Sprintf("%x", h), nil
}

func parseDepFile(data []byte) []string {
	content := string(data)
	var headers []string
	var current string

	for i := 0; i < len(content); i++ {
		ch := content[i]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' || ch == ':' {
			if current != "" && current != "\\" {
				headers = append(headers, current)
			}
			current = ""
		} else if ch == '\\' {
			// Skip escaped character
			i++
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		headers = append(headers, current)
	}

	return headers
}

func joinFlags(flags []string) string {
	sorted := make([]string, len(flags))
	copy(sorted, flags)
	sort.Strings(sorted)
	result := ""
	for _, f := range sorted {
		result += f + " "
	}
	return result
}

func joinStrings(strs []string) string {
	result := ""
	for _, s := range strs {
		result += s
	}
	return result
}

func GetNumCPU() int {
	return runtime.NumCPU()
}

func GetBuildJobs(jobsFlag int) int {
	if jobsFlag > 0 {
		return jobsFlag
	}
	numCPU := runtime.NumCPU()
	// Use all CPUs but leave one for OS
	if numCPU > 1 {
		return numCPU - 1
	}
	return 1
}

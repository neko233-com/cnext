package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Registry interface {
	GetPackage(name string) (*PackageInfo, error)
	GetVersions(name string) ([]string, error)
	Download(name, version string) (string, error)
}

type PackageInfo struct {
	Name    string            `json:"name"`
	Version string            `json:"version"`
	Deps    map[string]string `json:"dependencies"`
	URL     string            `json:"url"`
	Hash    string            `json:"hash"`
}

type HTTPRegistry struct {
	BaseURL string
	Client  *http.Client
}

func NewHTTPRegistry(baseURL string) *HTTPRegistry {
	return &HTTPRegistry{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (r *HTTPRegistry) GetPackage(name string) (*PackageInfo, error) {
	resp, err := r.Client.Get(fmt.Sprintf("%s/api/v1/packages/%s", r.BaseURL, name))
	if err != nil {
		return nil, fmt.Errorf("failed to get package %s: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("package %s not found", name)
	}

	var info PackageInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode package info: %w", err)
	}

	return &info, nil
}

func (r *HTTPRegistry) GetVersions(name string) ([]string, error) {
	resp, err := r.Client.Get(fmt.Sprintf("%s/api/v1/packages/%s/versions", r.BaseURL, name))
	if err != nil {
		return nil, fmt.Errorf("failed to get versions for %s: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("versions not found for package %s", name)
	}

	var versions []string
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("failed to decode versions: %w", err)
	}

	return versions, nil
}

func (r *HTTPRegistry) Download(name, version string) (string, error) {
	resp, err := r.Client.Get(fmt.Sprintf("%s/api/v1/packages/%s/%s/download", r.BaseURL, name, version))
	if err != nil {
		return "", fmt.Errorf("failed to download package %s@%s: %w", name, version, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download package %s@%s", name, version)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	cacheDir := filepath.Join(home, ".cnext", "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	cachePath := filepath.Join(cacheDir, fmt.Sprintf("%s-%s.tar.gz", name, version))
	out, err := os.Create(cachePath)
	if err != nil {
		return "", fmt.Errorf("failed to create cache file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("failed to save package: %w", err)
	}

	return cachePath, nil
}

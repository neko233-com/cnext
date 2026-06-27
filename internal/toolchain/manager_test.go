package toolchain

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewManager(t *testing.T) {
	mgr := NewManager()
	if mgr.baseDir == "" {
		t.Error("Expected non-empty baseDir")
	}
	if mgr.baseDir != filepath.Join(os.Getenv("HOME"), ".cnext/toolchains") &&
		mgr.baseDir != filepath.Join(os.Getenv("USERPROFILE"), ".cnext/toolchains") {
		t.Logf("baseDir: %s", mgr.baseDir)
	}
}

func TestNewManagerWithDir(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManagerWithDir(dir)
	if mgr.baseDir != dir {
		t.Errorf("Expected baseDir %s, got %s", dir, mgr.baseDir)
	}
}

func TestIsInstalled_NotInstalled(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManagerWithDir(dir)
	if mgr.IsInstalled("gcc", "14.2.0", "linux-x64") {
		t.Error("Expected toolchain to not be installed")
	}
}

func TestIsInstalled_Installed(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManagerWithDir(dir)

	tcDir := filepath.Join(dir, "gcc", "14.2.0", "linux-x64")
	if err := os.MkdirAll(tcDir, 0755); err != nil {
		t.Fatal(err)
	}

	if !mgr.IsInstalled("gcc", "14.2.0", "linux-x64") {
		t.Error("Expected toolchain to be installed")
	}
}

func TestGetPath(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManagerWithDir(dir)
	expected := filepath.Join(dir, "gcc", "14.2.0", "linux-x64")
	actual := mgr.GetPath("gcc", "14.2.0", "linux-x64")
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestList_Empty(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManagerWithDir(dir)
	installed := mgr.List()
	if len(installed) != 0 {
		t.Errorf("Expected 0 installed toolchains, got %d", len(installed))
	}
}

func TestList_WithToolchains(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManagerWithDir(dir)

	dirs := []string{
		filepath.Join(dir, "gcc", "14.2.0", "linux-x64"),
		filepath.Join(dir, "clang", "18.1.0", "linux-x64"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatal(err)
		}
	}

	installed := mgr.List()
	if len(installed) != 2 {
		t.Errorf("Expected 2 installed toolchains, got %d", len(installed))
	}
}

func TestRemove(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManagerWithDir(dir)

	tcDir := filepath.Join(dir, "gcc", "14.2.0", "linux-x64")
	if err := os.MkdirAll(tcDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := mgr.Remove("gcc", "14.2.0", "linux-x64"); err != nil {
		t.Fatalf("Failed to remove toolchain: %v", err)
	}

	if mgr.IsInstalled("gcc", "14.2.0", "linux-x64") {
		t.Error("Expected toolchain to be removed")
	}
}

func TestRemove_NotInstalled(t *testing.T) {
	dir := t.TempDir()
	mgr := NewManagerWithDir(dir)

	err := mgr.Remove("gcc", "14.2.0", "linux-x64")
	if err == nil {
		t.Error("Expected error when removing non-installed toolchain")
	}
}

func TestFindInRegistry(t *testing.T) {
	tc, err := FindInRegistry("gcc", "14.2.0", "linux-x64")
	if err != nil {
		t.Fatalf("Failed to find toolchain: %v", err)
	}
	if tc.Name != "gcc" {
		t.Errorf("Expected name 'gcc', got '%s'", tc.Name)
	}
	if tc.Version != "14.2.0" {
		t.Errorf("Expected version '14.2.0', got '%s'", tc.Version)
	}
}

func TestFindInRegistry_NotFound(t *testing.T) {
	_, err := FindInRegistry("gcc", "99.99.99", "linux-x64")
	if err == nil {
		t.Error("Expected error for non-existent toolchain")
	}
}

func TestFindLatest(t *testing.T) {
	tc, err := FindLatest("gcc", "linux-x64")
	if err != nil {
		t.Fatalf("Failed to find latest: %v", err)
	}
	if tc.Name != "gcc" {
		t.Errorf("Expected name 'gcc', got '%s'", tc.Name)
	}
}

func TestFindLatest_NotFound(t *testing.T) {
	_, err := FindLatest("nonexistent", "linux-x64")
	if err == nil {
		t.Error("Expected error for non-existent toolchain")
	}
}

func TestParseToolchainArg(t *testing.T) {
	tests := []struct {
		arg     string
		name    string
		version string
	}{
		{"gcc", "gcc", ""},
		{"gcc@14.2.0", "gcc", "14.2.0"},
		{"clang@18.1.0", "clang", "18.1.0"},
	}

	for _, tt := range tests {
		name, version := ParseToolchainArg(tt.arg)
		if name != tt.name {
			t.Errorf("ParseToolchainArg(%s): expected name %s, got %s", tt.arg, tt.name, name)
		}
		if version != tt.version {
			t.Errorf("ParseToolchainArg(%s): expected version %s, got %s", tt.arg, tt.version, version)
		}
	}
}

func TestCurrentPlatform(t *testing.T) {
	p := CurrentPlatform()
	if p == "" {
		t.Error("Expected non-empty platform")
	}
	t.Logf("Current platform: %s", p)
}

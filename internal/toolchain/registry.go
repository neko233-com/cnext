package toolchain

import (
	"fmt"
	"runtime"
	"strings"
)

type ToolchainInfo struct {
	Name       string
	Version    string
	Platform   string
	URL        string
	SHA256     string
	Size       int64
	Components []string
}

var Registry = []ToolchainInfo{
	{
		Name:       "gcc",
		Version:    "14.2.0",
		Platform:   "linux-x64",
		URL:        "https://github.com/neko233-com/cnext-toolchains/releases/download/gcc-14.2.0-linux-x64/gcc-14.2.0-linux-x64.tar.gz",
		Components: []string{"gcc", "g++", "binutils"},
	},
	{
		Name:       "gcc",
		Version:    "14.2.0",
		Platform:   "linux-arm64",
		URL:        "https://github.com/neko233-com/cnext-toolchains/releases/download/gcc-14.2.0-linux-arm64/gcc-14.2.0-linux-arm64.tar.gz",
		Components: []string{"gcc", "g++", "binutils"},
	},
	{
		Name:       "gcc",
		Version:    "14.2.0",
		Platform:   "macos-x64",
		URL:        "https://github.com/neko233-com/cnext-toolchains/releases/download/gcc-14.2.0-macos-x64/gcc-14.2.0-macos-x64.tar.gz",
		Components: []string{"gcc", "g++", "binutils"},
	},
	{
		Name:       "gcc",
		Version:    "14.2.0",
		Platform:   "macos-arm64",
		URL:        "https://github.com/neko233-com/cnext-toolchains/releases/download/gcc-14.2.0-macos-arm64/gcc-14.2.0-macos-arm64.tar.gz",
		Components: []string{"gcc", "g++", "binutils"},
	},
	{
		Name:       "gcc",
		Version:    "14.2.0",
		Platform:   "windows-x64",
		URL:        "https://github.com/neko233-com/cnext-toolchains/releases/download/gcc-14.2.0-windows-x64/gcc-14.2.0-windows-x64.zip",
		Components: []string{"gcc", "g++", "binutils"},
	},
	{
		Name:       "clang",
		Version:    "18.1.0",
		Platform:   "linux-x64",
		URL:        "https://github.com/neko233-com/cnext-toolchains/releases/download/clang-18.1.0-linux-x64/clang-18.1.0-linux-x64.tar.gz",
		Components: []string{"clang", "clang++", "lld"},
	},
	{
		Name:       "clang",
		Version:    "18.1.0",
		Platform:   "linux-arm64",
		URL:        "https://github.com/neko233-com/cnext-toolchains/releases/download/clang-18.1.0-linux-arm64/clang-18.1.0-linux-arm64.tar.gz",
		Components: []string{"clang", "clang++", "lld"},
	},
	{
		Name:       "clang",
		Version:    "18.1.0",
		Platform:   "macos-x64",
		URL:        "https://github.com/neko233-com/cnext-toolchains/releases/download/clang-18.1.0-macos-x64/clang-18.1.0-macos-x64.tar.gz",
		Components: []string{"clang", "clang++", "lld"},
	},
	{
		Name:       "clang",
		Version:    "18.1.0",
		Platform:   "macos-arm64",
		URL:        "https://github.com/neko233-com/cnext-toolchains/releases/download/clang-18.1.0-macos-arm64/clang-18.1.0-macos-arm64.tar.gz",
		Components: []string{"clang", "clang++", "lld"},
	},
	{
		Name:       "clang",
		Version:    "18.1.0",
		Platform:   "windows-x64",
		URL:        "https://github.com/neko233-com/cnext-toolchains/releases/download/clang-18.1.0-windows-x64/clang-18.1.0-windows-x64.zip",
		Components: []string{"clang", "clang++", "lld"},
	},
}

func CurrentPlatform() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	return fmt.Sprintf("%s-%s", goos, goarch)
}

func FindInRegistry(name, version, platform string) (*ToolchainInfo, error) {
	for _, tc := range Registry {
		if tc.Name == name && tc.Version == version && tc.Platform == platform {
			return &tc, nil
		}
	}
	return nil, fmt.Errorf("toolchain %s@%s not found for platform %s", name, version, platform)
}

func FindLatest(name, platform string) (*ToolchainInfo, error) {
	var latest *ToolchainInfo
	for i := range Registry {
		tc := &Registry[i]
		if tc.Name == name && tc.Platform == platform {
			if latest == nil || tc.Version > latest.Version {
				latest = tc
			}
		}
	}
	if latest == nil {
		return nil, fmt.Errorf("no %s toolchain found for platform %s", name, platform)
	}
	return latest, nil
}

func ParseToolchainArg(arg string) (name, version string) {
	parts := strings.SplitN(arg, "@", 2)
	name = parts[0]
	if len(parts) > 1 {
		version = parts[1]
	}
	return
}

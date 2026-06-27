package compiler

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type TargetTriple struct {
	Arch   string
	Vendor string
	OS     string
	Env    string
}

var CommonTargets = map[string]*TargetTriple{
	"linux-x64":   {Arch: "x86_64", Vendor: "unknown", OS: "linux", Env: "gnu"},
	"linux-arm64": {Arch: "aarch64", Vendor: "unknown", OS: "linux", Env: "gnu"},
	"linux-arm":   {Arch: "arm", Vendor: "unknown", OS: "linux", Env: "gnueabihf"},
	"windows-x64": {Arch: "x86_64", Vendor: "pc", OS: "windows", Env: "gnu"},
	"macos-arm64": {Arch: "aarch64", Vendor: "apple", OS: "darwin", Env: "none"},
	"macos-x64":   {Arch: "x86_64", Vendor: "apple", OS: "darwin", Env: "none"},
}

func ParseTarget(triple string) (*TargetTriple, error) {
	if t, ok := CommonTargets[triple]; ok {
		return t, nil
	}

	parts := strings.Split(triple, "-")

	t := &TargetTriple{}

	switch len(parts) {
	case 4:
		t.Arch = parts[0]
		t.Vendor = parts[1]
		t.OS = parts[2]
		t.Env = parts[3]
	case 3:
		t.Arch = parts[0]
		t.Vendor = "unknown"
		t.OS = parts[1]
		t.Env = parts[2]
	case 2:
		t.Arch = parts[0]
		t.Vendor = "unknown"
		t.OS = parts[1]
		t.Env = "gnu"
	default:
		return nil, fmt.Errorf("invalid target triple: %s", triple)
	}

	if t.Arch == "" || t.OS == "" {
		return nil, fmt.Errorf("invalid target triple: %s (missing arch or OS)", triple)
	}

	return t, nil
}

func (t *TargetTriple) String() string {
	if t.Vendor != "" && t.Vendor != "unknown" {
		return fmt.Sprintf("%s-%s-%s-%s", t.Arch, t.Vendor, t.OS, t.Env)
	}
	if t.Env != "" && t.Env != "none" {
		return fmt.Sprintf("%s-%s-%s", t.Arch, t.OS, t.Env)
	}
	return fmt.Sprintf("%s-%s", t.Arch, t.OS)
}

func DetectHost() *TargetTriple {
	var arch, vendor, osName, env string

	switch runtime.GOARCH {
	case "amd64":
		arch = "x86_64"
	case "arm64":
		arch = "aarch64"
	case "arm":
		arch = "arm"
	default:
		arch = runtime.GOARCH
	}

	switch runtime.GOOS {
	case "linux":
		osName = "linux"
		vendor = "unknown"
		env = "gnu"
	case "darwin":
		osName = "darwin"
		vendor = "apple"
		env = "none"
	case "windows":
		osName = "windows"
		vendor = "pc"
		env = "gnu"
	default:
		osName = runtime.GOOS
		vendor = "unknown"
		env = "none"
	}

	return &TargetTriple{
		Arch:   arch,
		Vendor: vendor,
		OS:     osName,
		Env:    env,
	}
}

func DetectForTarget(target *TargetTriple) (Compiler, error) {
	host := DetectHost()

	if target.Arch == host.Arch && target.OS == host.OS {
		return AutoDetect()
	}

	crossPrefix := target.String() + "-"
	compilers := []string{"gcc", "g++", "clang", "clang++"}
	for _, comp := range compilers {
		path, err := exec.LookPath(crossPrefix + comp)
		if err == nil {
			switch comp {
			case "gcc":
				return NewGCC(path), nil
			case "g++":
				return NewGCC(path), nil
			case "clang":
				return NewClang(path), nil
			case "clang++":
				return NewClang(path), nil
			}
		}
	}

	return nil, fmt.Errorf("no cross-compiler found for %s", target.String())
}

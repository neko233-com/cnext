package cpu

import (
	"os/exec"
	"runtime"
	"strings"
)

type CPUFeatures struct {
	HasSSE2    bool
	HasSSE4    bool
	HasAVX     bool
	HasAVX2    bool
	HasAVX512  bool
	HasNEON    bool // ARM
	HasAES     bool
	HasBMI     bool
	HasBMI2    bool
	HasFMA     bool
	HasLZCNT   bool
	HasPOPCNT  bool
	CoreCount  int
	ModelName  string
}

func Detect() *CPUFeatures {
	features := &CPUFeatures{
		CoreCount: runtime.NumCPU(),
	}

	switch runtime.GOARCH {
	case "amd64":
		detectX86Features(features)
	case "arm64":
		detectARMFeatures(features)
	}

	return features
}

func detectX86Features(f *CPUFeatures) {
	// Try to detect CPU features via /proc/cpuinfo on Linux
	if runtime.GOOS == "linux" {
		data, err := exec.Command("grep", "-m1", "flags", "/proc/cpuinfo").Output()
		if err == nil {
			flags := string(data)
			f.HasSSE2 = strings.Contains(flags, "sse2")
			f.HasSSE4 = strings.Contains(flags, "sse4_1") || strings.Contains(flags, "sse4_2")
			f.HasAVX = strings.Contains(flags, " avx ")
			f.HasAVX2 = strings.Contains(flags, "avx2")
			f.HasAVX512 = strings.Contains(flags, "avx512")
			f.HasAES = strings.Contains(flags, "aes")
			f.HasBMI = strings.Contains(flags, "bmi1")
			f.HasBMI2 = strings.Contains(flags, "bmi2")
			f.HasFMA = strings.Contains(flags, "fma")
			f.HasLZCNT = strings.Contains(flags, "abm")
			f.HasPOPCNT = strings.Contains(flags, "popcnt")
		}

		// Get model name
		model, err := exec.Command("grep", "-m1", "model name", "/proc/cpuinfo").Output()
		if err == nil {
			parts := strings.SplitN(string(model), ":", 2)
			if len(parts) == 2 {
				f.ModelName = strings.TrimSpace(parts[1])
			}
		}
	} else if runtime.GOOS == "darwin" {
		// macOS - use sysctl
		out, err := exec.Command("sysctl", "-n", "machdep.cpu.features").Output()
		if err == nil {
			flags := strings.ToUpper(string(out))
			f.HasSSE2 = strings.Contains(flags, "SSE2")
			f.HasSSE4 = strings.Contains(flags, "SSE4")
			f.HasAVX = strings.Contains(flags, "AVX")
			f.HasAVX2 = strings.Contains(flags, "AVX2")
		}

		model, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
		if err == nil {
			f.ModelName = strings.TrimSpace(string(model))
		}
	}
}

func detectARMFeatures(f *CPUFeatures) {
	f.HasNEON = true // All ARM64 have NEON

	if runtime.GOOS == "linux" {
		data, err := exec.Command("grep", "-m1", "Features", "/proc/cpuinfo").Output()
		if err == nil {
			flags := string(data)
			f.HasAES = strings.Contains(flags, "aes")
			f.HasFMA = strings.Contains(flags, "fphp")
		}

		model, err := exec.Command("grep", "-m1", "CPU implementer", "/proc/cpuinfo").Output()
		if err == nil {
			f.ModelName = strings.TrimSpace(string(model))
		}
	}
}

func (f *CPUFeatures) GetOptimizationFlags(std string) []string {
	var flags []string

	// Base optimization flags
	switch {
	case f.HasAVX512:
		flags = append(flags, "-march=x86-64-v4")
	case f.HasAVX2:
		flags = append(flags, "-march=x86-64-v3")
	case f.HasAVX:
		flags = append(flags, "-march=x86-64-v2")
	case f.HasSSE4:
		flags = append(flags, "-msse4.2")
	case f.HasSSE2:
		// Default for x86_64
	}

	// ARM NEON is always available on ARM64
	if f.HasNEON && runtime.GOARCH == "arm64" {
		// ARM64 has NEON by default, no special flags needed
	}

	// Enable specific instruction sets
	if f.HasAES {
		flags = append(flags, "-maes")
	}
	if f.HasFMA {
		flags = append(flags, "-mfma")
	}
	if f.HasBMI {
		flags = append(flags, "-mbmi")
	}
	if f.HasBMI2 {
		flags = append(flags, "-mbmi2")
	}
	if f.HasLZCNT {
		flags = append(flags, "-mlzcnt")
	}
	if f.HasPOPCNT {
		flags = append(flags, "-mpopcnt")
	}

	return flags
}

func (f *CPUFeatures) GetLTOFlags() []string {
	return []string{"-flto=auto"}
}

func (f *CPUFeatures) GetPGOFlags() []string {
	return []string{"-fprofile-generate", "-fprofile-use"}
}

func (f *CPUFeatures) GetNativeFlags() []string {
	if runtime.GOARCH == "amd64" {
		return []string{"-march=native"}
	}
	return nil
}

func (f *CPUFeatures) String() string {
	var parts []string
	parts = append(parts, f.ModelName)
	parts = append(parts, runtime.GOARCH)
	parts = append(parts, "cores:"+itoa(f.CoreCount))

	if f.HasAVX512 {
		parts = append(parts, "AVX512")
	} else if f.HasAVX2 {
		parts = append(parts, "AVX2")
	} else if f.HasAVX {
		parts = append(parts, "AVX")
	}

	if f.HasNEON {
		parts = append(parts, "NEON")
	}

	return strings.Join(parts, " | ")
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

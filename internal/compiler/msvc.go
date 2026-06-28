package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type MSVC struct {
	path string
}

func NewMSVC(path string) *MSVC {
	return &MSVC{path: path}
}

func (m *MSVC) Name() string {
	return "msvc"
}

func (m *MSVC) Version() string {
	return "MSVC (version unknown)"
}

func (m *MSVC) Compile(opts CompileOptions) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("msvc only available on Windows")
	}

	args := []string{}

	if opts.CompileOnly {
		args = append(args, "/c")
	}

	if opts.Std != "" {
		std := strings.Replace(opts.Std, "c++", "/std:c++", 1)
		std = strings.Replace(std, "c1", "/std:c1", 1)
		args = append(args, std)
	}

	switch opts.Optimization {
	case "debug":
		args = append(args, "/Od", "/Zi")
	case "release":
		args = append(args, "/O2")
	case "size":
		args = append(args, "/O1")
	case "speed":
		args = append(args, "/O2", "/Ot")
	}

	for _, dir := range opts.IncludeDirs {
		args = append(args, "/I"+dir)
	}

	for _, def := range opts.Defines {
		args = append(args, "/D"+def)
	}

	args = append(args, opts.Flags...)
	args = append(args, "/Fo:"+opts.Output)
	args = append(args, opts.Sources...)

	cmd := exec.Command(m.path, args...)
	cmd.Stdout = nil

	if err := runCommand(cmd); err != nil {
		return fmt.Errorf("msvc compilation failed: %w", err)
	}

	return nil
}

func (m *MSVC) Link(opts LinkOptions) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("msvc only available on Windows")
	}

	args := []string{}
	args = append(args, opts.Objects...)
	args = append(args, "/Fe:"+opts.Output)

	for _, dir := range opts.LibDirs {
		args = append(args, "/LIBPATH:"+dir)
	}

	for _, lib := range opts.Libraries {
		args = append(args, lib+".lib")
	}

	args = append(args, opts.Flags...)

	cmd := exec.Command(m.path, args...)
	cmd.Stdout = nil

	if err := runCommand(cmd); err != nil {
		return fmt.Errorf("msvc linking failed: %w", err)
	}

	return nil
}

func (m *MSVC) Archive(objects []string, output string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("msvc only available on Windows")
	}

	// Use lib.exe to create static library
	args := []string{"/OUT:" + output}
	args = append(args, objects...)

	// Find lib.exe
	libPath := "lib.exe"
	if dir := filepath.Dir(m.path); dir != "" {
		candidate := filepath.Join(dir, "lib.exe")
		if _, err := os.Stat(candidate); err == nil {
			libPath = candidate
		}
	}

	cmd := exec.Command(libPath, args...)
	cmd.Stdout = nil

	if err := runCommand(cmd); err != nil {
		return fmt.Errorf("msvc archive failed: %w", err)
	}

	return nil
}

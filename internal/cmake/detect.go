package cmake

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func IsInstalled() bool {
	_, err := exec.LookPath("cmake")
	return err == nil
}

func Version() (string, error) {
	cmd := exec.Command("cmake", "--version")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get cmake version: %w", err)
	}

	output := strings.TrimSpace(stdout.String())
	lines := strings.Split(output, "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}

	return output, nil
}

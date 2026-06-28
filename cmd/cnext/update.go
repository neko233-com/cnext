package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var (
	updateCheck bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update cnext to latest version",
	Long: `Self-update cnext to the latest version.

Checks GitHub releases for new versions and downloads the update.

Examples:
  cnext update              # Update to latest
  cnext update --check      # Check for updates without installing
  cnext update --version v1.2.0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		currentVersion := "0.1.0" // This would be injected at build time
		latestVersion := ""

		// Check for latest version
		fmt.Println("Checking for updates...")
		resp, err := http.Get("https://api.github.com/repos/neko233-com/cnext/releases/latest")
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		// Simple JSON parsing for tag_name
		bodyStr := string(body)
		start := 0
		for i := 0; i < len(bodyStr)-10; i++ {
			if bodyStr[i:i+10] == `"tag_name":"` {
				start = i + 10
				break
			}
		}
		if start > 0 {
			end := start
			for end < len(bodyStr) && bodyStr[end] != '"' {
				end++
			}
			latestVersion = bodyStr[start:end]
		}

		if latestVersion == "" {
			return fmt.Errorf("failed to parse latest version")
		}

		if updateCheck {
			if latestVersion == "v"+currentVersion || latestVersion == currentVersion {
				fmt.Printf("✓ Already up to date (%s)\n", currentVersion)
			} else {
				fmt.Printf("Update available: %s → %s\n", currentVersion, latestVersion)
			}
			return nil
		}

		if latestVersion == "v"+currentVersion || latestVersion == currentVersion {
			fmt.Printf("✓ Already up to date (%s)\n", currentVersion)
			return nil
		}

		fmt.Printf("Updating %s → %s...\n", currentVersion, latestVersion)

		// Determine download URL
		osName := runtime.GOOS
		arch := runtime.GOARCH
		asset := fmt.Sprintf("cnext-%s-%s", osName, arch)
		if osName == "windows" {
			asset += ".exe"
		}

		url := fmt.Sprintf("https://github.com/neko233-com/cnext/releases/download/%s/%s", latestVersion, asset)

		// Download
		fmt.Printf("Downloading %s...\n", url)
		resp, err = http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to download update: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("download failed with status %d", resp.StatusCode)
		}

		// Get current executable path
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}
		exePath, _ = filepath.EvalSymlinks(exePath)

		// Create temporary file
		tmpFile, err := os.CreateTemp("", "cnext-update-*")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		defer os.Remove(tmpFile.Name())

		// Write to temp file
		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			tmpFile.Close()
			return fmt.Errorf("failed to write update: %w", err)
		}
		tmpFile.Close()

		// Make executable (Unix)
		if osName != "windows" {
			if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
				return fmt.Errorf("failed to make executable: %w", err)
			}
		}

		// Replace current executable
		if osName == "windows" {
			// Windows: rename current, then rename new
			backupPath := exePath + ".bak"
			os.Rename(exePath, backupPath)
			if err := os.Rename(tmpFile.Name(), exePath); err != nil {
				os.Rename(backupPath, exePath)
				return fmt.Errorf("failed to replace executable: %w", err)
			}
			os.Remove(backupPath)
		} else {
			// Unix: write to same location
			if err := os.Rename(tmpFile.Name(), exePath); err != nil {
				return fmt.Errorf("failed to replace executable: %w", err)
			}
		}

		fmt.Printf("✓ Updated to %s\n", latestVersion)
		return nil
	},
}

func init() {
	updateCmd.Flags().BoolVar(&updateCheck, "check", false, "Check for updates without installing")
}

func init() {
	// Suppress unused import
	_ = time.Now
}

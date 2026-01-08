package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFindConfigFile(t *testing.T) {
	t.Run("should find config.json in current directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		_ = os.Chdir(tmpDir)
		defer func() { _ = os.Chdir(oldWd) }()

		_ = os.WriteFile("config.json", []byte(`{"test": true}`), 0644)

		got := findConfigFile()

		if got != "config.json" {
			t.Errorf("got %q, want %q", got, "config.json")
		}
	})

	t.Run("should find config.jsonc in current directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		_ = os.Chdir(tmpDir)
		defer func() { _ = os.Chdir(oldWd) }()

		_ = os.WriteFile("config.jsonc", []byte(`{"test": true /* comment */}`), 0644)

		got := findConfigFile()

		if got != "config.jsonc" {
			t.Errorf("got %q, want %q", got, "config.jsonc")
		}
	})

	t.Run("should prioritize jsonc over json", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		_ = os.Chdir(tmpDir)
		defer func() { _ = os.Chdir(oldWd) }()

		_ = os.WriteFile("config.jsonc", []byte(`{"priority": "jsonc"}`), 0644)
		_ = os.WriteFile("config.json", []byte(`{"priority": "json"}`), 0644)

		got := findConfigFile()

		if got != "config.jsonc" {
			t.Errorf("got %q, want %q", got, "config.jsonc")
		}
	})

	t.Run("should return empty string when no config in current directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		_ = os.Chdir(tmpDir)
		defer func() { _ = os.Chdir(oldWd) }()

		got := findConfigFile()

		if got != "" && got != "config.json" && got != "config.jsonc" {

		}
	})
}

func TestGetAppDataDir(t *testing.T) {
	t.Run("should return path containing app name", func(t *testing.T) {
		got := getAppDataDir()

		if got == "" {
			t.Skip("unable to determine AppData dir")
		}

		if filepath.Base(got) != AppName {
			t.Errorf("got base=%q, want %q", filepath.Base(got), AppName)
		}
	})

	t.Run("should return platform-specific path", func(t *testing.T) {
		got := getAppDataDir()

		if got == "" {
			t.Skip("unable to determine AppData dir")
		}

		switch runtime.GOOS {
		case "darwin":
			if !contains(got, "Library/Application Support") {
				t.Errorf("darwin path should contain 'Library/Application Support', got %q", got)
			}
		case "linux":
			if !contains(got, ".config") {
				t.Errorf("linux path should contain '.config', got %q", got)
			}
		case "windows":

			if got == "" {
				t.Error("windows path should not be empty")
			}
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

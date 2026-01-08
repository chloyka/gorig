package config

import (
	"os"
	"path/filepath"
	"runtime"
)

var configFileNames = []string{"config.jsonc", "config.json"}

func findConfigFile() string {

	for _, name := range configFileNames {
		if _, err := os.Stat(name); err == nil {
			return name
		}
	}

	appDataDir := getAppDataDir()
	if appDataDir != "" {
		for _, name := range configFileNames {
			path := filepath.Join(appDataDir, name)
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	return ""
}

func getAppDataDir() string {
	var baseDir string

	switch runtime.GOOS {
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		baseDir = filepath.Join(homeDir, "Library", "Application Support")
	case "windows":
		baseDir = os.Getenv("APPDATA")
	default:
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		baseDir = filepath.Join(homeDir, ".config")
	}

	if baseDir == "" {
		return ""
	}

	return filepath.Join(baseDir, AppName)
}

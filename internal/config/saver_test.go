package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	configTypes "github.com/chloyka/gorig/internal/config/types"
	"github.com/chloyka/gorig/internal/logger"
	"go.uber.org/zap"
)

func TestGetDefaultSavePath(t *testing.T) {
	t.Run("should return path ending with config.json", func(t *testing.T) {
		got := getDefaultSavePath()

		if got == nil {
			t.Fatal("expected non-nil path")
		}
		if filepath.Base(got.String()) != "config.json" {
			t.Errorf("got base=%q, want %q", filepath.Base(got.String()), "config.json")
		}
	})

	t.Run("should return path containing app name when AppData available", func(t *testing.T) {
		got := getDefaultSavePath()

		if got == nil {
			t.Fatal("expected non-nil path")
		}

		appDataDir := getAppDataDir()
		if appDataDir != "" {
			if !containsAt(got.String(), AppName) {
				t.Errorf("path %q should contain app name %q", got.String(), AppName)
			}
		}
	})
}

func TestConfigManagerSave(t *testing.T) {
	t.Run("Save", func(t *testing.T) {
		t.Run("should save config as JSON", func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := configTypes.ConfigPath(filepath.Join(tmpDir, "config.json"))

			sut := &configManager{
				loadedPath: &configPath,
				audio: &configTypes.AudioConfig{
					SampleRate: 48000,
				},
				logCfg:  &configTypes.LoggerConfig{Level: "debug"},
				effects: &configTypes.EffectsConfig{EffectsDir: "./fx"},
				state:   &configTypes.StateConfig{InputDevice: "USB"},
				presets: &configTypes.PresetsConfig{},
				logger:  newTestLogger(),
			}

			err := sut.save()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			data, _ := os.ReadFile(configPath.String())
			var raw configTypes.RawConfig
			if err := json.Unmarshal(data, &raw); err != nil {
				t.Errorf("saved file is not valid JSON: %v", err)
			}
			if raw.Audio.SampleRate != 48000 {
				t.Errorf("got SampleRate=%d, want 48000", raw.Audio.SampleRate)
			}
		})

		t.Run("should create directory if not exists", func(t *testing.T) {
			tmpDir := t.TempDir()
			nestedDir := filepath.Join(tmpDir, "nested", "dir")
			configPath := configTypes.ConfigPath(filepath.Join(nestedDir, "config.json"))

			sut := &configManager{
				loadedPath: &configPath,
				audio:      &configTypes.AudioConfig{},
				logCfg:     &configTypes.LoggerConfig{},
				effects:    &configTypes.EffectsConfig{},
				state:      &configTypes.StateConfig{},
				presets:    &configTypes.PresetsConfig{},
				logger:     newTestLogger(),
			}

			err := sut.save()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
				t.Error("expected directory to be created")
			}
		})
	})

	t.Run("GetLoadedPath", func(t *testing.T) {
		t.Run("should return loaded path", func(t *testing.T) {
			configPath := configTypes.ConfigPath("/some/path/config.json")
			sut := &configManager{
				loadedPath: &configPath,
			}

			got := sut.getLoadedPath()

			if got == nil {
				t.Fatal("expected non-nil path")
			}
			if got.String() != "/some/path/config.json" {
				t.Errorf("got %q, want %q", got.String(), "/some/path/config.json")
			}
		})

		t.Run("should return nil when no path loaded", func(t *testing.T) {
			sut := &configManager{
				loadedPath: nil,
			}

			got := sut.getLoadedPath()

			if got != nil {
				t.Errorf("expected nil, got %q", got.String())
			}
		})
	})

	t.Run("Stop", func(t *testing.T) {
		t.Run("should not panic when cancel is nil", func(t *testing.T) {
			sut := &configManager{cancel: nil}

			sut.stop()
		})
	})
}

func newTestLogger() *logger.Logger {
	return &logger.Logger{Logger: zap.NewNop()}
}

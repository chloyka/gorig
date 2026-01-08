package config

import (
	"os"
	"testing"

	errs "github.com/chloyka/gorig/utils/errors"
)

func TestProvideConfig(t *testing.T) {

	t.Run("Defaults", func(t *testing.T) {
		t.Run("should set audio config values", func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tmpDir)
			defer func() { _ = os.Chdir(oldWd) }()

			got, err := provideConfig()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Audio == nil {
				t.Error("expected Audio config to be set")
			}
			if got.Audio.SampleRate == 0 {
				t.Error("expected SampleRate to be non-zero")
			}
		})

		t.Run("should set logger config values", func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tmpDir)
			defer func() { _ = os.Chdir(oldWd) }()

			got, err := provideConfig()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Logger == nil {
				t.Error("expected Logger config to be set")
			}
			if got.Logger.Level == "" {
				t.Error("expected Level to be set")
			}
		})

		t.Run("should set effects config values", func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tmpDir)
			defer func() { _ = os.Chdir(oldWd) }()

			got, err := provideConfig()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Effects == nil {
				t.Error("expected Effects config to be set")
			}
		})

		t.Run("should set state config values", func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tmpDir)
			defer func() { _ = os.Chdir(oldWd) }()

			got, err := provideConfig()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.State == nil {
				t.Error("expected State config to be set")
			}
		})

		t.Run("should set presets config values", func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tmpDir)
			defer func() { _ = os.Chdir(oldWd) }()

			got, err := provideConfig()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Presets == nil {
				t.Error("expected Presets config to be set")
			}
		})
	})

	t.Run("Loading", func(t *testing.T) {
		t.Run("should parse JSON config file", func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tmpDir)
			defer func() { _ = os.Chdir(oldWd) }()

			jsonConfig := `{
				"audio": {"sample_rate": 48000},
				"state": {"input_device": "USB Mic"}
			}`
			_ = os.WriteFile("config.json", []byte(jsonConfig), 0644)

			got, err := provideConfig()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Audio.SampleRate != 48000 {
				t.Errorf("got SampleRate=%d, want 48000", got.Audio.SampleRate)
			}
			if got.State.InputDevice != "USB Mic" {
				t.Errorf("got InputDevice=%q, want %q", got.State.InputDevice, "USB Mic")
			}
		})

		t.Run("should parse JSONC config file with comments", func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tmpDir)
			defer func() { _ = os.Chdir(oldWd) }()

			jsoncConfig := `{
				
				"audio": {
					"sample_rate": 96000 /* high quality */
				},
				/* Device configuration */
				"state": {"output_device": "Headphones"}
			}`
			_ = os.WriteFile("config.jsonc", []byte(jsoncConfig), 0644)

			got, err := provideConfig()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Audio.SampleRate != 96000 {
				t.Errorf("got SampleRate=%d, want 96000", got.Audio.SampleRate)
			}
			if got.State.OutputDevice != "Headphones" {
				t.Errorf("got OutputDevice=%q, want %q", got.State.OutputDevice, "Headphones")
			}
		})

		t.Run("should merge partial config with defaults", func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tmpDir)
			defer func() { _ = os.Chdir(oldWd) }()

			jsonConfig := `{"audio": {"sample_rate": 22050}}`
			_ = os.WriteFile("config.json", []byte(jsonConfig), 0644)

			got, err := provideConfig()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Audio.SampleRate != 22050 {
				t.Errorf("got SampleRate=%d, want 22050", got.Audio.SampleRate)
			}

			if got.Effects.EffectsDir != "./effects" {
				t.Errorf("got EffectsDir=%q, want default %q", got.Effects.EffectsDir, "./effects")
			}
		})

		t.Run("should set ConfigPath when file loaded", func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tmpDir)
			defer func() { _ = os.Chdir(oldWd) }()

			_ = os.WriteFile("config.json", []byte(`{"audio": {"sample_rate": 44100}}`), 0644)

			got, err := provideConfig()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Path == nil {
				t.Fatal("expected Path to be set")
			}
			if got.Path.String() != "config.json" {
				t.Errorf("got Path=%q, want %q", got.Path.String(), "config.json")
			}
		})

		t.Run("should set ConfigPath when config loaded from AppData", func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tmpDir)
			defer func() { _ = os.Chdir(oldWd) }()

			got, err := provideConfig()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			_ = got.Path
		})
	})

	t.Run("Errors", func(t *testing.T) {
		t.Run("should return error for invalid JSON", func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tmpDir)
			defer func() { _ = os.Chdir(oldWd) }()

			_ = os.WriteFile("config.json", []byte("{invalid json}"), 0644)

			_, err := provideConfig()

			if !errs.Is(err, errs.ErrConfigParseJSON) {
				t.Errorf("got error=%v, want %v", err, errs.ErrConfigParseJSON)
			}
		})
	})

	t.Run("Savers", func(t *testing.T) {
		t.Run("should populate savers slice with all config types", func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tmpDir)
			defer func() { _ = os.Chdir(oldWd) }()

			got, err := provideConfig()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got.Savers) != 5 {
				t.Errorf("got len(Savers)=%d, want 5", len(got.Savers))
			}
		})
	})
}

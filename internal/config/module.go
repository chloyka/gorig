package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	configTypes "github.com/chloyka/gorig/internal/config/types"
	errs "github.com/chloyka/gorig/utils/errors"
	"github.com/tidwall/jsonc"
	"go.uber.org/fx"
)

type AppConfig struct {
	fx.Out

	Audio   *configTypes.AudioConfig
	Logger  *configTypes.LoggerConfig
	Effects *configTypes.EffectsConfig

	State   *configTypes.StateConfig
	Presets *configTypes.PresetsConfig

	Savers []configTypes.ConfigSaver `group:"savers,flatten"`

	Path *configTypes.ConfigPath
}

func provideConfig() (AppConfig, error) {
	cfg := AppConfig{
		Audio: &configTypes.AudioConfig{
			SampleRate:      44100,
			FramesPerBuffer: 64,
			NumChannels:     1,
			TargetLatency:   10 * time.Millisecond,
		},
		State: &configTypes.StateConfig{
			InputDevice:    "",
			OutputDevice:   "",
			EffectsEnabled: true,
		},
		Logger: &configTypes.LoggerConfig{
			MaxLogFiles:   30,
			LogsDir:       "./logs",
			BufferSize:    256 * 1024,
			FlushInterval: time.Second * 5,
			Level:         "info",
		},
		Effects: &configTypes.EffectsConfig{
			EffectsDir: "./effects",
		},
		Presets: &configTypes.PresetsConfig{
			Presets:      []configTypes.Preset{},
			ActivePreset: "",
		},
	}

	cfg.Savers = []configTypes.ConfigSaver{
		cfg.Audio, cfg.Effects, cfg.State, cfg.Logger, cfg.Presets,
	}

	configPath := findConfigFile()
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return AppConfig{}, errs.Wrap(errs.ErrConfigRead, err)
		}

		ext := strings.ToLower(filepath.Ext(configPath))
		if ext == ".jsonc" {
			data = jsonc.ToJSON(data)
		}

		var raw configTypes.RawConfig
		if err := json.Unmarshal(data, &raw); err != nil {
			return AppConfig{}, errs.Wrap(errs.ErrConfigParseJSON, err)
		}

		if raw.Audio != nil {
			cfg.Audio.SampleRate = raw.Audio.SampleRate
			cfg.Audio.FramesPerBuffer = raw.Audio.FramesPerBuffer
			cfg.Audio.NumChannels = raw.Audio.NumChannels
			cfg.Audio.TargetLatency = raw.Audio.TargetLatency
		}

		if raw.Logger != nil {
			cfg.Logger.MaxLogFiles = raw.Logger.MaxLogFiles
			cfg.Logger.LogsDir = raw.Logger.LogsDir
			cfg.Logger.BufferSize = raw.Logger.BufferSize
			cfg.Logger.FlushInterval = raw.Logger.FlushInterval
			cfg.Logger.Level = raw.Logger.Level
		}

		if raw.Effects != nil {
			cfg.Effects.EffectsDir = raw.Effects.EffectsDir
		}

		if raw.State != nil {
			cfg.State.InputDevice = raw.State.InputDevice
			cfg.State.OutputDevice = raw.State.OutputDevice
			cfg.State.EffectsEnabled = raw.State.EffectsEnabled
		}

		if raw.Presets != nil {
			cfg.Presets.Presets = raw.Presets.Presets
			cfg.Presets.ActivePreset = raw.Presets.ActivePreset
		}

		path := configTypes.ConfigPath(configPath)
		cfg.Path = &path
	}

	return cfg, nil
}

func ProvideConfig() fx.Option {
	return fx.Provide(provideConfig)
}

func ProvideManager() fx.Option {
	return fx.Provide(provideConfigManager)
}

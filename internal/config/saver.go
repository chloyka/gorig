package config

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	configTypes "github.com/chloyka/gorig/internal/config/types"
	"github.com/chloyka/gorig/internal/logger"
	"github.com/chloyka/gorig/internal/logger/keys"
	"github.com/chloyka/gorig/utils"
	errs "github.com/chloyka/gorig/utils/errors"
	"go.uber.org/fx"
)

type configManager struct {
	mu         sync.RWMutex
	loadedPath *configTypes.ConfigPath
	saveChan   chan struct{}
	logger     *logger.Logger
	cancel     context.CancelFunc

	audio   *configTypes.AudioConfig
	logCfg  *configTypes.LoggerConfig
	effects *configTypes.EffectsConfig
	state   *configTypes.StateConfig
	presets *configTypes.PresetsConfig
}

type newConfigManagerParams struct {
	fx.In

	Logger     *logger.Logger
	ConfigPath *configTypes.ConfigPath
	Configs    []configTypes.ConfigSaver `group:"savers"`

	Audio   *configTypes.AudioConfig
	LogCfg  *configTypes.LoggerConfig
	Effects *configTypes.EffectsConfig
	State   *configTypes.StateConfig
	Presets *configTypes.PresetsConfig
}

func provideConfigManager(in newConfigManagerParams) *configManager {
	saveChan := make(chan struct{}, 1)

	for _, config := range in.Configs {
		config.SetSaveChan(saveChan)
	}

	ctx, cancel := context.WithCancel(context.Background())

	m := &configManager{
		loadedPath: in.ConfigPath,
		saveChan:   saveChan,
		logger:     in.Logger,
		cancel:     cancel,
		audio:      in.Audio,
		logCfg:     in.LogCfg,
		effects:    in.Effects,
		state:      in.State,
		presets:    in.Presets,
	}

	loadedPathStr := ""
	if in.ConfigPath != nil {
		loadedPathStr = in.ConfigPath.String()
	}
	in.Logger.Info("configManager created",
		keys.PathConfigLoaded(loadedPathStr),
	)

	go m.saveLoop(ctx)

	return m
}

func (m *configManager) saveLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			m.logger.Debug("save loop stopped")
			return
		case <-m.saveChan:
			m.logger.Debug("received save signal")
			if err := m.save(); err != nil {
				m.logger.Error("failed to save config", keys.Error(err))
			}
		}
	}
}

func (m *configManager) stop() {
	if m.cancel != nil {
		m.cancel()
	}
}

func (m *configManager) getLoadedPath() *configTypes.ConfigPath {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.loadedPath
}

func (m *configManager) save() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	savePath := m.loadedPath

	if savePath == nil || savePath.String() == "" {
		savePath = getDefaultSavePath()
		m.logger.Debug("no config file loaded, using default path", keys.PathConfig(savePath.String()))
	}

	m.logger.Debug("saving config", keys.PathConfig(savePath.String()))

	dir := filepath.Dir(savePath.String())
	if err := os.MkdirAll(dir, 0755); err != nil {
		m.logger.Error("failed to create config directory", keys.Error(err), keys.PathDir(dir))
		return errs.Wrap(errs.ErrConfigCreateDir, err)
	}

	rawConfig := &configTypes.RawConfig{
		Audio:   m.audio,
		Logger:  m.logCfg,
		Effects: m.effects,
		State:   m.state,
		Presets: m.presets,
	}

	data, err := json.MarshalIndent(rawConfig, "", "  ")
	if err != nil {
		m.logger.Error("failed to marshal config", keys.Error(err))
		return errs.Wrap(errs.ErrConfigMarshal, err)
	}

	if err = os.WriteFile(savePath.String(), data, 0644); err != nil {
		m.logger.Error("failed to write config file", keys.Error(err), keys.PathConfig(savePath.String()))
		return errs.Wrap(errs.ErrConfigWrite, err)
	}

	m.loadedPath = savePath

	m.logger.Info("config saved successfully",
		keys.PathConfig(savePath.String()),
		keys.DeviceInputName(m.state.InputDevice),
		keys.DeviceOutputName(m.state.OutputDevice),
		keys.EffectEnabled(m.state.EffectsEnabled),
	)

	return nil
}

func getDefaultSavePath() *configTypes.ConfigPath {
	appDataDir := getAppDataDir()
	if appDataDir == "" {
		return utils.Ptr(configTypes.ConfigPath("config.json"))
	}

	return utils.Ptr(configTypes.ConfigPath(filepath.Join(appDataDir, "config.json")))
}

package effects

import (
	"sync"

	configTypes "github.com/chloyka/gorig/internal/config/types"
	"github.com/chloyka/gorig/internal/logger"
	"github.com/chloyka/gorig/internal/logger/keys"
)

type Chain struct {
	mu            sync.RWMutex
	registry      *EffectRegistry
	activeChain   []*InterpretedEffect
	effectsDir    string
	logger        *logger.Logger
	enabled       bool
	stateConfig   *configTypes.StateConfig
	presetsConfig *configTypes.PresetsConfig
}

func NewChain(log *logger.Logger, effectsDir string, stateConfig *configTypes.StateConfig, presetsConfig *configTypes.PresetsConfig) *Chain {
	c := &Chain{
		effectsDir:    effectsDir,
		logger:        log,
		enabled:       stateConfig.EffectsEnabled,
		stateConfig:   stateConfig,
		presetsConfig: presetsConfig,
	}

	if err := c.reload(); err != nil {
		log.Error("failed to load effects", keys.Error(err))
	}

	c.applyActivePreset()

	return c
}

func (c *Chain) reload() error {
	registry, err := loadEffectsFromDirRecursive(c.effectsDir)
	if err != nil {
		return err
	}

	c.registry = registry

	names := registry.GetAvailableEffectNames()
	c.logger.Info("effects loaded",
		keys.EffectList(names),
	)

	return nil
}

func (c *Chain) applyActivePreset() {
	preset := c.presetsConfig.GetActivePresetConfig()
	if preset == nil {
		c.activeChain = nil
		c.logger.Debug("no active preset, chain empty")
		return
	}

	c.activeChain = make([]*InterpretedEffect, 0)
	var missingEffects []string

	for _, effectName := range preset.EffectChain {
		effect := c.registry.GetEffect(effectName)
		if effect != nil {
			c.activeChain = append(c.activeChain, effect)
		} else {
			missingEffects = append(missingEffects, effectName)
		}
	}

	if len(missingEffects) > 0 {
		c.logger.Warn("preset has missing effects",
			keys.EffectName(preset.Name),
			keys.EffectList(missingEffects),
		)
	}

	c.logger.Info("preset chain applied",
		keys.EffectName(preset.Name),
	)
}

func (c *Chain) SetPresetChain(effectNames []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.activeChain = make([]*InterpretedEffect, 0)

	for _, name := range effectNames {
		effect := c.registry.GetEffect(name)
		if effect != nil {
			c.activeChain = append(c.activeChain, effect)
		}
	}

	c.logger.Debug("chain updated from preset")
}

func (c *Chain) Reload() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger.Debug("reloading effects from disk")
	if err := c.reload(); err != nil {
		return err
	}

	c.applyActivePreset()
	return nil
}

func (c *Chain) Process(samples []float32) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.enabled {
		return
	}

	for _, effect := range c.activeChain {
		effect.Process(samples)
	}
}

func (c *Chain) GetAvailableEffectNames() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.registry == nil {
		return nil
	}
	return c.registry.GetAvailableEffectNames()
}

func (c *Chain) HasActiveEffects() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.activeChain) > 0
}

type EffectInfo struct {
	Name      string
	Available bool
}

func (c *Chain) GetActiveChainInfo() []EffectInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var infos []EffectInfo
	for _, effect := range c.activeChain {
		infos = append(infos, EffectInfo{
			Name:      effect.Name(),
			Available: true,
		})
	}
	return infos
}

func (c *Chain) GetEffects() []EffectInfo {
	return c.GetActiveChainInfo()
}

func (c *Chain) ToggleChain() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.enabled = !c.enabled
	c.logger.Debug("chain toggled", keys.EffectEnabled(c.enabled))

	c.stateConfig.SetEffectsEnabled(c.enabled)

	return c.enabled
}

func (c *Chain) IsChainEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.enabled
}

func (c *Chain) GetRegistry() *EffectRegistry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.registry
}

func (c *Chain) GetMissingEffectsForPreset(preset *configTypes.Preset) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if preset == nil || c.registry == nil {
		return nil
	}

	var missing []string
	for _, name := range preset.EffectChain {
		if c.registry.GetEffect(name) == nil {
			missing = append(missing, name)
		}
	}
	return missing
}

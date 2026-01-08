package pedal

import (
	"github.com/chloyka/gorig/internal/effects"
	"github.com/chloyka/gorig/internal/logger"
	"github.com/chloyka/gorig/internal/logger/keys"
	"github.com/chloyka/gorig/internal/preset"
)

type State struct {
	chain         *effects.Chain
	presetManager *preset.Manager
	logger        *logger.Logger
}

func New(logger *logger.Logger, chain *effects.Chain, presetManager *preset.Manager) *State {
	return &State{
		chain:         chain,
		presetManager: presetManager,
		logger:        logger,
	}
}

func (s *State) ToggleEffects() bool {
	wasEnabled := s.chain.IsChainEnabled()

	if !wasEnabled {
		if !s.presetManager.HasPresets() {

			created, err := s.presetManager.EnsureActivePreset()
			if err != nil {
				s.logger.Error("failed to create default preset", keys.Error(err))
			}

			s.chain.ToggleChain()
			s.logger.Info("effects enabled, created default preset")
			return created
		}

		status := s.presetManager.GetActivePresetStatus()
		if status != nil && status.MissingCount > 0 {
			s.logger.Warn("active preset has missing effects")
		}
	}

	enabled := s.chain.ToggleChain()
	s.logger.Info("effects toggled", keys.EffectEnabled(enabled))
	return false
}

func (s *State) ReloadEffects() {
	if err := s.chain.Reload(); err != nil {
		s.logger.Error("failed to reload effects", keys.Error(err))
		return
	}
	s.logger.Info("effects reloaded", keys.EffectHasEnabled(s.chain.HasActiveEffects()))
}

func (s *State) IsEffectsOn() bool {
	return s.chain.IsChainEnabled()
}

func (s *State) GetEffects() []effects.EffectInfo {
	return s.chain.GetEffects()
}

func (s *State) GetPresetManager() *preset.Manager {
	return s.presetManager
}

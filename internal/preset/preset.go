package preset

import (
	"sync"

	configTypes "github.com/chloyka/gorig/internal/config/types"
	"github.com/chloyka/gorig/internal/logger"
	errs "github.com/chloyka/gorig/utils/errors"
)

type EffectStatus struct {
	Name      string
	Available bool
}

type PresetStatus struct {
	Name           string
	EffectChain    []EffectStatus
	MissingCount   int
	MissingEffects []string
}

type Manager struct {
	mu            sync.RWMutex
	presetsConfig *configTypes.PresetsConfig
	logger        *logger.Logger

	getAvailableEffects func() []string

	onPresetChanged func(chain []string)
}

func NewManager(
	logger *logger.Logger,
	presetsConfig *configTypes.PresetsConfig,
	getAvailableEffects func() []string,
	onPresetChanged func(chain []string),
) *Manager {
	return &Manager{
		presetsConfig:       presetsConfig,
		logger:              logger,
		getAvailableEffects: getAvailableEffects,
		onPresetChanged:     onPresetChanged,
	}
}

func (m *Manager) SetCallbacks(getAvailableEffects func() []string, onPresetChanged func(chain []string)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.getAvailableEffects = getAvailableEffects
	m.onPresetChanged = onPresetChanged
}

func (m *Manager) GetActivePreset() *configTypes.Preset {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.presetsConfig.GetActivePresetConfig()
}

func (m *Manager) GetActivePresetStatus() *PresetStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	preset := m.presetsConfig.GetActivePresetConfig()
	if preset == nil {
		return nil
	}
	return m.getPresetStatusUnlocked(*preset)
}

func (m *Manager) SetActivePreset(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	preset := m.presetsConfig.GetPreset(name)
	if preset == nil {
		return errs.Wrap(errs.ErrPresetNotFound, name)
	}

	m.presetsConfig.SetActivePreset(name)

	if m.onPresetChanged != nil {
		m.onPresetChanged(preset.EffectChain)
	}

	m.logger.Info("active preset changed")
	return nil
}

func (m *Manager) CreatePreset(name string, chain []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.presetsConfig.GetPreset(name) != nil {
		return errs.Wrap(errs.ErrPresetExists, name)
	}

	preset := configTypes.Preset{
		Name:        name,
		EffectChain: chain,
	}

	m.presetsConfig.AddPreset(preset)
	m.logger.Info("preset created")
	return nil
}

func (m *Manager) CreateDefaultPreset() error {
	return m.CreatePreset("Default", []string{})
}

func (m *Manager) EnsureActivePreset() (created bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.presetsConfig.Presets) == 0 {

		preset := configTypes.Preset{
			Name:        "Default",
			EffectChain: []string{},
		}
		m.presetsConfig.AddPreset(preset)
		m.presetsConfig.SetActivePreset("Default")
		m.logger.Info("created default preset")
		return true, nil
	}

	if m.presetsConfig.ActivePreset == "" {

		m.presetsConfig.SetActivePreset(m.presetsConfig.Presets[0].Name)
	}

	return false, nil
}

func (m *Manager) GetAllPresetsStatus() []PresetStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var statuses []PresetStatus
	for _, preset := range m.presetsConfig.Presets {
		statuses = append(statuses, *m.getPresetStatusUnlocked(preset))
	}
	return statuses
}

func (m *Manager) GetPresetNames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, len(m.presetsConfig.Presets))
	for i, p := range m.presetsConfig.Presets {
		names[i] = p.Name
	}
	return names
}

func (m *Manager) getPresetStatusUnlocked(preset configTypes.Preset) *PresetStatus {
	var available []string
	if m.getAvailableEffects != nil {
		available = m.getAvailableEffects()
	}

	availableSet := make(map[string]bool)
	for _, name := range available {
		availableSet[name] = true
	}

	status := &PresetStatus{
		Name:        preset.Name,
		EffectChain: make([]EffectStatus, len(preset.EffectChain)),
	}

	for i, effectName := range preset.EffectChain {
		isAvailable := availableSet[effectName]
		status.EffectChain[i] = EffectStatus{
			Name:      effectName,
			Available: isAvailable,
		}
		if !isAvailable {
			status.MissingCount++
			status.MissingEffects = append(status.MissingEffects, effectName)
		}
	}

	return status
}

func (m *Manager) UpdatePresetChain(presetName string, chain []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.presetsConfig.UpdatePreset(presetName, chain) {
		return errs.Wrap(errs.ErrPresetNotFound, presetName)
	}

	if presetName == m.presetsConfig.ActivePreset && m.onPresetChanged != nil {
		m.onPresetChanged(chain)
	}

	return nil
}

func (m *Manager) AddEffectToPreset(presetName, effectName string, position int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	preset := m.presetsConfig.GetPreset(presetName)
	if preset == nil {
		return errs.Wrap(errs.ErrPresetNotFound, presetName)
	}

	newChain := make([]string, 0, len(preset.EffectChain)+1)
	if position < 0 || position >= len(preset.EffectChain) {

		newChain = append(preset.EffectChain, effectName)
	} else {

		newChain = append(newChain, preset.EffectChain[:position]...)
		newChain = append(newChain, effectName)
		newChain = append(newChain, preset.EffectChain[position:]...)
	}

	m.presetsConfig.UpdatePreset(presetName, newChain)

	if presetName == m.presetsConfig.ActivePreset && m.onPresetChanged != nil {
		m.onPresetChanged(newChain)
	}

	return nil
}

func (m *Manager) RemoveEffectFromPreset(presetName, effectName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	preset := m.presetsConfig.GetPreset(presetName)
	if preset == nil {
		return errs.Wrap(errs.ErrPresetNotFound, presetName)
	}

	newChain := make([]string, 0)
	for _, name := range preset.EffectChain {
		if name != effectName {
			newChain = append(newChain, name)
		}
	}

	m.presetsConfig.UpdatePreset(presetName, newChain)

	if presetName == m.presetsConfig.ActivePreset && m.onPresetChanged != nil {
		m.onPresetChanged(newChain)
	}

	return nil
}

func (m *Manager) ReorderEffectInPreset(presetName string, fromIdx, toIdx int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	preset := m.presetsConfig.GetPreset(presetName)
	if preset == nil {
		return errs.Wrap(errs.ErrPresetNotFound, presetName)
	}

	if fromIdx < 0 || fromIdx >= len(preset.EffectChain) ||
		toIdx < 0 || toIdx >= len(preset.EffectChain) {
		return errs.Wrap(errs.ErrPresetInvalidIndex, []int{fromIdx, toIdx, len(preset.EffectChain)})
	}

	newChain := make([]string, len(preset.EffectChain))
	copy(newChain, preset.EffectChain)

	effect := newChain[fromIdx]
	newChain = append(newChain[:fromIdx], newChain[fromIdx+1:]...)

	if toIdx > fromIdx {
		toIdx--
	}
	newChain = append(newChain[:toIdx], append([]string{effect}, newChain[toIdx:]...)...)

	m.presetsConfig.UpdatePreset(presetName, newChain)

	if presetName == m.presetsConfig.ActivePreset && m.onPresetChanged != nil {
		m.onPresetChanged(newChain)
	}

	return nil
}

func (m *Manager) DeletePreset(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	wasActive := m.presetsConfig.ActivePreset == name

	if !m.presetsConfig.DeletePreset(name) {
		return errs.Wrap(errs.ErrPresetNotFound, name)
	}

	if wasActive {
		if len(m.presetsConfig.Presets) > 0 {
			newActive := m.presetsConfig.Presets[0].Name
			m.presetsConfig.SetActivePreset(newActive)
			if m.onPresetChanged != nil {
				m.onPresetChanged(m.presetsConfig.Presets[0].EffectChain)
			}
		} else {
			m.presetsConfig.ActivePreset = ""
			if m.onPresetChanged != nil {
				m.onPresetChanged(nil)
			}
		}
	}

	return nil
}

func (m *Manager) GetAvailableEffects() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.getAvailableEffects == nil {
		return nil
	}
	return m.getAvailableEffects()
}

func (m *Manager) HasPresets() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.presetsConfig.Presets) > 0
}

func (m *Manager) GetActivePresetName() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.presetsConfig.ActivePreset
}

func (m *Manager) GetPreset(name string) *configTypes.Preset {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.presetsConfig.GetPreset(name)
}

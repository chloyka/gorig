package configTypes

type Preset struct {
	Name        string   `json:"name" yaml:"name"`
	EffectChain []string `json:"effect_chain" yaml:"effect_chain"`
}

type PresetsConfig struct {
	configSaver

	Presets      []Preset `json:"presets" yaml:"presets"`
	ActivePreset string   `json:"active_preset" yaml:"active_preset"`
}

func (p *PresetsConfig) SetActivePreset(name string) {
	p.ActivePreset = name
	p.Save()
}

func (p *PresetsConfig) AddPreset(preset Preset) {
	p.Presets = append(p.Presets, preset)
	p.Save()
}

func (p *PresetsConfig) UpdatePreset(name string, chain []string) bool {
	for i, preset := range p.Presets {
		if preset.Name == name {
			p.Presets[i].EffectChain = chain
			p.Save()
			return true
		}
	}
	return false
}

func (p *PresetsConfig) DeletePreset(name string) bool {
	for i, preset := range p.Presets {
		if preset.Name == name {
			p.Presets = append(p.Presets[:i], p.Presets[i+1:]...)
			p.Save()
			return true
		}
	}
	return false
}

func (p *PresetsConfig) GetPreset(name string) *Preset {
	for i := range p.Presets {
		if p.Presets[i].Name == name {
			return &p.Presets[i]
		}
	}
	return nil
}

func (p *PresetsConfig) GetActivePresetConfig() *Preset {
	return p.GetPreset(p.ActivePreset)
}

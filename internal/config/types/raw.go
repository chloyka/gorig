package configTypes

type RawConfig struct {
	Audio   *AudioConfig   `json:"audio" yaml:"audio"`
	Logger  *LoggerConfig  `json:"logger" yaml:"logger"`
	Effects *EffectsConfig `json:"effects" yaml:"effects"`
	State   *StateConfig   `json:"state" yaml:"state"`
	Presets *PresetsConfig `json:"presets" yaml:"presets"`
}

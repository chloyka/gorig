package configTypes

type EffectsConfig struct {
	configSaver

	EffectsDir string `json:"effects_dir" yaml:"effects_dir"`
}

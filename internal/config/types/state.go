package configTypes

type StateConfig struct {
	configSaver

	InputDevice    string `json:"input_device" yaml:"input_device"`
	OutputDevice   string `json:"output_device" yaml:"output_device"`
	EffectsEnabled bool   `json:"effects_enabled" yaml:"effects_enabled"`

	RhythmBPM         float64 `json:"rhythm_bpm" yaml:"rhythm_bpm"`
	RhythmSubdivision int     `json:"rhythm_subdivision" yaml:"rhythm_subdivision"`
}

func (s *StateConfig) SetInputDevice(name string) {
	s.InputDevice = name
	s.Save()
}

func (s *StateConfig) SetOutputDevice(name string) {
	s.OutputDevice = name
	s.Save()
}

func (s *StateConfig) SetEffectsEnabled(enabled bool) {
	s.EffectsEnabled = enabled
	s.Save()
}

func (s *StateConfig) SetRhythmBPM(bpm float64) {
	s.RhythmBPM = bpm
	s.Save()
}

func (s *StateConfig) SetRhythmSubdivision(subdivision int) {
	s.RhythmSubdivision = subdivision
	s.Save()
}

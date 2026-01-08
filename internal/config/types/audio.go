package configTypes

import "time"

type AudioConfig struct {
	configSaver

	SampleRate      int           `json:"sample_rate" yaml:"sample_rate"`
	FramesPerBuffer int           `json:"frames_per_buffer" yaml:"frames_per_buffer"`
	NumChannels     int           `json:"num_channels" yaml:"num_channels"`
	TargetLatency   time.Duration `json:"target_latency" yaml:"target_latency"`
}

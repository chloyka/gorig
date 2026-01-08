package configTypes

import "time"

type LoggerConfig struct {
	configSaver

	MaxLogFiles   int           `json:"max_log_files" yaml:"max_log_files"`
	LogsDir       string        `json:"logs_dir" yaml:"logs_dir"`
	BufferSize    int           `json:"buffer_size" yaml:"buffer_size"`
	FlushInterval time.Duration `json:"flush_interval" yaml:"flush_interval"`
	Level         string        `json:"level" yaml:"level"`
}

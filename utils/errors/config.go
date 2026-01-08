package errors

var (
	ErrConfigRead              = New("config: failed to read file")
	ErrConfigParseJSON         = New("config: failed to parse JSON")
	ErrConfigParseYAML         = New("config: failed to parse YAML")
	ErrConfigUnsupportedFormat = New("config: unsupported format")
	ErrConfigCreateDir         = New("config: failed to create directory")
	ErrConfigMarshal           = New("config: failed to marshal")
	ErrConfigWrite             = New("config: failed to write file")
)

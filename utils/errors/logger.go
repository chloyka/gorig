package errors

var (
	ErrLoggerCreateDir = New("logger: failed to create directory")
	ErrLoggerOpenFile  = New("logger: failed to open file")
)

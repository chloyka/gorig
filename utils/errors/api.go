package errors

var (
	ErrAPIEffectNotFound = New("api: effect not found")
	ErrAPIPresetNotFound = New("api: preset not found")
	ErrAPINotImplemented = New("api: not implemented")
)

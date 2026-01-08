package errors

var (
	ErrPresetNotFound     = New("preset: not found")
	ErrPresetExists       = New("preset: already exists")
	ErrPresetInvalidIndex = New("preset: invalid index")
)

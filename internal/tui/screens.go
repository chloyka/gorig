package tui

type Screen int

const (
	ScreenMain Screen = iota
	ScreenPresetList
	ScreenPresetCreate
	ScreenPresetEdit
	ScreenEffectAdd
)

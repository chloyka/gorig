package tui

import (
	"github.com/chloyka/gorig/internal/keyboard"
)

const (
	ActionQuit      = "quit"
	ActionToggle    = "toggle"
	ActionReload    = "reload"
	ActionPresets   = "presets"
	ActionInput     = "input"
	ActionOutput    = "output"
	ActionUp        = "up"
	ActionDown      = "down"
	ActionEnter     = "enter"
	ActionSpace     = "space"
	ActionEsc       = "esc"
	ActionTab       = "tab"
	ActionNew       = "new"
	ActionEdit      = "edit"
	ActionDelete    = "delete"
	ActionSave      = "save"
	ActionAdd       = "add"
	ActionMoveUp    = "moveUp"
	ActionMoveDown  = "moveDown"
	ActionBackspace = "backspace"

	ActionTapTempo      = "tapTempo"
	ActionBpmUp         = "bpmUp"
	ActionBpmDown       = "bpmDown"
	ActionSubdivisionUp = "subdivisionUp"
	ActionSubdivisionDn = "subdivisionDn"
)

var keyMap = map[string][]string{

	ActionQuit:    {"q", "ctrl+c"},
	ActionToggle:  {"d"},
	ActionReload:  {"r"},
	ActionPresets: {"p"},
	ActionInput:   {"i"},
	ActionOutput:  {"o"},

	ActionUp:    {"up", "k"},
	ActionDown:  {"down", "j"},
	ActionEnter: {"enter"},
	ActionSpace: {" "},
	ActionEsc:   {"esc"},
	ActionTab:   {"tab"},

	ActionNew:    {"n"},
	ActionEdit:   {"e"},
	ActionDelete: {"x"},
	ActionSave:   {"s"},
	ActionAdd:    {"a"},

	ActionMoveUp:   {"K"},
	ActionMoveDown: {"J"},

	ActionBackspace: {"backspace"},

	ActionTapTempo:      {"t"},
	ActionBpmUp:         {".", ">"},
	ActionBpmDown:       {",", "<"},
	ActionSubdivisionUp: {"]", "}"},
	ActionSubdivisionDn: {"[", "{"},
}

func MatchKey(key, action string) bool {

	normalizedKey := keyboard.Normalize(key)

	keys, ok := keyMap[action]
	if !ok {
		return false
	}
	for _, k := range keys {
		if k == normalizedKey || k == key {
			return true
		}
	}
	return false
}

func MatchAnyKey(key string, actions ...string) bool {
	for _, action := range actions {
		if MatchKey(key, action) {
			return true
		}
	}
	return false
}

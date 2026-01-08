package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chloyka/gorig/internal/preset"
)

type presetCreateModel struct {
	name             string
	availableEffects []string
	selectedEffects  []bool
	cursor           int
	focusOnName      bool
	presetManager    *preset.Manager
}

func newPresetCreateModel(pm *preset.Manager) presetCreateModel {
	effects := pm.GetAvailableEffects()
	return presetCreateModel{
		name:             "",
		availableEffects: effects,
		selectedEffects:  make([]bool, len(effects)),
		cursor:           0,
		focusOnName:      true,
		presetManager:    pm,
	}
}

func (m presetCreateModel) Update(msg tea.Msg) (presetCreateModel, tea.Cmd, Screen) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		if m.focusOnName {
			switch {
			case MatchKey(key, "enter"), MatchKey(key, "tab"):
				if m.name != "" {
					m.focusOnName = false
				}
			case MatchKey(key, "backspace"):
				if len(m.name) > 0 {
					m.name = m.name[:len(m.name)-1]
				}
			case MatchKey(key, "esc"):
				return m, nil, ScreenPresetList
			default:

				if len(key) == 1 {
					m.name += key
				}
			}
		} else {
			switch {
			case MatchKey(key, "up"):
				if m.cursor > 0 {
					m.cursor--
				}
			case MatchKey(key, "down"):
				if m.cursor < len(m.availableEffects)-1 {
					m.cursor++
				}
			case MatchKey(key, "space"), MatchKey(key, "enter"):
				if m.cursor < len(m.selectedEffects) {
					m.selectedEffects[m.cursor] = !m.selectedEffects[m.cursor]
				}
			case MatchKey(key, "tab"):
				m.focusOnName = true
			case MatchKey(key, "save"):

				var chain []string
				for i, selected := range m.selectedEffects {
					if selected {
						chain = append(chain, m.availableEffects[i])
					}
				}
				if m.name != "" {
					m.presetManager.CreatePreset(m.name, chain)
					m.presetManager.SetActivePreset(m.name)
				}
				return m, nil, ScreenPresetList
			case MatchKey(key, "esc"):
				return m, nil, ScreenPresetList
			}
		}
	}
	return m, nil, ScreenPresetCreate
}

func (m presetCreateModel) View() string {
	var b strings.Builder

	b.WriteString("\n Create New Preset\n")
	b.WriteString(" ==================\n\n")

	nameFocus := ""
	if m.focusOnName {
		nameFocus = "_"
	}
	b.WriteString(fmt.Sprintf(" Name: %s%s\n\n", m.name, nameFocus))

	if len(m.availableEffects) == 0 {
		b.WriteString(" No effects available.\n")
		b.WriteString(" Add .go effect files to your effects directory.\n")
	} else {
		b.WriteString(" Select Effects:\n")
		for i, effect := range m.availableEffects {
			cursor := "  "
			if i == m.cursor && !m.focusOnName {
				cursor = "> "
			}

			check := "[ ]"
			if m.selectedEffects[i] {
				check = "[x]"
			}

			b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, check, effect))
		}
	}

	b.WriteString("\n [tab] Switch focus  [space] Toggle  [s] Save  [esc] Cancel\n")

	return b.String()
}

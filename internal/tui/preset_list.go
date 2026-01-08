package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chloyka/gorig/internal/preset"
)

type presetListModel struct {
	presets       []preset.PresetStatus
	cursor        int
	activePreset  string
	presetManager *preset.Manager
}

func newPresetListModel(pm *preset.Manager) presetListModel {
	statuses := pm.GetAllPresetsStatus()
	active := pm.GetActivePresetName()

	return presetListModel{
		presets:       statuses,
		cursor:        0,
		activePreset:  active,
		presetManager: pm,
	}
}

func (m presetListModel) Update(msg tea.Msg) (presetListModel, tea.Cmd, Screen, string) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch {
		case MatchKey(key, "up"):
			if m.cursor > 0 {
				m.cursor--
			}
		case MatchKey(key, "down"):
			if m.cursor < len(m.presets)-1 {
				m.cursor++
			}
		case MatchKey(key, "enter"):

			if len(m.presets) > 0 {
				selected := m.presets[m.cursor]
				m.presetManager.SetActivePreset(selected.Name)
				m.activePreset = selected.Name
				return m, nil, ScreenMain, ""
			}
		case MatchKey(key, "new"):
			return m, nil, ScreenPresetCreate, ""
		case MatchKey(key, "edit"):
			if len(m.presets) > 0 {
				return m, nil, ScreenPresetEdit, m.presets[m.cursor].Name
			}
		case MatchKey(key, "delete"):
			if len(m.presets) > 0 {
				selected := m.presets[m.cursor]
				m.presetManager.DeletePreset(selected.Name)
				m.presets = m.presetManager.GetAllPresetsStatus()
				m.activePreset = m.presetManager.GetActivePresetName()
				if m.cursor >= len(m.presets) && m.cursor > 0 {
					m.cursor--
				}
			}
		case MatchKey(key, "esc"):
			return m, nil, ScreenMain, ""
		}
	}
	return m, nil, ScreenPresetList, ""
}

func (m presetListModel) View() string {
	var b strings.Builder

	b.WriteString("\n Presets\n")
	b.WriteString(" ========\n\n")

	if len(m.presets) == 0 {
		b.WriteString(" No presets found. Press [n] to create one.\n")
	} else {
		for i, p := range m.presets {
			cursor := "  "
			if i == m.cursor {
				cursor = "> "
			}

			active := ""
			if p.Name == m.activePreset {
				active = " *"
			}

			missing := ""
			if p.MissingCount > 0 {
				missing = fmt.Sprintf(" (%d missing)", p.MissingCount)
			}

			effectCount := fmt.Sprintf(" [%d effects]", len(p.EffectChain))

			b.WriteString(fmt.Sprintf("%s%s%s%s%s\n", cursor, p.Name, active, effectCount, missing))
		}
	}

	if len(m.presets) > 0 && m.cursor < len(m.presets) {
		selected := m.presets[m.cursor]
		if selected.MissingCount > 0 {
			b.WriteString(fmt.Sprintf("\n Missing effects in '%s':\n", selected.Name))
			for _, name := range selected.MissingEffects {
				b.WriteString(fmt.Sprintf("   - %s\n", name))
			}
		}
	}

	b.WriteString("\n [enter] Select  [n] New  [e] Edit  [x] Delete  [esc] Back\n")

	return b.String()
}

package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chloyka/gorig/internal/preset"
)

type editMode int

const (
	editModeChain editMode = iota
	editModeAdd
)

type presetEditModel struct {
	presetName       string
	chain            []string
	availableEffects []string
	cursor           int
	mode             editMode
	addCursor        int
	presetManager    *preset.Manager
}

func newPresetEditModel(pm *preset.Manager, presetName string) presetEditModel {
	p := pm.GetPreset(presetName)
	var chain []string
	if p != nil {
		chain = make([]string, len(p.EffectChain))
		copy(chain, p.EffectChain)
	}

	return presetEditModel{
		presetName:       presetName,
		chain:            chain,
		availableEffects: pm.GetAvailableEffects(),
		cursor:           0,
		mode:             editModeChain,
		presetManager:    pm,
	}
}

func (m presetEditModel) Update(msg tea.Msg) (presetEditModel, tea.Cmd, Screen) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		if m.mode == editModeChain {
			switch {
			case MatchKey(key, "up"):
				if m.cursor > 0 {
					m.cursor--
				}
			case MatchKey(key, "down"):
				if m.cursor < len(m.chain)-1 {
					m.cursor++
				}
			case MatchKey(key, "moveUp"):
				if m.cursor > 0 && len(m.chain) > 0 {
					m.chain[m.cursor], m.chain[m.cursor-1] = m.chain[m.cursor-1], m.chain[m.cursor]
					m.cursor--
				}
			case MatchKey(key, "moveDown"):
				if m.cursor < len(m.chain)-1 {
					m.chain[m.cursor], m.chain[m.cursor+1] = m.chain[m.cursor+1], m.chain[m.cursor]
					m.cursor++
				}
			case MatchKey(key, "delete"), MatchKey(key, "backspace"):
				if len(m.chain) > 0 && m.cursor < len(m.chain) {
					m.chain = append(m.chain[:m.cursor], m.chain[m.cursor+1:]...)
					if m.cursor >= len(m.chain) && m.cursor > 0 {
						m.cursor--
					}
				}
			case MatchKey(key, "add"):
				m.mode = editModeAdd
				m.addCursor = 0
			case MatchKey(key, "save"):
				m.presetManager.UpdatePresetChain(m.presetName, m.chain)
				return m, nil, ScreenPresetList
			case MatchKey(key, "esc"):
				return m, nil, ScreenPresetList
			}
		} else {
			switch {
			case MatchKey(key, "up"):
				if m.addCursor > 0 {
					m.addCursor--
				}
			case MatchKey(key, "down"):
				if m.addCursor < len(m.availableEffects)-1 {
					m.addCursor++
				}
			case MatchKey(key, "enter"), MatchKey(key, "space"):

				if len(m.availableEffects) > 0 {
					effectToAdd := m.availableEffects[m.addCursor]
					m.chain = append(m.chain, effectToAdd)
					m.mode = editModeChain
					m.cursor = len(m.chain) - 1
				}
			case MatchKey(key, "esc"):
				m.mode = editModeChain
			}
		}
	}
	return m, nil, ScreenPresetEdit
}

func (m presetEditModel) View() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("\n Edit Preset: %s\n", m.presetName))
	b.WriteString(" ====================\n\n")

	if m.mode == editModeChain {
		b.WriteString(" Effect Chain (IN -> OUT):\n")
		if len(m.chain) == 0 {
			b.WriteString("   (empty chain - press [a] to add effects)\n")
		} else {
			for i, effect := range m.chain {
				cursor := "  "
				if i == m.cursor {
					cursor = "> "
				}
				b.WriteString(fmt.Sprintf("%s%d. %s\n", cursor, i+1, effect))
			}
		}
		b.WriteString("\n [j/k] Navigate  [J/K] Reorder  [a] Add  [x] Delete  [s] Save  [esc] Cancel\n")
	} else {
		b.WriteString(" Select effect to add:\n")
		if len(m.availableEffects) == 0 {
			b.WriteString("   (no effects available)\n")
		} else {
			for i, effect := range m.availableEffects {
				cursor := "  "
				if i == m.addCursor {
					cursor = "> "
				}
				b.WriteString(fmt.Sprintf("%s%s\n", cursor, effect))
			}
		}
		b.WriteString("\n [enter] Add  [esc] Cancel\n")
	}

	return b.String()
}

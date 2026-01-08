package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chloyka/gorig/internal/audio"
	"github.com/chloyka/gorig/internal/logger"
	"github.com/chloyka/gorig/internal/logger/keys"
	"github.com/chloyka/gorig/internal/pedal"
	"github.com/chloyka/gorig/internal/preset"
	"github.com/chloyka/gorig/internal/rhythm"
)

const hotkeysHelp = `
 Hotkeys:
   [d] Toggle Effects   [t] Tap Tempo
   [r] Reload Effects   [,/.] BPM -/+
   [p] Presets Menu     [/]] Subdivision
   [i/o] Input/Output   [q] Quit
`

type LampDecayMsg struct{}

type RhythmTickMsg struct{}

type QuantizedOnsetMsg struct {
	Onset rhythm.QuantizedOnset
}

const lampDecayDuration = 100 * time.Millisecond

const rhythmTickInterval = 20 * time.Millisecond

type model struct {
	pedalState    *pedal.State
	audioEngine   *audio.Engine
	presetManager *preset.Manager
	logger        *logger.Logger
	quitting      bool

	lampOn          bool
	lastOnsetEnergy float32

	rhythmViz rhythmVisualizer

	currentScreen Screen
	presetList    presetListModel
	presetCreate  presetCreateModel
	presetEdit    presetEditModel
}

func NewModel(pedalState *pedal.State, audioEngine *audio.Engine, presetManager *preset.Manager, logger *logger.Logger) model {
	logger.Debug("tui model created")

	rhythmViz := newRhythmVisualizer()
	if re := audioEngine.RhythmEngine(); re != nil {
		rhythmViz.bpm = re.GetBPM()
		rhythmViz.subdivision = re.GetSubdivision()
	}

	return model{
		pedalState:    pedalState,
		audioEngine:   audioEngine,
		presetManager: presetManager,
		logger:        logger,
		currentScreen: ScreenMain,
		rhythmViz:     rhythmViz,
	}
}

func listenForQuantizedOnsets(engine *rhythm.Engine) tea.Cmd {
	return func() tea.Msg {
		if engine == nil {
			return nil
		}
		onset, ok := <-engine.QuantizedOnsets()
		if !ok {
			return nil
		}
		return QuantizedOnsetMsg{Onset: onset}
	}
}

func rhythmTickCmd() tea.Cmd {
	return tea.Tick(rhythmTickInterval, func(t time.Time) tea.Msg {
		return RhythmTickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	m.logger.Debug("tui initialized")

	var cmds []tea.Cmd

	if re := m.audioEngine.RhythmEngine(); re != nil {
		cmds = append(cmds, listenForQuantizedOnsets(re))
	}

	cmds = append(cmds, rhythmTickCmd())

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case RhythmTickMsg:

		if re := m.audioEngine.RhythmEngine(); re != nil {
			m.rhythmViz.Update(
				re.GetBPM(),
				re.GetSubdivision(),
				re.GetBeatPhase(),
				re.GetBeatCount(),
			)
		}
		return m, rhythmTickCmd()

	case QuantizedOnsetMsg:

		m.lampOn = true
		m.lastOnsetEnergy = msg.Onset.OriginalEvent.Energy

		beatInDisplay := int(m.rhythmViz.beatCount) % beatsToShow
		m.rhythmViz.AddHitMarker(beatInDisplay, msg.Onset.SlotIndex)

		var cmds []tea.Cmd
		cmds = append(cmds, tea.Tick(lampDecayDuration, func(time.Time) tea.Msg {
			return LampDecayMsg{}
		}))
		if re := m.audioEngine.RhythmEngine(); re != nil {
			cmds = append(cmds, listenForQuantizedOnsets(re))
		}
		return m, tea.Batch(cmds...)

	case LampDecayMsg:
		m.lampOn = false
		return m, nil
	}

	switch m.currentScreen {
	case ScreenPresetList:
		var cmd tea.Cmd
		var nextScreen Screen
		var editPresetName string
		m.presetList, cmd, nextScreen, editPresetName = m.presetList.Update(msg)
		if nextScreen != ScreenPresetList {
			m.currentScreen = nextScreen
			if nextScreen == ScreenPresetCreate {
				m.presetCreate = newPresetCreateModel(m.presetManager)
			} else if nextScreen == ScreenPresetEdit && editPresetName != "" {
				m.presetEdit = newPresetEditModel(m.presetManager, editPresetName)
			}
		}
		return m, cmd

	case ScreenPresetCreate:
		var cmd tea.Cmd
		var nextScreen Screen
		m.presetCreate, cmd, nextScreen = m.presetCreate.Update(msg)
		if nextScreen != ScreenPresetCreate {
			m.currentScreen = nextScreen
			if nextScreen == ScreenPresetList {
				m.presetList = newPresetListModel(m.presetManager)
			}
		}
		return m, cmd

	case ScreenPresetEdit:
		var cmd tea.Cmd
		var nextScreen Screen
		m.presetEdit, cmd, nextScreen = m.presetEdit.Update(msg)
		if nextScreen != ScreenPresetEdit {
			m.currentScreen = nextScreen
			if nextScreen == ScreenPresetList {
				m.presetList = newPresetListModel(m.presetManager)
			}
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		m.logger.Debug("key pressed", keys.UIKey(key))

		switch {
		case MatchKey(key, ActionQuit):
			m.logger.Info("quit requested", keys.UIKey(key))
			m.quitting = true
			return m, tea.Quit
		case MatchKey(key, ActionToggle):
			m.logger.Debug("toggle effects requested")
			needsPresetScreen := m.pedalState.ToggleEffects()
			if needsPresetScreen {

				m.presetCreate = newPresetCreateModel(m.presetManager)
				m.currentScreen = ScreenPresetCreate
			}
			return m, nil
		case MatchKey(key, ActionReload):
			m.logger.Debug("reload effects requested")
			m.pedalState.ReloadEffects()
			return m, nil
		case MatchKey(key, ActionPresets):
			m.logger.Debug("preset menu requested")
			m.presetList = newPresetListModel(m.presetManager)
			m.currentScreen = ScreenPresetList
			return m, nil
		case MatchKey(key, ActionInput):
			m.logger.Debug("next input device requested")
			m.audioEngine.NextInputDevice()
			return m, nil
		case MatchKey(key, ActionOutput):
			m.logger.Debug("next output device requested")
			m.audioEngine.NextOutputDevice()
			return m, nil

		case MatchKey(key, ActionTapTempo):
			if re := m.audioEngine.RhythmEngine(); re != nil {
				re.RegisterTap(time.Now())
				m.rhythmViz.bpm = re.GetBPM()
			}
			return m, nil
		case MatchKey(key, ActionBpmUp):
			if re := m.audioEngine.RhythmEngine(); re != nil {
				re.AdjustBPM(1)
				m.rhythmViz.bpm = re.GetBPM()
			}
			return m, nil
		case MatchKey(key, ActionBpmDown):
			if re := m.audioEngine.RhythmEngine(); re != nil {
				re.AdjustBPM(-1)
				m.rhythmViz.bpm = re.GetBPM()
			}
			return m, nil
		case MatchKey(key, ActionSubdivisionUp):
			if re := m.audioEngine.RhythmEngine(); re != nil {
				re.NextSubdivision()
				m.rhythmViz.subdivision = re.GetSubdivision()
			}
			return m, nil
		case MatchKey(key, ActionSubdivisionDn):
			if re := m.audioEngine.RhythmEngine(); re != nil {
				re.PrevSubdivision()
				m.rhythmViz.subdivision = re.GetSubdivision()
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.logger.Debug("window resized", keys.UIWidth(msg.Width), keys.UIHeight(msg.Height))
		m.rhythmViz.SetWidth(msg.Width)
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	switch m.currentScreen {
	case ScreenPresetList:
		return m.presetList.View()
	case ScreenPresetCreate:
		return m.presetCreate.View()
	case ScreenPresetEdit:
		return m.presetEdit.View()
	}

	amp := getAmpArt(m.pedalState.IsEffectsOn(), m.lampOn)

	presetInfo := ""
	if p := m.presetManager.GetActivePreset(); p != nil {
		presetInfo = fmt.Sprintf("\n Preset: %s", p.Name)

		status := m.presetManager.GetActivePresetStatus()
		if status != nil && status.MissingCount > 0 {
			presetInfo += fmt.Sprintf(" (%d effects missing!)", status.MissingCount)
		}
		presetInfo += "\n"
	} else {
		presetInfo = "\n Preset: (none - press [p] to create)\n"
	}

	effects := m.pedalState.GetEffects()
	var chainParts []string
	for _, e := range effects {
		chainParts = append(chainParts, fmt.Sprintf("[%s]", e.Name))
	}

	chainDisplay := ""
	if len(effects) > 0 {
		chainDisplay = fmt.Sprintf("\n Effects Chain:\n   IN -> %s -> OUT\n", strings.Join(chainParts, " -> "))
	} else {
		chainDisplay = "\n Effects Chain: (empty - add effects in preset menu)\n"
	}

	devices := fmt.Sprintf(`
 Devices:
   IN:  %s
   OUT: %s
`, m.audioEngine.CurrentInputDevice(), m.audioEngine.CurrentOutputDevice())

	rhythmDisplay := "\n" + m.rhythmViz.View()

	return fmt.Sprintf("%s%s%s%s%s%s", amp, presetInfo, chainDisplay, devices, rhythmDisplay, hotkeysHelp)
}

package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chloyka/gorig/internal/audio"
	"github.com/chloyka/gorig/internal/logger"
	"github.com/chloyka/gorig/internal/pedal"
	"github.com/chloyka/gorig/internal/preset"
	"go.uber.org/fx"
)

type TUI struct {
	program *tea.Program
	logger  *logger.Logger
}

func NewTUI(pedalState *pedal.State, audioEngine *audio.Engine, presetManager *preset.Manager, log *logger.Logger) *TUI {
	log.Debug("creating TUI")
	m := NewModel(pedalState, audioEngine, presetManager, log)
	p := tea.NewProgram(m, tea.WithAltScreen())
	return &TUI{program: p, logger: log}
}

func (t *TUI) Run() error {
	t.logger.Info("TUI started")
	_, err := t.program.Run()
	t.logger.Info("TUI stopped")
	return err
}

func (t *TUI) Quit() {
	t.logger.Debug("TUI quit called")
	t.program.Quit()
}

var Module = fx.Module("tui",
	fx.Provide(NewTUI),
	fx.Invoke(registerHooks),
)

func registerHooks(lc fx.Lifecycle, tui *TUI, shutdowner fx.Shutdowner, log *logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Debug("TUI OnStart hook")
			go func() {
				_ = tui.Run()
				log.Debug("TUI exited, triggering shutdown")
				_ = shutdowner.Shutdown()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Debug("TUI OnStop hook")
			tui.Quit()
			return nil
		},
	})
}

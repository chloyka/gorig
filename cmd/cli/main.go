package main

import (
	"github.com/chloyka/gorig/internal/audio"
	"github.com/chloyka/gorig/internal/config"
	"github.com/chloyka/gorig/internal/effects"
	"github.com/chloyka/gorig/internal/logger"
	"github.com/chloyka/gorig/internal/onset"
	"github.com/chloyka/gorig/internal/pedal"
	"github.com/chloyka/gorig/internal/preset"
	"github.com/chloyka/gorig/internal/rhythm"
	"github.com/chloyka/gorig/internal/tui"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	fx.New(
		config.ProvideConfig(),

		logger.Module,

		config.ProvideManager(),
		fx.WithLogger(func(log *logger.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Logger}
		}),
		effects.Module,
		onset.Module,
		rhythm.Module,
		preset.Module,
		audio.Module,
		pedal.Module,
		tui.Module,
	).Run()
}

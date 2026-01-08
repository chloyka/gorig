package audio

import (
	"context"

	configTypes "github.com/chloyka/gorig/internal/config/types"
	"github.com/chloyka/gorig/internal/effects"
	"github.com/chloyka/gorig/internal/logger"
	"github.com/chloyka/gorig/internal/onset"
	"github.com/chloyka/gorig/internal/rhythm"
	"go.uber.org/fx"
)

type newParams struct {
	fx.In

	Logger        *logger.Logger
	Chain         *effects.Chain
	OnsetDetector *onset.Detector
	RhythmEngine  *rhythm.Engine
	AudioConfig   *configTypes.AudioConfig
	StateConfig   *configTypes.StateConfig
}

var Module = fx.Module("audio",
	fx.Provide(func(p newParams) (*Engine, error) {
		return newEngine(p.Logger, p.Chain, p.OnsetDetector, p.RhythmEngine, p.AudioConfig, p.StateConfig)
	}),
	fx.Invoke(registerHooks),
)

func registerHooks(lc fx.Lifecycle, engine *Engine) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return engine.Start()
		},
		OnStop: func(ctx context.Context) error {
			return engine.Stop()
		},
	})
}

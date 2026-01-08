package preset

import (
	configTypes "github.com/chloyka/gorig/internal/config/types"
	"github.com/chloyka/gorig/internal/effects"
	"github.com/chloyka/gorig/internal/logger"
	"go.uber.org/fx"
)

type newManagerParams struct {
	fx.In

	Logger        *logger.Logger
	PresetsConfig *configTypes.PresetsConfig
	Chain         *effects.Chain
}

var Module = fx.Module("preset",
	fx.Provide(provideManager),
)

func provideManager(p newManagerParams) *Manager {
	m := NewManager(
		p.Logger,
		p.PresetsConfig,
		p.Chain.GetAvailableEffectNames,
		p.Chain.SetPresetChain,
	)
	return m
}

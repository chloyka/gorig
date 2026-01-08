package effects

import (
	configTypes "github.com/chloyka/gorig/internal/config/types"
	"github.com/chloyka/gorig/internal/logger"
	"github.com/chloyka/gorig/internal/logger/keys"
	"go.uber.org/fx"
)

type newChainParams struct {
	fx.In

	Logger        *logger.Logger
	EffectsConfig *configTypes.EffectsConfig
	StateConfig   *configTypes.StateConfig
	PresetsConfig *configTypes.PresetsConfig
}

var Module = fx.Module("effects",
	fx.Provide(newEffectsChain),
)

func newEffectsChain(p newChainParams) *Chain {
	p.Logger.Debug("loading effects from directory", keys.PathEffectsDir(p.EffectsConfig.EffectsDir))
	return NewChain(p.Logger, p.EffectsConfig.EffectsDir, p.StateConfig, p.PresetsConfig)
}

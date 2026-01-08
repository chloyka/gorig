package onset

import (
	"go.uber.org/fx"

	configTypes "github.com/chloyka/gorig/internal/config/types"
)

var Module = fx.Module("onset",
	fx.Provide(NewDetectorFromConfig),
)

func NewDetectorFromConfig(audioCfg *configTypes.AudioConfig) *Detector {
	cfg := DefaultConfig(float32(audioCfg.SampleRate))
	return NewDetector(cfg)
}

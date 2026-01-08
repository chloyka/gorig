package rhythm

import (
	"go.uber.org/fx"

	configTypes "github.com/chloyka/gorig/internal/config/types"
	"github.com/chloyka/gorig/internal/onset"
)

var Module = fx.Module("rhythm",
	fx.Provide(NewEngineFromConfig),
)

func NewEngineFromConfig(
	audioCfg *configTypes.AudioConfig,
	stateCfg *configTypes.StateConfig,
	detector *onset.Detector,
) *Engine {

	bpm := stateCfg.RhythmBPM
	if bpm < MinBPM || bpm > MaxBPM {
		bpm = DefaultBPM
	}

	sub := SubdivisionFromInt(stateCfg.RhythmSubdivision)

	return NewEngine(EngineConfig{
		SampleRate:  float32(audioCfg.SampleRate),
		InitialBPM:  bpm,
		Subdivision: sub,
		OnsetEvents: detector.Events(),
		OnStateChange: func(newBPM float64, newSub int) {
			stateCfg.SetRhythmBPM(newBPM)
			stateCfg.SetRhythmSubdivision(newSub)
		},
	})
}

package pedal

import "go.uber.org/fx"

var Module = fx.Module("pedal",
	fx.Provide(New),
)

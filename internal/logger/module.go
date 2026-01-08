package logger

import (
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module("logger",
	fx.Provide(newLogger),
	fx.Invoke(registerHooks),
)

func registerHooks(lc fx.Lifecycle, logger *Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			_ = logger.Sync()

			return logger.Close()
		},
	})
}

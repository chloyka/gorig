package config

import (
	"testing"

	"github.com/chloyka/gorig/internal/logger"
	"go.uber.org/fx"
)

func TestFXIntegration(t *testing.T) {
	t.Run("should validate complete config module dependencies", func(t *testing.T) {
		err := fx.ValidateApp(
			ProvideConfig(),
			logger.Module,
			ProvideManager(),
		)
		if err != nil {
			t.Fatalf("fx validation failed: %v", err)
		}
	})
}

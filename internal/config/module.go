package config

import (
	"auth/internal/config/env"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"config",
	fx.Provide(
		env.NewConfig,
	),
)

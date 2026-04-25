package application

import (
	"auth/internal/application/usecase"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"application",
	fx.Provide(
		usecase.NewRegister,
		usecase.NewLogin,
	),
)

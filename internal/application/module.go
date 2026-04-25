package application

import (
	"auth/internal/application/service"
	"auth/internal/application/usecase"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"application",
	fx.Provide(
		service.NewSession,
		usecase.NewRegister,
		usecase.NewLogin,
	),
)

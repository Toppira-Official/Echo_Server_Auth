package infrastructure

import (
	"auth/internal/domain/contract"
	"auth/internal/infrastructure/uuid"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"infrastructure",
	fx.Provide(
		fx.Annotate(
			uuid.NewKsuidIdGenerator,
			fx.As(new(contract.UuidGenerator)),
		),
	),
)

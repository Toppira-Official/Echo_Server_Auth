package ui

import (
	"auth/internal/ui/doc/router"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"ui",
	fx.Invoke(
		router.RegisterSwaggerRoutes,
	),
)

package ui

import (
	docsrouter "auth/internal/ui/doc/router"
	"auth/internal/ui/register/controller"
	registerrouter "auth/internal/ui/register/router"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"ui",
	fx.Provide(
		controller.NewRegister,
	),
	fx.Invoke(
		docsrouter.RegisterSwaggerRoutes,
		registerrouter.RegisterAuthRegisterRoutes,
	),
)

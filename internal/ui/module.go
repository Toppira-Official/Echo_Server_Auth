package ui

import (
	docsrouter "auth/internal/ui/doc/router"
	logincontroller "auth/internal/ui/login/controller"
	loginrouter "auth/internal/ui/login/router"
	"auth/internal/ui/middlewares"
	registercontroller "auth/internal/ui/register/controller"
	registerrouter "auth/internal/ui/register/router"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"ui",
	fx.Provide(
		registercontroller.NewRegister,
		logincontroller.NewLogin,
		middlewares.NewError,
	),
	fx.Invoke(
		docsrouter.RegisterSwaggerRoutes,
		registerrouter.RegisterAuthRegisterRoutes,
		middlewares.RegisterErrorMiddleware,
		loginrouter.RegisterAuthLoginRoutes,
	),
)

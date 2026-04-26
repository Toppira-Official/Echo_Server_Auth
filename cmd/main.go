package main

import (
	"auth/internal/application"
	"auth/internal/config"
	"auth/internal/infrastructure"
	"auth/internal/ui"

	"go.uber.org/fx"

	_ "auth/docs/swagger"
)

//	@title			Echo Swagger Document
//	@version		0.0.1
//	@description	Echo is a social media platform where people can freely share their feelings without being judged.

// @contact.name	Ali Moradi
// @contact.url	toppira.com/support
// @contact.email	info@toppira.com
func main() {
	fx.
		New(
			config.Module,
			application.Module,
			infrastructure.Module,
			ui.Module,
		).
		Run()
}

package main

import (
	"auth/internal/application"
	"auth/internal/config"
	"auth/internal/infrastructure"

	"go.uber.org/fx"
)

func main() {
	fx.
		New(
			config.Module,
			application.Module,
			infrastructure.Module,
		).
		Run()
}

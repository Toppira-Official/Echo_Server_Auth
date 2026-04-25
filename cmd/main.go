package main

import (
	"auth/internal/application"
	"auth/internal/config"

	"go.uber.org/fx"
)

func main() {
	fx.
		New(
			config.Module,
			application.Module,
		).
		Run()
}

package main

import (
	"auth/internal/config"

	"go.uber.org/fx"
)

func main() {
	fx.
		New(
			config.Module,
		).
		Run()
}

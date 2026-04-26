package main

import (
	"auth/internal/infrastructure/db/gorm/model"

	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/infrastructure/db/gorm/dao",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.ApplyBasic(
		model.Credential{},
	)

	g.Execute()
}

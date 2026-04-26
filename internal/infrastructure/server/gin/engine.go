package gin

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ali127Dev/xerr"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type GinEngineConfig struct {
	Mode    string
	Port    int
	Origins []string
}

func NewGinEngine(lc fx.Lifecycle, cfg GinEngineConfig) *gin.Engine {
	switch cfg.Mode {
	case "prod":
	case "production":
	case "release":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Port),
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil {
					log.Printf("http server listen error: %v\n", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := srv.Shutdown(ctx); err != nil {
				return xerr.Wrap(
					err,
					xerr.CodeInternalError,
					xerr.WithMessage("failed to gracefully shutdown http server"),
					xerr.WithMeta("address", srv.Addr),
				)
			}
			return nil
		},
	})

	return r
}

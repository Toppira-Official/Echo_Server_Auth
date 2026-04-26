package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Ali127Dev/xerr"
	"github.com/avast/retry-go"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Mode     string
	Host     string
	Port     int
	User     string
	Password string
	DB       string
}

func NewDB(lc fx.Lifecycle, cfg DBConfig) (*gorm.DB, *sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DB,
		cfg.Port,
	)

	var db *gorm.DB
	var err error

	var loggerMode logger.LogLevel
	switch cfg.Mode {
	case "prod":
	case "production":
	case "release":
		loggerMode = logger.Silent
	default:
		loggerMode = logger.Warn
	}

	err = retry.Do(func() error {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:         logger.Default.LogMode(loggerMode),
			PrepareStmt:    true,
			TranslateError: true,
		})
		if err != nil {
			return xerr.Wrap(
				err,
				xerr.CodeServiceUnavailable,
				xerr.WithMessage("failed to open postgres connection"),
				xerr.WithMeta("host", cfg.Host),
				xerr.WithMeta("db", cfg.DB),
			)
		}
		return nil
	}, retry.Attempts(5), retry.Delay(2*time.Second))

	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, xerr.Wrap(
			err,
			xerr.CodeInternalError,
			xerr.WithMessage("failed to retrieve underlying sql.DB from gorm"),
		)
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return retry.Do(func() error {
				if pingErr := sqlDB.PingContext(ctx); pingErr != nil {
					return xerr.Wrap(
						pingErr,
						xerr.CodeServiceUnavailable,
						xerr.WithMessage("postgres ping failed"),
						xerr.WithMeta("host", cfg.Host),
						xerr.WithMeta("db", cfg.DB),
					)
				}
				return nil
			}, retry.Attempts(3), retry.Delay(1*time.Second), retry.Context(ctx))
		},
		OnStop: func(ctx context.Context) error {
			if closeErr := sqlDB.Close(); closeErr != nil {
				return xerr.Wrap(
					closeErr,
					xerr.CodeInternalError,
					xerr.WithMessage("failed to close postgres connection"),
				)
			}
			return nil
		},
	})

	return db, sqlDB, nil
}

package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresGormDBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DB       string
}

func NewPostgresGormDB(lc fx.Lifecycle, cfg PostgresGormDBConfig) (*gorm.DB, *sql.DB, error) {
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

	err = retry.Do(func() error {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:      logger.Default.LogMode(logger.Warn),
			PrepareStmt: true,
		})
		return err
	}, retry.Attempts(5), retry.Delay(2*time.Second))

	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return retry.Do(func() error {
				return sqlDB.PingContext(ctx)
			}, retry.Attempts(3), retry.Delay(1*time.Second), retry.Context(ctx))
		},
		OnStop: func(ctx context.Context) error {
			return sqlDB.Close()
		},
	})

	return db, sqlDB, nil
}

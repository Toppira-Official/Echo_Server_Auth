package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapLoggerConfig struct {
	Mode       string
	LogPath    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

func NewZapLogger(cfg ZapLoggerConfig) (*zap.Logger, error) {
	var encoder zapcore.Encoder
	var core zapcore.Core
	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	// -------------------------
	//  DEV MODE
	// -------------------------
	if cfg.Mode == "dev" {
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		encoder = zapcore.NewConsoleEncoder(encoderCfg)

		core = zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			level,
		)

		return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
	}

	// -------------------------
	//  PRODUCTION MODE
	// -------------------------

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoder = zapcore.NewJSONEncoder(encoderCfg)

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.LogPath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	})

	consoleWriter := zapcore.Lock(os.Stdout)

	core = zapcore.NewTee(
		zapcore.NewCore(encoder, fileWriter, level),

		zapcore.NewCore(consoleEncoder, consoleWriter, level),
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}

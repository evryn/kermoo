package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func MustInitLogger(level string) {
	var err error

	if level == "" {
		level = "info"
	}

	config := zap.NewProductionConfig()

	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	levelObject, err := zapcore.ParseLevel(level)

	if err != nil {
		panic(err)
	}

	config.Level = zap.NewAtomicLevelAt(levelObject)

	Log, err = config.Build()

	if err != nil {
		panic(err)
	}
}

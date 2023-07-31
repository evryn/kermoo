package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func MustInitLogger(verbosity string) {
	var err error

	config := zap.NewProductionConfig()

	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	if verbosity == "debug" {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	if verbosity == "fatal" {
		config.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	}

	Log, err = config.Build()

	if err != nil {
		panic(err)
	}
}

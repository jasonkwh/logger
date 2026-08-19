package logger

import (
	"errors"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Interface = (*logger)(nil)

type logger struct {
	zl       *zap.Logger
	logLevel zapcore.Level
}

func New(opts ...loggerOption) (Interface, error) {
	logr := &logger{
		logLevel: zapcore.InfoLevel, // default log level
	}

	zlCfg := zap.NewProductionConfig()

	for _, opt := range opts {
		opt(logr)
	}

	zlCfg.Level.SetLevel(logr.logLevel)

	zl, err := zlCfg.Build()
	if err != nil {
		return nil, err
	}

	logr.zl = zl

	return logr, nil
}

func (logr *logger) Logger() *zap.Logger {
	return logr.zl
}

func (logr *logger) Close() error {
	if err := logr.zl.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
		return err
	}

	return nil
}

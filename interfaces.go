package logger

import (
	"context"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Interface interface {
	Logger() *zap.Logger
	ShiftLevel(ctx context.Context,
		level zapcore.Level) error

	io.Closer
}

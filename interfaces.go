package logger

import (
	"io"

	"go.uber.org/zap"
)

type Interface interface {
	Logger() *zap.Logger

	io.Closer
}

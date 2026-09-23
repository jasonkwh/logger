package logger

import (
	"errors"
	"fmt"

	"go.uber.org/zap/zapcore"
)

var (
	ErrMarshalJSON = func(err error) error {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}
	ErrInvalidLogLevel = func(level zapcore.Level) error {
		//nolint:err113
		return fmt.Errorf("invalid log level: %d", level)
	}
	ErrFailedServerListenerInit = func(err error) error {
		return fmt.Errorf("failed to initialise server listener: %w", err)
	}
	ErrShiftLevelHTTPStatus = func(statusCode int) error {
		//nolint:err113
		return fmt.Errorf("failed to shift log level: HTTP %d", statusCode)
	}

	ErrShiftLevelWhenHTTPServerEnabled = errors.New("cannot shift level when HTTP server is enabled")
)

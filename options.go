package logger

import "go.uber.org/zap/zapcore"

type loggerOption func(*logger)

func WithLogLevel(logLevel int8) loggerOption {
	return func(c *logger) {
		c.logLevel = zapcore.Level(logLevel)
	}
}

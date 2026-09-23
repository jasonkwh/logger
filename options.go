package logger

import "go.uber.org/zap/zapcore"

type loggerOption func(*logger)

func WithLogLevel(logLevel int8) loggerOption {
	return func(c *logger) {
		c.logLevel = zapcore.Level(logLevel)
	}
}

func WithHTTPServerPort(port int) loggerOption {
	return func(c *logger) {
		c.httpServerPort = port
	}
}

func WithHTTPServerPath(path string) loggerOption {
	return func(c *logger) {
		c.httpServerPath = path
	}
}

func DisableHTTPServer() loggerOption {
	return func(c *logger) {
		c.disableHTTPServer = true
	}
}

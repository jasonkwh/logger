package logger

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"syscall"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Interface = (*logger)(nil)

type logger struct {
	httpServer *http.Server
	httpClient *http.Client
	zl         *zap.Logger
	logLevel   zapcore.Level

	disableHTTPServer bool
	httpServerPort    int
	httpServerPath    string
}

func New(opts ...loggerOption) (Interface, error) {
	logr := &logger{
		httpServerPort: defaultHTTPServerPort,
		httpServerPath: defaultHTTPServerPath,
	}

	zlCfg := zap.NewProductionConfig()

	for _, opt := range opts {
		opt(logr)
	}

	if logr.disableHTTPServer {
		zlCfg.Level.SetLevel(logr.logLevel)

		zl, err := zlCfg.Build()
		if err != nil {
			return nil, err
		}

		logr.zl = zl
		logr.httpClient = defaultHTTPClient()

		return logr, nil
	}

	atomicLevel := zap.NewAtomicLevelAt(logr.logLevel)
	zlCfg.Level = atomicLevel

	zl, err := zlCfg.Build()
	if err != nil {
		return nil, err
	}

	logr.zl = zl

	if err := logr.startHTTPServer(atomicLevel); err != nil {
		return nil, err
	}

	return logr, nil
}

func (logr *logger) Logger() *zap.Logger {
	return logr.zl
}

// ShiftLevel - shifts the log level of the main container,
// only for Log Shifter to use.
func (logr *logger) ShiftLevel(
	ctx context.Context,
	level zapcore.Level,
) error {
	if !logr.disableHTTPServer {
		return ErrShiftLevelWhenHTTPServerEnabled
	}

	if level < zapcore.DebugLevel || level > zapcore.FatalLevel {
		return ErrInvalidLogLevel(level)
	}

	payloadBytes, err := loggerEndpointPayload{
		Level: level.String(),
	}.ToJSON()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		logr.shiftLevelURL(),
		bytes.NewReader(payloadBytes),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := logr.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return ErrShiftLevelHTTPStatus(resp.StatusCode)
	}

	return nil
}

func (logr *logger) Close() error {
	var errs error

	if err := logr.zl.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
		errs = multierr.Append(errs, err)
	}

	if logr.httpServer != nil {
		errs = multierr.Append(errs, logr.httpServer.Close())
	}

	if logr.httpClient != nil {
		logr.httpClient.CloseIdleConnections()
	}

	return errs
}

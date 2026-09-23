package logger

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func (logr *logger) startHTTPServer(
	atomicLevel zap.AtomicLevel,
) error {
	mux := http.NewServeMux()

	// logger level change endpoint
	mux.Handle(logr.httpServerPath, atomicLevel)

	logr.httpServer = &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,  //nolint:mnd
		ReadTimeout:       5 * time.Second,  //nolint:mnd
		WriteTimeout:      5 * time.Second,  //nolint:mnd
		IdleTimeout:       30 * time.Second, //nolint:mnd
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", logr.httpServerPort))
	if err != nil {
		return ErrFailedServerListenerInit(err)
	}

	go func() {
		if err := logr.httpServer.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logr.zl.Error("logger server failed", zap.Error(err))
		}
	}()

	logr.zl.Info("logger server started", zap.Int("port", logr.httpServerPort))

	return nil
}

func (logr *logger) shiftLevelURL() string {
	u := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort("localhost", strconv.Itoa(logr.httpServerPort)),
		Path:   logr.httpServerPath,
	}

	return u.String()
}

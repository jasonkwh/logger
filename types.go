package logger

import (
	"encoding/json"
	"net/http"
	"time"
)

type loggerEndpointPayload struct {
	Level string `json:"level"`
}

func (pl loggerEndpointPayload) ToJSON() ([]byte, error) {
	payloadBytes, err := json.Marshal(pl)
	if err != nil {
		return nil, ErrMarshalJSON(err)
	}

	return payloadBytes, nil
}

func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second, //nolint:mnd
		Transport: &http.Transport{
			DisableKeepAlives:     true,
			MaxConnsPerHost:       1,
			TLSHandshakeTimeout:   2 * time.Second, //nolint:mnd
			ResponseHeaderTimeout: 3 * time.Second, //nolint:mnd
		},
	}
}

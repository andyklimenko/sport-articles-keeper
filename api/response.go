package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type statusResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func (a *API) respondNotOK(w http.ResponseWriter, statusCode int, err error) {
	resp := statusResponse{
		Code: statusCode,
		Text: err.Error(),
	}
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("marshal response", slog.Any("error", err))
		return
	}
}

func (a *API) respondOK(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if resp == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("marshal response", slog.Any("error", err))
	}
}

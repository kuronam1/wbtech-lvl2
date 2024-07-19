package handlers

import (
	"encoding/json"
	"events/task/internal/appErrors"
	"events/task/internal/event/repository"
	"events/task/internal/parser"
	"log/slog"
	"net/http"
)

type errResponse struct {
	Error string `json:"error"`
}

type handler struct {
	logger *slog.Logger
	repo   repository.Repository
}

func (h *handler) CreationHandler(w http.ResponseWriter, r *http.Request) {
	request, err := parser.ParseRequest(r.Body)
	if err != nil {
		w.WriteHeader(400)
		errResp, err := errWrapper(err)
		if err != nil {
			http.Error(w, err.Error(), 503)
		}
		w.Write(errResp)
	}

	if err = h.repo.CreateEvent(request); err != nil {
		w.WriteHeader(503)
		errResp, err := errWrapper(err)
		if err != nil {
			http.Error(w, err.Error(), 503)
		}
		w.Write(errResp)
	}
}

func errWrapper(err error) ([]byte, error) {
	const functionName = "handler.errWrapper"
	response := errResponse{Error: err.Error()}

	result, err := json.Marshal(response)
	if err != nil {
		return nil, appErrors.WrapErr(functionName, "error while wrapping error:", err)
	}

	return result, nil
}

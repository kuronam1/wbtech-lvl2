package router

import (
	"events/task/internal/event/repository"
	"events/task/internal/router/handler"
	"log/slog"
	"net/http"
)

type AppRouter struct {
	Router *http.ServeMux
}

func New(logger *slog.Logger, repo repository.Repository) http.Handler {
	router := http.NewServeMux()

	eventHandler := handler.New(logger, repo)

	eventHandler.Register(router)

	return router
}

func (r AppRouter) Handle(pattern string, handler http.Handler) {
	r.Router.Handle(pattern, handler)
}

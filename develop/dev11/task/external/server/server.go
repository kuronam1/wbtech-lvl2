package server

import (
	"context"
	"events/task/internal/config"
	"fmt"
	"net/http"
	"time"
)

type App struct {
	Server       *http.Server
	Notify       chan error
	ShutdownTime time.Duration
}

func New(handler http.Handler, config config.Config) App {
	server := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", config.Server.Host, config.Server.Port),
		Handler:      handler,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
	}

	app := App{
		Server:       server,
		Notify:       make(chan error),
		ShutdownTime: config.Server.ShutdownTime,
	}

	app.Start()

	return app
}

func (a *App) Start() {
	go func() {
		a.Notify <- a.Server.ListenAndServe()
		close(a.Notify)
	}()
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), a.ShutdownTime)
	defer cancel()

	return a.Server.Shutdown(ctx)
}

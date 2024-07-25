package main

import (
	"events/task/external/client"
	"events/task/external/logger"
	"events/task/external/server"
	"events/task/internal/config"
	"events/task/internal/event/repository"
	"events/task/internal/router"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var pathToCfgFile string
	flag.StringVar(&pathToCfgFile, "cfg", "internal/config/config.yaml", "set path to config file")

	cfg := config.MustNew(pathToCfgFile)

	lg := logger.New()
	lg.Info("[Logger] initialized")

	lg.Info("[File client] initializing")
	fileClient, err := client.GetFSClient(lg, cfg.PathToData)
	if err != nil {
		lg.Error(err.Error())
	}

	lg.Info("[Repository] initializing")
	repo, err := repository.NewRepository(fileClient, lg)
	if err != nil {
		lg.Error(err.Error())
	}

	lg.Info("[Router] initializing")
	route := router.New(lg, repo)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	lg.Info("[Server] starting server", "host", cfg.Server.Host, cfg.Server.Port)
	app := server.New(route, cfg)
	lg.Info("[Server] server started", "address", app.Server.Addr)

	select {
	case s := <-interrupt:
		lg.Error("[Run] interrupt signal", "signal", s.String())
	case err := <-app.Notify:
		lg.Error("[Server] run error signal", "error", err.Error())
	}

	lg.Error("[Server] shutting down server", "shutdown time", cfg.Server.ShutdownTime.Seconds())
	if err = app.Shutdown(); err != nil {
		lg.Error("[Server] stopped", "error", err.Error())
	}
}

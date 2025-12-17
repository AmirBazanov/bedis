package main

import (
	"bedis/internal/app"
	"bedis/internal/config"
	"bedis/internal/server"
	"bedis/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	bedis := app.New(cfg.Server, cfg.Logger)
	go func() {
		err := bedis.Server.Start()
		if err != nil {
			panic(err)
		}
	}()
	GracefulShutdown(bedis.Server)
}

func GracefulShutdown(server *server.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	sign := <-stop
	server.Stop()
	logger.GetLogger().Info("gracefully shutting down", sign.String())
}

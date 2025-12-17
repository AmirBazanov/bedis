package app

import (
	"bedis/internal/config"
	"bedis/internal/handler"
	"bedis/internal/server"
	"bedis/internal/storage"
	"bedis/pkg/logger"
	"log/slog"
)

type App struct {
	Logger  *slog.Logger
	Server  *server.Server
	Storage *storage.Storage
	Handler *handler.Handler
}

func New(serverCfg config.Server, loggerCfg config.Logger) *App {
	log := logger.InitLogger(loggerCfg.Service, loggerCfg.Level, loggerCfg.Logfile)
	s := storage.New(log)
	h := handler.New(s, log)
	srv := server.New(serverCfg.Address+":"+serverCfg.Port, h, log)

	if err := srv.Start(); err != nil {
		panic(err)
	}
	return &App{
		Logger:  log,
		Server:  srv,
		Storage: s,
		Handler: h,
	}

}

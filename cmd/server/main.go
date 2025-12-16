package main

import (
	"bedis/internal/handler"
	"bedis/internal/server"
	"bedis/internal/storage"
	logger2 "bedis/pkg/logger"
)

func main() {
	logger := logger2.InitLogger("bedis", "info", "./logs")
	logger2.GetLogger()
	s := storage.New(logger)
	h := handler.New(s, logger)
	srv := server.New(":6380", h, logger)
	if err := srv.Start(); err != nil {
		panic(err)
	}
}

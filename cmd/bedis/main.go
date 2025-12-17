package main

import (
	"bedis/internal/app"
	"bedis/internal/config"
)

func main() {
	cfg := config.MustLoad()
	bedis := app.New(cfg.Server, cfg.Logger)
	defer bedis.Server.Stop()
}

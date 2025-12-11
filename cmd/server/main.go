package main

import (
	"bedis/internal/handler"
	"bedis/internal/server"
	"bedis/internal/storage"
	"log"
)

func main() {
	s := storage.New()
	h := handler.New(s)
	srv := server.New(":6380", h)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}

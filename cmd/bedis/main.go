package main

import (
	"bedis/internal/sender"
)

func main() {

	send := sender.New(":6380")
	send.Send("SET KEY VALUE\n")
}

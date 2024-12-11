package main

import (
	"os"
	"prior-chat-bot/internal/adapter/api"
)

func main() {
	framework := os.Args[1]

	switch framework {
	case "echo":
		api.StartEchoServer()
	default:
		panic("Unknown framework")
	}
}

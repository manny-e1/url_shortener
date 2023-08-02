package main

import (
	"github.com/manny-e1/url_shortener/model"
	"github.com/manny-e1/url_shortener/server"
)

func main() {
	model.Setup()
	server.SetupAndListen()
}

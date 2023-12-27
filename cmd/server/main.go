package main

import (
	"flag"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/pinghandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/userhandler"

	"github.com/labstack/echo/v4"
)

var (
	port   int = *flag.Int("port", 8080, "Which port to run.")
	server httpserver.Server
)

func init() {
	flag.Parse()
}

func main() {
	handlers := []httpserver.Handler{
		userhandler.New(),
		pinghandler.New(),
	}

	config := httpserver.Config{
		HTTPServer: httpserver.HTTPServer{
			Port: port,
		},
	}

	server = httpserver.New(config, echo.New(), handlers)

	server.Serve()
}

package main

import (
	"flag"
	"fmt"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/pinghandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/userhandler"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongouser"
	"mdhesari/kian-quiz-golang-game/service/userservice"
	"time"

	"github.com/hellofresh/janus/pkg/plugin/basic/encrypt"
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
	cfg := config.Load("config.yml")
	fmt.Println(cfg)

	cli, err := mongorepo.New(cfg.Database.MongoDB, 30*time.Second, encrypt.Hash{})
	if err != nil {

		panic("could not connect to mongodb.")
	}

	repo := mongouser.New(cli)

	userSrv := userservice.New(repo)

	handlers := []httpserver.Handler{
		userhandler.New(userSrv),
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

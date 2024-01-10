package main

import (
	"flag"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/pinghandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/userhandler"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongouser"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/userservice"
	"time"

	"github.com/hellofresh/janus/pkg/plugin/basic/encrypt"
)

var (
	cfg    config.Config
	server httpserver.Server
)

func init() {
	cfg = config.Load("config.yml")
	flag.Parse()
}

func main() {
	// Todo: duration should be in config
	cli, err := mongorepo.New(cfg.Database.MongoDB, 5*time.Second, encrypt.Hash{})
	if err != nil {

		panic("could not connect to mongodb.")
	}

	repo := mongouser.New(cli)

	authSrv := authservice.New(cfg.JWT.Secret)

	userSrv := userservice.New(&authSrv, repo)

	handlers := []httpserver.Handler{
		userhandler.New(userSrv, authSrv),
		pinghandler.New(),
	}

	config := httpserver.Config{
		HTTPServer: cfg.Server.HTTPServer,
	}

	server = httpserver.New(config, handlers)

	server.Serve()
}

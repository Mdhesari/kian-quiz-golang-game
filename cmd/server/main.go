package main

import (
	"flag"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/backpanelhandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/pinghandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/userhandler"
	"mdhesari/kian-quiz-golang-game/delivery/validator/uservalidator"
	"mdhesari/kian-quiz-golang-game/repository/migrator"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongorbac"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongouser"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/rbacservice"
	"mdhesari/kian-quiz-golang-game/service/userservice"

	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	cli, err := mongorepo.New(cfg.Database.MongoDB, encrypt.Hash{})
	if err != nil {

		panic("could not connect to mongodb.")
	}

	migrator, err := migrator.New(cli.Conn().Client(), &mongodb.Config{
		DatabaseName:         cfg.Database.MongoDB.DBName,
		MigrationsCollection: cfg.Database.MongoDB.Migrations,
		TransactionMode:      false,
		Locking:              mongodb.Locking{},
	})
	if err != nil {

		panic(err)
	}
	migrator.Up()

	userRepo := mongouser.New(cli)
	rbacRepo := mongorbac.New(cli)

	authConfig := authservice.Config{
		Secret: []byte(cfg.JWT.Secret),
	}
	authSrv := authservice.New(authConfig)

	rbacSrv := rbacservice.New(rbacRepo)

	userSrv := userservice.New(&authSrv, userRepo)

	userValidator := uservalidator.New(userRepo)

	handlers := []httpserver.Handler{
		userhandler.New(&userSrv, &authSrv, &rbacSrv, authConfig, userValidator),
		pinghandler.New(),
		backpanelhandler.New(&userSrv, &rbacSrv, &authSrv, authConfig),
	}

	config := httpserver.Config{
		HTTPServer: cfg.Server.HTTPServer,
	}

	server = httpserver.New(config, handlers)

	server.Serve()
}

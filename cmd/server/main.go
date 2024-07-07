package main

import (
	"flag"
	"mdhesari/kian-quiz-golang-game/adapter/redisadapter"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/backpanelhandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/matchinghandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/pinghandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/userhandler"
	"mdhesari/kian-quiz-golang-game/delivery/validator/matchingvalidator"
	"mdhesari/kian-quiz-golang-game/delivery/validator/uservalidator"
	"mdhesari/kian-quiz-golang-game/repository/migrator"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongocategory"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongorbac"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongouser"
	"mdhesari/kian-quiz-golang-game/repository/redisrepo/redismatching"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/matchingservice"
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
	cli, err := mongorepo.New(cfg.Database.MongoDB, encrypt.Hash{})
	if err != nil {

		panic("could not connect to mongodb.")
	}

	migrator, err := migrator.New(cli.Conn().Client(), &mongodb.Config{
		DatabaseName:         cfg.Database.MongoDB.DBName,
		MigrationsCollection: cfg.Database.MongoDB.Migrations,
		TransactionMode:      false,
		Locking:              mongodb.Locking{},
	}, cfg.Database.Migrations)
	if err != nil {

		panic(err)
	}
	migrator.Up()

	authConfig := authservice.Config{
		Secret: []byte(cfg.JWT.Secret),
	}
	authSrv := authservice.New(authConfig)

	rbacRepo := mongorbac.New(cli)
	rbacSrv := rbacservice.New(rbacRepo)

	userRepo := mongouser.New(cli)
	userSrv := userservice.New(&authSrv, userRepo)
	userValidator := uservalidator.New(userRepo)

	categoryRepo := mongocategory.New(cli)
	redisAdap := redisadapter.New(cfg.Redis)
	matchingRepo := redismatching.New(redisAdap)
	matchingSrv := matchingservice.New(matchingRepo)
	matchingValidator := matchingvalidator.New(categoryRepo)

	handlers := []httpserver.Handler{
		userhandler.New(&userSrv, &authSrv, &rbacSrv, authConfig, userValidator),
		pinghandler.New(),
		backpanelhandler.New(&userSrv, &rbacSrv, &authSrv, authConfig),
		matchinghandler.New(authConfig, &authSrv, matchingSrv, matchingValidator),
	}

	config := httpserver.Config{
		HTTPServer: cfg.Server.HTTPServer,
	}

	server = httpserver.New(config, handlers)

	server.Serve()
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mdhesari/kian-quiz-golang-game/adapter/presenceadapter"
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
	"mdhesari/kian-quiz-golang-game/repository/redisrepo/redispresence"
	"mdhesari/kian-quiz-golang-game/scheduler"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/matchingservice"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
	"mdhesari/kian-quiz-golang-game/service/rbacservice"
	"mdhesari/kian-quiz-golang-game/service/userservice"
	"os"
	"os/signal"
	"sync"

	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
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
	cli, err := mongorepo.New(cfg.Database.MongoDB)
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

	grpConn, err := grpc.Dial("172.18.0.5:8089", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Grpc could not dial %v\n", err)
	}
	presenceCli := presenceadapter.New(grpConn)
	matchingSrv := matchingservice.New(matchingRepo, categoryRepo, presenceCli)
	matchingValidator := matchingvalidator.New()

	presenceRepo := redispresence.New(redisAdap)
	presenceSrv := presenceservice.New(cfg.Presence, presenceRepo)

	handlers := []httpserver.Handler{
		userhandler.New(&userSrv, &authSrv, &rbacSrv, &presenceSrv, authConfig, userValidator),
		pinghandler.New(),
		backpanelhandler.New(&userSrv, &rbacSrv, &authSrv, authConfig),
		matchinghandler.New(authConfig, &authSrv, matchingSrv, matchingValidator, &presenceSrv),
	}

	config := httpserver.Config{
		HTTPServer: cfg.Server.HTTPServer,
	}

	server = httpserver.New(config, handlers)

	go server.Serve()

	var wg sync.WaitGroup
	scheduler := scheduler.New(cfg.Scheduler, &matchingSrv)
	go scheduler.Start(&wg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Application.GracefulShutdownTimeout)
	defer cancel()
	if err := server.Router.Shutdown(ctx); err != nil {
		fmt.Println("Err: ", err)
	}
	<-ctx.Done()

	wg.Wait()

	// done := make(chan bool, 1)
	// done <- true

	// fmt.Println("Shutting down gracefully")

	// time.Sleep(5 * time.Second)
}

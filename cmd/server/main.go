package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mdhesari/kian-quiz-golang-game/adapter/presenceadapter"
	"mdhesari/kian-quiz-golang-game/adapter/redisadapter"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/delivery/grpcserver"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/backpanelhandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/categoryhandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/matchinghandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/pinghandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/userhandler"
	"mdhesari/kian-quiz-golang-game/delivery/validator/matchingvalidator"
	"mdhesari/kian-quiz-golang-game/delivery/validator/uservalidator"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/pkg/protobufdecoder"
	"mdhesari/kian-quiz-golang-game/repository/migrator"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongocategory"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongorbac"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongouser"
	"mdhesari/kian-quiz-golang-game/repository/redisrepo/redismatching"
	"mdhesari/kian-quiz-golang-game/repository/redisrepo/redispresence"
	"mdhesari/kian-quiz-golang-game/scheduler"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/categoryservice"
	"mdhesari/kian-quiz-golang-game/service/matchingservice"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
	"mdhesari/kian-quiz-golang-game/service/rbacservice"
	"mdhesari/kian-quiz-golang-game/service/userservice"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	_ "net/http/pprof"
)

var (
	cfg    config.Config
	server httpserver.Server
)

type services struct {
	authSrv     *authservice.Service
	userSrv     *userservice.Service
	matchingSrv *matchingservice.Service
	categorySrv *categoryservice.Service
	presenceSrv *presenceservice.Service
	rbacSrv     *rbacservice.Service
}

func init() {
	cfg = config.Load("config.yml")
	flag.Parse()
}

func setupServices(cfg config.Config) services {
	authConfig := authservice.Config{
		Secret: []byte(cfg.JWT.Secret),
	}
	authSrv := authservice.New(authConfig)

	cli, err := mongorepo.New(cfg.Database.MongoDB)
	if err != nil {

		panic("could not connect to mongodb.")
	}

	rbacRepo := mongorbac.New(cli)
	rbacSrv := rbacservice.New(rbacRepo)

	userRepo := mongouser.New(cli)
	userSrv := userservice.New(&authSrv, userRepo)
	userValidator := uservalidator.New(userRepo)

	categoryRepo := mongocategory.New(cli)
	redisAdap := redisadapter.New(cfg.Redis)
	matchingRepo := redismatching.New(redisAdap)

	presenceRepo := redispresence.New(redisAdap)
	presenceSrv := presenceservice.New(cfg.Presence, presenceRepo)

	matchingSrv := matchingservice.New(matchingRepo, categoryRepo, presenceCli, redisAdap)

	categorySrv := categoryservice.New(categoryRepo)
}

func main() {
	log.Println("goroutines: ", runtime.NumGoroutine())

	var srvs services = setupServices(cfg)

	go func() {
		// Start the HTTP server
		server := http.Server{
			Addr: ":6060",
		}

		go func() {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			<-sigChan

			if err := server.Close(); err != nil {
				logger.L().Error("Could not close server.", zap.Error(err))
			}
		}()

		if err := server.ListenAndServe(); err != nil {
			logger.L().Error("Http server error.", zap.Error(err))

			return
		}
	}()

	// TODO - Sepearte cmd
	go func() {
		redisAdap := redisadapter.New(cfg.Redis)
		subscriber := redisAdap.Cli().Subscribe(context.Background(), string(entity.UsersMatched))
		for {
			msg, err := subscriber.ReceiveMessage(context.Background())
			if err != nil {

				logger.L().Error("Redis sub: Could not recieve message.", zap.Error(err))
			}

			playersMatched := protobufdecoder.DecodeUsersMatchedEvent(msg.Payload)
			fmt.Printf("new message %v\n", playersMatched)
		}
	}()

	// TODO - this should be removed later it's just for development purposes
	fmt.Println("Starting presence server...")
	presenceserver := grpcserver.New(srvs.presenceSrv)
	go presenceserver.Start()

	grpConn, err := grpc.Dial("127.0.0.1:8089", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Grpc could not dial %v\n", err)
	}
	presenceCli := presenceadapter.New(grpConn)

	handlers := []httpserver.Handler{
		userhandler.New(srvs.userSrv, srvs.authSrv, srvs.rbacSrv, srvs.presenceSrv, authConfig, userValidator),
		pinghandler.New(),
		backpanelhandler.New(srvs.userSrv, srvs.rbacSrv, srvs.authSrv, authConfig),
		matchinghandler.New(authConfig, srvs.authSrv, *srvs.matchingSrv, matchingValidator, &presenceSrv),
		categoryhandler.New(&categorySrv, &presenceSrv, &authSrv, authConfig),
	}

	// TODO - refactor
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8090", nil)

	config := httpserver.Config{
		HTTPServer: cfg.Server.HTTPServer,
	}

	server = httpserver.New(config, handlers)

	go server.Serve()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		scheduler := scheduler.New(cfg.Scheduler, &matchingSrv)
		scheduler.Start()
	}()

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

	log.Print("Application shutdown gracefully.")
	log.Println("goroutines: ", runtime.NumGoroutine())
}

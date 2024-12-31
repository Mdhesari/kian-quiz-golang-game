package main

import (
	"context"
	"flag"
	"fmt"
	"mdhesari/kian-quiz-golang-game/adapter/presenceadapter"
	"mdhesari/kian-quiz-golang-game/adapter/redisadapter"
	"mdhesari/kian-quiz-golang-game/adapter/websocketadapter"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/delivery/grpcserver"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/backpanelhandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/categoryhandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/gamehandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/matchinghandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/pinghandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/userhandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/websockethandler"
	"mdhesari/kian-quiz-golang-game/delivery/validator/matchingvalidator"
	"mdhesari/kian-quiz-golang-game/delivery/validator/uservalidator"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/matchmaking"
	"mdhesari/kian-quiz-golang-game/pubsub"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongocategory"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongogame"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongoquestion"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongorbac"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongouser"
	"mdhesari/kian-quiz-golang-game/repository/redisrepo/redismatching"
	"mdhesari/kian-quiz-golang-game/repository/redisrepo/redispresence"
	"mdhesari/kian-quiz-golang-game/scheduler"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/categoryservice"
	"mdhesari/kian-quiz-golang-game/service/gameservice"
	"mdhesari/kian-quiz-golang-game/service/matchingservice"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
	"mdhesari/kian-quiz-golang-game/service/questionservice"
	"mdhesari/kian-quiz-golang-game/service/rbacservice"
	"mdhesari/kian-quiz-golang-game/service/userservice"

	"os"
	"os/signal"
	"sync"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"

	_ "net/http/pprof"
)

var (
	cfg config.Config
)

type services struct {
	questionSrv *questionservice.Service
	authSrv     *authservice.Service
	userSrv     *userservice.Service
	matchingSrv *matchingservice.Service
	categorySrv *categoryservice.Service
	presenceSrv *presenceservice.Service
	rbacSrv     *rbacservice.Service
	gameSrv     *gameservice.Service

	websocketAdap *websocketadapter.Adapter

	pubsubManager *pubsub.PubSubManager

	userValidator     *uservalidator.Validator
	matchingValidator *matchingvalidator.Validator
}

func init() {
	cfg = config.Load("config.yml")

	flag.Parse()
}

func main() {
	logger.L().Info("Welcome to KianQuiz.")

	var srvs services = setupServices(&cfg)

	mm := matchmaking.New(srvs.pubsubManager, srvs.gameSrv, srvs.userSrv, srvs.questionSrv)
	mm.SubscribeEventHandlers()

	presenceserver := grpcserver.New(cfg.Server.GrpcServer, srvs.presenceSrv)
	go presenceserver.Start()

	go srvs.websocketAdap.Run()

	handlers := []httpserver.Handler{
		pinghandler.New(),
		websockethandler.New(srvs.websocketAdap, srvs.authSrv, &cfg.Auth),
		gamehandler.New(srvs.gameSrv, srvs.presenceSrv, srvs.authSrv, cfg.Auth),
		userhandler.New(srvs.userSrv, srvs.authSrv, srvs.rbacSrv, srvs.presenceSrv, cfg.Auth, *srvs.userValidator),
		backpanelhandler.New(srvs.userSrv, srvs.rbacSrv, srvs.authSrv, cfg.Auth),
		matchinghandler.New(cfg.Auth, srvs.authSrv, *srvs.matchingSrv, *srvs.matchingValidator, srvs.presenceSrv),
		categoryhandler.New(srvs.categorySrv, srvs.presenceSrv, srvs.authSrv, cfg.Auth),
	}

	httpSvr := httpserver.New(cfg.Server.HttpServer, handlers)
	go httpSvr.Serve()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		scheduler := scheduler.New(cfg.Scheduler, srvs.matchingSrv)
		scheduler.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Gracefully shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Application.GracefulShutdownTimeout)
	defer cancel()

	if err := httpSvr.Router.Shutdown(ctx); err != nil {
		logger.L().Error("Could not shutdown http server.", zap.Error(err))
	}

	<-ctx.Done()

	wg.Wait()

	logger.L().Info("Shutdown services gracefully.")
}

func setupServices(cfg *config.Config) services {
	authConfig := cfg.Auth
	authSrv := authservice.New(authConfig)

	mongoCli := mongorepo.New(cfg.Database.MongoDB)

	rbacRepo := mongorbac.New(mongoCli)
	rbacSrv := rbacservice.New(rbacRepo)

	userRepo := mongouser.New(mongoCli)
	userSrv := userservice.New(&authSrv, userRepo)

	categoryRepo := mongocategory.New(mongoCli)
	redisAdap := redisadapter.New(cfg.Redis)

	questionRepo := mongoquestion.New(mongoCli)
	questionSrv := questionservice.New(cfg.Application.Question, questionRepo)

	pubsubManager := pubsub.NewPubSubManager(redisAdap)

	presenceRepo := redispresence.New(redisAdap)
	presenceSrv := presenceservice.New(cfg.Presence, presenceRepo)

	address := fmt.Sprintf(":%d", cfg.Server.GrpcServer.Port)
	presenceCli := presenceadapter.New(address)
	matchingRepo := redismatching.New(redisAdap)
	matchingSrv := matchingservice.New(cfg.Matching, matchingRepo, categoryRepo, presenceCli, pubsubManager)

	userValidator := uservalidator.New(userRepo)
	matchingValidator := matchingvalidator.New(categoryRepo)

	categorySrv := categoryservice.New(categoryRepo)

	gameRepo := mongogame.New(mongoCli)
	gameSrv := gameservice.New(gameRepo)

	fmt.Println(cfg.Server.Websocket)
	websocketAdap := websocketadapter.New(cfg.Server.Websocket)

	return services{
		questionSrv:       &questionSrv,
		authSrv:           &authSrv,
		userSrv:           &userSrv,
		matchingSrv:       &matchingSrv,
		categorySrv:       &categorySrv,
		presenceSrv:       &presenceSrv,
		rbacSrv:           &rbacSrv,
		gameSrv:           &gameSrv,
		websocketAdap:     websocketAdap,
		pubsubManager:     pubsubManager,
		userValidator:     &userValidator,
		matchingValidator: &matchingValidator,
	}
}

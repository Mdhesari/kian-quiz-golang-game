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
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/gamehandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/matchinghandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/pinghandler"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/handler/userhandler"
	"mdhesari/kian-quiz-golang-game/delivery/validator/matchingvalidator"
	"mdhesari/kian-quiz-golang-game/delivery/validator/uservalidator"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/protobufdecoder"
	"mdhesari/kian-quiz-golang-game/pkg/protobufencoder"
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
	"google.golang.org/grpc"

	_ "net/http/pprof"
)

var (
	cfg config.Config
)

type services struct {
	authSrv     *authservice.Service
	userSrv     *userservice.Service
	matchingSrv *matchingservice.Service
	categorySrv *categoryservice.Service
	presenceSrv *presenceservice.Service
	rbacSrv     *rbacservice.Service
	gameSrv     *gameservice.Service
}

func init() {
	cfg = config.Load("config.yml")
	flag.Parse()
}

func main() {
	logger.L().Info("Welcome to KianQuiz.")

	mongoCli, err := mongorepo.New(cfg.Database.MongoDB)
	if err != nil {

		panic("could not connect to mongodb.")
	}

	var srvs services = setupServices(&cfg, mongoCli)

	// TODO - Sepearte cmd for subscription
	go func() {
		gameRepo := mongogame.New(mongoCli)
		gameSrv := gameservice.New(gameRepo)

		questionRepo := mongoquestion.New(mongoCli)
		questionSrv := questionservice.New(questionRepo)
		redisAdap := redisadapter.New(cfg.Redis)
		subscriber := redisAdap.Cli().Subscribe(context.Background(), string(entity.UsersMatchedEvent))
		for {
			msg, err := subscriber.ReceiveMessage(context.Background())
			if err != nil {

				logger.L().Error("Redis sub: Could not recieve message.", zap.Error(err))
			}

			playersMatched := protobufdecoder.DecodeUsersMatchedEvent(msg.Payload)

			questionRes, err := questionSrv.GetRandomQuestions(context.Background(), param.QuestionGetRequest{
				CategoryId: playersMatched.Category.ID,
				Count:      cfg.Application.Game.QuestionsCount,
			})
			if err != nil {

				logger.L().Error("Could not get random questions for creating game.", zap.Error(err), zap.Any("Event", playersMatched))
			}

			game, err := gameSrv.Create(context.Background(), param.GameCreateRequest{
				Players:   playersMatched.Players,
				Category:  playersMatched.Category,
				Questions: questionRes.Items,
			})
			if err != nil {

				logger.L().Error(err.Error(), zap.Error(err), zap.Any("game", game))
			}

			logger.L().Info("A new game created.", zap.Any("game", game.Game.ID))

			payload := protobufencoder.EncodeGameStartedEvent(entity.GameStarted{
				PlayerIds: game.Game.PlayerIDs,
			})
			redisAdap.Publish(context.Background(), string(entity.GameStartedEvent), payload)
		}
	}()

	// TODO - Seperate cmd for presence server
	presenceserver := grpcserver.New(cfg.Server.GrpcServer, srvs.presenceSrv)
	go presenceserver.Start()

	userRepo := mongouser.New(mongoCli)
	userValidator := uservalidator.New(userRepo)

	categoryRepo := mongocategory.New(mongoCli)
	matchingValidator := matchingvalidator.New(categoryRepo)

	handlers := []httpserver.Handler{
		pinghandler.New(),
		gamehandler.New(srvs.gameSrv, srvs.presenceSrv, srvs.authSrv, cfg.Auth),
		userhandler.New(srvs.userSrv, srvs.authSrv, srvs.rbacSrv, srvs.presenceSrv, cfg.Auth, userValidator),
		backpanelhandler.New(srvs.userSrv, srvs.rbacSrv, srvs.authSrv, cfg.Auth),
		matchinghandler.New(cfg.Auth, srvs.authSrv, *srvs.matchingSrv, matchingValidator, srvs.presenceSrv),
		categoryhandler.New(srvs.categorySrv, srvs.presenceSrv, srvs.authSrv, cfg.Auth),
	}

	httpSvr := httpserver.New(cfg.Server.HttpServer, handlers)
	go httpSvr.Serve()

	// TODO - Separate cmd for schedule
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
	logger.L().Info("Gracefully shutting down http server...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Application.GracefulShutdownTimeout)
	defer cancel()
	if err := httpSvr.Router.Shutdown(ctx); err != nil {
		logger.L().Error("Could not shutdown http server.", zap.Error(err))
	}
	<-ctx.Done()

	wg.Wait()

	logger.L().Info("Shutdown services gracefully.")
}

func setupServices(cfg *config.Config, mongoCli *mongorepo.MongoDB) services {
	authConfig := cfg.Auth
	authSrv := authservice.New(authConfig)

	rbacRepo := mongorbac.New(mongoCli)
	rbacSrv := rbacservice.New(rbacRepo)

	userRepo := mongouser.New(mongoCli)
	userSrv := userservice.New(&authSrv, userRepo)

	categoryRepo := mongocategory.New(mongoCli)
	redisAdap := redisadapter.New(cfg.Redis)

	presenceRepo := redispresence.New(redisAdap)
	presenceSrv := presenceservice.New(cfg.Presence, presenceRepo)
	address := fmt.Sprintf(":%d", cfg.Server.GrpcServer.Port)
	grpConn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {

		log.Fatalf("Grpc could not dial %v\n", err)
	}
	presenceCli := presenceadapter.New(grpConn)
	matchingRepo := redismatching.New(redisAdap)
	matchingSrv := matchingservice.New(cfg.Matching, matchingRepo, categoryRepo, presenceCli, redisAdap)

	categorySrv := categoryservice.New(categoryRepo)

	gameRepo := mongogame.New(mongoCli)
	gameSrv := gameservice.New(gameRepo)

	return services{
		authSrv:     &authSrv,
		userSrv:     &userSrv,
		matchingSrv: &matchingSrv,
		categorySrv: &categorySrv,
		presenceSrv: &presenceSrv,
		rbacSrv:     &rbacSrv,
		gameSrv:     &gameSrv,
	}
}

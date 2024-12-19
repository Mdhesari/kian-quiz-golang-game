package main

import (
	"flag"
	"log"
	"mdhesari/kian-quiz-golang-game/adapter/presenceadapter"
	"mdhesari/kian-quiz-golang-game/adapter/redisadapter"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongocategory"
	"mdhesari/kian-quiz-golang-game/repository/redisrepo/redismatching"
	"mdhesari/kian-quiz-golang-game/scheduler"
	"mdhesari/kian-quiz-golang-game/service/matchingservice"
	"sync"

	"google.golang.org/grpc"
)

var (
	cfg config.Config
)

func init() {
	cfg = config.Load("config.yml")
	flag.Parse()
}

func main() {
	var wg sync.WaitGroup

	redisAdap := redisadapter.New(cfg.Redis)
	matchingRepo := redismatching.New(redisAdap)

	cli, err := mongorepo.New(cfg.Database.MongoDB)
	if err != nil {

		panic("could not connect to mongodb.")
	}
	mongocategory := mongocategory.New(cli)

	grpConn, err := grpc.Dial(":8089", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Grpc could not dial %v\n", err)
	}
	presenceCli := presenceadapter.New(grpConn)
	matchingSrv := matchingservice.New(cfg.Matching, matchingRepo, mongocategory, presenceCli, redisAdap)

	scheduler := scheduler.New(cfg.Scheduler, &matchingSrv)

	wg.Add(1)
	go func() {
		defer wg.Done()

		scheduler.Start()
	}()

	wg.Wait()
}

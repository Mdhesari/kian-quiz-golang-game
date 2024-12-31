package main

import (
	"flag"
	"fmt"
	"mdhesari/kian-quiz-golang-game/adapter/presenceadapter"
	"mdhesari/kian-quiz-golang-game/adapter/redisadapter"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo/mongocategory"
	"mdhesari/kian-quiz-golang-game/repository/redisrepo/redismatching"
	"mdhesari/kian-quiz-golang-game/scheduler"
	"mdhesari/kian-quiz-golang-game/service/matchingservice"
	"sync"
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

	cli := mongorepo.New(cfg.Database.MongoDB)
	mongocategory := mongocategory.New(cli)

	address := fmt.Sprintf(":%d", cfg.Server.GrpcServer.Port)
	presenceCli := presenceadapter.New(address)
	matchingSrv := matchingservice.New(cfg.Matching, matchingRepo, mongocategory, presenceCli, redisAdap)

	scheduler := scheduler.New(cfg.Scheduler, &matchingSrv)

	wg.Add(1)
	go func() {
		defer wg.Done()

		scheduler.Start()
	}()

	wg.Wait()
}

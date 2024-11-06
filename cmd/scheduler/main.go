package main

import (
	"flag"
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
	cfg    config.Config
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

	matchingSrv := matchingservice.New(matchingRepo, mongocategory)

	scheduler := scheduler.New(cfg.Scheduler, &matchingSrv)

	wg.Add(1)
	go scheduler.Start(&wg)

	wg.Wait()
}
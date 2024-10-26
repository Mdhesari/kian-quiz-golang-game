package main

import (
	"flag"
	"mdhesari/kian-quiz-golang-game/adapter/redisadapter"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/delivery/httpserver"
	"mdhesari/kian-quiz-golang-game/repository/redisrepo/redismatching"
	"mdhesari/kian-quiz-golang-game/scheduler"
	"mdhesari/kian-quiz-golang-game/service/matchingservice"
	"sync"
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
	var wg sync.WaitGroup
	
	redisAdap := redisadapter.New(cfg.Redis)
	matchingRepo := redismatching.New(redisAdap)
	matchingSrv := matchingservice.New(matchingRepo)
	scheduler := scheduler.New(&matchingSrv)

	wg.Add(1)
	go scheduler.Start(&wg)

	wg.Wait()
}
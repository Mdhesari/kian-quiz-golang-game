package main

import (
	"context"
	"flag"
	"fmt"
	"mdhesari/kian-quiz-golang-game/adapter/redisadapter"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/protobufdecoder"
)

var cfg config.Config

func init() {
	cfg = config.Load("config.yml")
	flag.Parse()
}

func main() {
	redisAdap := redisadapter.New(cfg.Redis)
	subscriber := redisAdap.Cli().Subscribe(context.Background(), string(entity.UsersMatched))

	for {
		msg, err := subscriber.ReceiveMessage(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		
		playersMatched := protobufdecoder.DecodeUsersMatchedEvent(msg.Payload)
		fmt.Println("new message %v\n", playersMatched)
	}
}

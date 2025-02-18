package main

import (
	"context"
	"flag"
	"mdhesari/kian-quiz-golang-game/adapter/redisadapter"
	"mdhesari/kian-quiz-golang-game/config"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/pkg/protobufdecoder"

	"go.uber.org/zap"
)

var cfg config.Config

func init() {
	cfg = config.Load("config.yml")
	flag.Parse()
}

func main() {
	redisAdap := redisadapter.New(cfg.Redis)
	subscriber := redisAdap.Cli().Subscribe(context.Background(), string(entity.PlayersMatchedEvent))

	for {
		msg, err := subscriber.ReceiveMessage(context.Background())
		if err != nil {
			logger.L().Error("Could not recieve pusub msg.", zap.Error(err))
		}

		playersMatched := protobufdecoder.DecodePlayersMatchedEvent(msg.Payload)
		logger.L().Info("Players matched.", zap.Any("playersMatched", playersMatched))
	}
}

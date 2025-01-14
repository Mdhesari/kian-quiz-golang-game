package main

import (
	"context"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/protobuf/golang/presence"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8089", grpc.WithInsecure())
	if err != nil {

		logger.L().Error("Could not dial grpc.")
	}
	defer conn.Close()

	cli := presence.NewPresenceServiceClient(conn)
	_, err = cli.GetPresence(context.Background(), &presence.GetPresenceRequest{
		UserId: []string{"000000000000000000000000"},
	})
	if err != nil {
		logger.L().Error("Clould not get presence: %v\n", zap.Error(err))
	}
}

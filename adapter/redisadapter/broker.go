package redisadapter

import (
	"context"
	"mdhesari/kian-quiz-golang-game/logger"

	"go.uber.org/zap"
)

func (a Adapter) Publish(ctx context.Context, topic string, payload string) {
	logger.L().Info("Publishing topic to queue.", zap.String("topic", topic))

	res, err := a.cli.Publish(ctx, topic, payload).Result()
	if err != nil {

		panic(res)
	}

	logger.L().Info("Topic published.", zap.Int64("result", res))
}

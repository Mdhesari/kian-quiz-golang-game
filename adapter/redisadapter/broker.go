package redisadapter

import (
	"context"
	"mdhesari/kian-quiz-golang-game/logger"

	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
)

func (a Adapter) Publish(ctx context.Context, topic string, payload string) {
	logger.L().Info("Publishing topic to queue.", zap.String("topic", topic))

	res, err := a.cli.Publish(ctx, topic, payload).Result()
	if err != nil {

		panic(res)
	}

	logger.L().Info("Topic published.", zap.Int64("result", res))
}

func (a Adapter) Subscribe(ctx context.Context, topics ...string) *redis.PubSub {
	logger.L().Info("Subscribing to queue.", zap.Any("topics", topics))

	return a.cli.Subscribe(ctx, topics...)
}

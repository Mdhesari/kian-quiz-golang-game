package pubsub

import (
	"context"

	"mdhesari/kian-quiz-golang-game/adapter/redisadapter"
	"mdhesari/kian-quiz-golang-game/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type EventHandler func(ctx context.Context, topic string, payload string) error

type PubSubManager struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewPubSubManager(redisAdap redisadapter.Adapter) *PubSubManager {
	return &PubSubManager{
		redisClient: redisAdap.Cli(),
		ctx:         context.Background(),
	}
}

func (p *PubSubManager) Subscribe(topic string, handler EventHandler) {
	go func() {
		subscriber := p.redisClient.Subscribe(p.ctx, topic)
		ch := subscriber.Channel()

		for msg := range ch {
			if err := handler(p.ctx, topic, msg.Payload); err != nil {
				logger.L().Error("Could not call subscribed handler.", zap.Error(err), zap.String("topic", topic), zap.String("payload", msg.Payload))
			}
		}
	}()
}

func (p *PubSubManager) Publish(ctx context.Context, topic string, payload string) {
	if err := p.redisClient.Publish(ctx, topic, payload).Err(); err != nil {
		logger.L().Error("Pubsub: could not publish.", zap.Error(err), zap.String("topic", topic), zap.String("payload", payload))
	}
}

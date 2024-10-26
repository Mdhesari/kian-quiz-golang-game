package redismatching

import (
	"context"
	"fmt"
	"mdhesari/kian-quiz-golang-game/pkg/timestamp"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db DB) AddToWaitingList(ctx context.Context, userId primitive.ObjectID, categoryId primitive.ObjectID) error {
	scoreKey := fmt.Sprintf("%s:%s", "scores", categoryId)
	_, err := db.adapter.Cli().ZAdd(ctx, scoreKey, redis.Z{
		Score:  float64(timestamp.Now()),
		Member: fmt.Sprintf(`%s`, userId),
	}).Result()
	if err != nil {

		return err
	}

	return nil
}


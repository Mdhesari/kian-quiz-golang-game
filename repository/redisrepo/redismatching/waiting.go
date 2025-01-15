package redismatching

import (
	"context"
	"fmt"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/timestamp"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *DB) AddToWaitingList(ctx context.Context, userId primitive.ObjectID, category entity.Category) error {
	categoryKey := db.getCategoryKey(category)
	_, err := db.adapter.Cli().ZAdd(ctx, categoryKey, redis.Z{
		Score:  float64(timestamp.Now()),
		Member: userId.Hex(),
	}).Result()
	if err != nil {

		return err
	}

	return nil
}

func (db *DB) GetWaitingListByCategory(ctx context.Context, category entity.Category, maxWaitingTime time.Duration) ([]entity.WaitingMember, error) {
	waitingMembers := []entity.WaitingMember{}

	min := timestamp.Add(-1 * maxWaitingTime)
	max := timestamp.Now()

	categoryKey := db.getCategoryKey(category)
	list, err := db.adapter.Cli().ZRangeByScoreWithScores(ctx, categoryKey, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", min),
		Max: fmt.Sprintf("%d", max),
	}).Result()
	if err != nil {

		return waitingMembers, err
	}

	for _, item := range list {
		member := item.Member.(string)
		userId, err := primitive.ObjectIDFromHex(member)
		if err != nil {

			return waitingMembers, err
		}

		waitingMembers = append(waitingMembers, entity.WaitingMember{
			UserId:    userId,
			Category:  category,
			Timestamp: int64(item.Score),
		})
	}

	return waitingMembers, nil
}

func (db *DB) RemoveUsersFromWaitingList(ctx context.Context, category entity.Category, userIds []string) error {
	categoryKey := db.getCategoryKey(category)
	// TODO - do we need to check deleted count?
	_, err := db.adapter.Cli().ZRem(ctx, categoryKey, userIds).Result()
	if err != nil {

		return err
	}

	return nil
}

func (db *DB) getCategoryKey(category entity.Category) string {
	return fmt.Sprintf("%s:%s", db.waitingListPrefix, category.ID.Hex())
}

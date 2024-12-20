package redismatching

import (
	"context"
	"fmt"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/timestamp"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const WaitingListPrefix = "waitinglist"

func (db DB) AddToWaitingList(ctx context.Context, userId primitive.ObjectID, category entity.Category) error {
	categoryKey := getCategoryKey(category)
	_, err := db.adapter.Cli().ZAdd(ctx, categoryKey, redis.Z{
		Score:  float64(timestamp.Now()),
		Member: userId.Hex(),
	}).Result()
	if err != nil {

		return err
	}

	return nil
}

func (db DB) GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error) {
	waitingMembers := []entity.WaitingMember{}

	categoryKey := getCategoryKey(category)
	list, err := db.adapter.Cli().ZRangeByScoreWithScores(ctx, categoryKey, &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
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

func (db DB) RemoveUsersFromWaitingList(ctx context.Context, category entity.Category, userIds []string) error {
	categoryKey := getCategoryKey(category)
	// TODO - do we need to check deleted count?
	_, err := db.adapter.Cli().ZRem(ctx, categoryKey, userIds).Result()
	if err != nil {

		return err
	}

	return nil
}

func getCategoryKey(category entity.Category) string {
	return fmt.Sprintf("%s:%s", WaitingListPrefix, category.ID.Hex())
}

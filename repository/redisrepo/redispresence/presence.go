package redispresence

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db DB) GetPresence(ctx context.Context, prefix string, userIds []primitive.ObjectID) (map[primitive.ObjectID]int64, error) {
	var keys []string
	presenceList := make(map[primitive.ObjectID]int64, len(userIds))

	for _, userId := range userIds {
		keys = append(keys, fmt.Sprintf("%s:%s", prefix, userId.Hex()))
	}
	res, err := db.adapter.Cli().MGet(ctx, keys...).Result()
	if err != nil {

		return presenceList, err
	}

	for i, userId := range userIds {
		if res[i] == nil {

			continue
		}

		if val, ok := res[i].(string); ok {
			presenceList[userId], _ = strconv.ParseInt(val, 10, 64)
		}
	}

	return presenceList, nil
}

func (db DB) Upsert(ctx context.Context, key string, timestamp int64, exp time.Duration) error {
	_, err := db.adapter.Cli().Set(ctx, key, timestamp, exp).Result()
	if err != nil {

		return err
	}

	return nil
}

package redisleaderboard

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/pkg/mongoutils"

	"go.uber.org/zap"
)

const (
	LeaderboardKey = "leaderboard"
)

func (db *DB) GetLeaderboard(ctx context.Context, limit int) ([]entity.UserScore, error) {
	lb, err := db.adapter.Cli().ZRevRangeWithScores(context.Background(), LeaderboardKey, 0, int64(limit-1)).Result()
	if err != nil {

		return nil, err
	}

	if len(lb) < 1 {

		return []entity.UserScore{}, nil
	}

	var userIDs []string
	for _, us := range lb {
		userIDs = append(userIDs, us.Member.(string))
	}

	var userData []interface{}
	userData, err = db.adapter.Cli().HMGet(ctx, "user:data", userIDs...).Result()
	if err != nil {

		return nil, err
	}

	logger.L().Info("User data retrieved.", zap.Any("user_data", userData))

	leaderboard := make([]entity.UserScore, 0)
	var displayName string
	var ok bool
	for i, entry := range lb {
		displayName, ok = userData[i].(string)
		if ok {
			leaderboard = append(leaderboard, entity.UserScore{
				UserId:      mongoutils.HexToObjectID(entry.Member.(string)),
				DisplayName: displayName,
				Score:       entity.Score(entry.Score),
			})
		}
	}

	logger.L().Info("Leaderboard retrieved.", zap.Any("leaderboard", len(lb)))

	return leaderboard, nil
}

func (db *DB) UpsertLeaderboardUserScore(ctx context.Context, us entity.UserScore) error {
	err := db.adapter.Cli().ZIncrBy(ctx, LeaderboardKey, float64(us.Score), us.UserId.Hex()).Err()
	if err != nil {

		return err
	}

	if err = db.adapter.Cli().HSet(ctx, "user:data", us.UserId.Hex(), us.DisplayName).Err(); err != nil {

		return err
	}

	return nil
}

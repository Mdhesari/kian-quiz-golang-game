package redisleaderboard

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/mongoutils"
)

func (db *DB) GetLeaderboard(ctx context.Context, limit int) ([]entity.UserScore, error) {
	lb, err := db.adapter.Cli().ZRevRangeWithScores(context.Background(), "leaderboard", 0, int64(limit-1)).Result()
	if err != nil {

		return nil, err
	}

	leaderboard := make([]entity.UserScore, len(lb))
	for i, entry := range lb {
		leaderboard[i].UserId = mongoutils.HexToObjectID(entry.Member.(string))
		leaderboard[i].Score = entity.Score(entry.Score)
	}

	return leaderboard, nil
}

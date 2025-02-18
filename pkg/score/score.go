package score

import (
	"math"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"time"

	"go.uber.org/zap"
)

func ClaculateScore(sMax entity.Score, t time.Duration, tMax time.Duration) entity.Score {
	var current, max float64 = t.Seconds(), tMax.Seconds()
	logger.L().Info("zap", zap.Any("current", current), zap.Any("max", max))
	if current > max {

		return 0
	}
	if current < 4 {

		return entity.Score(float64(sMax))
	}

	du := 1 - current/max

	return entity.Score(math.Ceil(float64(sMax) * du))
}

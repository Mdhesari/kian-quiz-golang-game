package score

import (
	"math"
	"mdhesari/kian-quiz-golang-game/entity"
	"time"
)

func ClaculateScore(sMax entity.Score, t time.Duration, tMax time.Duration) float64 {
	var current, max float64 = t.Seconds(), tMax.Seconds()
	if current > max {

		return 0
	}
	if current < 4 {

		return float64(sMax)
	}

	du := 1 - current/max

	return math.Ceil(float64(sMax) * du)
}

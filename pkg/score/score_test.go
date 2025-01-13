package score

import (
	"mdhesari/kian-quiz-golang-game/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		Expected float64
		MaxScore entity.Score
		TimeTook time.Duration
		TimeMax  time.Duration
	}{
		{
			Expected: 1,
			MaxScore: 5,
			TimeTook: time.Second * 8,
			TimeMax:  time.Second * 10,
		},
		{
			Expected: 0,
			MaxScore: 5,
			TimeTook: time.Second * 16,
			TimeMax:  time.Second * 15,
		},
		{
			Expected: 1,
			MaxScore: 5,
			TimeTook: time.Second * 14,
			TimeMax:  time.Second * 15,
		},
		{
			Expected: 4,
			MaxScore: 5,
			TimeTook: time.Second * 5,
			TimeMax:  time.Second * 15,
		},
		{
			Expected: 5,
			MaxScore: 5,
			TimeTook: time.Second * 3,
			TimeMax:  time.Second * 15,
		},
		{
			Expected: 0,
			MaxScore: 5,
			TimeTook: time.Second * 100,
			TimeMax:  time.Second * 15,
		},
	}

	t.Run("Test CalculateScore", func(t *testing.T) {
		for _, test := range tests {
			assert.Equal(t, test.Expected, ClaculateScore(test.MaxScore, test.TimeTook, test.TimeMax))
		}
	})
}

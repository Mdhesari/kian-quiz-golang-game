package param

import (
	"mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserIncrementScoreRequest struct {
	UserId primitive.ObjectID `json:"user_id"`
	Score  entity.Score       `json:"score"`
}

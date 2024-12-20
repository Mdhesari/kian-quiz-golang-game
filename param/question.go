package param

import (
	"mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionGetRequest struct {
	CategoryId primitive.ObjectID `json:"category_id"`
	Count      int                `json:"count"`
}

type QuestionGetResponse struct {
	Items []entity.Question `json:"questions"`
}

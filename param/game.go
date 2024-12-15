package param

import (
	"mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameCreateRequest struct {
	Players  []primitive.ObjectID `json:"players"`
	Category entity.Category      `json:"category"`
}

type GameCreateResponse struct {
	Game entity.Game `json:"game"`
}

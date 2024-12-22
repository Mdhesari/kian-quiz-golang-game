package param

import (
	"mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameCreateRequest struct {
	PlayerIDs []primitive.ObjectID `json:"player_ids"`
	Questions []entity.Question    `json:"questions"`
	Category  entity.Category      `json:"category"`
}

type GameCreateResponse struct {
	Game entity.Game `json:"game"`
}

type GameGetRequest struct {
	GameId primitive.ObjectID `param:"game_id"`
}

type GameGetResponse struct {
	Game entity.Game `json:"game"`
}

package param

import (
	"mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameGetAllRequest struct {
	UserID     primitive.ObjectID `json:"user_id"`
	CategoryID primitive.ObjectID `json:"category_id,omitempty"`
}

type GameGetAllResponse struct {
	Items []entity.Game `json:"games"`
}

type GameCreateRequest struct {
	Players   []entity.Player   `json:"players,omitempty"`
	Questions []entity.Question `json:"questions,omitempty"`
	Category  entity.Category   `json:"category,omitempty"`
}

type GameCreateResponse struct {
	Game entity.Game `json:"game"`
}

type GameGetRequest struct {
	GameId primitive.ObjectID `param:"game_id,omitempty"`
}

type GameGetResponse struct {
	Game entity.Game `json:"game"`
}

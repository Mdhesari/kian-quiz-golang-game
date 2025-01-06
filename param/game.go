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
	Players   map[string]entity.Player `json:"players,omitempty"`
	Questions []entity.Question        `json:"questions,omitempty"`
	Category  entity.Category          `json:"category,omitempty"`
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

type GameAnswerQuestionRequest struct {
	UserId       primitive.ObjectID  `json:"user_id"`
	GameId       primitive.ObjectID  `json:"game_id"`
	PlayerAnswer entity.PlayerAnswer `json:"player_answer"`
}

type GameAnswerQuestionResponse struct {
	PlayerAnswer entity.PlayerAnswer `json:"player_answer"`
}

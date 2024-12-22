package param

import (
	"mdhesari/kian-quiz-golang-game/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameCreateRequest struct {
	PlayerIDs []primitive.ObjectID `json:"player_ids,omitempty"`
	Questions []entity.Question    `json:"questions,omitempty"`
	Category  entity.Category      `json:"category,omitempty"`
}

type GameUpdateRequest struct {
	ID          primitive.ObjectID   `json:"id"`
	PlayerIDs   []primitive.ObjectID `json:"player_ids,omitempty"`
	Category    entity.Category      `json:"category,omitempty"`
	QuestionIDs []primitive.ObjectID `json:"question_ids,omitempty"`
	StartTime   time.Time            `json:"start_time,omitempty"`
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

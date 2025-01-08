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
	GameId primitive.ObjectID `json:"game_id,omitempty" param:"game_id"`
}

type GameGetResponse struct {
	Game entity.Game `json:"game"`
}

type GameAnswerQuestionRequest struct {
	UserId     primitive.ObjectID `json:"user_id"`
	GameId     primitive.ObjectID `json:"game_id" param:"game_id"`
	QuestionId primitive.ObjectID `json:"question_id" form:"question_id"`
	Answer     string             `json:"answer" form:"answer"`
}

type GameAnswerQuestionResponse struct {
	Answer        entity.PlayerAnswer `json:"answer"`
	CorrectAnswer entity.Answer       `json:"correct_answer"`
}

type GameGetNextQuestionRequest struct {
	UserId primitive.ObjectID `json:"user_id"`
	GameId primitive.ObjectID `json:"game_id" param:"game_id"`
}

type GameGetNextQuestionResponse struct {
	Question entity.Question `json:"question"`
}

type GameFinishRequest struct {
	GameId primitive.ObjectID `json:"game_id" param:"game_id"`
}

type GameFinishResponse struct {
	//
}

package param

import (
	"mdhesari/kian-quiz-golang-game/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayerCreateRequest struct {
	UserID    primitive.ObjectID `json:"user_id"`
	GameID    primitive.ObjectID `json:"game_id"`
	CreatedAt time.Time          `json:"created_at"`
}

type PlayerCreateResponse struct {
	Player entity.Player `json:"player"`
}

type PlayerGetRequest struct {
	ID primitive.ObjectID `json:"id"`
}

type PlayerGetResponse struct {
	Player entity.Player `json:"player"`
}

type PlayerUpdateRequest struct {
	ID        primitive.ObjectID    `json:"id"`
	Answers   []entity.PlayerAnswer `json:"answers"`
	Score     int                   `json:"score"`
	UpdatedAt time.Time             `json:"updated_at"`
}

type PlayerDeleteRequest struct {
	ID primitive.ObjectID `json:"id"`
}

type PlayerStatusUpdateRequest struct {
	GameId primitive.ObjectID  `json:"game_id"`
	UserId primitive.ObjectID  `json:"user_id"`
	Status entity.PlayerStatus `json:"status"`
}

type PlayerStatusUpdateResponse struct {
	//
}

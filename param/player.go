package param

import (
	"mdhesari/kian-quiz-golang-game/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayerCreateRequest struct {
	UserID    primitive.ObjectID
	GameID    primitive.ObjectID
	CreatedAt time.Time
}

type PlayerCreateResponse struct {
	Player entity.Player
}

type PlayerGetRequest struct {
	ID primitive.ObjectID
}

type PlayerGetResponse struct {
	Player entity.Player
}

type PlayerUpdateRequest struct {
	ID        primitive.ObjectID
	Answers   []entity.PlayerAnswer
	Score     int
	UpdatedAt time.Time
}

type PlayerDeleteRequest struct {
	ID primitive.ObjectID
}

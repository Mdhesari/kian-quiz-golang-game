package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Game struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty"`
	CategoryID     primitive.ObjectID   `bson:"category_id,omitempty"`
	QuestionIDs    []primitive.ObjectID `bson:"question_ids,omitempty"`
	PlayerIDs      []primitive.ObjectID `bson:"player_ids,omitempty"`
	WinnerPlayerID primitive.ObjectID   `bson:"winner_player_id,omitempty"`
	StartTime      time.Time            `bson:"start_time"`
	EndTime        time.Time            `bson:"end_time"`
	CreatedAt      time.Time            `bson:"created_at"`
	UpdatedAt      time.Time            `bson:"updated_at"`
}

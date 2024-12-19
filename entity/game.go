package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Game struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	CategoryID  primitive.ObjectID   `bson:"category_id"`
	QuestionIDs []primitive.ObjectID `bson:"question_ids"`
	PlayerIDs   []primitive.ObjectID `bson:"player_ids"`
	WinnerID    primitive.ObjectID   `bson:"winner_id"`
	StartTime   time.Time            `bson:"start_time"`
	ExpiresTime time.Time            `bson:"expires_time"`
}

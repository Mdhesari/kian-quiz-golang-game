package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Round struct {
	QuestionId primitive.ObjectID `bson:"question_id"`
	PlayerId   primitive.ObjectID `bson:"player_id"`
	Answer     Answer             `bson:"answer"`
	StartTime  time.Time          `bson:"start_time"`
	EndTime    time.Time          `bson:"end_time"`
}

type Game struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	CategoryID  primitive.ObjectID   `bson:"category_id"`
	QuestionIDs []primitive.ObjectID `bson:"question_ids"`
	PlayerIDs   []primitive.ObjectID `bson:"player_ids"`
	Rounds      []Round              `bson:"rounds"`
	WinnerID    primitive.ObjectID   `bson:"winner_id"`
	StartTime   time.Time            `bson:"start_time"`
	ExpiresTime time.Time            `bson:"expires_time"`
}

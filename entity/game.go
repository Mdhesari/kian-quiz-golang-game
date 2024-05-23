package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Game struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CategoryID  primitive.ObjectID
	QuestionIDs []primitive.ObjectID
	PlayerIDs   []primitive.ObjectID
	WinnerID    primitive.ObjectID
	StartTime   time.Time
	ExpiresTime time.Time
}

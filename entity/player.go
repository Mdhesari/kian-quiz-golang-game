package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayerAnswer struct {
	QuestionId primitive.ObjectID `bson:"question_id"`
	Answer     Answer             `bson:"answer"`
	Score      int                `bson:"score"`
	StartTime  time.Time          `bson:"start_time"`
	EndTime    time.Time          `bson:"end_time"`
}

type Player struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	GameID    primitive.ObjectID `bson:"game_id"`
	Answers   []PlayerAnswer     `bson:"answers"`
	Score     int                `bson:"score"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Player struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	UserID    primitive.ObjectID   `bson:"user_id"`
	GameID    primitive.ObjectID   `bson:"game_id"`
	AnswerIDs []primitive.ObjectID `bson:"answer_ids"`
	Score     int                  `bson:"score"`
}

type PlayerAnswer struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	PlayerID primitive.ObjectID `bson:"player_id"`
	AnswerID primitive.ObjectID `bson:"answer_id"`
}

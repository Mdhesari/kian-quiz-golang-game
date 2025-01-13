package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type PlayerFinished struct {
	UserId primitive.ObjectID `json:"user_id"`
	GameId primitive.ObjectID `json:"game_id"`
}

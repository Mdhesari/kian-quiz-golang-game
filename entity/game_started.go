package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type GameStarted struct {
	GameID primitive.ObjectID `json:"game_id"`
}

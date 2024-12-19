package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type GameStarted struct {
	PlayerIds []primitive.ObjectID `json:"player_ids"`
}

package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type PlayersMatched struct {
	PlayerIDs []primitive.ObjectID `json:"player_ids"`
	Category  Category             `json:"category"`
}
